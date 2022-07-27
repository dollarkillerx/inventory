package storage

import (
	"github.com/dollarkillerx/inventory/internal/conf"
	"github.com/dollarkillerx/inventory/internal/storage/simple"
	"log"
	"testing"
)

func TestStorage(t *testing.T) {
	s, err := simple.NewSimple(&conf.CONF.PgSQLConfig)
	if err != nil {
		log.Fatalln(err)
	}

	good, err := s.Good("6921168509256", "10086")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(good)
}
