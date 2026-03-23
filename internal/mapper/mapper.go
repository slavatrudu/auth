package mapper

import (
	"github.com/slavatrudu/auth/internal/model"
	authpb "github.com/slavatrudu/contracts/auth/go"
)

func PbRegisterToRegisterModel(user *authpb.RegisterRequest) model.Register {
	return model.Register{
		Login:    user.Login,
		Email:    user.Email,
		Password: user.Password,
	}
}

func PbLoginToLoginModel(user *authpb.LoginRequest) model.Login {
	return model.Login{
		LoginOrEmail: user.LoginOrEmail,
		Password:     user.Password,
	}
}
