package auth

import (
	"auth-service/internal/services/auth"
	"context"
	"errors"

	ssov1 "github.com/MOONLAYT400/Proto_sso/gen/go/sso"

	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer // for debug without full struct realization
	auth Auth
}

type Auth interface {
	Login(ctx context.Context, email,password string, appId int) (token string, err error)
	Register(ctx context.Context,email,password string) (userId int, err error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}


func RegisterServer(gRPC *grpc.Server,auth Auth)  {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

const emptyValue =0

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse,error) {
	// validation refactor to lib
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty email")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty password")
	}

		if req.GetAppId() == emptyValue {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}

	// service layer
	token,err:=s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId())); 
	if err != nil {

		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ssov1.LoginResponse{
		Token: token,
	},nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse,error) {
		// validation refactor to lib
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty email")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty password")
	}

	// service layer
	userId,err:=s.auth.Register(ctx, req.GetEmail(), req.GetPassword()); 
	if err != nil {

		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ssov1.RegisterResponse{
		UserId: fmt.Sprint(userId),
	},nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse,error) {
	// validation refactor to lib
	if req.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty user_id")
	}

	// service layer
	isAdmin,err:=s.auth.IsAdmin(ctx, req.GetUserId()); 
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	},nil
}