package models

import (
	"database/sql"
	"errors"
)

type User struct {
	ID    int
	Name  string
	Phone string
	Email string
}

type UserService struct {
	DB *sql.DB
}

// Create is a method on a *UserService that takes an email, and password string and returns a *User, and an error. The function changes the email to lowercase, hashes the password, creates the user, then queries the DB to store the email and passwordhash
func (us *UserService) Read(name string) (*[]User, error) {
	rows, err := us.DB.Query(`
	SELECT *
	FROM contacts
	WHERE name % $1;`, name)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, ErrNoResults
		}
		return nil, err
	}
	defer rows.Close()
	users := []User{}
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return &users, nil
}

func (us *UserService) ReadAll() (*[]User, error) {
	rows, err := us.DB.Query(`
	SELECT *
	FROM contacts;`)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, ErrNoResults
		}
		return nil, err
	}
	defer rows.Close()
	users := []User{}
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return &users, nil
}
