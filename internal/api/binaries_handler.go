package api

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type binaryResponse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	OriginalName string `json:"original_name"`
	Path         string `json:"path"`
	Size         int64  `json:"size"`
	Version      string `json:"version"`
}

func (s *Server) ListBinaries(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query("SELECT id, name, original_name, path, size, version FROM binaries ORDER BY id DESC")
	if err != nil {
		http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	binaries := make([]binaryResponse, 0)
	for rows.Next() {
		var b binaryResponse
		if err := rows.Scan(&b.ID, &b.Name, &b.OriginalName, &b.Path, &b.Size, &b.Version); err != nil {
			continue
		}
		binaries = append(binaries, b)
	}
	writeJSON(w, http.StatusOK, binaries)
}

func (s *Server) UploadBinary(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		http.Error(w, `{"error":"failed to parse multipart form"}`, http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, `{"error":"file is required"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	version := r.FormValue("version")
	if version == "" {
		version = "1.0.0"
	}

	if err := os.MkdirAll(s.cfg.BinariesDir, 0755); err != nil {
		http.Error(w, `{"error":"failed to create binaries directory"}`, http.StatusInternalServerError)
		return
	}

	name := filepath.Base(header.Filename)
	dstPath := filepath.Join(s.cfg.BinariesDir, name)
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, `{"error":"failed to save file"}`, http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	written, err := io.Copy(dst, file)
	if err != nil {
		http.Error(w, `{"error":"failed to write file"}`, http.StatusInternalServerError)
		return
	}

	if err := os.Chmod(dstPath, 0755); err != nil {
		http.Error(w, `{"error":"failed to set permissions"}`, http.StatusInternalServerError)
		return
	}

	_, err = s.db.Exec(
		`INSERT INTO binaries (name, original_name, path, size, version) VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(name) DO UPDATE SET original_name=excluded.original_name, path=excluded.path, size=excluded.size, version=excluded.version`,
		name, header.Filename, dstPath, written, version,
	)
	if err != nil {
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}

	var id int64
	s.db.QueryRow("SELECT id FROM binaries WHERE name = ?", name).Scan(&id)
	writeJSON(w, http.StatusCreated, binaryResponse{
		ID:           int(id),
		Name:         name,
		OriginalName: header.Filename,
		Path:         dstPath,
		Size:         written,
		Version:      version,
	})
}

func (s *Server) DeleteBinary(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}
	var path string
	err = s.db.QueryRow("SELECT path FROM binaries WHERE id = ?", id).Scan(&path)
	if err != nil {
		http.Error(w, `{"error":"binary not found"}`, http.StatusNotFound)
		return
	}
	os.Remove(path)
	_, _ = s.db.Exec("DELETE FROM binaries WHERE id = ?", id)
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}


