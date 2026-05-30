package api

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/lhqua/gopanel/internal/caddy"
	"github.com/lhqua/gopanel/internal/publicshare"
)

func isSafeShareID(id string) bool {
	if id == "" || len(id) > 64 || id == "." || id == ".." {
		return false
	}
	for _, char := range id {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-' || char == '_') {
			return false
		}
	}
	return true
}

type CreateShareRequest struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (s *Server) ListShares(w http.ResponseWriter, r *http.Request) {
	entries, err := os.ReadDir(s.cfg.PublicDir)
	if err != nil {
		writeJSON(w, http.StatusOK, []publicshare.Share{})
		return
	}
	var shares []publicshare.Share
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		shares = append(shares, publicshare.Share{
			ID:     e.Name(),
			Folder: e.Name(),
			Type:   "album",
			Title:  e.Name(),
		})
	}
	writeJSON(w, http.StatusOK, shares)
}

func (s *Server) CreateShare(w http.ResponseWriter, r *http.Request) {
	var req CreateShareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	shareID := req.ID
	if shareID == "" {
		shareID = randomID()
	} else if !isSafeShareID(shareID) {
		http.Error(w, `{"error":"invalid share id"}`, http.StatusBadRequest)
		return
	}
	title := req.Title
	if title == "" {
		title = shareID
	}

	folderPath := filepath.Join(s.cfg.PublicDir, shareID)
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		http.Error(w, `{"error":"failed to create share"}`, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, publicshare.Share{
		ID:          shareID,
		Folder:      shareID,
		Type:        "album",
		Title:       title,
		Description: req.Description,
	})
}

func (s *Server) DeleteShare(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !isSafeShareID(id) {
		http.Error(w, `{"error":"invalid share id"}`, http.StatusBadRequest)
		return
	}
	folderPath := filepath.Join(s.cfg.PublicDir, id)
	if err := os.RemoveAll(folderPath); err != nil {
		http.Error(w, `{"error":"failed to delete share"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (s *Server) ListShareFiles(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !isSafeShareID(id) {
		http.Error(w, `{"error":"invalid share id"}`, http.StatusBadRequest)
		return
	}
	pm := publicshare.New(s.cfg.PublicDir)
	files, err := pm.ListFiles(id)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, files)
}

func (s *Server) UploadShareFile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !isSafeShareID(id) {
		http.Error(w, `{"error":"invalid share id"}`, http.StatusBadRequest)
		return
	}
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		http.Error(w, `{"error":"failed to parse form"}`, http.StatusBadRequest)
		return
	}

	folderPath := filepath.Join(s.cfg.PublicDir, id)
	os.MkdirAll(folderPath, 0755)

	files := r.MultipartForm.File["files"]
	var uploaded []string
	for _, fh := range files {
		f, err := fh.Open()
		if err != nil {
			continue
		}
		filename := filepath.Base(fh.Filename)
		dst, err := os.Create(filepath.Join(folderPath, filename))
		if err != nil {
			f.Close()
			continue
		}
		if _, err := dst.ReadFrom(f); err != nil {
			dst.Close()
			f.Close()
			continue
		}
		dst.Close()
		f.Close()
		uploaded = append(uploaded, fh.Filename)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"uploaded": uploaded,
	})
}

func (s *Server) DeleteShareFile(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}
	if !isSafeShareID(req.ID) {
		http.Error(w, `{"error":"invalid share id"}`, http.StatusBadRequest)
		return
	}
	filename := filepath.Base(req.Name)
	filePath := filepath.Join(s.cfg.PublicDir, req.ID, filename)
	if err := os.Remove(filePath); err != nil {
		http.Error(w, `{"error":"failed to delete file"}`, http.StatusInternalServerError)
		return
	}
	pm := publicshare.New(s.cfg.PublicDir)
	pm.DeleteCachedThumbnail(req.ID, req.Name)
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (s *Server) PublicDomainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		port := 3637
		if s.cfg.PublicPort != 0 {
			port = s.cfg.PublicPort
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"public_domain": s.cfg.PublicDomain,
			"public_port":   port,
		})
	case "PUT":
		var req struct {
			PublicDomain string `json:"public_domain"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		oldDomain := s.cfg.PublicDomain
		s.cfg.PublicDomain = req.PublicDomain

		if oldDomain != "" && req.PublicDomain == "" {
			s.caddy.RemoveSite(oldDomain)
		}

		if req.PublicDomain != "" {
			s.caddy.AddSite(caddy.Site{
				Domain:      req.PublicDomain,
				ServicePort: s.cfg.PublicPort,
				TLSEnabled:  false,
				Type:        "proxy",
			})
		}

		s.cfg.Save()
		writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
	}
}

var cacheImmutable = "public, max-age=31536000, immutable"
var cacheCDN = "public, s-maxage=31536000, max-age=86400, immutable"

// PublicServer returns an http.Handler for the public-facing server on port 3637.
func (s *Server) PublicServer() http.Handler {
	r := chi.NewRouter()

	r.Get("/p/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if !isSafeShareID(id) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		s.servePublicPage(w, r, id)
	})

	r.HandleFunc("/p/{id}/*", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if !isSafeShareID(id) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		subPath := strings.TrimPrefix(r.URL.Path, "/p/"+id+"/")
		subPath = filepath.Clean(subPath)
		if strings.HasPrefix(subPath, "..") || strings.Contains(subPath, "/") || strings.Contains(subPath, "\\") {
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}

		// Thumbnail — serve from disk cache if available
		if r.URL.Query().Get("thumb") == "1" {
			pm := publicshare.New(s.cfg.PublicDir)
			data, err := pm.GenerateThumbnail(id, subPath)
			if err != nil {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "image/jpeg")
			w.Header().Set("Cache-Control", cacheImmutable)
			w.Header().Set("CDN-Cache-Control", cacheImmutable)
			w.Header().Set("Cloudflare-CDN-Cache-Control", cacheImmutable)
			w.Write(data)
			return
		}

		// Original file
		filePath := filepath.Join(s.cfg.PublicDir, id, subPath)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Cache-Control", cacheImmutable)
		w.Header().Set("CDN-Cache-Control", cacheImmutable)
		w.Header().Set("Cloudflare-CDN-Cache-Control", cacheImmutable)
		http.ServeFile(w, r, filePath)
	})

	return r
}

func (s *Server) servePublicPage(w http.ResponseWriter, r *http.Request, id string) {
	pm := publicshare.New(s.cfg.PublicDir)
	files, err := pm.ListFiles(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	baseURL := ""
	if s.cfg.PublicDomain != "" {
		baseURL = "https://" + s.cfg.PublicDomain
	}

	share := &publicshare.Share{
		Title:  id,
		Folder: id,
	}

	html := publicshare.GalleryPage(share, files, baseURL)

	etag := fmt.Sprintf(`"%s-%d"`, id, len(files))

	if match := r.Header.Get("If-None-Match"); match == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("ETag", etag)
	w.Header().Set("Cache-Control", "public, s-maxage=300, max-age=0, must-revalidate")
	w.Header().Set("CDN-Cache-Control", "public, s-maxage=300")
	w.Header().Set("Cloudflare-CDN-Cache-Control", "public, s-maxage=300")
	w.Write([]byte(html))
}

func randomID() string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 8)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			b[i] = chars[0]
		} else {
			b[i] = chars[num.Int64()]
		}
	}
	return string(b)
}


