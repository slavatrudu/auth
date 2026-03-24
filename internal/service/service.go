package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"github.com/slavatrudu/auth/internal/config"
	"github.com/slavatrudu/auth/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo   Repository
	cfg    *config.Config
	logger *zerolog.Logger
}

type Repository interface {
	CreateUser(context.Context, model.User) error
	GetUserByID(context.Context, uint64) (model.User, error)
	GetUserByLoginOrEmail(context.Context, string) (model.User, error)
	UpdateRefreshToken(context.Context, model.RefreshToken) error
	GetRefreshToken(context.Context, string) (model.RefreshToken, error)
	RevokeRefreshToken(context.Context, string) error
	DeleteUser(context.Context, uint64) error
}

func New(repo Repository, cfg *config.Config, logger *zerolog.Logger) *AuthService {
	return &AuthService{repo: repo, cfg: cfg, logger: logger}
}

func (s *AuthService) Register(ctx context.Context, req model.Register) error {
	// Валидация входных данных
	if req.Login == "" {
		return fmt.Errorf("login cannot be empty")
	}
	if req.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	if req.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	if len(req.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}

	// Проверяем, не существует ли уже пользователь с таким логином
	_, err := s.repo.GetUserByLoginOrEmail(ctx, req.Login)
	if err == nil {
		return fmt.Errorf("user with login %s already exists", req.Login)
	}

	// Проверяем, не существует ли уже пользователь с таким email
	_, err = s.repo.GetUserByLoginOrEmail(ctx, req.Email)
	if err == nil {
		return fmt.Errorf("user with email %s already exists", req.Email)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := model.NewUser(req.Login, req.Email, string(hash))

	return s.repo.CreateUser(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, req model.Login) (model.TokenPair, error) {
	user, err := s.repo.GetUserByLoginOrEmail(ctx, req.LoginOrEmail)
	if err != nil {
		return model.TokenPair{}, fmt.Errorf("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return model.TokenPair{}, fmt.Errorf("invalid credentials")
	}
	return s.issueTokens(ctx, user.ID)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (model.TokenPair, error) {
	rt, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return model.TokenPair{}, fmt.Errorf("invalid refresh token")
	}
	if rt.RevokedAt != nil || time.Now().After(rt.ExpiresAt) {
		return model.TokenPair{}, fmt.Errorf("refresh token expired or revoked")
	}
	return s.issueTokens(ctx, rt.UserID)
}

func (s *AuthService) Validate(ctx context.Context, accessToken string) (uint64, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JwtSecret), nil
	})
	if err != nil {
		return 0, fmt.Errorf("invalid token: %w", err)
	}
	uidFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid subject")
	}
	return uint64(uidFloat), nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.repo.RevokeRefreshToken(ctx, refreshToken)
}

func (s *AuthService) issueTokens(ctx context.Context, userID uint64) (model.TokenPair, error) {
	now := time.Now()
	accessExp := now.Add(time.Duration(s.cfg.AccessTokenTTLMinutes) * time.Minute)
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     accessExp.Unix(),
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessStr, err := access.SignedString([]byte(s.cfg.JwtSecret))
	if err != nil {
		return model.TokenPair{}, fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshRaw := fmt.Sprintf("%d:%d:%s", userID, now.UnixNano(), s.cfg.JwtSecret)
	h := sha256.Sum256([]byte(refreshRaw))
	refreshStr := hex.EncodeToString(h[:])
	refresh := model.NewRefreshToken(userID, refreshStr, now.Add(time.Duration(s.cfg.RefreshTokenTTLDays)*24*time.Hour))
	if err := s.repo.UpdateRefreshToken(ctx, refresh); err != nil {
		return model.TokenPair{}, err
	}
	return model.TokenPair{AccessToken: accessStr, RefreshToken: refreshStr}, nil
}

func (s *AuthService) DeleteUser(ctx context.Context, userID uint64) error {
	return s.repo.DeleteUser(ctx, userID)
}
