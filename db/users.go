package db

import "database/sql"

type User struct {
	ID        int    `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}

func (d *Database) QueryUserByID(tx *sql.Tx, userID int) (*User, error) {
	row := tx.QueryRow(`
		SELECT *
		FROM users
		WHERE id = $1
	`, userID)

	var user User
	if err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}

	return &user, nil
}
