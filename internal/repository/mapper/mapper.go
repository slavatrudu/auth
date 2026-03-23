package mapper

import (
	"github.com/slavatrudu/auth/internal/model"
	repomodel "github.com/slavatrudu/auth/internal/repository/model"
)

func UserToRepoUser(user model.User) repomodel.User {
	return repomodel.User{
		ID:           user.ID,
		Login:        user.Login,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

func RepoUserToUser(user repomodel.User) model.User {
	return model.User{
		ID:           user.ID,
		Login:        user.Login,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

func RefreshTokenToRepoRefresh(refreshToken model.RefreshToken) repomodel.RefreshToken {
	return repomodel.RefreshToken{
		ID:        refreshToken.ID,
		UserID:    refreshToken.UserID,
		Token:     refreshToken.Token,
		ExpiresAt: refreshToken.ExpiresAt,
		RevokedAt: refreshToken.RevokedAt,
		CreatedAt: refreshToken.CreatedAt,
	}
}

func RepoRefreshTokenToRefresh(refreshToken repomodel.RefreshToken) model.RefreshToken {
	return model.RefreshToken{
		ID:        refreshToken.ID,
		UserID:    refreshToken.UserID,
		Token:     refreshToken.Token,
		ExpiresAt: refreshToken.ExpiresAt,
		RevokedAt: refreshToken.RevokedAt,
		CreatedAt: refreshToken.CreatedAt,
	}
}
