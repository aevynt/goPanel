package api

import (
	"fmt"
	"net/http"

	"github.com/aevynt/goPanel/internal/auth"
	"github.com/aevynt/goPanel/internal/middleware"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Code     string `json:"code,omitempty"`
}

type loginResponse struct {
	Token    string `json:"token,omitempty"`
	UserID   int    `json:"user_id,omitempty"`
	Username string `json:"username,omitempty"`
	Role     string `json:"role,omitempty"`
	Status   string `json:"status,omitempty"`
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := decodeJSON(r, &req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	u, err := s.db.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	if !auth.CheckPassword(req.Password, u.PasswordHash) {
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	// 2FA Verification Flow
	if u.TOTPEnabled == 1 {
		if req.Code == "" {
			writeJSON(w, http.StatusOK, loginResponse{
				Status: "require_2fa",
			})
			return
		}
		if !auth.VerifyTOTP(u.TOTPSecret, req.Code) {
			http.Error(w, `{"error":"invalid 2FA code"}`, http.StatusUnauthorized)
			return
		}
	}

	token, err := s.auth.Generate(u.ID, u.Username, u.Role)
	if err != nil {
		http.Error(w, `{"error":"failed to generate token"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, loginResponse{
		Token:    token,
		UserID:   u.ID,
		Username: u.Username,
		Role:     u.Role,
	})
}

func (s *Server) Me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	
	u, err := s.db.GetUserByID(claims.UserID)
	if err != nil {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user_id":      claims.UserID,
		"username":     claims.Username,
		"role":         claims.Role,
		"totp_enabled": u.TOTPEnabled == 1,
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

func (s *Server) ChangePassword(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := decodeJSON(r, &req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}
	if req.CurrentPassword == "" || req.NewPassword == "" {
		http.Error(w, `{"error":"both fields are required"}`, http.StatusBadRequest)
		return
	}

	u, err := s.db.GetUserByID(claims.UserID)
	if err != nil {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		return
	}

	if !auth.CheckPassword(req.CurrentPassword, u.PasswordHash) {
		http.Error(w, `{"error":"incorrect current password"}`, http.StatusUnauthorized)
		return
	}

	newHash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		http.Error(w, `{"error":"failed to hash password"}`, http.StatusInternalServerError)
		return
	}

	err = s.db.UpdatePassword(claims.UserID, newHash)
	if err != nil {
		http.Error(w, `{"error":"failed to update password"}`, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (s *Server) Setup2FA(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	secret, err := auth.GenerateSecret()
	if err != nil {
		http.Error(w, `{"error":"failed to generate 2FA secret"}`, http.StatusInternalServerError)
		return
	}

	err = s.db.SetTOTPSecret(claims.UserID, secret)
	if err != nil {
		http.Error(w, `{"error":"failed to save 2FA secret"}`, http.StatusInternalServerError)
		return
	}

	issuer := "goPanel"
	uri := fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s", issuer, claims.Username, secret, issuer)

	writeJSON(w, http.StatusOK, map[string]string{
		"secret":           secret,
		"provisioning_uri": uri,
	})
}

func (s *Server) Enable2FA(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	var req struct {
		Code string `json:"code"`
	}
	if err := decodeJSON(r, &req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	u, err := s.db.GetUserByID(claims.UserID)
	if err != nil || u.TOTPSecret == "" {
		http.Error(w, `{"error":"2FA not configured"}`, http.StatusBadRequest)
		return
	}

	if !auth.VerifyTOTP(u.TOTPSecret, req.Code) {
		http.Error(w, `{"error":"invalid verification code"}`, http.StatusUnauthorized)
		return
	}

	err = s.db.Enable2FA(claims.UserID)
	if err != nil {
		http.Error(w, `{"error":"failed to enable 2FA"}`, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "enabled"})
}

func (s *Server) Disable2FA(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	var req struct {
		Password string `json:"password"`
	}
	if err := decodeJSON(r, &req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	u, err := s.db.GetUserByID(claims.UserID)
	if err != nil {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		return
	}

	if !auth.CheckPassword(req.Password, u.PasswordHash) {
		http.Error(w, `{"error":"incorrect password"}`, http.StatusUnauthorized)
		return
	}

	err = s.db.Disable2FA(claims.UserID)
	if err != nil {
		http.Error(w, `{"error":"failed to disable 2FA"}`, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "disabled"})
}
