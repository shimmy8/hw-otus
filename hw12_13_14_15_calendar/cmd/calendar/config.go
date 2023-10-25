package main

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  LoggerConf  `toml:"logger"`
	Storage StorageConf `toml:"storage"`
}

type LoggerConf struct {
	Level string `toml:"level"`
}

type StorageConf struct {
	db string `toml:"db"`
}

func NewConfig(fileName string) Config {
	c := Config{}

	data, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Can't load condig file")
		os.Exit(1)
	}

	err = toml.Unmarshal(data, &c)
	if err != nil {
		fmt.Printf("Can't load condig file")
		os.Exit(1)
	}
	return c
}
