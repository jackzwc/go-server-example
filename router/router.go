package router

import (
	"github.com/gin-gonic/gin"
	"go-server-example/service"
)

func Routers(group *gin.RouterGroup) {
	Hello(group)
}

func Hello(group *gin.RouterGroup) {
	group.GET("/hello", service.Hello)
}
