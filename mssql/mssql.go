package mssql

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	// mssql driver
	_ "github.com/denisenkom/go-mssqldb"
)

var (
	// Conn mssql DB connection
	Conn   *sql.DB
	config *Config
)

// Config mssql配置
type Config struct {
	ConnectionTimeout  int    `json:"ConnectionTimeout"`
	MaxOpenConnections int    `json:"MaxOpenConnections"`
	Host               string `json:"Host"`
	Usr                string `json:"Usr"`
	Pass               string `json:"Pass"`
	DB                 string `json:"DB"`
	EncryptEnabled     bool   `json:"EncryptEnabled"`
}

// Startup 启动
func Startup(config *Config) bool {
	query := url.Values{}
	query.Add("connection timeout", fmt.Sprintf("%d", config.ConnectionTimeout))
	query.Add("database", config.DB)

	var connectionString string
	if config.EncryptEnabled {
		u := &url.URL{
			Scheme: "mssql",
			User:   url.UserPassword(config.Usr, config.Pass),
			Host:   config.Host,
			// Path:  instance, // if connecting to an instance instead of a port
			RawQuery: query.Encode(),
		}

		connectionString = u.String()
	} else {
		connectionString = fmt.Sprintf("sqlserver://%s:%s@%s?database=%s&connection+timeout=%d&log=63&encrypt=disable", config.Usr, config.Pass, config.Host, config.DB, config.ConnectionTimeout)
	}

	dbConn, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		log.Fatalln(config.DB, "Cannot connect: ", err.Error())
		return false
	}

	err = dbConn.Ping()
	if err != nil {
		log.Fatalln(config.DB, "Cannot ping: ", err.Error())
		return false
	}

	dbConn.SetMaxOpenConns(config.MaxOpenConnections)
	Conn = dbConn
	log.Printf("DB %s connected.\n", config.DB)
	return true
}

func printValue(pval *interface{}) {
	switch v := (*pval).(type) {
	case nil:
		fmt.Print("NULL")
	case bool:
		if v {
			fmt.Print("1")
		} else {
			fmt.Print("0")
		}
	case []byte:
		fmt.Print(string(v))
	case time.Time:
		fmt.Print(v.Format("2006-01-02 15:04:05.999"))
	default:
		fmt.Print(v)
	}
}
