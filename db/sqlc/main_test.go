package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/bytepharoh/simplebank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config , err :=util.LoadConfig("../..")
	if err != nil {
		log.Fatal("can not load config")
	}
	testDB , err = sql.Open(config.DBDriver , config.DBSource)
	if err != nil {
		log.Fatal("Can not establish a connection:", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}