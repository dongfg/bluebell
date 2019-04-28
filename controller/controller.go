package controller

import "github.com/gin-gonic/gin"

func Setup(r *gin.Engine) {
	newHealth(r.Group("/health"))
}
