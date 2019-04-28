package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Health struct {
}

type healthEndpoint struct {
	Status string
}

func newHealth(g *gin.RouterGroup) {
	h := &Health{}
	g.GET("", h.Get)
}

func (h *Health) Get(c *gin.Context) {
	c.JSON(http.StatusOK, healthEndpoint{
		Status: "Normal",
	})
}
