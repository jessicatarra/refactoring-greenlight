package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	pb "github.com/jessicatarra/greenlight/api/proto"
	"github.com/jessicatarra/greenlight/ms/auth/internal/domain"
)

type Service interface {
	ValidateAuthToken(ctx context.Context, request *pb.ValidateAuthTokenRequest) (*pb.User, error)
	UserPermission(ctx context.Context, request *pb.UserPermissionRequest) (*empty.Empty, error)
}

//
//func (s Service) mustEmbedUnimplementedAuthGRPCServiceServer() {
//	//TODO implement me
//	panic("implement me")
//}

type Server struct {
	Appl domain.Appl
	pb.UnimplementedAuthGRPCServiceServer
}

func NewGRPCServer(appl domain.Appl) *Server {
	return &Server{
		Appl: appl,
	}
}

func (s Server) ValidateAuthToken(ctx context.Context, request *pb.ValidateAuthTokenRequest) (*pb.User, error) {
	user, err := s.Appl.ValidateAuthTokenUseCase(request.Token)
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

func (s Server) UserPermission(ctx context.Context, request *pb.UserPermissionRequest) (*empty.Empty, error) {
	err := s.Appl.UserPermissionUseCase(request.Code, request.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
