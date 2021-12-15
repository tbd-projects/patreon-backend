package push_server

//easyjson:json
type PushResponse struct {
	Type string `json:"type"`
	Push interface{} `json:"push"`
}
