package db

import (
	"fmt"
	"strings"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pedrocelso/go-rest-service/lib/config"
)

// ConnectToDB generates an connection string and connect ot MySQL
func ConnectToDB(conn config.Mysql) *sqlx.DB {
	options := []string{
		"parseTime=true",
		"interpolateParams=true",
		"multiStatements=true",
	}
	user := conn.User
	pass := conn.Password
	host := conn.Host
	databaseName := conn.Database
	port := conn.Port
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", user, pass, host, port, databaseName, strings.Join(options, "&"))
	return sqlx.MustConnect("mysql", connectionString)
}
