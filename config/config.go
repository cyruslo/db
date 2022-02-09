package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	HttpClient "github.com/cyruslo/library/pkg/httpclient"
)

type config struct {
	HTTPUrl               			string  `json:"HTTPUrl"`
	HTTPSUrl               			string  `json:"HTTPSUrl"`
	ChatUrl               			string  `json:"ChatUrl"`
	SQLConnectionTimeout  			int     `json:"SQLConnectionTimeout"`
	LogFile              		 	string  `json:"LogFile"`
	SQLUsr                			string  `json:"SQLUsr"`
	SQLPass               			string  `json:"SQLPass"`
	SQLProcPass           			string  `json:"SQLProcPass"`
	SYSQLServerHost       			string  `json:"SYSQLServerHost"`
	CYSQLServerHost       			string  `json:"CYSQLServerHost"`
	PassportSQLServerHost 			string  `json:"PassportSQLServerHost"`
	SYGameCenterSQLServerHost 		string  `json:"SYGameCenterSQLServerHost"`
	GameCenterSQLServerHost 		string 	`json:"GameCenterSQLServerHost"`
	SYRecperGameSQLServerHost 		string  `json:"SYRecperGameSQLServerHost"`
	SYAgentSystemSQLServerHost 		string  `json:"SYAgentSystemSQLServerHost"`
	SignEnabled           			bool  	`json:"SignEnabled"`
	SignKey               			string 	`json:"SignKey"`
	EncryptEnabled					bool   	`json:"EncryptEnabled"`
	TLSEnabled               		bool    `json:"TLSEnabled"`
	TLSServerKey                   	string  `json:"TLSServerKey"`
	TLSServerCert                  	string  `json:"TLSServerCert"`
	TLSClientCert					string	`json:"TLSClientCert"`
}

var (
	configFilePath string
	cfgFilepath    string
	// Params qp配置
	Params *config = &config{}
)

func HttpUrl() string {
	http := Params.HTTPUrl
	// 检测切割Http字符串
	httpSlice := strings.Split(http, "//")
	if len(httpSlice) == 2 {
		httpSlice1 := strings.Split(httpSlice[1], "/")
		http = httpSlice1[0]
	}

	return http
}

func HttpsUrl() string {
	http := Params.HTTPSUrl
	// 检测切割Http字符串
	httpSlice := strings.Split(http, "//")
	if len(httpSlice) == 2 {
		httpSlice1 := strings.Split(httpSlice[1], "/")
		http = httpSlice1[0]
	}

	return http
}

// ReloadConfig 重新加载配置
func ReloadConfig(client bool) bool {
	if configFilePath == `` || cfgFilepath == `` {
		log.Printf("configFilePath == `` || cfgFilepath == ``")
		return false
	}

	return LoadConfig(cfgFilepath, configFilePath, client, false)
}

// LoadConfig 解析配置
func LoadConfig(configpath, filepath string, client, encode bool) bool {
	var params = &config{}

	f, err := os.Open(filepath)
	if err != nil {
		log.Println("failed to open config file:", filepath)
		return false
	}

	buf := make([]byte, 4096)
	n, err := f.Read(buf)
	if err != nil {
		log.Println("failed to read config file:", filepath)
		return false
	}

	json.Unmarshal(buf[:n], params)
	log.Println("-------------------Configure are:-------------------")
	log.Println(params)

	Params = params
	configFilePath = filepath
	cfgFilepath = configpath

	if encode == false {
		Params.TLSServerKey = fmt.Sprintf("%s%s", cfgFilepath, Params.TLSServerKey)
		Params.TLSServerCert = fmt.Sprintf("%s%s", cfgFilepath, Params.TLSServerCert)
		Params.TLSClientCert = fmt.Sprintf("%s%s", cfgFilepath, Params.TLSClientCert)
		log.Printf("LoadConfig cert=%s,key=%s.", Params.TLSClientCert, Params.TLSClientCert)
	}
	if client == true && Params.TLSClientCert != "" && Params.TLSEnabled == true {
		HttpClient.OnRead(Params.TLSClientCert)
	}
	return true
}

// DecodeConfig 解密
func DecodeConfig() {
	Params.SQLUsr, _ = Base64Decode(Params.SQLUsr)
	Params.SQLPass, _ = Base64Decode(Params.SQLPass)
	Params.SQLProcPass, _ = Base64Decode(Params.SQLProcPass)
}

// EncodeConfig 加密配置（只加密用户）
func EncodeConfig(filepath string, outputpath string) bool {
	if !LoadConfig(filepath, filepath, false, true) {
		return false
	}

	fout, err := os.Create(outputpath)
	defer fout.Close()
	if err != nil {
		fmt.Println(outputpath, err)
		return false
	}

	// 加密
	Params.SQLUsr = Base64Encode(Params.SQLUsr)
	Params.SQLPass = Base64Encode(Params.SQLPass)
	Params.SQLProcPass = Base64Encode(Params.SQLProcPass)

	b, err := json.Marshal(Params)
	if err != nil {
		log.Println("json marshal error:", err)
		return false
	}

	//log.Println("-------------------Encoded config are:-------------------")
	//log.Println(Params)

	fout.Write(b)
	return true
}

const (
	//BASE64字符表,不要有重复
	base64Table = "<>:;',./?~!@#$CDVWX%^&*ABYZabcghijklmnopqrstuvwxyz01EFGHIJKLMNOP"
)

var coder = base64.NewEncoding(base64Table)

/**
 * base64加密
 */
func Base64Encode(str string) string {
	src := []byte(str)
	return string(coder.EncodeToString(src))
}

/**
 * base64解密
 */
func Base64Decode(str string) (string, error) {
	src := []byte(str)
	by, err := coder.DecodeString(string(src))
	return string(by), err
}
