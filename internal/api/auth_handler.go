package api

import (
	"net/http"

	"github.com/lhqua/gopanel/internal/auth"
	"github.com/lhqua/gopanel/internal/middleware"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token    string `json:"token"`
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := decodeJSON(r, &req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	var id int
	var username, passwordHash, role string
	err := s.db.QueryRow("SELECT id, username, password_hash, role FROM users WHERE username = ?", req.Username).
		Scan(&id, &username, &passwordHash, &role)
	if err != nil {
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}
	if !auth.CheckPassword(req.Password, passwordHash) {
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}
	token, err := s.auth.Generate(id, username, role)
	if err != nil {
		http.Error(w, `{"error":"failed to generate token"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, loginResponse{
		Token:    token,
		UserID:   id,
		Username: username,
		Role:     role,
	})
}

func (s *Server) Me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user_id":  claims.UserID,
		"username": claims.Username,
		"role":     claims.Role,
	})
}

func (s *Server) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}
	if err := decodeJSON(r, &req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}
	claims, err := s.auth.Validate(req.Token)
	if err != nil {
		http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
		return
	}
	newToken, err := s.auth.Generate(claims.UserID, claims.Username, claims.Role)
	if err != nil {
		http.Error(w, `{"error":"failed to refresh"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": newToken})
}
