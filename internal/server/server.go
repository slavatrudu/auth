package server

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/slavatrudu/auth/internal/mapper"
	"github.com/slavatrudu/auth/internal/model"
	authpb "github.com/slavatrudu/contracts/auth/go"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	authpb.UnimplementedAuthServer

	authService AuthService
	logger      *zerolog.Logger
}

type AuthService interface {
	Register(context.Context, model.Register) error
	Login(context.Context, model.Login) (model.TokenPair, error)
	Refresh(context.Context, string) (model.TokenPair, error)
	Validate(context.Context, string) (uint64, error)
	Logout(context.Context, string) error
}

func New(authService AuthService, logger *zerolog.Logger) *Server {
	return &Server{authService: authService, logger: logger}
}

func (s *Server) Register(ctx context.Context, req *authpb.RegisterRequest) (*emptypb.Empty, error) {
	r := mapper.PbRegisterToRegisterModel(req)
	if err := s.authService.Register(ctx, r); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	pair, err := s.authService.Login(ctx, mapper.PbLoginToLoginModel(req))
	if err != nil {
		return nil, err
	}
	return &authpb.LoginResponse{TokenPair: &authpb.TokenPair{AccessToken: pair.AccessToken, RefreshToken: pair.RefreshToken}}, nil
}

func (s *Server) Refresh(ctx context.Context, req *authpb.RefreshRequest) (*authpb.RefreshResponse, error) {
	pair, err := s.authService.Refresh(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &authpb.RefreshResponse{TokenPair: &authpb.TokenPair{AccessToken: pair.AccessToken, RefreshToken: pair.RefreshToken}}, nil
}

func (s *Server) Validate(ctx context.Context, req *authpb.ValidateRequest) (*authpb.ValidateResponse, error) {
	uid, err := s.authService.Validate(ctx, req.AccessToken)
	if err != nil {
		return nil, err
	}
	return &authpb.ValidateResponse{UserId: uid}, nil
}

func (s *Server) Logout(ctx context.Context, req *authpb.LogoutRequest) (*emptypb.Empty, error) {
	if err := s.authService.Logout(ctx, req.RefreshToken); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
