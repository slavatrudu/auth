package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	Login        string    `gorm:"column:login"`
	Email        string    `gorm:"column:email"`
	PasswordHash string    `gorm:"column:password_hash"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "users"
}

// BeforeCreate - хук GORM для установки ID в 0 перед созданием
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = 0
	return nil
}

type RefreshToken struct {
	ID        uint64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    uint64     `gorm:"column:user_id"`
	Token     string     `gorm:"column:token"`
	ExpiresAt time.Time  `gorm:"column:expires_at"`
	RevokedAt *time.Time `gorm:"column:revoked_at"`
	CreatedAt time.Time  `gorm:"column:created_at"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// BeforeCreate - хук GORM для установки ID в 0 перед созданием
func (rt *RefreshToken) BeforeCreate(tx *gorm.DB) error {
	rt.ID = 0
	return nil
}
