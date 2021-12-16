package utils

//go:generate easyjson -all -disallow_unknown_fields response_models.go

//easyjson:json
type PushResponse struct {
	Type string      `json:"type"`
	Push interface{} `json:"push"`
}
