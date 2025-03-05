package util

import (
	"fmt"
	"path"
	"runtime"

	"github.com/spf13/viper"
)

var (
	ProjectRootPath = getoncurrentPath() + "/../"
)

func getoncurrentPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Dir(filename)
}

func CreateConfig(file string) *viper.Viper {
	config := viper.New()
	config.AddConfigPath(ProjectRootPath + "config/")
	config.SetConfigName(file)
	config.SetConfigType("yaml")

	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("找不到配置文件:%s", ProjectRootPath+"config/"+file+".yaml"))
		} else {
			panic(fmt.Errorf("解析配置文件%s出错:%s", ProjectRootPath+"config/"+file+".yaml", err))
		}
	}

	return config
}
