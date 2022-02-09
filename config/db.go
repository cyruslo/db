package config

import (
	"encoding/json"
	LOGGER "github.com/cyruslo/util/logger"
	"os"
)

type GameDbConfig struct {
	GameID               int32  `json:"GameID"`
	Host                 string `json:"Host"`
	Database             string `json:"Database"`
	UserID               string `json:"UserID"`
	Password             string `json:"Password"`
	EncryptEnabled       bool   `json:"EncryptEnabled"`
	SQLConnectionTimeout int32  `json:"SQLConnectionTimeout"`
}

type redis_c struct {
	Addr			string		`json:"addr"`
	Pwd				string		`json:"pwd"`
}

type redis_s struct {
	Temporary      redis_c      `json:"temporary"`     
	Persistent     redis_c      `json:"persistent"`     
} 

type DbConfig struct {
	Games []GameDbConfig `json:"Games"`
	Redis redis_s `json:"redis"`
}


var (
	DBParams *DbConfig = &DbConfig{}
)

func GetRedisConnectionString() string {
    return DBParams.Redis.Temporary.Addr
}

func GetRedisConnectionPwd() string {
    return DBParams.Redis.Temporary.Pwd
}

func loadDBConfig(filepath string) bool {
	var params = &DbConfig{}

	f, err := os.Open(filepath)
	if err != nil {
		LOGGER.Error("failed to open config file:", filepath)
		return false
	}

	buf := make([]byte, 4096)
	n, err := f.Read(buf)
	if err != nil {
		LOGGER.Error("failed to read config file:", filepath)
		return false
	}

	json.Unmarshal(buf[:n], params)
	//LOGGER.Info("-------------------Configure are:-------------------")
	//LOGGER.Info(params)

	DBParams = params

	return true
}
