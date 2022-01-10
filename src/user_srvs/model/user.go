package model

import (
	"time"

	base "github.com/ahlemarg/shop-srvs/src/user_srvs/model/base"
)

type User struct {
	base.BaseModel
	Mobile   string     `gorm:"column:mobile;type:varchar(11);index:idx_mobile;unique;not null;"`
	PassWord string     `gorm:"column:pass_word;type:varchar(128);not null;"`
	NickName string     `gorm:"column:nick_name;type:varchar(20);"`
	Birthday *time.Time `gorm:"column:birthday;"`
	Sex      string     `gorm:"column:sex;type:varchar(6);default:mela;comment:mela:男性 femela:女性;"`
	Role     uint       `gorm:"column:role;type:int;default:1;comment:1:普通用户 2:管理员用户;"`
}

func (User) TableName() string {
	return "shop_user_srvs"
}
