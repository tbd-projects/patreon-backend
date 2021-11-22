package server

import (
	"context"
	"net"
	"os"
	proto "patreon/internal/microservices/auth/delivery/grpc/protobuf"
	"patreon/internal/microservices/auth/sessions"
	prometheus_monitoring "patreon/pkg/monitoring/prometheus-monitoring"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
)

type AuthServer struct {
	grpcServer     *grpc.Server
	sessionManager sessions.SessionsManager
	logger         *logrus.Logger
}

func NewAuthGRPCServer(logger *logrus.Logger, grpcServer *grpc.Server, sessionManager sessions.SessionsManager) *AuthServer {
	server := &AuthServer{
		sessionManager: sessionManager,
		grpcServer:     grpcServer,
		logger:         logger,
	}
	return server
}

func (server *AuthServer) StartGRPCServer(listenUrl string) error {
	lis, err := net.Listen("tcp", listenUrl)
	server.logger.Infof("my listen url %s \n", listenUrl)

	if err != nil {
		server.logger.Errorf("AUTHSERVER\n")
		server.logger.Errorf("can not listen url: %s err :%v\n", listenUrl, err)
		return err
	}
	proto.RegisterAuthCheckerServer(server.grpcServer, server)

	go prometheus_monitoring.CreateNewMonitoringRouter(os.Getenv("sessions-service"))

	server.logger.Info("Start session service\n")
	return server.grpcServer.Serve(lis)
}

func (s *AuthServer) Check(ctx context.Context, sessionID *proto.SessionID) (*proto.Result, error) {
	res, err := s.sessionManager.Check(sessionID.ID)
	if err != nil {
		s.logger.Errorf("AUTHSERVER\n")
		s.logger.Errorf("can not check session with sessionID = %s, err = %v", sessionID.ID,
			err)
		return nil, err
	}

	return &proto.Result{
		UserID:    res.UserID,
		SessionID: res.UniqID,
	}, nil
}

func (s *AuthServer) Create(ctx context.Context, userID *proto.UserID) (*proto.Result, error) {
	res, err := s.sessionManager.Create(userID.ID)
	if err != nil {
		s.logger.Errorf("AUTHSERVER\n")
		s.logger.Errorf("can not create session with userID = %d, err = %v", userID.ID,
			err)
		return nil, err
	}
	return &proto.Result{
		UserID:    res.UserID,
		SessionID: res.UniqID,
	}, nil
}
func (s *AuthServer) Delete(ctx context.Context, sessionID *proto.SessionID) (*proto.Nothing, error) {
	err := s.sessionManager.Delete(sessionID.ID)
	if err != nil {
		s.logger.Errorf("AUTHSERVER\n")
		s.logger.Errorf("can not delete session with sessionID = %s, err = %v", sessionID.ID,
			err)
		return &proto.Nothing{Dummy: false}, err
	}
	return &proto.Nothing{
		Dummy: true,
	}, nil
}
