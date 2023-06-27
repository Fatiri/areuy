package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	dbr "github.com/abiewardani/dbr/v2"
	"github.com/abiewardani/dbr/v2/dialect"
	_ "go.elastic.co/apm/module/apmsql/mysql"
)

// DbrDatabaseConfig config
type DbrDatabaseConfig struct {
	Host         string
	Port         string
	UserName     string
	Password     string
	DatabaseName string
	DatabaseType string
}

type dbrInstance struct {
	dbr *dbr.Session
}

func (c *dbrInstance) DB() *dbr.Session {
	return c.dbr
}

// DbrDatabase abstraction
type DbrDatabase interface {
	DB() *dbr.Session
}

// InitDbr ...
func (dbc *DbrDatabaseConfig) InitDbr() DbrDatabase {
	inst := new(dbrInstance)

	dnsMaster := ""
	databaseType := strings.ToLower(dbc.DatabaseType)
	if databaseType == "mysql" {
		// username, password, host, port, database
		dnsMaster = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
			dbc.UserName, dbc.Password, dbc.Host,
			dbc.Port, dbc.DatabaseName)
		dnsMaster += `&loc=Asia%2FJakarta&charset=utf8`
	} else if databaseType == "postgress" {
		dnsMaster = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbc.Host, dbc.Port, dbc.UserName, dbc.Password, dbc.DatabaseName,
		)
	} else {
		log.Panic(errors.New(fmt.Sprintf("database : %s not support", databaseType)))
	}

	sqlMaster, errMaster := sql.Open(databaseType, dnsMaster)
	connMaster := &dbr.Connection{
		DB:            sqlMaster, // <- underlying database/sql.DB is instrumented
		EventReceiver: &dbr.NullEventReceiver{},
		Dialect:       dialect.MySQL,
	}

	if errMaster != nil {
		log.Panic(errMaster)
	}

	if errPing := connMaster.Ping(); errPing != nil {
		log.Panic(errPing)
	}

	connMaster.SetMaxOpenConns(100)
	connMaster.SetMaxIdleConns(10)
	inst.dbr = connMaster.NewSession(nil)

	return inst
}
