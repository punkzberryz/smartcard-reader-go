package server

import (
	"net/http"

	smc "github.com/punkzberryz/smartcard-reader-go/pkg/smc"
	"github.com/somprasongd/go-thai-smartcard/pkg/model"
)

type ServerConfig struct {
	Port      string
	ApiKey    string
	Broadcast chan model.Message
}
type Server struct {
	ServerConfig ServerConfig
	Card         *smc.SmartCard
	CardConfig   *smc.Options
}

func (s *Server) RunServer() error {
	mux := http.NewServeMux()
	mux.Handle("/health", CORS(http.HandlerFunc(s.handleHome)))
	//auth has CORS included
	mux.Handle("/read", s.auth(http.HandlerFunc(s.handleRead)))

	return http.ListenAndServe(":"+s.ServerConfig.Port, mux)
}
