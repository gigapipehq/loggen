package _default

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	progressCh chan int
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

func (s *DefaultSender) AddProgress(count int) {
	s.progressCh <- count
}

func (s *DefaultSender) Progress() <-chan int {
	return s.progressCh
}

func (s *DefaultSender) Send(batch []byte) error {
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

func New() *DefaultSender {
	return &DefaultSender{
		httpClient: &http.Client{},
		progressCh: make(chan int),
	}
}
