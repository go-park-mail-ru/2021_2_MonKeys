package grpcserver

import (
	"context"
	_userModels "dripapp/internal/dripapp/models"
	proto "dripapp/internal/microservices/auth/delivery/grpc/protobuff"
	_sessionModels "dripapp/internal/microservices/auth/models"
	"net"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type server struct {
	sessionManager _sessionModels.SessionRepository
	userRepo       _userModels.UserRepository
}

func NewAuthServerGRPC(gserver *grpc.Server, su _sessionModels.SessionRepository, uu _userModels.UserRepository) {
	authServer := &server{
		sessionManager: su,
		userRepo:       uu,
	}
	proto.RegisterAuthGrpcHandlerServer(gserver, authServer)
	reflection.Register(gserver)
}

func StartStaffGrpcServer(su _sessionModels.SessionRepository, uu _userModels.UserRepository, url string) {
	list, err := net.Listen("tcp", url)
	if err != nil {
		log.Err(err)
	}
	server := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
		}),
	)
	NewAuthServerGRPC(server, su, uu)
	_ = server.Serve(list)
}

func (s *server) GetFromSession(ctx context.Context, cookie *proto.Cookie) (*proto.Session, error) {
	userSession, err := s.sessionManager.GetSessionByCookie(cookie.Cookie)
	if err != nil {
		userSession = _sessionModels.Session{
			UserID: 0,
			Cookie: "",
		}
	}
	return transSessionIntoRPC(&userSession), err
}

func (s *server) GetById(ctx context.Context, session *proto.Session) (*proto.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, session.UserID)
	return transUserIntoRPC(&user), err
}

func transSessionIntoRPC(session *_sessionModels.Session) *proto.Session {
	if session == nil {
		return nil
	}
	res := &proto.Session{
		Cookie: session.Cookie,
		UserID: session.UserID,
	}
	return res
}

func transUserIntoRPC(user *_userModels.User) *proto.User {
	if user == nil {
		return nil
	}
	res := &proto.User{
		ID:           user.ID,
		Email:        user.Email,
		Password:     user.Password,
		Name:         user.Name,
		Gender:       user.Gender,
		Prefer:       user.Prefer,
		FromAge:      user.FromAge,
		ToAge:        user.ToAge,
		Date:         user.Date,
		Age:          user.Age,
		Description:  user.Description,
		Imgs:         user.Imgs,
		Tags:         user.Tags,
		ReportStatus: user.ReportStatus,
	}
	return res
}
