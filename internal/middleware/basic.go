package middleware

import (
	"github.com/dollarkillerx/inventory/internal/pkg/enum"
	"github.com/dollarkillerx/inventory/internal/pkg/errs"
	"github.com/dollarkillerx/inventory/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"

	"log"
)

// HttpRecover recover
func HttpRecover() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("recover panic: ", err)
				response.Return(ctx, errs.SystemError)
			}
		}()
	}
}

// SetBasicInformation 设置基础信息
func SetBasicInformation() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(enum.RequestID.String(), xid.New().String())
	}
}
