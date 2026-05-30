package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lhqua/gopanel/internal/auth"
	"github.com/lhqua/gopanel/internal/middleware"
)

type userResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (s *Server) ListUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query("SELECT id, username, role FROM users ORDER BY id")
	if err != nil {
		http.Error(w, `{"error":"query users"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	users := make([]userResponse, 0)
	for rows.Next() {
		var u userResponse
		if err := rows.Scan(&u.ID, &u.Username, &u.Role); err != nil {
			continue
		}
		users = append(users, u)
	}
	writeJSON(w, http.StatusOK, users)
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}
	if req.Username == "" || req.Password == "" {
		http.Error(w, `{"error":"username and password required"}`, http.StatusBadRequest)
		return
	}
	if req.Role == "" {
		req.Role = "viewer"
	}
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, `{"error":"failed to hash password"}`, http.StatusInternalServerError)
		return
	}
	result, err := s.db.Exec("INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?)",
		req.Username, hash, req.Role)
	if err != nil {
		http.Error(w, `{"error":"username already exists"}`, http.StatusConflict)
		return
	}
	id, _ := result.LastInsertId()
	writeJSON(w, http.StatusCreated, userResponse{
		ID:       int(id),
		Username: req.Username,
		Role:     req.Role,
	})
}

func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid user id"}`, http.StatusBadRequest)
		return
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}
	if req.Username != "" {
		_, err = s.db.Exec("UPDATE users SET username = ? WHERE id = ?", req.Username, id)
		if err != nil {
			http.Error(w, `{"error":"failed to update username: `+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
	}
	if req.Password != "" {
		hash, err := auth.HashPassword(req.Password)
		if err != nil {
			http.Error(w, `{"error":"failed to hash password"}`, http.StatusInternalServerError)
			return
		}
		_, err = s.db.Exec("UPDATE users SET password_hash = ? WHERE id = ?", hash, id)
		if err != nil {
			http.Error(w, `{"error":"failed to update password: `+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
	}
	if req.Role != "" {
		_, err = s.db.Exec("UPDATE users SET role = ? WHERE id = ?", req.Role, id)
		if err != nil {
			http.Error(w, `{"error":"failed to update role: `+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid user id"}`, http.StatusBadRequest)
		return
	}
	claims := middleware.GetClaims(r)
	if claims.UserID == id {
		http.Error(w, `{"error":"cannot delete yourself"}`, http.StatusBadRequest)
		return
	}
	_, _ = s.db.Exec("DELETE FROM users WHERE id = ?", id)
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
