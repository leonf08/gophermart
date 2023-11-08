package http

import (
	"github.com/leonf08/gophermart.git/internal/config"
	"net/http"
)

type Server struct {
	server *http.Server
	config *config.Config
}

func NewServer(h http.Handler, cfg *config.Config) {

}
