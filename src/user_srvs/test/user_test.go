package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/ahlemarg/shop-srvs/src/user_srvs/proto"
	"google.golang.org/grpc"
)

var (
	conn       *grpc.ClientConn
	userClient proto.UserClient
	err        error
)

var (
	ids       []uint32
	mobiles   []string
	passwords []string
)

func init() {
	conn, err = grpc.Dial("127.0.0.1:50052", grpc.WithInsecure())
	if err != nil {
		str := fmt.Sprintf("failed to connect grpc server. %v\n", err)
		panic(str)
	}
	userClient = proto.NewUserClient(conn)
}

func TestCreateUser(t *testing.T) {
	for i := 1; i <= 9; i++ {
		resp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
			Mobile:   fmt.Sprintf("1859759682%d", i),
			NickName: fmt.Sprintf("testUser%d", i),
			Password: "admin123",
		})
		if err != nil {
			panic(err)
		}
		t.Errorf("ID: %d nickName: %s\n", resp.ID, resp.NickName)
	}
}

func TestGetUserList(t *testing.T) {
	resp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 4,
	})
	if err != nil {
		panic(err)
	}

	for _, user := range resp.Data {
		t.Errorf("mobile: %s nickName: %s PassWord: %s\n", user.Mobile, user.NickName, user.PassWord)
		ids = append(ids, user.ID)
		mobiles = append(mobiles, user.Mobile)
		passwords = append(passwords, user.PassWord)
	}
	t.Error("GetUserList over...\n\n")
}

func TestGetUserByID(t *testing.T) {
	t.Run("testing function GetUserList", TestGetUserList)
	for _, id := range ids {
		resp, err := userClient.QueryUserByID(context.Background(), &proto.GetUserID{
			ID: id,
		})
		if err != nil {
			panic(err)
		}

		t.Errorf("ID: %d nickName: %s\n", resp.ID, resp.NickName)
	}
}

func TestGetUserByMobile(t *testing.T) {
	t.Run("testing function GetUserList", TestGetUserList)
	for _, mobile := range mobiles {
		resp, err := userClient.QueryUserByMobile(context.Background(), &proto.GetUserMobile{
			Mobile: mobile,
		})
		if err != nil {
			panic(err)
		}

		t.Errorf("ID: %d nickName: %s\n", resp.ID, resp.NickName)
	}
}

func TestCheckPassWord(t *testing.T) {
	t.Run("testing function GetUserList", TestGetUserList)
	for _, password := range passwords {
		resp, err := userClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
			Password:         "admin123",
			EncrytedPassword: password,
		})
		if err != nil {
			panic(err)
		}

		t.Errorf("password verityed result: %v\n", resp.Success)
	}
}
