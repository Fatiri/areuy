package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// GormDatabaseConfig config
type GormDatabaseConfig struct {
	Host         string
	Port         string
	UserName     string
	Password     string
	DatabaseName string
	DatabaseType string
}

type gormInstance struct {
	db *gorm.DB
}

// GormDatabase abstraction
type GormDatabase interface {
	DB() *gorm.DB
	Close()
}

// InitGorm ...
func (dbc *GormDatabaseConfig) InitGorm() *gorm.DB {

	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Silent, // Log level
			Colorful:      true,          // Disable color
		},
	)

	gormConfig := &gorm.Config{
		// enhance performance config
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 dbLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	dnsMaster := ""
	databaseType := strings.ToLower(dbc.DatabaseType)
	if databaseType == "mysql" {
		// username, password, host, port, database
		dnsMaster = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
			dbc.UserName, dbc.Password, dbc.Host,
			dbc.Port, dbc.DatabaseName)
		dnsMaster += `&charset=utf8mb4&parseTime=True&loc=Local`
		sqlMaster, errMaster := sql.Open("mysql", dnsMaster)
		if errMaster != nil {
			log.Panic(errMaster)
		}
		// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
		sqlMaster.SetMaxIdleConns(10)
		// SetMaxOpenConns sets the maximum number of open connections to the database.
		sqlMaster.SetMaxOpenConns(100)
		// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
		sqlMaster.SetConnMaxLifetime(time.Hour)
		dbMaster, errMaster := gorm.Open(mysql.New(mysql.Config{
			Conn: sqlMaster,
		}), gormConfig)
		if errMaster != nil {
			log.Panic(errMaster)
		}

		return dbMaster
	} else if databaseType == "postgres" {
		dnsMaster = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbc.Host, dbc.Port, dbc.UserName, dbc.Password, dbc.DatabaseName)
		sqlMaster, errMaster := sql.Open("postgres", dnsMaster)
		if errMaster != nil {
			log.Panic(errMaster)
		}
		// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
		sqlMaster.SetMaxIdleConns(10)
		// SetMaxOpenConns sets the maximum number of open connections to the database.
		sqlMaster.SetMaxOpenConns(100)
		// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
		sqlMaster.SetConnMaxLifetime(time.Hour)
		dbMaster, errMaster := gorm.Open(postgres.New(postgres.Config{
			Conn: sqlMaster,
		}), gormConfig)
		if errMaster != nil {
			log.Panic(errMaster)
		}
		return dbMaster
	} else {
		log.Panic(fmt.Errorf(fmt.Sprintf("database : %s not support", databaseType)))
	}

	return nil
}
