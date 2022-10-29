package db

import (
	"database/sql"

	_ "github.com/lib/pq"

	ytrssilConfig "gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/config"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/lib/log"
)

type postgresDB struct {
	l  log.Logger
	db *sql.DB
}

func NewPostgresDB(log log.Logger, dbCfg ytrssilConfig.DB) (*postgresDB, error) {
	db, err := sql.Open("postgres", dbCfg.DBURI)
	if err != nil {
		return nil, err
	}

	return &postgresDB{
		l:  log,
		db: db,
	}, nil
}
