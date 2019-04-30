package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type healthController struct {
}

type healthEndpoint struct {
	Status string
}

func newHealthController(g *gin.RouterGroup) {
	c := &healthController{}
	g.GET("", c.healthCheck)
}

func (controller *healthController) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, healthEndpoint{
		Status: "Normal",
	})
}
