package file

import (
	"fmt"
	"io"
	"log"
	"os"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/gigapipehq/loggen/internal/otel"
)

type FileSender struct {
	file       *os.File
	progressCh chan int
}

func (s *FileSender) Send(batch []byte) error {
	batch = append(batch, []byte(",")...)
	if _, err := s.file.Write(batch); err != nil {
		return err
	}
	return nil
}

func (s *FileSender) AddProgress(count int) {
	s.progressCh <- count
}

func (s *FileSender) Progress() <-chan int {
	return s.progressCh
}

func (s *FileSender) SupportsMetrics() bool {
	return false
}

func (s *FileSender) TracesExporter() sdktrace.SpanExporter {
	f, err := os.Create(fmt.Sprintf("%s.trace.json", s.file.Name()))
	if err != nil {
		log.Printf("Unable to create file for traces export: %v. Discarding traces", err)
		return otel.NewDiscardExporter()
	}
	return otel.NewFileExporter(f)
}

func (s *FileSender) Close() error {
	close(s.progressCh)
	if _, err := s.file.Seek(-1, io.SeekEnd); err != nil {
		return fmt.Errorf("unable to seek last file byte: %v", err)
	}
	if _, err := s.file.Write([]byte("]")); err != nil {
		return fmt.Errorf("unable to write final byte to file: %v", err)
	}
	return s.file.Close()
}

func New(filename string) (*FileSender, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	if _, err := f.Write([]byte("[")); err != nil {
		return nil, fmt.Errorf("unable to write initial byte to file: %v", err)
	}
	return &FileSender{
		file:       f,
		progressCh: make(chan int),
	}, nil
}
