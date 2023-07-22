package postgres

import (
	"fmt"
	

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/yervsil/onlineShop/internal/config"
)


func NewPostgresDB(c *config.Config) (*sqlx.DB, error){
	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", c.Postgres.Host,
	c.Postgres.Port, c.Postgres.Dbname, c.Postgres.Username, c.Postgres.Password, c.Postgres.Sslmode)

	db, err := sqlx.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}