package router

import (
	"github.com/gin-gonic/gin"
	"go-server-example/config"
	"go-server-example/log"
	"go-server-example/middleware"
	"strconv"
)

func Server() {
	// Init log
	log.InitLog()

	config.InitConfig()

	gin.SetMode(gin.DebugMode)
	router := gin.New()
	// Use middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Cors())

	// Route Prefix
	apiGroup := router.Group("test/")
	// Register Routes
	Routers(apiGroup)
	// Start Server
	err := router.Run(":" + strconv.FormatInt(config.AppConfig.Server.Port, 10))
	if err != nil {
		panic("start server err:" + err.Error())
	}
}
