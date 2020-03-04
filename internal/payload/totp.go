package payload

// TotpValidation request
type TotpValidation struct {
	Secret string `json:"secret"`
	Code   string `json:"code"`
}
