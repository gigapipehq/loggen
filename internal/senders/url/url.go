package url

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/gigapipehq/loggen/internal/otel"
)

type roundTripper struct {
	headers http.Header
	url     *url.URL
}

func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header = rt.headers
	return http.DefaultTransport.RoundTrip(req)
}

type URLSender struct {
	httpClient *http.Client
	transport  *roundTripper
	progressCh chan int
}

func (s *URLSender) Client() *http.Client {
	return s.httpClient
}

func (s *URLSender) WithHeaders(headers map[string]string) *URLSender {
	httpHeaders := http.Header{}
	for k, v := range headers {
		httpHeaders.Add(k, v)
	}
	if s.transport == nil {
		s.transport = &roundTripper{}
	}
	s.transport.headers = httpHeaders
	s.httpClient.Transport = s.transport
	return s
}

func (s *URLSender) AddProgress(count int) {
	s.progressCh <- count
}

func (s *URLSender) Progress() <-chan int {
	return s.progressCh
}

func (s *URLSender) Send(batch []byte) error {
	httpMethod := "POST"
	reqSize := int64(len(batch))
	body := io.NopCloser(bytes.NewReader(batch))
	req := &http.Request{
		URL:           s.transport.url,
		Method:        httpMethod,
		Body:          body,
		ContentLength: reqSize,
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode >= http.StatusBadRequest {
		rbody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d received with error: %s", resp.StatusCode, string(rbody))
	}
	return nil
}

func (s *URLSender) SupportsMetrics() bool {
	return true
}

func (s *URLSender) TracesExporter() sdktrace.SpanExporter {
	return otel.NewZipkinExporter(s.transport.url.String(), s.httpClient)
}

func (s *URLSender) Close() error {
	s.httpClient.CloseIdleConnections()
	close(s.progressCh)
	return nil
}

func New(httpURL *url.URL) *URLSender {
	return &URLSender{
		httpClient: &http.Client{},
		transport: &roundTripper{
			url: httpURL,
		},
		progressCh: make(chan int),
	}
}
