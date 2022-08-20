package server

import (
	"github.com/dollarkillerx/inventory/internal/pkg/errs"
	"github.com/dollarkillerx/inventory/internal/pkg/response"
	"github.com/dollarkillerx/inventory/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) statistics(ctx *gin.Context) {
	model := utils.GetAuthModel(ctx)

	scs, err := s.storage.Statistics(model.Account)
	if err != nil {
		response.Return(ctx, errs.BadRequest)
		return
	}

	response.Return(ctx, scs)
}
