package client

import (
	"context"
	"patreon/internal/microservices/auth/sessions/models"

	proto "patreon/internal/microservices/auth/delivery/grpc/protobuf"

	"google.golang.org/grpc"
)

type SessionClient struct {
	sessionClient proto.AuthCheckerClient
}

func NewSessionClient(con *grpc.ClientConn) *SessionClient {
	client := proto.NewAuthCheckerClient(con)
	return &SessionClient{
		sessionClient: client,
	}
}

// Check Errors:
//		Status 401 "not authorized user"
func (c *SessionClient) Check(ctx context.Context, sessionID string) (models.Result, error) {
	protoSessionID := &proto.SessionID{ID: sessionID}
	res, err := c.sessionClient.Check(ctx, protoSessionID)
	if err != nil {
		return models.Result{}, err
	}
	return ConvertAuthServerRespond(res), err
}
func (c *SessionClient) Create(ctx context.Context, userID int64) (models.Result, error) {
	protoUserID := &proto.UserID{
		ID: userID,
	}
	res, err := c.sessionClient.Create(ctx, protoUserID)
	if err != nil {
		return models.Result{}, err
	}
	return ConvertAuthServerRespond(res), nil
}

func (c *SessionClient) Delete(ctx context.Context, sessionID string) error {
	protoSessionID := &proto.SessionID{
		ID: sessionID,
	}
	_, err := c.sessionClient.Delete(ctx, protoSessionID)
	if err != nil {
		return err
	}

	return nil
}
