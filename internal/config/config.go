package config

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"`
}

// env-default:"production"
type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

func MustLoad() *Config {
	var configPath string

	// First try to get config path from environment
	configPath = os.Getenv("CONFIG_PATH")

	// If not set in environment, try command line flag
	if configPath == "" {
		flags := flag.String("config", "", "path to config file")
		flag.Parse()

		configPath = *flags
	}

	// If still not set, use default config path for local development
	if configPath == "" {
		// Get the current working directory
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal("Failed to get working directory: ", err)
		}
		// Default to config/local.yaml in the project root
		configPath = filepath.Join(wd, "config", "production.yaml")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Cannot read config file: %s", err.Error())
	}

	return &cfg
}
