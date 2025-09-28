package main

import (
	consumer "QiNiuCloud/QiNiuCloud/internal/events/comsumer"
	"github.com/gin-gonic/gin"
)

type App struct {
	server    *gin.Engine
	consumers []consumer.Consumer
}
