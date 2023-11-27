package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger     LoggerConf     `mapstructure:"logger"`
	Storage    StorageConf    `mapstructure:"storage"`
	HTTPServer HTTPServerConf `mapstructure:"httpserver"`
	GRPCServer GRPCServerConf `mapstructure:"grpcserver"`
	AMQPConf   AMQPConf       `mapstructure:"amqp"`
	Scheduler  SchedulerConf  `mapstructure:"scheduler"`
}

type LoggerConf struct {
	Level string `mapstructure:"level"`
}

type StorageConf struct {
	DB      string `mapstructure:"db"`
	URL     string `mapstructure:"url"`
	Timeout int    `mapstructure:"timeout"`
}

type HTTPServerConf struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type GRPCServerConf struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type AMQPConf struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	QueueName string `mapstructure:"queue_name"`
}

type SchedulerConf struct {
	NotifyCheckPeriod int `mapstructure:"notify_check_period_s"`
	RemoveCheckPeriod int `mapstructure:"remove_check_period_s"`
}

func NewConfig(fileName string) Config {
	v := viper.New()
	v.SetConfigFile(fileName)

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Can't load config file")
		os.Exit(1)
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		fmt.Printf("Can't unmarshall config file")
		os.Exit(1)
	}

	return c
}
