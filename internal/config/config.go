package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct{
	Addr string
}

type Config struct{
	Env string `yaml:"env"`
	StoragePath string `yaml:"storage_path"`
	HTTPServer `yaml:"http_server"`
}

func MustLoad() *Config{
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == ""{
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		configPath = *flags

		if configPath == ""{
			log.Fatalf("config path is not set")
		}
	}

	_, err := os.Stat(configPath)
	if os.IsNotExist(err){
		log.Fatalf("config file does not exists %s", configPath)
	}

	var cfg Config

	error := cleanenv.ReadConfig(configPath, &cfg)
	if error != nil{
		log.Fatalf("cannot read config file %s", error.Error())
	}

	return &cfg
}