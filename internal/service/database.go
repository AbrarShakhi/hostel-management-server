package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Database struct {
	psql *sql.DB
}

var (
	database   = os.Getenv("POSTGRES_DB_DATABASE")
	password   = os.Getenv("POSTGRES_DB_PASSWORD")
	username   = os.Getenv("POSTGRES_DB_USERNAME")
	port       = os.Getenv("POSTGRES_DB_PORT")
	host       = os.Getenv("POSTGRES_DB_HOST")
	schema     = os.Getenv("POSTGRES_DB_SCHEMA")
	dbInstance *Database
)

func DbInstance() *Database {
	once.Do(func() {
		psql, err := sql.Open("pgx",
			fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
				username, password, host, port, database, schema))
		if err != nil {
			log.Fatal(err)
		}
		dbInstance = &Database{
			psql: psql,
		}
	})
	return dbInstance
}

func (s *Database) Query(query string, args ...any) (*sql.Rows, error) {
	return s.psql.Query(query, args...)
}

func (s *Database) QueryRow(query string, args ...any) *sql.Row {
	return s.psql.QueryRow(query, args...)
}

func (s *Database) Exec(query string, args ...any) (sql.Result, error) {
	return s.psql.Exec(query, args...)
}

func (s *Database) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	err := s.psql.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"

	dbStats := s.psql.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	if dbStats.OpenConnections > 40 {
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

func (s *Database) Close() error {
	log.Printf("Disconnected from database: %s", database)
	return s.psql.Close()
}
