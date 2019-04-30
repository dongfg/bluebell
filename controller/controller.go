package controller

import "github.com/gin-gonic/gin"

func Register(r *gin.Engine) {
	newHealthController(r.Group("/health"))
}
