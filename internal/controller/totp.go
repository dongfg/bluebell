package controller

import (
	"bytes"
	"github.com/dongfg/bluebell/internal/payload"
	"github.com/dongfg/bluebell/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type totpController struct {
	opts *totpControllerOptions
}

type totpControllerOptions struct {
	routerGroup *gin.RouterGroup
	totpService *service.TotpService
}

func newTotpController(opts *totpControllerOptions) {
	ctrl := &totpController{
		opts,
	}
	routerGroup := ctrl.opts.routerGroup
	routerGroup.POST("/generate", ctrl.generate)
	routerGroup.POST("/validate", ctrl.validate)
}

func (ctrl *totpController) generate(c *gin.Context) {
	totpService := ctrl.opts.totpService
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "error read request body")
		return
	}

	secret := buf.String()
	code, err := totpService.Generate(secret)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.String(http.StatusOK, code)
}

func (ctrl *totpController) validate(c *gin.Context) {
	totpService := ctrl.opts.totpService
	var validation payload.TotpValidation
	err := c.BindJSON(&validation)
	if err != nil {
		c.JSON(http.StatusBadRequest, failed(err.Error()))
		return
	}

	c.JSON(http.StatusOK, data(totpService.Validate(validation)))
}
