package repo

import (
	"database/sql"

	"steadylearner.com/sqlite/models"
)

// UsersRepository thing
type UsersRepository struct {
	db *sql.DB
}

// NewUsersRepository returns a new user repository with db.
func NewUsersRepository(db *sql.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

// List list users
func (r *UsersRepository) List() ([]*models.User, error) {

	rows, err := r.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	users := []*models.User{}
	for rows.Next() {
		u := models.User{}
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

// Get fetches a user by it's ID
func (r *UsersRepository) Get(id int64) (*models.User, error) {
	// No need tx in this operation

	ret := models.User{}
	res := r.db.QueryRow("SELECT * FROM users WHERE id=?", id)
	if err := res.Scan(&ret.ID, &ret.Name); err != nil {
		return nil, err
	}

	return &ret, nil
}

// Create creates a user by name.
func (r *UsersRepository) Create(user *models.User) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	res, err := tx.Exec(`INSERT INTO users (name) values(?)`, user.Name)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id

	return tx.Commit()
}

// Update updates the user
func (r *UsersRepository) Update(user *models.User) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec("UPDATE users SET name=? where id=?", user.Name, user.ID); err != nil {
		return err
	}

	return tx.Commit()
}

// Delete deletes an user
func (r *UsersRepository) Delete(id int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec("DELETE FROM users WHERE id=?", id); err != nil {
		return err
	}
	return tx.Commit()
}
