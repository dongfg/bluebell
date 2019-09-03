package controller

import (
	"github.com/gin-gonic/gin"
	"time"
)

// Register all controller
func Register(r *gin.Engine) {
	newHealthController(r.Group("/health"))
	newSeriesController(r.Group("/series"))
	newTotpController(r.Group("/totp"))
}

func success() gin.H {
	return gin.H{
		"msg":       "success",
		"timestamp": time.Now().Unix(),
	}
}

func failed(msg string) gin.H {
	return gin.H{
		"msg":       msg,
		"timestamp": time.Now().Unix(),
	}
}

func data(data interface{}) gin.H {
	return gin.H{
		"msg":       "success",
		"data":      data,
		"timestamp": time.Now().Unix(),
	}
}
