package grpc_app

import (
	authGrpc "auth-service/internal/grpc/auth"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type App struct {
	log *slog.Logger
	gRPCServer *grpc.Server
	port int
}

func New(log *slog.Logger,  port int) *App {
	gRPCServer := grpc.NewServer()

	authGrpc.RegisterServer(gRPCServer)

	return &App{log: log, gRPCServer: gRPCServer, port: port}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpc_app.Run"

	log := a.log.With( slog.String("op", op),slog.Int("port", a.port))

	log.Info("Starting gRPC server")

	l,err :=net.Listen("tcp",fmt.Sprintf(":%d",a.port))

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop()  {
	const op = "grpc_app.Stop"

	a.log.With(slog.String("op", op)).Info("Stopping gRPC server")

	a.gRPCServer.GracefulStop()
}