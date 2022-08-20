package server

import (
	"github.com/dollarkillerx/inventory/internal/pkg/errs"
	"github.com/dollarkillerx/inventory/internal/pkg/models"
	"github.com/dollarkillerx/inventory/internal/pkg/request"
	"github.com/dollarkillerx/inventory/internal/pkg/response"
	"github.com/dollarkillerx/inventory/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"

	"log"
	"strings"
)

func (s *Server) Goods(ctx *gin.Context) {
	//model := utils.GetAuthModel(ctx)
}

func (s *Server) Search(ctx *gin.Context) {
	model := utils.GetAuthModel(ctx)

	keyword := strings.TrimSpace(ctx.Query("keyword"))
	search, err := s.storage.Search(keyword, model.Account)
	if err != nil {
		response.Return(ctx, errs.BadRequest)
		return
	}

	response.Return(ctx, search)
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

func (s *Server) AddGood(ctx *gin.Context) {
	model := utils.GetAuthModel(ctx)

	var good request.AddGoods
	if err := ctx.ShouldBindJSON(&good); err != nil {
		log.Println(err)
		response.Return(ctx, errs.BadRequest)
		return
	}

	if len(good.Barcode) != 13 {
		response.Return(ctx, errs.NewError("40001", "条形码不合法 长度不为 13"))
		return
	}

	err := s.storage.DB().Model(&models.Goods{}).Create(&models.Goods{
		BasicModel: models.BasicModel{ID: xid.New().String()},
		Barcode:    good.Barcode,
		Name:       good.Name,
		Spec:       good.Spec,
		Cost:       good.Cost,
		Price:      good.Price,
		Brand:      good.Brand,
		MadeIn:     good.MadeIn,
		Img:        good.Img,
		ByAccount:  model.Account,
	}).Error
	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			response.Return(ctx, errs.NewError("4002", "商品以存在"))
			return
		}
		log.Println(err)
		response.Return(ctx, errs.SqlSystemError)
		return
	}

	response.Return(ctx, gin.H{})
}
