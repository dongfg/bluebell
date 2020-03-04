package controller

import (
	"github.com/dongfg/bluebell/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type healthController struct {
	opts *healthControllerOptions
}

type healthControllerOptions struct {
	routerGroup   *gin.RouterGroup
	healthService *service.HealthService
}

func newHealthController(opts *healthControllerOptions) {
	ctrl := &healthController{
		opts,
	}
	routerGroup := ctrl.opts.routerGroup
	routerGroup.GET("", ctrl.check)
}

func (ctrl *healthController) check(c *gin.Context) {
	healthService := ctrl.opts.healthService
	c.JSON(http.StatusOK, healthService.Check())
}
