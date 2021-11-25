package client

import (
	proto "patreon/internal/microservices/auth/delivery/grpc/protobuf"
	"patreon/internal/microservices/auth/sessions/models"
)

func ConvertAuthServerRespond(result *proto.Result) models.Result {
	if result == nil {
		return models.Result{}
	}
	res := models.Result{
		UserID: result.UserID,
		UniqID: result.SessionID,
	}
	return res
}
