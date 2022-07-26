package server

import "github.com/dollarkillerx/inventory/internal/middleware"

func (s *Server) router() {
	s.app.POST("/api/v1/login", s.userCenter)

	v1 := s.app.Group("/api/v1", middleware.UAAuthorization())
	{
		v1.POST("/goods")
	}
}
