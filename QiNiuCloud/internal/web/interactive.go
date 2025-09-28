package web

import (
	"QiNiuCloud/QiNiuCloud/internal/service"
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *InteractiveHandler) RegisiterInteractiveRoutes(server *gin.Engine) {
	server.POST("/interactive")
}

type InteractiveHandler struct {
	server *gin.Engine
	svc    service.InteractiveService
	l      logger.LoggerV1
}

func NewInteractiveHandler(svc service.InteractiveService) *InteractiveHandler {
	return &InteractiveHandler{
		svc: svc,
	}
}

type Link struct {
	op    InteractiveOp
	token string
	hash  string
}
type InteractiveOp string

var (
	IncrLinkCntOp                 InteractiveOp = "intr_like_cnt"
	IncrDownloadCntOp             InteractiveOp = "intr_download_cnt"
	IncrCloseAfterDownloadedCntOp InteractiveOp = "intr_close_after_downloaded_cnt"
)
var (
	ErrUnknownInteractiveOp = errors.New("unknown InteractiveOp")
)

func (h *InteractiveHandler) Interactive(ctx *gin.Context) {
	var r Link
	if err := ctx.ShouldBind(&r); err != nil {
		h.l.Debug(err.Error())
		ctx.JSON(http.StatusOK, Result{
			Code: http.StatusBadRequest,
			Msg:  "Internal Server Error",
		})
		return
	}
	var err error
	switch r.op {
	case IncrLinkCntOp:
		err = h.svc.IncrLinkCnt(ctx, r.token, r.hash)
	case IncrDownloadCntOp:
		err = h.svc.IncrDownloadCnt(ctx, r.token, r.hash)
	case IncrCloseAfterDownloadedCntOp:
		err = h.svc.IncrCloseAfterDownloadedCnt(ctx, r.token, r.hash)
	default:
		err = ErrUnknownInteractiveOp
	}
	if err != nil {
		h.l.Debug(err.Error(), logger.String("op", string(r.op)), logger.String("token", r.token), logger.String("hash", r.hash))
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
