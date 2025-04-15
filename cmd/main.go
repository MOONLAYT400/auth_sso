package main

import (
	"auth-service/internal/app"
	"auth-service/internal/config"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// init config
	cfg := config.MustLoad()

	fmt.Println(cfg)
	//init logger
	log := config.SetupLogger(cfg.Env)	
	log.Info("Custom logger enabled in",slog.String("env",cfg.Env) )

	//init app
	application  := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)
	go application.GRPCSrv.MustRun()
	//init grpc server

	//gracefull shutdown
	stop :=make(chan os.Signal,1)
	signal.Notify(stop,syscall.SIGTERM,syscall.SIGINT)//TODO:what is this?
	receivedSignal:=<-stop

	log.Info("Received signal",slog.String("signal",receivedSignal.String()))
	application.GRPCSrv.Stop()
	log.Info("Application stopped")
}