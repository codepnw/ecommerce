package database

import (
	"log"

	"github.com/codepnw/ecommerce/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func DbConnect(cfg config.IDbConfig) *sqlx.DB {
	db, err := sqlx.Connect("pgx", cfg.Url())
	if err != nil {
		log.Fatalf("connect to db failed: %v\n", err)
	}
	db.DB.SetMaxOpenConns(cfg.MaxConnections())
	return db
}
