package web

import (
	"QiNiuCloud/QiNiuCloud/internal/service"
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *InteractiveHandler) RegisiterInteractiveRoutes(server *gin.Engine) {
	ug := server.Group("/interactive")
	ug.POST("/link", h.IncrLinkCnt)
	ug.POST("/download", h.IncrDownloadCnt)
}

type InteractiveHandler struct {
	server *gin.Engine
	svc    service.InteractiveService
	l      logger.ZapLogger
}

func NewInteractiveHandler(svc service.InteractiveService) *InteractiveHandler {
	return &InteractiveHandler{
		svc: svc,
	}
}

type Link struct {
	token string
	hash  string
}

func (h *InteractiveHandler) IncrLinkCnt(ctx *gin.Context) {
	var r Link
	if err := ctx.ShouldBind(&r); err != nil {
		h.l.Debug(err.Error())
		ctx.JSON(http.StatusOK, Result{
			Code: http.StatusBadRequest,
			Msg:  "Internal Server Error",
		})
		return
	}
	err := h.svc.IncrLinkCnt(ctx, r.token, r.hash)
	if err != nil {
		h.l.Debug(err.Error())
		ctx.JSON(http.StatusOK, Result{
			Code: http.StatusBadRequest,
			Msg:  "Internal Server Error",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: http.StatusOK,
		Msg:  "Request Success",
	})
	return
}
func (h *InteractiveHandler) IncrDownloadCnt(ctx *gin.Context) {

}
