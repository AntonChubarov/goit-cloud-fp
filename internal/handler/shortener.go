package handler

import (
	"encoding/json"
	"io/fs"
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

	// API
	r.Post("/api/links", h.create)
	r.Get("/{code}", h.redirect)

	// Serve static files
	fileServer := http.StripPrefix("/", http.FileServer(http.FS(h.fs)))
	r.Handle("/*", withMimeFix(fileServer, h))

	return r
}

// --- Handlers ---

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	link, err := h.svc.Create(r.Context(), req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := createRes{Short: link.ShortCode}
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	url, err := h.svc.Resolve(r.Context(), code)
	if err != nil || url == "" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

// --- Static Middleware ---

func withMimeFix(next http.Handler, h *Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ext := path.Ext(r.URL.Path)

		if ext != "" {
			// Try to set the correct MIME type
			if ctype := mime.TypeByExtension(ext); ctype != "" {
				w.Header().Set("Content-Type", ctype)
			}
			next.ServeHTTP(w, r)
			return
		}

		// Not a file request – serve index.html (for React router)
		h.serveIndex(w)
	})
}

func (h *Handler) serveIndex(w http.ResponseWriter) {
	data, err := fs.ReadFile(h.fs, "index.html")
	if err != nil {
		http.Error(w, "index.html not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
