package main

import (
	"auth-service/internal/config"
	"fmt"
	"log/slog"
)

func main() {
	// init config
	cfg := config.MustLoad()

	fmt.Println(cfg)
	//init logger
	log := config.SetupLogger(cfg.Env)	
	log.Info("Custom logger enabled in",slog.String("env",cfg.Env) )

	//init app

	//init grpc server
}