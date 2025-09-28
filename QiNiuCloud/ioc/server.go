package ioc

import "github.com/gin-gonic/gin"

func InitServer() *gin.Engine {
	server := gin.Default()
	return server
}
