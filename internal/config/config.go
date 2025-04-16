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
	Port    int `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path ==""{
		panic("config path is empty")
	}



	return MustLoadByPath(path)
}

func MustLoadByPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic(err)
	}

	var cfg Config

	cleanenv.ReadConfig(configPath, &cfg)

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