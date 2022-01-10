package global

import (
	"github.com/ahlemarg/shop-srvs/src/user_srvs/config"
	"gorm.io/gorm"
)

var (
	DB         *gorm.DB
	ServerInfo config.ServerConfig
)
