package config

import (
	"embed"
	"gopkg.in/yaml.v3"
	"io/fs"
	"pixivImages/build"
)

//go:embed *.yaml
var configFiles embed.FS

var config *Config

type Config struct {
	Redis      Redis      `yaml:"redis"`
	Mysql      Mysql      `yaml:"mysql"`
	Pixiv      Pixiv      `yaml:"pixiv"`
	MsGraph    MsGraph    `yaml:"ms-graph"`
	HttpServer HttpServer `yaml:"http-server"`
	Socks5     Socks5     `yaml:"socks5"`
	Logger     Logger     `yaml:"logger"`
}

type HttpServer struct {
	Port int `yaml:"port"`
}

type Redis struct {
	Pool     int    `yaml:"pool"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

type Mysql struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Db       string `yaml:"database"`
	Pool     int    `yaml:"pool"`
}

type MsGraph struct {
	ClientId   string   `yaml:"client_id"`
	SecretId   string   `yaml:"secret_id"`
	ReceiveUrl string   `yaml:"receive-url"`
	Scopes     []string `yaml:"scopes"`
}

type Socks5 struct {
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Logger struct {
	MaxSize    int    `yaml:"max-size"`
	MaxAge     int    `yaml:"max-age"`
	MaxBackups int    `yaml:"max-backups"`
	Path       string `yaml:"path"`
}

type Pixiv struct {
}

func Get() *Config {
	return config
}

func LoadConfig() {
	config = &Config{}

	file, err := fs.ReadFile(configFiles, "config.yaml")
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(file, config); err != nil {
		panic(err)
	}

	if build.Version == build.Prod {
		ProductionLoad()
	}
}

func ProductionLoad() {
	file, err := fs.ReadFile(configFiles, "config.prod.yaml")
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(file, config); err != nil {
		panic(err)
	}
}
