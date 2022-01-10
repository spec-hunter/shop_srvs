package handle

import (
	"context"
	"crypto/sha512"
	"fmt"
	"strings"
	"time"

	"github.com/ahlemarg/shop-srvs/src/user_srvs/global"
	"github.com/ahlemarg/shop-srvs/src/user_srvs/model"
	"github.com/ahlemarg/shop-srvs/src/user_srvs/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UserServer struct{}

func ModelInResponse(user model.User) (userInfo proto.UserInfoResponse) {
	userInfo = proto.UserInfoResponse{
		ID:       uint32(user.ID),
		Mobile:   user.Mobile,
		PassWord: user.PassWord,
		NickName: user.NickName,
		Sex:      user.Sex,
		Role:     uint32(user.Role),
	}

	// grpc的 message 类型不允许传递 nil 值
	if user.Birthday != nil {
		userInfo.Birthday = uint64(user.Birthday.Unix())
	}
	return
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (u *UserServer) GetUserList(ctx context.Context, reply *proto.PageInfo) (*proto.UserListResponse, error) {
	var users []model.User
	res := global.DB.Find(&users)
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "不存在任何的用户信息")
	}
	if res.Error != nil {
		return nil, status.Error(codes.Internal, res.Error.Error())
	}

	resp := &proto.UserListResponse{}
	resp.Total = uint32(res.RowsAffected)

	global.DB.Scopes(Paginate(int(reply.Pn), int(reply.PSize))).Find(&users)

	for _, user := range users {
		userInfoResp := ModelInResponse(user)
		resp.Data = append(resp.Data, &userInfoResp)
	}

	return resp, nil
}

func (u *UserServer) QueryUserByMobile(ctx context.Context, reply *proto.GetUserMobile) (*proto.UserInfoResponse, error) {
	var user model.User
	res := global.DB.Where(&model.User{Mobile: reply.Mobile}).Find(&user)
	// 用户无法查询
	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户 %v 不存在.", reply.Mobile)
	}
	// 用户查询出现错误
	if res.Error != nil {
		return nil, res.Error
	}
	userInfo := ModelInResponse(user)
	return &userInfo, nil
}

func (u *UserServer) QueryUserByID(ctx context.Context, reply *proto.GetUserID) (*proto.UserInfoResponse, error) {
	var user model.User
	res := global.DB.Find(&user, reply.ID)
	// 用户无法查询
	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户无法通过id进行查询.", reply.ID)
	}
	// 用户查询出现错误
	if res.Error != nil {
		return nil, res.Error
	}
	userInfo := ModelInResponse(user)
	return &userInfo, nil
}

func (u *UserServer) CreateUser(ctx context.Context, reply *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	res := global.DB.Where(&model.User{Mobile: reply.Mobile}).Find(&user)
	// 数据库中已存在用户
	if res.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在, 无法继续添加")
	}
	// 查询结果出现问题
	if res.Error != nil {
		return nil, status.Errorf(codes.Internal, res.Error.Error())
	}

	options := global.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha512.New,
	}
	salt, encoderPwd := global.Encode(reply.Password, &options)
	// 存储到数据库中的password字段格式
	// $[加密算法]$[盐值]$[加密后的密码]
	password := fmt.Sprintf("$%s$%s$%s", "pbkdf2-sha512", salt, encoderPwd)

	user.Mobile = reply.Mobile
	user.PassWord = password
	user.NickName = reply.NickName

	res = global.DB.Create(&user)
	if res.Error != nil {
		return nil, status.Errorf(codes.Internal, res.Error.Error())
	}

	responseInfo := ModelInResponse(user)
	return &responseInfo, nil
}

func (u *UserServer) UpdateUser(ctx context.Context, reply *proto.UpdateUserInfo) (*empty.Empty, error) {
	var user model.User
	res := global.DB.Find(&user, reply.ID)
	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "更新的用户不存在")
	}
	if res.Error != nil {
		return nil, status.Errorf(codes.Internal, res.Error.Error())
	}

	// 将传入的birthday进行转换
	birthDay := time.Unix(int64(reply.Birthday), 0)

	user.NickName = reply.NickName
	user.Sex = reply.Sex
	user.Birthday = &birthDay

	res = global.DB.Save(&user)
	if res.Error != nil {
		return nil, status.Errorf(codes.Internal, res.Error.Error())
	}

	return &empty.Empty{}, nil
}

func (u *UserServer) CheckPassword(ctx context.Context, reply *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	options := global.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha512.New,
	}
	passwordInfo := strings.Split(reply.EncrytedPassword, "$")
	check := global.Verify(reply.Password, passwordInfo[2], passwordInfo[3], &options)

	return &proto.CheckResponse{Success: check}, nil
}
