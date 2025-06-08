package handler

import (
	"encoding/json"
	"io/fs"
	"log/slog"
	"mime"
	"net/http"
	"path"

	"github.com/AntonChubarov/goit-cloud-fp/internal/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *service.Shortener
	fs  fs.FS // embedded dist/ filesystem
}

func New(svc *service.Shortener, webFS fs.FS) *Handler {
	return &Handler{svc: svc, fs: webFS}
}

type createReq struct {
	URL string `json:"url"`
}

type createRes struct {
	Short string `json:"short"`
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Post("/api/links", h.create)
	r.Get("/r/{code}", h.redirect)

	fileServer := http.StripPrefix("/", http.FileServer(http.FS(h.fs)))
	r.Handle("/*", withMimeFix(fileServer, h))

	slog.Debug("routes configured for handler")
	return r
}

// --- Handlers ---

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("unable to decode request body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slog.Debug("create request received", "url", req.URL)

	link, err := h.svc.Create(r.Context(), req.URL)
	if err != nil {
		slog.Error("unable to create short link", "url", req.URL, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := createRes{Short: link.ShortCode}
	slog.Debug("short link created", "short", res.Short)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		slog.Error("unable to encode response", "error", err)
	}
}

func (h *Handler) redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	slog.Debug("redirect request received", "code", code)

	url, err := h.svc.Resolve(r.Context(), code)
	if err != nil || url == "" {
		slog.Warn("unable to resolve short code", "code", code, "error", err)
		http.NotFound(w, r)
		return
	}
	slog.Debug("short code resolved", "code", code, "url", url)
	http.Redirect(w, r, url, http.StatusFound)
}

// --- Static Middleware ---

func withMimeFix(next http.Handler, h *Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ext := path.Ext(r.URL.Path)

		if ext != "" {
			if ctype := mime.TypeByExtension(ext); ctype != "" {
				w.Header().Set("Content-Type", ctype)
				slog.Debug("MIME type set for static file", "path", r.URL.Path, "type", ctype)
			}
			next.ServeHTTP(w, r)
			return
		}

		slog.Debug("serving index.html for non-file route", "path", r.URL.Path)
		h.serveIndex(w)
	})
}

func (h *Handler) serveIndex(w http.ResponseWriter) {
	data, err := fs.ReadFile(h.fs, "index.html")
	if err != nil {
		slog.Error("unable to read index.html from embedded FS", "error", err)
		http.Error(w, "index.html not found", http.StatusInternalServerError)
		return
	}
	slog.Debug("index.html served")
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
