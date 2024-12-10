package data_test

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"testing"
	"user/internal/data"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// 测试 data 方法
func TestData(t *testing.T) {
	//  Ginkgo 测试通过调用 Fail(description string) 功能来表示失败
	// 使用 RegisterFailHandler 将此函数传递给 Gomega 。这是 Ginkgo 和 Gomega 之间的唯一连接点
	RegisterFailHandler(Fail)
	// 通知 Ginkgo 启动测试套件。如果您的任何 specs 失败，Ginkgo 将自动使 testing.T 失败。
	RunSpecs(t, "biz data test user")
}

var cleaner func()      // 定义删除 mysql 容器的回调函数
var Db *data.Data       // 用于测试的 data
var ctx context.Context // 上下文

// initialize  AutoMigrate gorm自动建表
func initialize(db *gorm.DB) error {
	err := db.AutoMigrate(
		&data.User{},
	)
	return errors.WithStack(err)
}
