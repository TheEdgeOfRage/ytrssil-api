package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/alexedwards/argon2id"
	"github.com/lib/pq"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/models"
)

var authenticateUserQuery = `SELECT password FROM users WHERE username = $1`

func (d *postgresDB) AuthenticateUser(ctx context.Context, user models.User) (bool, error) {
	row := d.db.QueryRowContext(ctx, authenticateUserQuery, user.Username)
	var hashedPassword string
	err := row.Scan(&hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		d.l.Log("level", "ERROR", "function", "db.AuthenticateUser", "error", err)
		return false, err
	}

	match, err := argon2id.ComparePasswordAndHash(user.Password, hashedPassword)
	if err != nil {
		d.l.Log("level", "ERROR", "function", "db.AuthenticateUser", "error", err)
		return false, err
	}

	return match, nil
}

var createUserQuery = `INSERT INTO users (username, password) VALUES ($1, $2)`

func (d *postgresDB) CreateUser(ctx context.Context, user models.User) error {
	_, err := d.db.ExecContext(ctx, createUserQuery, user.Username, user.Password)
	if err != nil {
		if pgerr, ok := err.(*pq.Error); ok {
			if pgerr.Code == "23505" {
				return ErrUserExists
			}
		}
		d.l.Log("level", "ERROR", "function", "db.CreateUser", "error", err)
		return err
	}

	return nil
}

var deleteUserQuery = `DELETE FROM users WHERE username = $1`

func (d *postgresDB) DeleteUser(ctx context.Context, username string) error {
	_, err := d.db.ExecContext(ctx, deleteUserQuery, username)
	if err != nil {
		d.l.Log("level", "ERROR", "function", "db.DeleteUser", "error", err)
		return err
	}

	return nil
}
