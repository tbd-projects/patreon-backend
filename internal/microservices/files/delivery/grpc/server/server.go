package server

import (
	"bytes"
	"context"
	"net"
	grpc2 "patreon/internal/microservices/files/delivery/grpc"
	proto "patreon/internal/microservices/files/delivery/grpc/protobuf"
	repository_files "patreon/internal/microservices/files/files/repository/files"
	prometheus_monitoring "patreon/pkg/monitoring/prometheus-monitoring"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
)

type FileServer struct {
	grpcServer     *grpc.Server
	fileRepository repository_files.Repository
	logger         *logrus.Logger
}

func (server *FileServer) MakeUnused(ctx context.Context, path *proto.FilePath) (*proto.Nothing, error) {
	return &proto.Nothing{}, nil
}

func (server *FileServer) SaveFile(ctx context.Context, file *proto.SaveArgs) (*proto.FilePath, error) {
	fileReader := bytes.NewReader(file.GetFile().File)

	res, err := server.fileRepository.SaveFile(fileReader, repository_files.FileName(file.Name),
		repository_files.TypeFiles(file.Type))

	if err != nil {
		server.logger.Errorf("FILESERVER\n")
		server.logger.Errorf("can not save file with name =  %s, err = %v", file.Name,
			err)
		return nil, err
	}
	return &proto.FilePath{
		Path: res,
	}, nil
}

func (server *FileServer) LoadFile(ctx context.Context, path *proto.FilePath) (*proto.File, error) {
	res, err := server.fileRepository.LoadFile(path.Path)
	if err != nil {
		server.logger.Errorf("FILESERVER\n")
		server.logger.Errorf("can not load file with path =  %s, err = %v", path.Path, err)
		return nil, err
	}
	convertRes := grpc2.StreamToByte(res)

	return &proto.File{
		File: convertRes,
	}, nil
}

func NewFilesGRPCServer(logger *logrus.Logger, grpcServer *grpc.Server, repo repository_files.Repository) *FileServer {
	server := &FileServer{
		fileRepository: repo,
		grpcServer:     grpcServer,
		logger:         logger,
	}
	return server
}

func (server *FileServer) StartGRPCServer(listenUrl string) error {
	lis, err := net.Listen("tcp", listenUrl)
	server.logger.Infof("my listen url %s \n", listenUrl)

	if err != nil {
		server.logger.Errorf("FILESERVER\n")
		server.logger.Errorf("can not listen url: %s err :%v\n", listenUrl, err)
		return err
	}
	proto.RegisterFileServiceServer(server.grpcServer, server)

	go prometheus_monitoring.CreateNewMonitoringRouter("files-service")

	server.logger.Info("Start file service\n")
	return server.grpcServer.Serve(lis)
}
