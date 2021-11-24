package client

import (
	"context"
	_userModels "dripapp/internal/dripapp/models"
	proto "dripapp/internal/microservices/auth/delivery/grpc/protobuff"
	_sessionModels "dripapp/internal/microservices/auth/models"
	"fmt"

	"google.golang.org/grpc"
)

type SessionClient struct {
	client proto.AuthGrpcHandlerClient
}

func NewStaffClient(conn *grpc.ClientConn) *SessionClient {
	c := proto.NewAuthGrpcHandlerClient(conn)
	return &SessionClient{
		client: c,
	}
}

func (s *SessionClient) GetFromSession(ctx context.Context, cookie string) (_sessionModels.Session, error) {
	cook := proto.Cookie{
		Cookie: cookie,
	}
	userSession, err := s.client.GetFromSession(ctx, &cook)
	if err != nil {
		fmt.Println("Unexpected Error", err)
		return _sessionModels.Session{}, err
	}
	return transformSessionFromRPC(userSession), nil
}

func (s *SessionClient) GetById(ctx context.Context, session _sessionModels.Session) (_userModels.User, error) {
	sess := proto.Session{
		Cookie: session.Cookie,
		UserID: session.UserID,
	}
	user, err := s.client.GetById(ctx, &sess)
	return transformUserFromRPC(user), err
}

func transformSessionFromRPC(session *proto.Session) _sessionModels.Session {
	if session == nil {
		return _sessionModels.Session{}
	}
	res := _sessionModels.Session{
		Cookie: session.Cookie,
		UserID: session.UserID,
	}
	return res
}

func transformUserFromRPC(user *proto.User) _userModels.User {
	if user == nil {
		return _userModels.User{}
	}
	res := _userModels.User{
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
