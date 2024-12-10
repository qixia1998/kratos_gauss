package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
	v1 "user/api/user/v1"
	"user/internal/biz"
)

type UserService struct {
	v1.UnimplementedUserServer

	uc  *biz.UserUsecase
	log *log.Helper
}

// NewUserService new a greeter service.
func NewUserService(uc *biz.UserUsecase, logger log.Logger) *UserService {
	return &UserService{uc: uc, log: log.NewHelper(logger)}
}

// CreateUser create a user
func (u *UserService) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.CreateUserReply, error) {
	user, err := u.uc.Create(ctx, &biz.User{
		Mobile:   req.Mobile,
		Password: req.Password,
		Name:     req.Name,
	})
	if err != nil {
		return nil, err
	}

	userInfoRsp := v1.CreateUserReply{
		Id:       user.ID,
		Mobile:   user.Mobile,
		Password: user.Password,
		Name:     user.Name,
		Gender:   user.Gender,
	}

	return &userInfoRsp, nil
}

func (u *UserService) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*v1.UpdateUserReply, error) {
	user, err := u.uc.UpdateUser(ctx, &biz.User{
		ID:       req.Id,
		Mobile:   req.Mobile,
		Password: req.Password,
		Name:     req.Name,
		Gender:   req.Gender,
	})
	if err != nil {
		return nil, err
	}
	if user == false {
		return nil, err
	}
	return &v1.UpdateUserReply{Success: true}, nil
}

func UserResponse(user *biz.User) v1.GetUserReply {
	userInfoRsp := v1.GetUserReply{
		Id:       user.ID,
		Mobile:   user.Mobile,
		Password: user.Password,
		Name:     user.Name,
		Gender:   user.Gender,
	}
	return userInfoRsp
}

func (u *UserService) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.GetUserReply, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "get user info by Id")
	defer span.End()
	user, err := u.uc.UserById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	rsp := UserResponse(user)
	return &rsp, nil
}

func (u *UserService) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*v1.DeleteUserReply, error) {
	_, err := u.uc.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteUserReply{Success: true}, nil
}
