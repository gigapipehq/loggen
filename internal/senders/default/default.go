package _default

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

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

type DefaultSender struct {
	httpClient *http.Client
	transport  *roundTripper
}

func (s *DefaultSender) Client() *http.Client {
	return s.httpClient
}

func (s *DefaultSender) WithHeaders(headers map[string]string) *DefaultSender {
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

func (s *DefaultSender) WithURL(httpURL string) (*DefaultSender, error) {
	u, err := url.Parse(httpURL)
	if err != nil {
		return nil, err
	}
	if s.transport == nil {
		s.transport = &roundTripper{}
	}
	s.transport.url = u
	return s, nil
}

func (s *DefaultSender) Send(ctx context.Context, batch []byte) error {
	httpMethod := "POST"
	reqSize := int64(len(batch))
	ctx, span := otel.Tracer.Start(ctx, "send log batch", trace.WithAttributes(
		attribute.Key("http.url").String(s.transport.url.String()),
		attribute.Key("http.method").String(httpMethod),
		attribute.Key("http.request_size").Int64(reqSize),
	))
	defer span.End()

	body := io.NopCloser(bytes.NewReader(batch))
	req := &http.Request{
		URL:           s.transport.url,
		Method:        httpMethod,
		Body:          body,
		ContentLength: reqSize,
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	ctx, span = otel.Tracer.Start(ctx, "batch logs sent", trace.WithAttributes(
		attribute.Key("http.status").Int(resp.StatusCode),
	))
	defer span.End()
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= http.StatusBadRequest {
		rbody, _ := io.ReadAll(resp.Body)
		span.RecordError(err)
		span.SetStatus(codes.Error, string(rbody))
		return fmt.Errorf("unexpected status code %d received with error: %s", resp.StatusCode, string(rbody))
	}
	return nil
}

func New() *DefaultSender {
	return &DefaultSender{
		httpClient: &http.Client{},
	}
}
