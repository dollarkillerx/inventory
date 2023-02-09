package server

import (
	"github.com/dollarkillerx/inventory/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (s *Server) router() {
	s.app.Use(gin.Logger())
	s.app.POST("/api/v1/login", s.userCenter)
	s.app.POST("/api/v1/user_info", s.userInfo)

	s.app.GET("/export", s.export)
	s.app.GET("/export_core/:account/:password", s.coreExport)
	//s.app.GET("/good/:store")

	v1 := s.app.Group("/api/v1", middleware.UAAuthorization())
	{
		v1.GET("/search", s.search)                  // 搜索商品
		v1.GET("/goods", s.goods)                    // 所有商品
		v1.GET("/good/:barcodes", s.good)            // 更具条码查询单个商品
		v1.POST("/good", s.addGood)                  // 添加商品
		v1.POST("/good/update", s.upGood)            // 更新商品
		v1.POST("/upload", s.uploadFile)             // 上传图片
		v1.POST("/warehousing", s.wareHousing)       // 入库
		v1.POST("/out_stock", s.outStock)            // 出库
		v1.GET("/io_history/:goods_id", s.ioHistory) // io history 出入庫記錄
		v1.POST("/io_revoke", s.iORevoke)            // 撤銷出入庫記錄

		v1.GET("/statistics", s.statistics)
		v1.POST("/io_list", s.ioList)
	}
}
