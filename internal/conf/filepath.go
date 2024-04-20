package conf

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Binding(config interface{}, filePath string) error {
	v := viper.New()
	v.SetConfigFile(filePath)
	if err := v.ReadInConfig(); err != nil {
		return err
	}
	// 把读取到的配置信息反序列化到 SysConfig 变量中
	if err := v.Unmarshal(&config); err != nil {
		return err
	}
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		viper.Unmarshal(&config)
	})
	return nil
}
