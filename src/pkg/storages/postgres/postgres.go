package postgres

import (
	"ecomdream/src/pkg/configs"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"sync"
)

var database *sqlx.DB
var once sync.Once
var databaseUrl string

func Connection() *sqlx.DB {
	if database == nil {
		once.Do(Init)
	}

	return database
}

func Init() {
	var err error

	databaseUrl, _ = configs.ConnectionURLBuilder("postgres")

	database, err = sqlx.Connect("pgx", databaseUrl)
	if err != nil {
		logrus.Error(err)
		return
	}

	if err = database.Ping(); err != nil {
		logrus.Error(err)
		return
	}

	database.SetMaxOpenConns(32)
	database.SetMaxIdleConns(32)

	logrus.Printf("Connected to %s", databaseUrl)
}

func init() {
	Init()
}
