package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	ServerAddress   string `env:"RUN_ADDRESS"`
	DatabaseAddress string `env:"DATABASE_URI"`
	AccrualAddress  string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	JWTSecret       string `toml:"jwt_secret" env:"JWT_SECRET" env-required:"true"`
}

func MustLoadConfig() *Config {
	f := flag.NewFlagSet(os.Args[0], flag.PanicOnError)
	serverAddress := f.String("a", "localhost:8080", "host address")
	dsn := f.String("d", "", "database uri")
	accrualAddress := f.String("r", "", "accrual system address")
	_ = f.Parse(os.Args[1:])

	cfg := &Config{
		ServerAddress:   *serverAddress,
		DatabaseAddress: *dsn,
		AccrualAddress:  *accrualAddress,
	}

	if err := cleanenv.ReadConfig("./secret.toml", cfg); err != nil {
		panic(err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		panic(err)
	}

	if cfg.DatabaseAddress == "" {
		panic("database address must be not empty")
	}

	return cfg
}
