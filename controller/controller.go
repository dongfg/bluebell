package controller

import "github.com/gin-gonic/gin"

// Register all controller
func Register(r *gin.Engine) {
	newHealthController(r.Group("/health"))
	newSeriesController(r.Group("/series"))
}
