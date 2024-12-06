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
	DB			`yaml:"db" env-required:"true"`
	HttpServer `yaml:"http_server" env-required:"true"`
}

type HttpServer struct {
	Address 		string `yaml:"address" env-required:"true"`
	Timeout 		time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout	time.Duration `yaml:"idle_timeout" env-default:"60s"`
}
type DB struct {
	Username 		string `yaml:"username" env-default:"postgres"`
	Host 			string `yaml:"host" env-default:"localhost"`
	Port			string `yaml:"port" env-default:"5432"`
	Dbname			string `yaml:"dbname" env-default:"postgres"`
	Sslmode			string `yaml:"sslmode" env-default:"disable"`
	Password		string `env-default:""`
}

func MustLoad() *Config {
	var configPath string
	flag.StringVar(&configPath, "config", "", "config_path")
	flag.Parse()
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	if err := godotenv.Load(); err != nil {
		log.Fatal("env variables loading error")
	}
	cfg.DB.Password = os.Getenv("DB_PASSWORD")
	return &cfg
}