package postgres

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type PostgresSettings struct {
	Host     string `xml:"psqlHost"`
	Port     string `xml:"psqlPort"`
	User     string `xml:"psqlUser"`
	Password string `xml:"psqlPassword"`
	DbName   string `xml:"psqlDbName"`
}

func (settings *PostgresSettings) Valid() error {
	if settings.Host == "" {
		settings.Host = "postgres"
	}
	if settings.Port == "" {
		settings.Port = "5432"
	}
	if settings.User == "" {
		settings.User = "postgres"
	}
	if settings.Password == "" {
		settings.Password = "postgres"
	}
	if settings.DbName == "" {
		settings.DbName = "hezzl"
	}
	return nil
}

type PostgresDb struct {
	Db *sql.DB

	PingLoopDelay int
}

func InitPostgresDb(s *PostgresSettings, pingLoopDelay int) (db *PostgresDb, err error) {
	if err = s.Valid(); err != nil {
		return nil, fmt.Errorf("InitPostgresDb: %s", err)
	}

	db = &PostgresDb{}
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", s.Host, s.Port, s.User, s.Password, s.DbName)
	if sqlDb, err := sql.Open("postgres", psqlconn); err == nil {
		db.Db = sqlDb
	} else {
		return nil, fmt.Errorf("InitPostgresDb: %s", err)
	}
	if err := db.Db.Ping(); err != nil {
		return nil, fmt.Errorf("InitPostgresDb: %s", err)
	}

	if pingLoopDelay > 0 {
		db.PingLoopDelay = pingLoopDelay
		go db.PingLoop()
	}
	return db, nil
}

func (db *PostgresDb) Close() error {
	db.PingLoopDelay = 0
	if err := db.Db.Close(); err != nil {
		return fmt.Errorf("PostgresClient.Close: %s", err)
	}
	return nil
}

func (db *PostgresDb) PingLoop() {
	for db.PingLoopDelay > 0 {
		db.Db.Ping()
		time.Sleep(time.Duration(db.PingLoopDelay) * time.Second)
	}
}
