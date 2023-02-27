package senders

import (
	"context"
	"log"
	"math"
	"sync"
	"time"

	"github.com/gigapipehq/loggen/internal/prom"
)

type Sender interface {
	Send(batch []byte) error
	AddProgress(int)
	Progress() <-chan int
}

type Generator interface {
	Generate() ([]byte, error)
	Rate() int
}

func Start(ctx context.Context, sender Sender, generator Generator) {
	deadline, _ := ctx.Deadline()
	batchMax := int(math.Ceil(time.Until(deadline).Seconds()))
	batchChannel := make(chan []byte, 5)

	go func() {
		batchesCreated := 0
		for {
			select {
			case <-ctx.Done():
				return
			default:
				batch, err := generator.Generate()
				if err != nil {
					log.Printf("Error generating batch: %v", err)
					continue
				}
				batchChannel <- batch
				batchesCreated++
				if batchesCreated >= batchMax {
					return
				}
			}
		}
	}()

	t := time.NewTicker(time.Second)
	wg := &sync.WaitGroup{}
	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return
		case <-t.C:
			batch := <-batchChannel
			go func() {
				go func() {
					sender.AddProgress(generator.Rate())
				}()
				wg.Add(1)
				defer wg.Done()

				prom.AddLines(generator.Rate())
				prom.AddBytes(len(batch))

				if err := sender.Send(batch); err != nil {
					prom.AddErrors(1)
					log.Printf("Error sending request: %v", err)
					return
				}
			}()
		}
	}
}
