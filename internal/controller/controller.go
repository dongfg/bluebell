package controller

import (
	"github.com/dongfg/bluebell/internal/service"
	"github.com/gin-gonic/gin"
	"time"
)

// Controller ref
type Controller struct {
	opts *Options
}

// Options dependency
type Options struct {
	Router        *gin.Engine
	HealthService *service.HealthService
	TotpService   *service.TotpService
	SeriesService *service.SeriesService
}

// New controller
func New(opts *Options) *Controller {
	return &Controller{opts}
}

// Register all controller
func (ctrl *Controller) Register() {
	router := ctrl.opts.Router

	newHealthController(&healthControllerOptions{
		routerGroup:   router.Group("/health"),
		healthService: ctrl.opts.HealthService,
	})
	newTotpController(&totpControllerOptions{
		routerGroup: router.Group("/totp"),
		totpService: ctrl.opts.TotpService,
	})
	newSeriesController(&seriesControllerOptions{
		routerGroup:   router.Group("/series"),
		seriesService: ctrl.opts.SeriesService,
	})
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
