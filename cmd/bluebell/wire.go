//+build wireinject

package main

import (
	"github.com/dongfg/bluebell/internal/config"
	"github.com/dongfg/bluebell/internal/controller"
	"github.com/dongfg/bluebell/internal/middleware"
	"github.com/dongfg/bluebell/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
	"github.com/google/wire"
	"net/http"
)

func initConfig() (*config.Config, error) {
	wire.Build(config.New)
	return &config.Config{}, nil
}

// initRouter init gin router and register controller
func initRouter(conf *config.Config) *gin.Engine {
	r := gin.Default()
	box := packr.New("static", "../../static")
	r.GET("/", func(c *gin.Context) {
		src, _ := box.Find("index.html")
		c.Data(http.StatusOK, "text/html; charset=utf-8", src)
	})
	r.GET("/swagger.yml", func(c *gin.Context) {
		src, _ := box.Find("swagger.yml")
		c.Data(http.StatusOK, "text/vnd.yaml; charset=utf-8", src)
	})
	r.Use(middleware.RateLimit())
	return r
}

func initController(conf *config.Config, r *gin.Engine) (*controller.Controller, error) {
	wire.Build(
		controller.WireSet,
		service.WireSet,
		service.OptionsSet,
	)
	return &controller.Controller{}, nil
}
