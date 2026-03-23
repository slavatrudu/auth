package model

import "time"

type User struct {
	ID           uint64
	Login        string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewUser создает нового пользователя с правильной инициализацией
func NewUser(login, email, passwordHash string) User {
	now := time.Now()
	return User{
		ID:           0,
		Login:        login,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

type RefreshToken struct {
	ID        uint64
	UserID    uint64
	Token     string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

// NewRefreshToken создает новый refresh token с правильной инициализацией
func NewRefreshToken(userID uint64, token string, expiresAt time.Time) RefreshToken {
	return RefreshToken{
		ID:        0,
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		RevokedAt: nil,
		CreatedAt: time.Now(),
	}
}

type Register struct {
	Login    string
	Email    string
	Password string
}

type Login struct {
	LoginOrEmail string
	Password     string
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
