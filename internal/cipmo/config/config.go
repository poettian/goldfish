package config

import (
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

type Config struct {
	Tapd   Tapd   `mapstructure:"cipmo"`
	Feishu Feishu `mapstructure:"feishu"`
	Sqlite Sqlite `mapstructure:"sqlite"`
}

type Tapd struct {
	AppId       string `mapstructure:"app_id"`
	AppSecret   string `mapstructure:"app_secret"`
	WorkspaceId string `mapstructure:"workspace_id"`
	PageSize    int    `mapstructure:"page_size"`
}

type Feishu struct {
	AppId             string `mapstructure:"app_id"`
	AppSecret         string `mapstructure:"app_secret"`
	TenantAccessToken string `mapstructure:"tenant_access_token"`
	DocsToken         string `mapstructure:"docs_token"`
	StoryTableId      string `mapstructure:"story_table_id"`
	BugTableId        string `mapstructure:"bug_table_id"`
	PageSize          int    `mapstructure:"page_size"`
}

type Sqlite struct {
	file string `mapstructure:"file"`
}

var c Config
var once sync.Once

// GetTapdConfig 获取一个配置的单例
func GetTapdConfig() *Config {
	once.Do(func() {
		initConfig()
	})
	return &c
}

// initConfig 初始化配置
func initConfig() {
	viper.SetConfigName("cipmo")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./etc")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error read config file: %s \n", err))
	}
	err = viper.Unmarshal(&c)
	if err != nil {
		panic(fmt.Errorf("error Unmarshal config: %s \n", err))
	}
}
