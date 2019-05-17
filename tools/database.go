package tools

import (
	"database/sql"
	"net/url"
	"os"
	"strings"

	// _ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	// _ "github.com/jinzhu/gorm/dialects/sqlite"
)

// DB ...
type DB struct {
	db       *sql.DB
	Username string
	Password string
	Hostname string
	Database string
}

// NewDB ...
func NewDB(db *sql.DB, username, password, hostname, database string) *DB {
	v := DB{db, username, password, hostname, database}
	return &v
}

// NewDBWithEnv ...
func NewDBWithEnv() *DB {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		panic("Missing DATABASE_URL environment variable")
	}
	return NewDBWithString(dbURL)
}

// NewDBWithString ...
func NewDBWithString(urlString string) *DB {
	u, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}

	urlString = getConnectionString(*u)

	db, err := sql.Open(u.Scheme, urlString)
	if err != nil {
		panic(err)
	}
	// db.LogMode(true)
	user := u.User
	username := user.Username()
	password, _ := user.Password()
	return NewDB(db, username, password, u.Hostname(), getDatabaseName(u))
}

func getConnectionString(u url.URL) string {
	if u.Scheme != "mysql" {
		panic("Only MySQL is currently supported")
	}
	// if u.Scheme == "postgres" {
	// 	password, _ := u.User.Password()
	// 	host := strings.Split(u.Host, ":")[0]
	// 	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, u.Port(), u.User.Username(), password, strings.TrimPrefix(u.Path, "/"))
	// }
	if u.Scheme != "sqlite3" {
		u.Host = "tcp(" + u.Host + ")"
	}
	// for db clear-tables
	u.Query().Set("multipleStatements", "true")
	return strings.Replace(u.String(), u.Scheme+"://", "", 1)
}
func getDatabaseName(u *url.URL) string {
	return strings.TrimPrefix(u.Path, "/")
}

// Client ...
func (db *DB) Client() *sql.DB {
	return db.db
}

// Close ...
func (db *DB) Close() error {
	return db.db.Close()
}
