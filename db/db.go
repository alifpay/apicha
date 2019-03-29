package db

import (
	"log"
	"runtime"

	"github.com/jackc/pgx"
)

var pgPool *pgx.ConnPool

// afterConnect creates the prepared statements that this application uses
func afterConnect(conn *pgx.Conn) (err error) {

	_, err = conn.Prepare("user", `SELECT name, password, token FROM users WHERE id = $1 AND status = 'active' AND ip ? $2`)
	if err != nil {
		return
	}
	_, err = conn.Prepare("updateUser", "UPDATE users SET status = $1, name = $2 WHERE id = $3")
	if err != nil {
		return
	}

	_, err = conn.Prepare("insertUser", `INSERT INTO users(name, password, token, ip, status) VALUES($1, $2, $3, $4, $5) RETURNING id`)
	return

}

// Connect connect to postgres and adds new connection pool
func Connect(dbhost string, dbuser string, dbpass string, dbname string) {
	var err error
	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     dbhost,
			User:     dbuser,
			Password: dbpass,
			Database: dbname,
		},
		MaxConnections: runtime.NumCPU() * 2,
		AfterConnect:   afterConnect,
	}
	pgPool, err = pgx.NewConnPool(connPoolConfig)
	if err != nil {
		log.Fatalln("Databese connection", err)
	}
}

// Close db connection
func Close() {
	pgPool.Close()
}
