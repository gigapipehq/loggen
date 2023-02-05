package senders

import (
	"context"
	"log"
	"sync"
	"time"

	"go.opentelemetry.io/otel/codes"

	"github.com/gigapipehq/loggen/internal/otel"
	"github.com/gigapipehq/loggen/internal/prom"
)

type Sender interface {
	Send(ctx context.Context, batch []byte) error
	AddProgress(int)
	Progress() <-chan int
}

type Generator interface {
	Generate(ctx context.Context) ([]byte, error)
	Rate() int
}

func Start(ctx context.Context, sender Sender, generator Generator) {
	batchChannel := make(chan []byte, 5)
	go func() {
		gctx, span := otel.Tracer.Start(ctx, "start generating")
		defer span.End()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				lctx, span := otel.Tracer.Start(gctx, "generate new batch")
				batch, err := generator.Generate(lctx)
				if err != nil {
					log.Printf("Error generating batch: %v", err)
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
					continue
				}
				span.End()
				batchChannel <- batch
			}
		}
	}()

	t := time.NewTicker(time.Second)
	wg := &sync.WaitGroup{}
	sctx, span := otel.Tracer.Start(ctx, "start sending")
	defer span.End()
	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return
		case <-t.C:
			batch := <-batchChannel
			go func() {
				lctx, span := otel.Tracer.Start(sctx, "receive new batch")
				defer span.End()

				go func() {
					sender.AddProgress(generator.Rate())
				}()
				wg.Add(1)
				defer wg.Done()

				prom.AddLines(generator.Rate())
				prom.AddBytes(len(batch))

				if err := sender.Send(lctx, batch); err != nil {
					prom.AddErrors(1)
					log.Printf("Error sending request: %v", err)
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
					return
				}
			}()
		}
	}
}
