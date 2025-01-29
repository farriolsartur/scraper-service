package cfg

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type CacheCfg struct {
	Host string `yaml:"host" env:"CACHE_DB_HOST"`
	Port string `yaml:"port" env:"CACHE_DB_PORT"`
	Psw  string `yaml:"psw" env:"CACHE_DB_PSW"`
}

func LoadConfig(file string, config interface{}) error {
	return cleanenv.ReadConfig(file, config)
}
