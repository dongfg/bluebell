package controller

import (
	"fmt"
	"github.com/dongfg/bluebell/internal/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

type healthController struct {
}

type healthEndpoint struct {
	Status string
	Series struct {
		Domain string
		Status string
	}
}

func newHealthController(g *gin.RouterGroup) {
	c := &healthController{}
	g.GET("", c.healthCheck)
}

func (ctrl *healthController) healthCheck(c *gin.Context) {
	health := healthEndpoint{
		Status: "Normal",
		Series: struct {
			Domain string
			Status string
		}{Domain: config.Basic.Series.Domain, Status: "Normal"},
	}

	r, err := http.Get(fmt.Sprintf("http://%s", config.Basic.Series.Domain))
	if err != nil {
		health.Series.Status = err.Error()
	} else {
		health.Series.Status = r.Status
		defer r.Body.Close()
	}

	c.JSON(http.StatusOK, health)
}
