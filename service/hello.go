package service

import (
	"github.com/gin-gonic/gin"
	"go-server-example/model"
	"net/http"
)

func Hello(c *gin.Context) {
	response := model.ResponseInfo{
		Code: 200,
		Data: 1,
		Msg:  "success",
	}
	c.JSON(http.StatusOK, response)
}
