package server

import (
	"github.com/dollarkillerx/inventory/internal/pkg/errs"
	"github.com/dollarkillerx/inventory/internal/pkg/response"
	"github.com/dollarkillerx/inventory/internal/utils"
	"github.com/gin-gonic/gin"

	"log"
)

func (s *Server) Goods(ctx *gin.Context) {

}

func (s *Server) Good(ctx *gin.Context) {
	model := utils.GetAuthModel(ctx)

	barcodes := ctx.Param("barcodes")
	good, err := s.storage.Good(barcodes, model.Account)
	if err != nil {
		log.Println(err)
		response.Return(ctx, errs.BadRequest)
		return
	}

	// https://jkl-1253341723.cos.ap-chengdu.myqcloud.com/default.png
	response.Return(ctx, good)
}
