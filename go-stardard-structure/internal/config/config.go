package configs

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"`
}

type Config struct {
	Env         string `yaml:"env" env-required:"true" `
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server" env-required:"true"`
}

func MustLoad() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "Path to the configuration file")
		flag.Parse()
		configPath = *flags

		if configPath == "" {
			log.Fatalf("CONFIG_PATH environment variable or --config flag must be set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Configuration file does not exist: %s", configPath)
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)

	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err.Error())
	}

	return &cfg
}
