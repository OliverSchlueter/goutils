package ratelimit

import "net/http"

func GetIP(r *http.Request) string {
	ip := "unknown"
	if r.RemoteAddr != "" {
		ip = r.RemoteAddr
	}
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ip = xff
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		ip = xri
	}
	return ip
}
