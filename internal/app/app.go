package app

import (
	grpcapp "auth-service/internal/app/grpc_app"
	"auth-service/internal/services/auth"
	"auth-service/internal/storage/sqllite"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	//db init
	storage,err := sqllite.New(storagePath)
	if err != nil {
		panic(err)
	}
	//init auth service

	authService := auth.New(log,storage,storage,storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService,grpcPort)

	return &App{GRPCSrv: grpcApp}
}