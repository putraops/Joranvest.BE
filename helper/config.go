package helper

import (
	"github.com/gin-gonic/gin"
)

type Config struct {
	DB *DBConfig
	C  *gin.Context
}

type DBConfig struct {
	ServerName string
	Dialect    string
	Username   string
	Password   string
	Name       string
	Charset    string
}

func GetConfig() *Config {
	//session := sessions.Default(Config.C)
	return &Config{
		DB: &DBConfig{
			ServerName: "localhost",
			Dialect:    "sqlserver",
			Username:   "sa",
			Password:   "P@ssw0rd",
			Name:       "jireh_laundry",
		},
	}
}
