package initialize

import (
	"github.com/ahlemarg/shop-srvs/src/user_srvs/global"
	"github.com/spf13/viper"
)

func GetEnv(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	debug := GetEnv("CONFIG_DEBUG")

	config_file := "src/user_srvs/config-pro.yaml"
	if debug {
		config_file = "src/user_srvs/config-debug.yaml"
	}

	v := viper.New()

	v.SetConfigFile(config_file)

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(&global.ServerInfo); err != nil {
		panic(err)
	}
}
