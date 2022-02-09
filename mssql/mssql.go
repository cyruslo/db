package mssql

import (
	"database/sql"
	"fmt"
	"net/url"

	LOGGER "github.com/cyruslo/util/logger"
	_ "github.com/denisenkom/go-mssqldb"

	gsCfg "github.com/cyruslo/db/config"
	_ "github.com/denisenkom/go-mssqldb"
)

// 错误码
const (
	Success = 0
	Failed  = -1
)

// 数据库返回
const (
	SQLSuccess = 1
)

// 充值来源(sourceType)
const (
	SourceIos     = 1
	SourceAndroid = 2
)

var (
	gameDBConns map[int32]*sql.DB
	//MProp表主键id,服务器实现自增<gameId, Id>
	mPropCurId map[int32]int64
)

// Startup 启动
func Startup() {
	if gameDBConns == nil {
		gameDBConns = make(map[int32]*sql.DB)
	}

	if nil == mPropCurId {
		mPropCurId = make(map[int32]int64)
	}

	arrlen := len(gsCfg.DBParams.Games)
	for i := 0; i < arrlen; i++ {
		dbCfg := gsCfg.DBParams.Games[i]

		dbConn := ConnectGameDB(dbCfg.Database, dbCfg.Host, dbCfg.UserID, dbCfg.Password, dbCfg.EncryptEnabled, dbCfg.SQLConnectionTimeout)
		if gameDBConns[dbCfg.GameID] == nil {
			gameDBConns[dbCfg.GameID] = new(sql.DB)
		}

		gameDBConns[dbCfg.GameID] = dbConn
	}
}

func GetSqlConn(gameID int32) *sql.DB {
	if _, ok := gameDBConns[gameID]; ok {
		return gameDBConns[gameID]
	}
	return nil
}

func GetMPropCurId() map[int32]int64 {

	return mPropCurId
}

func ConnectGameDB(database, host, user, pass string, EncryptEnabled bool, timeout int32) *sql.DB {
	query := url.Values{}
	query.Add("connection timeout", fmt.Sprintf("%d", timeout))
	query.Add("database", database)

	var connectionString string
	if EncryptEnabled == true {
		u := &url.URL{
			Scheme: "sqlserver",
			User:   url.UserPassword(user, pass),
			Host:   host,
			// Path:  instance, // if connecting to an instance instead of a port
			RawQuery: query.Encode(),
		}
		//	log.Printf("================= EncryptEnabled ")
		connectionString = u.String()
	} else {
		connectionString = fmt.Sprintf("sqlserver://%s:%s@%s?database=%s&connection+timeout=30&log=63&encrypt=disable", user, pass, host, database)
	}

	LOGGER.Info("connectionString:%s.", connectionString)
	dbConn, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		LOGGER.Error("Cannot connect: %s.", err.Error())
		return nil
	}

	LOGGER.Info("connectionString:%s.", connectionString)
	err = dbConn.Ping()
	if err != nil {
		LOGGER.Error("Cannot ping: %s.", err.Error())
		return nil
	}

	dbConn.SetMaxOpenConns(100)
	LOGGER.Info("DB %s (%s) connected.\n", database, host)
	return dbConn
}


func ConnectDB(database, host, user, pass string) *sql.DB {
	query := url.Values{}
	query.Add("connection timeout", fmt.Sprintf("%d", gsCfg.Params.SQLConnectionTimeout))
	query.Add("database", database)

	//log.Printf("================= connectDB %s,%s. ", database, host)
	var connectionString string
	if gsCfg.Params.EncryptEnabled {
		u := &url.URL{
			Scheme: "sqlserver",
			User:   url.UserPassword(user, pass),
			Host:   host,
			// Path:  instance, // if connecting to an instance instead of a port
			RawQuery: query.Encode(),
		}
	//	log.Printf("================= EncryptEnabled ")
		connectionString = u.String()
	} else {
		connectionString = fmt.Sprintf("sqlserver://%s:%s@%s?database=%s&connection+timeout=30&log=63&encrypt=disable", user, pass, host, database)
	}
	LOGGER.Info("connectionString:%s.", connectionString)
	dbConn, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		LOGGER.Error("Cannot connect: %s.", err.Error())
		return nil
	}

	err = dbConn.Ping()
	if err != nil {
		LOGGER.Error("Cannot ping: %s.", err.Error())
		return nil
	}

	dbConn.SetMaxOpenConns(100)
	LOGGER.Info("DB %s (%s) connected.\n", database, host)
	return dbConn
}
