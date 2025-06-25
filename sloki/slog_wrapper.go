package sloki

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

var ctxFuncs = map[string]func(context.Context) string{}

func RegisterContextFunc(key string, fn func(ctx context.Context) string) {
	ctxFuncs[key] = fn
}

func WrapContext(ctx context.Context) slog.Attr {
	var attrs []slog.Attr
	for key, fn := range ctxFuncs {
		if value := fn(ctx); value != "" {
			attrs = append(attrs, slog.Any(key, value))
		}
	}

	return slog.Group("context", unpackArray(attrs)...)
}

func unpackArray[S ~[]E, E any](s S) []any {
	r := make([]any, len(s))
	for i, e := range s {
		r[i] = e
	}
	return r
}

func WrapRequest(r *http.Request) slog.Attr {
	method := slog.Any("method", r.Method)
	url := slog.Any("url", r.URL.String())
	userAgent := slog.Any("userAgent", r.UserAgent())
	referrer := slog.Any("referrer", r.Referer())

	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		bodyData = make([]byte, 0) // If reading the body fails, use an empty byte slice
	} else {
		// Reset the body so it can be read again later
		r.Body = io.NopCloser(strings.NewReader(string(bodyData)))
	}
	body := slog.Any("body", string(bodyData))

	jsonHeaders, _ := json.Marshal(r.Header)
	jsonHeadersStr := strings.ReplaceAll(string(jsonHeaders), "\"", "")
	headers := slog.Any("headers", jsonHeadersStr)

	return slog.Group("request", method, url, userAgent, headers, referrer, body)
}

func WrapError(err error) slog.Attr {
	return slog.Any("error", err)
}
