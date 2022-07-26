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

type conf struct {
	ListenAddr string `json:"listen_addr"`

	Salt        string      `json:"salt"`
	JWTToken    string      `json:"jwt_token"`
	PgSQLConfig PgSQLConfig `json:"pgsql_config"`
}

var CONF *conf

func init() {
	initConf()
}

func initConf() {
	var cf conf

	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		marshal, _ := json.MarshalIndent(cf, " ", "  ")
		if err2 := ioutil.WriteFile("config.json", marshal, 00666); err2 != nil {
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
