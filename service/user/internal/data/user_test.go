package data_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"user/internal/biz"
	"user/internal/data"
	"user/internal/testdata"
)

var _ = Describe("User", func() {
	var ro biz.UserRepo
	var uD *biz.User
	BeforeEach(func() {
		ro = data.NewUserRepo(Db, nil)
		// 这里你可以不引入外部组装好的数据，可以在这里直接写
		uD = testdata.User()
	})
	// 设置 It 块来添加单个规格
	It("CreateUser", func() {
		u, err := ro.CreateUser(ctx, uD)
		Ω(err).ShouldNot(HaveOccurred())
		// 组装的数据 mobile 为 13509876789
		Ω(u.Mobile).Should(Equal("13509876789")) // 手机号应该为创建的时候写入的手机号
	})

	// 设置 It 块来添加单个规格
	It("UpdateUser", func() {
		uD.Name = "gyl"
		uD.Gender = "female"
		user, err := ro.UpdateUser(ctx, uD)
		Ω(err).ShouldNot(HaveOccurred()) // 更新不应该出现错误
		Ω(user).Should(BeTrue())         // 结果应该为 true
	})
})
