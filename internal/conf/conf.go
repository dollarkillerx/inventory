package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type PgSQLConfig struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"db_name"`
	Port     int    `json:"port"`
}

type OSSConf struct {
	Bucket    string `json:"bucket"`
	Region    string `json:"region"`
	Url       string `json:"url"`
	SecretID  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
}

type conf struct {
	ListenAddr string `json:"listen_addr"`

	Salt        string      `json:"salt"`
	JWTToken    string      `json:"jwt_token"`
	PgSQLConfig PgSQLConfig `json:"pgsql_config"`
	OSSConf     OSSConf     `json:"oss_conf"`
}

var CONF *conf

func init() {
	initConf()
}

func initConf() {
	var cf conf

	file, err := ioutil.ReadFile("configs/config.json")
	if err != nil {
		marshal, _ := json.MarshalIndent(cf, " ", "  ")
		if err2 := ioutil.WriteFile("configs/config.json", marshal, 00666); err2 != nil {
			log.Fatalln(err2)
		}
		log.Fatalln(err)
	}

	err = json.Unmarshal(file, &cf)
	if err != nil {
		log.Fatalln(err)
	}

	CONF = &cf
}
