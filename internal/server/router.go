package server

import "github.com/dollarkillerx/inventory/internal/middleware"

func (s *Server) router() {
	s.app.POST("/api/v1/login", s.userCenter)
	s.app.POST("/api/v1/user_info", s.userInfo)

	v1 := s.app.Group("/api/v1", middleware.UAAuthorization())
	{
		v1.GET("/goods", s.Goods)
		v1.GET("/good/:barcodes", s.Good)
		v1.POST("/good", s.AddGood)
		v1.POST("/upload", s.UploadFile)
		v1.POST("/warehousing", s.wareHousing) // 入库
		v1.POST("/out_stock", s.outStock)      // 出库
	}
}
