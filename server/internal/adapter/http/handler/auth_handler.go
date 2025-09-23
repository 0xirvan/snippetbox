package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/0xirvan/snippetbox/internal/adapter/http/util"
	"github.com/0xirvan/snippetbox/internal/dto"
	"github.com/0xirvan/snippetbox/internal/shared/validator"
	authsvc "github.com/0xirvan/snippetbox/internal/usecase/auth"
)

const (
	SessionCookieName = "session_id"
	CookiePath        = "/"
)

type AuthHandler struct {
	authService  authsvc.AuthService
	cookieTTL    time.Duration
	cookieSecure bool
}

func NewAuthHandler(
	authSvc authsvc.AuthService,
	ttl time.Duration,
	secure bool,
) *AuthHandler {
	return &AuthHandler{
		authService:  authSvc,
		cookieTTL:    ttl,
		cookieSecure: secure,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	req := dto.RegisterRequest{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	if err := validator.ValidateStruct(req); err != nil {
		util.WriteValidationError(w, err)
		return
	}

	if _, err := h.authService.Register(r.Context(), req.Email, req.Password); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":  "success",
		"message": "user registered",
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	req := dto.LoginRequest{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	if err := validator.ValidateStruct(req); err != nil {
		util.WriteValidationError(w, err)
		return
	}

	sid, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	cookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    sid,
		Path:     CookiePath,
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(h.cookieTTL),
		MaxAge:   int(h.cookieTTL.Seconds()),
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":  "success",
		"message": "login successful",
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err := h.authService.Logout(r.Context(), cookie.Value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     CookiePath,
		HttpOnly: true,
		Secure:   h.cookieSecure,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":  "success",
		"message": "logged out",
	})
}
