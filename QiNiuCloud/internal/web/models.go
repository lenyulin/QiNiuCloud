package web

import (
	"QiNiuCloud/QiNiuCloud/internal/service"
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ModelsHandler struct {
	server *gin.Engine
	svc    service.ModelsService
	l      logger.ZapLogger
}

func (h *ModelsHandler) RegisiterRoutes(server *gin.Engine) {
	ug := server.Group("/models")
	ug.GET("/generate", h.Generate)
}
func NewCommentHandler(svc service.ModelsService) *ModelsHandler {
	return &ModelsHandler{
		svc: svc,
	}
}

func (h *ModelsHandler) Generate(ctx *gin.Context) {
	type req struct {
		Description string `json:"description"`
	}
	var r req
	if err := ctx.ShouldBind(&r); err != nil {
		h.l.Debug(err.Error())
		ctx.JSON(http.StatusOK, Result{
			Code: http.StatusBadRequest,
			Msg:  "Internal Server Error",
		})
		return
	}
	if len(r.Description) == 0 || len(r.Description) > 200 {
		ctx.JSON(http.StatusOK, Result{
			Code: http.StatusBadRequest,
			Msg:  "Description Cannot less than 0 characters or longer than 200",
		})
		return
	}
	//usr := ctx.MustGet("user").(UserClaims)
	res, err := h.svc.Generate(ctx, r.Description)
	if err != nil {
		h.l.Debug(err.Error())
		ctx.JSON(http.StatusOK, Result{
			Code: http.StatusInternalServerError,
			Msg:  "Internal Server Error",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: http.StatusOK,
		Data: res,
	})
}
