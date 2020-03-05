package service

import (
	"fmt"
	"github.com/dongfg/bluebell/internal/config"
	"github.com/dongfg/bluebell/internal/payload"
	"net/http"
)

// HealthService ref
type HealthService struct {
	opts *HealthServiceOptions
}

// HealthServiceOptions service dependency
type HealthServiceOptions struct {
	Conf *config.Config
}

// NewHealthService instance
func NewHealthService(opts *HealthServiceOptions) *HealthService {
	return &HealthService{opts}
}

// Check health status
func (svc *HealthService) Check() payload.HealthEndpoint {
	domain := svc.opts.Conf.Series.Domain
	health := payload.HealthEndpoint{
		Status: "Normal",
		Series: payload.SeriesHealth{
			Domain: domain,
			Status: "Normal",
		},
	}

	r, err := http.Get(fmt.Sprintf("http://%s", domain))
	if err != nil {
		health.Series.Status = err.Error()
	} else {
		health.Series.Status = r.Status
		defer r.Body.Close()
	}
	return health
}
