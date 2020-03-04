package service

import (
	"github.com/dongfg/bluebell/internal/payload"
	"github.com/pquerna/otp/totp"
	"log"
	"time"
)

// TotpService ref
type TotpService struct {
}

// NewTotpService
func NewTotpService() *TotpService {
	return &TotpService{}
}

// Generate totp code by secret
func (svc *TotpService) Generate(secret string) (string, error) {
	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		log.Print(err)
	}
	return code, err
}

// Validate code by secret
func (svc *TotpService) Validate(validation payload.TotpValidation) bool {
	return totp.Validate(validation.Code, validation.Secret)
}
