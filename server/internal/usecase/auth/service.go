package authsvc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	domainUser "github.com/0xirvan/snippetbox/internal/domain/user"
	authPort "github.com/0xirvan/snippetbox/internal/port/auth"
	sessionPort "github.com/0xirvan/snippetbox/internal/port/session"
	"github.com/0xirvan/snippetbox/internal/shared/security"
)

type AuthService struct {
	userRepo    authPort.AuthUserRepository
	sessionRepo sessionPort.SessionRepository
	sessionTTL  time.Duration
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

func NewAuthService(
	userRepo authPort.AuthUserRepository,
	sessionRepo sessionPort.SessionRepository,
	sessionTTL time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		sessionTTL:  sessionTTL,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*domainUser.User, error) {
	if existing, _ := s.userRepo.FindByEmail(ctx, email); existing != nil {
		return nil, ErrUserExists
	}

	hashed, err := security.HashPassword(password)
	if err != nil {
		return nil, err
	}

	u := domainUser.NewUser(email, hashed)
	if err := s.userRepo.Save(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if s := security.CheckPassword(u.Password, password); s != true {
		return "", ErrInvalidCredentials
	}

	sid, err := generateSessionID(32)
	if err != nil {
		return "", err
	}

	if err := s.sessionRepo.Save(ctx, sid, u.ID, s.sessionTTL); err != nil {
		return "", err
	}

	return sid, nil
}

func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	return s.sessionRepo.Delete(ctx, sessionID)
}

func generateSessionID(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
