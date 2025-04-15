package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env_required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC        GRPCConfig `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    string `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path ==""{
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(err)
	}

	var cfg Config

	cleanenv.ReadConfig(path, &cfg)

	return &cfg
}

func fetchConfigPath() string {
	var res string

	
	
	flag.StringVar(&res, "config", "", "path to config file")
	
	flag.Parse()
	
	if(res == "") {
			err := godotenv.Load()
				if err != nil {
					log.Fatal("Error loading .env file")
				}
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}