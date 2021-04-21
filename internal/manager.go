package internal

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	Host     string
	Username string
	Pass     string
	Port     string
	DBName   string
}

// SQLManager - manage connect to db
type SQLManager struct {
	conn *sql.DB
}

// InitManager - init connect to db
func InitManager() *SQLManager {
	var m = &SQLManager{}

	m.open(&config{
		Host:     os.Getenv("MYSQL_HOST"),
		Username: os.Getenv("MYSQL_USER"),
		Pass:     os.Getenv("MYSQL_PASSWORD"),
		Port:     "3306",
		DBName:   os.Getenv("MYSQL_DATABASE"),
	})

	return m
}

// Close - close connect to db
func (m *SQLManager) Close() {
	_ = m.conn.Close()
}

func (m *SQLManager) GetConnection() *sql.DB {
	return m.conn
}

func (m *SQLManager) open(c *config) {
	var conn *sql.DB
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?collation=utf8_unicode_ci", c.Username, c.Pass, c.Host, c.Port, c.DBName)
	if conn, err = sql.Open("mysql", dsn); err != nil {
		log.Printf("on insert: on open connection to db: %s \n", err.Error())

		os.Exit(1)
	}

	m.conn = conn
}
