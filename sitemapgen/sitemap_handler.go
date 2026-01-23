package sitemapgen

import (
	"encoding/xml"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/OliverSchlueter/goutils/problems"
	"github.com/OliverSchlueter/goutils/ratelimit"
	"github.com/OliverSchlueter/goutils/sloki"
)

type Handler struct {
	provider      UrlProvider
	rl            *ratelimit.Service
	cacheDuration int
}

type Configuration struct {
	Provider      UrlProvider
	Ratelimit     *ratelimit.Service
	CacheDuration *int
}

func NewHandler(cfg Configuration) *Handler {
	if cfg.Ratelimit == nil {
		cfg.Ratelimit = ratelimit.NewService(ratelimit.Configuration{
			TokensPerSecond: 2,
			MaxTokens:       10,
		})
	}
	if cfg.CacheDuration == nil {
		defaultCacheDuration := 10800 // 3 hours
		cfg.CacheDuration = &defaultCacheDuration
	}

	return &Handler{
		provider:      cfg.Provider,
		rl:            cfg.Ratelimit,
		cacheDuration: *cfg.CacheDuration,
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/sitemap.xml", h.handle)
}

func (h *Handler) handle(w http.ResponseWriter, r *http.Request) {
	if err := h.rl.CheckRequest(r, "sitemap_xml"); err != nil {
		ratelimit.RateLimitExceededProblem().WriteToHTTP(w)
		return
	}

	urls := h.provider()

	sitemap := UrlSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Urls:  urls,
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(h.cacheDuration))

	w.WriteHeader(http.StatusOK)

	enc := xml.NewEncoder(w)
	defer enc.Close()

	enc.Indent("", "  ")

	// Write XML header
	w.Write([]byte(xml.Header))

	if err := enc.Encode(sitemap); err != nil {
		slog.Error("Failed to encode sitemap XML", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}
}
