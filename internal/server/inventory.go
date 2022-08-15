package server

import (
	"github.com/dollarkillerx/inventory/internal/pkg/errs"
	"github.com/dollarkillerx/inventory/internal/pkg/request"
	"github.com/dollarkillerx/inventory/internal/pkg/response"
	"github.com/dollarkillerx/inventory/internal/utils"
	"github.com/gin-gonic/gin"

	"log"
)

func (s *Server) wareHousing(ctx *gin.Context) {
	model := utils.GetAuthModel(ctx)

	var war request.Warehousing
	if err := ctx.ShouldBindJSON(&war); err != nil {
		log.Println(err)
		response.Return(ctx, errs.BadRequest)
		return
	}

	vg, err := s.storage.Good(war.Barcode, model.Account)
	if err != nil {
		response.Return(ctx, errs.NewError("4001", err.Error()))
		return
	}

	err = s.storage.WareHousing(vg.Goods.ID, vg.Barcode, model.Account, war.Cost, war.NumberProducts, war.Remark)
	if err != nil {
		log.Println(err)
		response.Return(ctx, errs.SqlSystemError)
		return
	}

	response.Return(ctx, gin.H{})
}

func (s *Server) outStock(ctx *gin.Context) {
	model := utils.GetAuthModel(ctx)

	var war request.OutStock
	if err := ctx.ShouldBindJSON(&war); err != nil {
		log.Println(err)
		response.Return(ctx, errs.BadRequest)
		return
	}

	vg, err := s.storage.Good(war.Barcode, model.Account)
	if err != nil {
		response.Return(ctx, errs.NewError("4001", err.Error()))
		return
	}

	err = s.storage.OutStock(vg.Goods.ID, vg.Barcode, model.Account, war.Cost, war.NumberProducts, war.Price, war.Remark)
	if err != nil {
		log.Println(err)
		response.Return(ctx, errs.SqlSystemError)
		return
	}

	response.Return(ctx, gin.H{})
}

func (s *Server) ioHistory(ctx *gin.Context) {
	model := utils.GetAuthModel(ctx)
	barcodes := ctx.Param("barcodes")
	history, err := s.storage.IOHistory(barcodes, model.Account)
	if err != nil {
		log.Println(err)
		response.Return(ctx, errs.BadRequest)
		return
	}

	response.Return(ctx, history)
}

func (s *Server) iORevoke(ctx *gin.Context) {

}
