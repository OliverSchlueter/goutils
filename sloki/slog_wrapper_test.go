package sloki_test

import (
	"bytes"
	"context"
	"errors"
	"github.com/OliverSchlueter/goutils/sloki"
	"io"
	"net/http"
	"testing"
)

func TestWrapContext(t *testing.T) {
	fn := func(ctx context.Context) string {
		return "test value"
	}

	sloki.RegisterContextFunc("testKey", fn)

	got := sloki.WrapContext(context.Background())
	if got.Key != "context" {
		t.Errorf("expected key 'context', got %s", got.Key)
	}
	if len(got.Value.Group()) != 1 {
		t.Errorf("expected 1 attribute, got %d", len(got.Value.Group()))
	}
}

func TestWrapRequest_Basic(t *testing.T) {
	body := "test body"
	req, err := http.NewRequest("POST", "http://example.com/test", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("X-Test-Header", "header-value")
	req.Header.Set("User-Agent", "test-agent")
	req.Header.Set("Referer", "http://referrer.com")

	attr := sloki.WrapRequest(req)
	if attr.Key != "request" {
		t.Errorf("expected key 'request', got %s", attr.Key)
	}
	group := attr.Value.Group()
	if len(group) != 6 {
		t.Errorf("expected 6 attributes, got %d", len(group))
	}

	for _, a := range group {
		if a.Key == "body" && a.Value.String() != "test body" {
			t.Errorf("expected body, got %q", a.Value.String())
		}
	}
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (int, error) {
	return 0, errors.New("read error")
}

func (e *errorReader) Close() error {
	return nil
}

func TestWrapRequest_BodyReadError(t *testing.T) {
	req, err := http.NewRequest("POST", "http://example.com/test", io.NopCloser(&errorReader{}))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	attr := sloki.WrapRequest(req)
	if attr.Key != "request" {
		t.Errorf("expected key 'request', got %s", attr.Key)
	}
	group := attr.Value.Group()
	if len(group) != 6 {
		t.Errorf("expected 6 attributes, got %d", len(group))
	}
	// Optionally, check that body is empty string
	for _, a := range group {
		if a.Key == "body" && a.Value.String() != "" {
			t.Errorf("expected empty body, got %q", a.Value.String())
		}
	}
}
