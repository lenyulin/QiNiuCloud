package ioc

import (
	"QiNiuCloud/QiNiuCloud/internal/web"
	"github.com/gin-gonic/gin"
)

func InitServer(interhdl *web.InteractiveHandler, model *web.ModelsHandler) *gin.Engine {
	server := gin.Default()
	interhdl.RegisiterInteractiveRoutes(server)
	model.RegisiterModelRoutes(server)
	return server
}
