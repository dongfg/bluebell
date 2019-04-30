package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthController struct {
}

type healthEndpoint struct {
	Status string
}

func newHealthController(g *gin.RouterGroup) {
	h := &HealthController{}
	g.GET("", h.Get)
}

func (h *HealthController) Get(c *gin.Context) {
	c.JSON(http.StatusOK, healthEndpoint{
		Status: "Normal",
	})
}
