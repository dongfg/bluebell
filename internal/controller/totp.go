package controller

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"net/http"
	"time"
)

type totpController struct {
}

type totpValidation struct {
	Secret string `json:"secret"`
	Code   string `json:"code"`
}

func newTotpController(g *gin.RouterGroup) {
	c := &totpController{}
	g.POST("/generate", c.generate)
	g.POST("/validate", c.validate)
}

func (ctrl *totpController) generate(c *gin.Context) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "error read request body")
		return
	}

	secret := buf.String()
	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.String(http.StatusOK, code)
}

func (ctrl *totpController) validate(c *gin.Context) {
	var validation totpValidation
	err := c.BindJSON(&validation)
	if err != nil {
		c.JSON(http.StatusBadRequest, failed(err.Error()))
		return
	}

	c.JSON(http.StatusOK, data(totp.Validate(validation.Code, validation.Secret)))
}
