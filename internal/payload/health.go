package payload

// HealthEndpoint health status
type HealthEndpoint struct {
	Status string       `json:"status"`
	Series SeriesHealth `json:"series"`
}

// Series module health
type SeriesHealth struct {
	Domain string `json:"domain"`
	Status string `json:"status"`
}
