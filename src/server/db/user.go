package db

import (
	"fmt"
	"time"

	"github.com/georgysavva/scany/pgxscan"
)

type User struct {
	ID        string    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	SpotifyId *string   `db:"spotify_id"`
	Name      string    `db:"name"`
	Token     *string   `db:"token"`
}

func (db *DB) PutUser(user *User) (userId string, err error) {
	query := `
		INSERT INTO users (spotify_id, name, token)
		VALUES ($1, $2, $3)
		ON CONFLICT (spotify_id) DO UPDATE SET token = $3
		RETURNING id;
	`

	err = db.pool.QueryRow(*db.ctx, query, user.SpotifyId, user.Name, user.Token).Scan(&userId)
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("unable to insert user")
	}

	return userId, nil
}

func (db *DB) CreateUser(user *User) (userId string, err error) {
	query := `
		INSERT INTO users (spotify_id, name, token)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	err = db.pool.QueryRow(*db.ctx, query, user.SpotifyId, user.Name, user.Token).Scan(&userId)
	if err != nil {
		return "", fmt.Errorf("unable to insert user")
	}

	return userId, nil
}

func (db *DB) GetUserById(userId string) (user User, err error) {
	fmt.Println("Getting user by id", userId)
	query := `
		SELECT * FROM users WHERE id = $1;
	`

	err = pgxscan.Get(*db.ctx, db.pool, &user, query, userId)
	if err != nil {
		fmt.Println(err)
		return user, fmt.Errorf("failed to get user")
	}

	fmt.Println(user)

	return user, nil
}
