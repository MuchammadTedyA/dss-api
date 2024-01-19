package data

import (
	"context"
)

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
