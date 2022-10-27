package db

import (
	"context"
	"fmt"

	"github.com/alexedwards/argon2id"
	"github.com/georgysavva/scany/v2/sqlscan"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/models"
)

var authenticateUserQuery = `SELECT username, password FROM users WHERE username = $1`

func (d *psqlDB) AuthenticateUser(ctx context.Context, username string, password string) (bool, error) {
	var user []*models.User
	err := sqlscan.Select(ctx, d.db, &user, authenticateUserQuery, username)
	if err != nil {
		d.l.Log("level", "ERROR", "function", "db.AuthenticateUser", "msg", "failed to query user", "error", err)
		return false, err
	}

	match, err := argon2id.ComparePasswordAndHash(password, user[0].Password)
	if err != nil {
		d.l.Log("level", "ERROR", "function", "db.AuthenticateUser", "msg", "failed to check hashed passsword", "error", err)
		return false, err
	}

	return match, nil
}

var createUserQuery = `INSERT INTO users (username, password) VALUES ($1, $2)`

func (d *psqlDB) CreateUser(ctx context.Context, user models.User) error {
	res, err := d.db.Exec(createUserQuery, user.Username, user.Password)
	if err != nil {
		d.l.Log("level", "ERROR", "function", "db.CreateUser", "msg", "failed to create user", "error", err)
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		d.l.Log("level", "ERROR", "function", "db.CreateUser", "msg", "failed to get affected row count", "error", err)
		return err
	}

	if affected != 1 {
		d.l.Log("level", "ERROR", "function", "db.CreateUser", "msg", "failed to get affected row count", "error", err)
		return fmt.Errorf("expected to insert one row, but %d were inserted", affected)
	}

	return nil
}
