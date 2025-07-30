package usecases

import (
	"context"
	"football-team-management/internal/domain/user"
)

/*
	TODO: This usecase should have its own repository to fetch users from wherever they are stored
	but, to keep this template simple lets return a hardcoded user always
*/

type GetUser interface {
	Execute(ctx context.Context, username string) (user.User, error)
}

type getUserImpl struct {
}

func (g getUserImpl) Execute(ctx context.Context, username string) (user.User, error) {
	return user.User{
		Username: "testusername",
		Password: "$2a$04$LC8tT2q56k4LIOgxCZXLkODy0i3G9wLFpp2Itm6JMnnrJ5fpZmMH.", // a hash for "pass" string
		Role:     "admin",
	}, nil
}

func NewGetUserImpl() *getUserImpl {
	return &getUserImpl{}
}

var _ GetUser = (*getUserImpl)(nil)
