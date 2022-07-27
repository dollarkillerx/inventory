package server

import (
	"github.com/dollarkillerx/inventory/internal/conf"
	"github.com/dollarkillerx/inventory/internal/middleware"
	"github.com/dollarkillerx/inventory/internal/storage"
	"github.com/dollarkillerx/inventory/internal/storage/simple"
	"github.com/gin-gonic/gin"
)

type Server struct {
	app     *gin.Engine
	storage storage.Interface
}

func NewServer() *Server {
	return &Server{
		app: gin.New(),
	}
}

func (s *Server) Run() error {
	newSimple, err := simple.NewSimple(&conf.CONF.PgSQLConfig)
	if err != nil {
		return err
	}

	s.storage = newSimple

	s.app.Use(middleware.SetBasicInformation())
	s.app.Use(middleware.Cors())
	s.app.Use(middleware.HttpRecover())

	s.router()

	return s.app.Run(conf.CONF.ListenAddr)
}
