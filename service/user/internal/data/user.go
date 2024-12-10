package data

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
	"user/internal/biz"
)

// User 定义数据表结构体
type User struct {
	ID        int64     `gorm:"primarykey"`
	Mobile    string    `gorm:"index:idx_mobile;unique;type:varchar(255);not null"`
	Password  string    `gorm:"type:varchar(255);not null "` // 用户密码的保存需要注意是否加密
	Name      string    `gorm:"type:varchar(255) comment '用户昵称'"`
	Gender    string    `gorm:"column:gender;default:male;type:char(10)"`
	CreatedAt time.Time `gorm:"column:create_at"`
	UpdatedAt time.Time `gorm:"column:update_at"`
}
type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo . 这里需要注意，上面 data 文件 wire 注入的是此方法，方法名不要写错了
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// CreateUser .
func (r *userRepo) CreateUser(ctx context.Context, u *biz.User) (*biz.User, error) {
	var user User
	// 验证是否已经创建
	result := r.data.db.Where(&biz.User{Mobile: u.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}

	user.Mobile = u.Mobile
	user.Name = u.Name
	user.Password = encrypt(u.Password) // 密码加密
	res := r.data.db.Create(&user)
	if res.Error != nil {
		return nil, status.Errorf(codes.Internal, res.Error.Error())
	}
	return &biz.User{
		ID:       user.ID,
		Mobile:   user.Mobile,
		Password: user.Password,
		Name:     user.Name,
		Gender:   user.Gender,
	}, nil
}

// UpdateUser .
func (r *userRepo) UpdateUser(ctx context.Context, user *biz.User) (bool, error) {
	var userInfo User
	result := r.data.db.Where(&User{ID: user.ID}).First(&userInfo)
	if result.RowsAffected == 0 {
		return false, status.Errorf(codes.NotFound, "用户不存在")
	}

	userInfo.Name = user.Name
	userInfo.Password = encrypt(user.Password)
	userInfo.Mobile = user.Mobile
	userInfo.UpdatedAt = time.Now()
	userInfo.Gender = user.Gender

	res := r.data.db.Save(&userInfo)
	if res.Error != nil {
		return false, status.Errorf(codes.Internal, res.Error.Error())
	}

	return true, nil
}

func (r *userRepo) DeleteUser(ctx context.Context, id int64) (bool, error) {
	var user User
	result := r.data.db.Where(&biz.User{ID: id}).First(&user)
	if result.RowsAffected == 0 {
		return false, status.Errorf(codes.NotFound, "用户不存在")
	}
	res := r.data.db.Delete(&user)
	if res.Error != nil {
		return false, status.Errorf(codes.Internal, res.Error.Error())
	}
	return true, nil
}

// GetUserById .
func (r *userRepo) GetUserById(ctx context.Context, Id int64) (*biz.User, error) {
	var user User
	result := r.data.db.Where(&User{ID: Id}).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	re := modelToResponse(user)
	return &re, nil
}

// ModelToResponse 转换 user 表中所有字段的值
func modelToResponse(user User) biz.User {
	userInfoRsp := biz.User{
		ID:        user.ID,
		Mobile:    user.Mobile,
		Password:  user.Password,
		Name:      user.Name,
		Gender:    user.Gender,
		CreatedAt: user.CreatedAt,
	}
	return userInfoRsp
}

// Password encryption
func encrypt(psd string) string {
	options := &password.Options{SaltLen: 16, Iterations: 10000, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode(psd, options)
	return fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
}
