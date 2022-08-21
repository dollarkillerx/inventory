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
	barcodes := ctx.Param("goods_id")
	history, err := s.storage.IOHistory(barcodes, model.Account)
	if err != nil {
		log.Println(err)
		response.Return(ctx, errs.BadRequest)
		return
	}

	response.Return(ctx, history)
}

func (s *Server) iORevoke(ctx *gin.Context) {
	model := utils.GetAuthModel(ctx)

	var payload request.IORevoke
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		log.Println(err)
		response.Return(ctx, errs.BadRequest)
		return
	}

	revoke, err := s.storage.IORevoke(payload.OrderID, model.Account)
	if err != nil {
		log.Println(err)
		response.Return(ctx, errs.NewError("403", "當前訂單可能存在子訂單 無權刪除"))
		return
	}

	response.Return(ctx, revoke)
}

func (s *Server) ioList(ctx *gin.Context) {
	model := utils.GetAuthModel(ctx)

	var payload request.IOList
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		response.Return(ctx, errs.BadRequest)
		return
	}

	count, ids, err := s.storage.IOList(model.Account, payload.Limit, payload.Offset)
	if err != nil {
		log.Println(err)
		response.Return(ctx, errs.SqlSystemError)
		return
	}

	response.Return(ctx, response.IOList{
		Count: count,
		Items: ids,
	})
}
