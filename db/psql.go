package db

import (
	"database/sql"

	_ "github.com/lib/pq"

	ytrssilConfig "gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/config"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/lib/log"
)

type psqlDB struct {
	l  log.Logger
	db *sql.DB
}

func NewPSQLDB(log log.Logger, dbCfg ytrssilConfig.DB) (*psqlDB, error) {
	db, err := sql.Open("postgres", dbCfg.DBURI)
	if err != nil {
		return nil, err
	}

	return &psqlDB{
		l:  log,
		db: db,
	}, nil
}
