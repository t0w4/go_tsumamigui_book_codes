package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func Init() (*sql.DB, error) {
	connectionString := getConnectionString()
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getParamString(param string, defaultValue string) string {
	env := os.Getenv(param)
	if env != "" {
		return env
	}
	return defaultValue
}

func getConnectionString() string {
	host := getParamString("MYSQL_DB_HOST", "localhost")
	port := getParamString("MYSQL_PORT", "3306")
	user := getParamString("MYSQL_USER", "root")
	pass := getParamString("MYSQL_PASSWORD", "")
	dbName := getParamString("MYSQL_DB", "test_database")
	protocol := getParamString("MYSQL_PROTOCOL", "tcp")
	dbArgs := getParamString("MYSQL_DBARGS", " ")

	if strings.Trim(dbArgs, " ") != "" {
		dbArgs = "?" + dbArgs
	} else {
		dbArgs = ""
	}
	return fmt.Sprintf("%s:%s@%s([%s]:%s)/%s%s",
		user, pass, protocol, host, port, dbName, dbArgs)
}
