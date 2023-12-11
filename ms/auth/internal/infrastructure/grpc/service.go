package grpc

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	pb "github.com/jessicatarra/greenlight/api/proto"
	"github.com/jessicatarra/greenlight/ms/auth/internal/domain"
)

type Service interface {
	validateAuthToken(req *pb.ValidateAuthTokenRequest) (*pb.User, error)
	userPermission(req *pb.UserPermissionRequest) (*pb.Empty, error)
}

type Server struct {
	Appl domain.Appl
	pb.UnimplementedMyServiceServer
}

func NewGRPCServer(appl domain.Appl) *Server {
	return &Server{
		Appl: appl,
	}
}

func (g *Server) validateAuthToken(req *pb.ValidateAuthTokenRequest) (*pb.User, error) {
	user, err := g.Appl.ValidateAuthTokenUseCase(req.Token)
	if err != nil {
		return nil, err
	}

	//createdAt, err := ptypes.TimestampProto(user.CreatedAt)
	createdAt := &timestamp.Timestamp{
		Seconds: user.CreatedAt.Unix(),
		Nanos:   int32(user.CreatedAt.Nanosecond()),
	}
	if err != nil {
		return nil, err
	}

	return &pb.User{
		Id:             user.ID,
		CreatedAt:      createdAt,
		Name:           user.Name,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		Activated:      user.Activated,
		Version:        int32(user.Version),
	}, nil
}

func (g *Server) userPermission(req *pb.UserPermissionRequest) (*pb.Empty, error) {
	err := g.Appl.UserPermissionUseCase(req.Code, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
