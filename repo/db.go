package repo

import (
	"context"
	"fmt"

	"main/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Dbstruct struct {
	*pgxpool.Pool
}

func NewDb(ctx context.Context, cfg config.Config) (*Dbstruct, error) {

	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable pool_max_conns=20",
		cfg.DBUsername,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBdatabase,
	)

	dbConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.New(ctx, dbConfig.ConnString())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return &Dbstruct{
		db,
	}, nil
}