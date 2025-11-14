package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText *string
	hash      []byte
}

func (p *password) Set(plainPasswordText string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPasswordText), 12)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	p.plainText = &plainPasswordText
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, fmt.Errorf("failed to compare passwords: %w", err)
		}
	}
	return true, nil
}

type User struct {
	ID           int       `json:"id"`
	Username        string    `json:"username"`
	PasswordHash password  `json:"-"`
	Email    string    `json:"email"`
	BIO     string    `json:"bio"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

type UserStore interface {
	CreateUser(user *User) error
	GetUserByID(id int64) (*User, error)
	GetUserByUsername(username string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id int64) error
}

func (pg *PostgresUserStore) CreateUser(user *User) error {
	query :=
		`INSERT INTO users (username, email, password_hash, bio, created_at, updated_at)
	VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id;
	`
	err := pg.db.QueryRow(query, user.Username, user.Email, user.PasswordHash.hash, user.BIO, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (pg *PostgresUserStore) GetUserByID(id int64) (*User, error) {
	user := &User{PasswordHash: password{}}
	query := `
	SELECT id, username, email, password_hash, bio, created_at, updated_at from users where id = $1;
	`
	err := pg.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.BIO, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (pg *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	user := &User{PasswordHash: password{}}
	query := `
	SELECT id, username, email, password_hash, bio, created_at, updated_at from users where username = $1;
	`
	err := pg.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.BIO, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (pg *PostgresUserStore) UpdateUser(user *User) error {
	query := `
	UPDATE users SET username = $1, email = $2, password_hash = $3, bio = $4, updated_at = NOW()
	WHERE id = $5;
	`
	result, err := pg.db.Exec(query, user.Username, user.Email, user.PasswordHash, user.BIO, user.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", user.ID)
	}
	return nil
}

func (pg *PostgresUserStore) DeleteUser(id int64) error {
	query := `
	DELETE FROM users WHERE id = $1;
	`
	result, err := pg.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("article with ID %d not found", id)
	}
	return nil
}
