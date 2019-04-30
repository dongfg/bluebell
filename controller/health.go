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
	h := &healthController{}
	g.GET("", h.Get)
}

func (h *healthController) Get(c *gin.Context) {
	c.JSON(http.StatusOK, healthEndpoint{
		Status: "Normal",
	})
}
