package server

import (
	"fmt"
	"time"

	"github.com/dollarkillerx/inventory/internal/pkg/errs"
	"github.com/dollarkillerx/inventory/internal/pkg/request"
	"github.com/dollarkillerx/inventory/internal/pkg/response"
	"github.com/dollarkillerx/inventory/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) userCenter(ctx *gin.Context) {
	var payload request.UserLogin
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		response.Return(ctx, errs.BadRequest)
		return
	}

	uc, err := s.storage.GetUserCenter(payload.Account)
	if err != nil {
		response.Return(ctx, errs.LoginFailed)
		return
	}

	if uc.Password != payload.Password {
		response.Return(ctx, errs.LoginFailed)
		return
	}

	token, err := utils.JWT.CreateToken(request.AuthJWT{
		Account:    uc.Account,
		Storehouse: uc.Storehouse,
	}, time.Now().Add(time.Hour*24*600).Unix())
	if err != nil {
		response.Return(ctx, errs.SystemError)
		return
	}

	fmt.Println(token)

	response.Return(ctx, response.UserLogin{
		JWT:        token,
		Storehouse: uc.Storehouse,
	})
}

func (s *Server) userInfo(ctx *gin.Context) {
	var payload request.UserInfo
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		response.Return(ctx, errs.BadRequest)
		return
	}

	uc, err := s.storage.GetUserCenter(payload.Account)
	if err != nil {
		response.Return(ctx, errs.BadRequest)
		return
	}

	response.Return(ctx, gin.H{
		"storehouse": uc.Storehouse,
	})
}
