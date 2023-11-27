package main

import (
	"github.com/leonf08/gophermart.git/internal/app"
	"github.com/leonf08/gophermart.git/internal/config"
)

func main() {
	cfg := config.MustLoadConfig()

	app.Run(cfg)
}
