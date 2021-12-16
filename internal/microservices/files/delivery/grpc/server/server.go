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

func (s *FileServer) SaveFile(ctx context.Context, file *proto.SaveArgs) (*proto.FilePath, error) {
	s.logger.Infof("FILESERVER - SaveFile: call with fileName = %v, fileType = %v\n", file.Name, file.Type)

	fileReader := bytes.NewReader(file.GetFile().File)

	res, err := s.fileRepository.SaveFile(fileReader, repository_files.FileName(file.Name),
		repository_files.TypeFiles(file.Type))

	if err != nil {
		s.logger.Errorf("FILESERVER\n")
		s.logger.Errorf("can not save file with name =  %s, err = %v", file.Name,
			err)
		return nil, err
	}
	s.logger.Infof("FILESERVER - SaveFile: correctly work resPath = %v\n", res)

	return &proto.FilePath{
		Path: res,
	}, nil
}

func (s *FileServer) LoadFile(ctx context.Context, path *proto.FilePath) (*proto.File, error) {
	s.logger.Infof("FILESERVER - LoadFile: call with path = %v\n", path.Path)

	res, err := s.fileRepository.LoadFile(path.Path)
	if err != nil {
		s.logger.Errorf("FILESERVER\n")
		s.logger.Errorf("can not load file with path = %s, err = %v", path.Path, err)
		return nil, err
	}
	convertRes := grpc2.StreamToByte(res)

	s.logger.Infof("FILESERVER - LoadFile: correctly work res = io.Reader \n")

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

func (s *FileServer) StartGRPCServer(listenUrl string) error {
	lis, err := net.Listen("tcp", listenUrl)
	s.logger.Infof("my listen url %s \n", listenUrl)

	if err != nil {
		s.logger.Errorf("FILESERVER\n")
		s.logger.Errorf("can not listen url: %s err :%v\n", listenUrl, err)
		return err
	}
	proto.RegisterFileServiceServer(s.grpcServer, s)

	go prometheus_monitoring.CreateNewMonitoringRouter("files-service")

	s.logger.Info("Start file service\n")
	return s.grpcServer.Serve(lis)
}
