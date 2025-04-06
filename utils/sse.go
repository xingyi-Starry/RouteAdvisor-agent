package utils

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// send
func SetWebStream(c *fiber.Ctx) {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")
}

func GenerateWebStreamMessage(data []byte) string {
	var builder strings.Builder
	builder.WriteString("data: ")
	builder.Write(data)
	builder.WriteString("\n\n")
	return builder.String()
}

func GenerateWebStreamMessageByString(data string) string {
	var builder strings.Builder
	builder.WriteString("data: ")
	builder.WriteString(data)
	builder.WriteString("\n\n")
	return builder.String()
}

// client
type SseClient struct {
	url       string
	body      io.Reader
	eventChan chan string
	ctx       context.Context
	cancel    context.CancelFunc
	header    http.Header
}

func NewSseClient(ctx context.Context, url string, body any) (*SseClient, error) {
	marshalBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	logrus.Debug(string(marshalBody))
	ctx, cancel := context.WithCancel(ctx)
	return &SseClient{
		url:       url,
		body:      bytes.NewReader(marshalBody),
		eventChan: make(chan string),
		ctx:       ctx,
		cancel:    cancel,
	}, nil
}

func (s *SseClient) AddHeader(key, value string) {
	if s.header == nil {
		s.header = make(http.Header)
	}
	s.header.Set(key, value)
}

func (s *SseClient) GetEventChan() <-chan string {
	return s.eventChan
}

func (s *SseClient) Cancel() {
	close(s.eventChan)
	s.cancel()
}

func (s *SseClient) Start() error {
	req, err := http.NewRequestWithContext(s.ctx, http.MethodPost, s.url, s.body)
	// logrus.Debug(req)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/event-stream")
	for key, value := range s.header {
		req.Header[key] = value
	}
	client := &http.Client{
		Timeout:   time.Minute * 5,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	reader := bufio.NewReader(resp.Body)

	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}

		if text != "" {
			select {
			case <-s.ctx.Done():
				return nil
			default:
				s.eventChan <- text
			}
		}
	}
	close(s.eventChan)
	return nil
}
