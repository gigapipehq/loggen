package progress

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/schollz/progressbar/v3"
)

type Bar struct {
	*progressbar.ProgressBar
}

type Server struct {
	current   int
	max       int
	ch        chan struct{}
	writeFunc func(*bufio.Writer)
}

func (b *Bar) Add(count int) {
	_ = b.ProgressBar.Add(count)
}

func (p *Server) marshal() []byte {
	return []byte(fmt.Sprintf("{\"current\":%d,\"max\":%d}", p.current, p.max))
}

func (p *Server) Add(count int) {
	p.current += count
	if p.current == p.max {
		defer close(p.ch)
	}
	p.ch <- struct{}{}
}

func (p *Server) WriteProgress(w *bufio.Writer) {
	p.writeFunc(w)
}

func NewBar(max int, writer io.Writer) *Bar {
	return &Bar{
		progressbar.NewOptions(max,
			progressbar.OptionSetDescription("Batches sent"),
			progressbar.OptionSetWriter(writer),
			progressbar.OptionSetWidth(10),
			progressbar.OptionThrottle(65*time.Millisecond),
			progressbar.OptionShowCount(),
			progressbar.OptionOnCompletion(func() {
				_, _ = fmt.Fprint(writer, "\nAll batches sent. Waiting for all requests to complete...\n")
			}),
			progressbar.OptionSpinnerType(14),
			progressbar.OptionFullWidth(),
			progressbar.OptionSetRenderBlankState(true),
		),
	}
}

func NewServer(max int) *Server {
	s := &Server{
		current: 0,
		max:     max,
		ch:      make(chan struct{}),
	}
	s.writeFunc = func(w *bufio.Writer) {
		for range s.ch {
			_, _ = w.Write(s.marshal())
			_ = w.Flush()
		}
	}
	return s
}
