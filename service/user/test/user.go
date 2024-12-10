package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	v1 "user/api/user/v1"
)

var userClient v1.UserClient
var conn *grpc.ClientConn

func main() {
	Init()

	TestCreateUser() // 创建用户

	conn.Close()
}

// Init 初始化 grpc 链接 注意这里链接的 端口
func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:9000", grpc.WithInsecure())
	if err != nil {
		panic("grpc link err" + err.Error())
	}
	userClient = v1.NewUserClient(conn)
}

func TestCreateUser() {

	rsp, err := userClient.CreateUser(context.Background(), &v1.CreateUserRequest{
		Mobile:   fmt.Sprintf("1388888888%d", 1),
		Password: "admin123",
		Name:     fmt.Sprintf("YWWW%d", 1),
	})
	if err != nil {
		panic("grpc 创建用户失败" + err.Error())
	}
	fmt.Println(rsp.Id)
}
