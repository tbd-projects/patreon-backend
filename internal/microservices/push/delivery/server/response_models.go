package push_server

import (
	"patreon/internal/microservices/push/push/repository"
	"time"
)

//go:generate easyjson -all -disallow_unknown_fields response_models.go

//easyjson:json
type PushResponse struct {
	Type   string      `json:"type"`
	Push   interface{} `json:"push"`
	Date   time.Time   `json:"date"`
	Viewed bool        `json:"viewed"`
	Id     int64       `json:"push_id"`
}

//easyjson:json
type PushesResponse struct {
	Pushes []PushResponse `json:"pushes"`
}

func ToPushesResponse(pushes []repository.Push) PushesResponse {
	res := PushesResponse{}
	for _, ph := range pushes {
		res.Pushes = append(res.Pushes, PushResponse{
			Type:   ph.Type,
			Id:     ph.Id,
			Date:   ph.Date,
			Push:   ph.Push.Push,
			Viewed: ph.Viewed,
		})
	}
	return res
}
