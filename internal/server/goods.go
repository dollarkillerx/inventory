package server

import (
	"github.com/dollarkillerx/inventory/internal/pkg/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) Goods(ctx *gin.Context) {
	s.storage.DB().Model(&models.Goods{})
}

func (s *Server) Good(ctx *gin.Context) {

}
