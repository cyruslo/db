package config

import (
	"encoding/json"
	"fmt"
	"github.com/cyruslo/helper/gtimer"
	"log"
	"os"

	LOGGER "github.com/cyruslo/util/logger"
)

type gsconfig struct {
	MonitorEnabled     bool   `json:"MonitorEnabled"`
	Log4NetConfigPath  string `json:"Log4NetConfigPath"`
	LoggerName         string `json:"LoggerName"`
	Http               string `json:"Http"`
	Tcp                string `json:"Tcp"`
	GameId             int32  `json:"GameId"`
	DistributeMode     int32  `json:"DistributeMode"`
	Capacity           int32  `json:"Capacity"`
	Capacitys		  []int32 `json:"Capacitys"`
	AliveCheckDuration int32  `json:"AliveCheckDuration"`
	ConfigPath         string `json:"ConfigPath"`
	CheckServiceTicks  int64 `json:"CheckServiceTicks"`
	ClearServiceTicks  int64 `json:"ClearServiceTicks"`
}

var (
	// Params qp配置
	GsParams   *gsconfig = &gsconfig{}
	cfigPath string
)

func polling() {
	//	LOGGER.Info("========================= polling")
	loadConfig(cfigPath)
}

func InvalGameID(gameID int32) bool {

	return GsParams.GameId != gameID
}

// LoadConfig 解析配置
func loadConfig(filepath string) bool {
	var params = &gsconfig{}

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
	log.Printf("-------------------Configure are:-------------------")
	fmt.Println(params)

	GsParams = params
	
	/*arrlen := len(Params.JackpotConfig)

	for i:=0;i<arrlen;i++ {
		arrlenj := len(Params.JackpotConfig[i].Items)
		log.Printf("gameId(%d),%d.", Params.JackpotConfig[i].GameID, arrlenj)
	}*/
	return true
}

func OnStart(timerID int32, filepath, dbfilepath string) {
	cfigPath = filepath
	loadConfig(filepath)
	loadDBConfig(dbfilepath)

	gtimer.SetTimer(timerID, 600, true, true, polling)
}
