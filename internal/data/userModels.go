package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// START CRUD USERS
func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `
	select id, username, email, first_name, last_name, password, status, level, created_at, updated_at,
	case
		when (select count(id) from tokens t where username = users.username and t.expiry > NOW()) > 0
		then 1
		else 0
	end as hash_token
	from users order by last_name
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Status,
			&user.Level,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Token.ID,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (u *User) GetOne(username string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, username, email, first_name, last_name, password, status, level, created_at, updated_at from users where username = $1`

	var user User
	row := db.QueryRowContext(ctx, query, username)

	err := row.Scan(
		&user.ID,
		&user.UserName,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Status,
		&user.Level,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, username, email, first_name, last_name, password, status, level, created_at, updated_at from users where email = $1`

	var user User
	row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.UserName,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Status,
		&user.Level,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) GetByUsername(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, username, email, first_name, last_name, password, status, level, created_at, updated_at from users where username = $1`

	var user User
	row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.UserName,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Status,
		&user.Level,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) Insert(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	var newID int
	stmt := `
	insert into users(username, email, first_name, last_name, password, status, level, created_at, updated_at)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id
	`

	err = db.QueryRowContext(ctx, stmt,
		user.UserName,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPassword,
		user.Status,
		user.Level,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (u *User) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `update users set
		username = $1,
		email = $2,
		first_name = $3,
		last_name = $4,
		password = $5,
		status = $6,
		level = $7,
		created_at = $8,
		updated_at = $9,
	`
	_, err := db.ExecContext(ctx, stmt,
		u.UserName,
		u.Email,
		u.FirstName,
		u.LastName,
		u.Password,
		u.Level,
		u.CreatedAt,
		u.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil

}

func (u *User) DeleteByUsername(username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `delete from users where username = $1`

	_, err := db.ExecContext(ctx, stmt, username)
	if err != nil {
		return err
	}

	return nil
}

// END CRUD USERS

// START ABOUT PASSWORD
// Matching password
func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

// Reset password
func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil
	}

	stmt := `update users set password = $1 where username = $2`
	_, err = db.ExecContext(ctx, stmt, hashedPassword, u.ID)
	if err != nil {
		return nil
	}

	return nil
}

// END ABOUT PASSWORD

// START GET TOKEN
func (t *Token) GetByToken(plainText string) (*Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, username, email, token, token_hash, created_at, updated_at, expiry from tokens where token = $1`

	var token Token
	row := db.QueryRowContext(ctx, query, plainText)
	err := row.Scan(
		&token.ID,
		&token.UserName,
		&token.Email,
		&token.Token,
		&token.TokenHash,
		&token.CreatedAt,
		&token.UpdatedAt,
		&token.Expiry,
	)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (t *Token) GetUserForToken(token Token) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, username, email, first_name, last_name, password, status, level, created_at, updated_at from users where username = $1`

	var user User
	row := db.QueryRowContext(ctx, query, token.UserName)

	err := row.Scan(
		&user.ID,
		&user.UserName,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Status,
		&user.Level,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (t *Token) GenerateToken(username string, ttl time.Duration) (*Token, error) {
	token := &Token{
		UserName: username,
		Expiry:   time.Now().Add(ttl),
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Token = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Token))
	token.TokenHash = hash[:]

	return token, nil
}

// END GET TOKEN
