package storage

import (
	"github.com/dollarkillerx/inventory/internal/conf"
	"github.com/dollarkillerx/inventory/internal/storage/simple"
	"github.com/dollarkillerx/inventory/internal/utils"
	"log"
	"testing"
)

func TestStorage(t *testing.T) {
	s, err := simple.NewSimple(&conf.CONF.PgSQLConfig)
	if err != nil {
		log.Fatalln(err)
	}

	good, err := s.Good("8801116016730", "10086")
	if err != nil {
		log.Fatalln(err)
	}
	utils.PrintObj(good)
}
