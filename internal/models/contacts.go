package models

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"unicode"
)

var rxEmail = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type Contact struct {
	ID      int
	Name    string
	Phone   string
	Email   string
	Errors  map[string]string
	Message string
}

func (c *Contact) Valid() bool {
	c.Errors = make(map[string]string)
	c.Email = strings.TrimSpace(c.Email)
	match := rxEmail.Match([]byte(c.Email))
	if !match {
		c.Errors["Email"] = "Please enter a valid email address"
	}

	if strings.TrimSpace(c.Name) == "" {
		c.Errors["Name"] = "Please enter a name"
	}
	c.Phone = nums(c.Phone)
	if len(c.Phone) < 10 || len(strings.TrimSpace(c.Phone)) > 11 {
		c.Errors["Phone"] = "Please enter a valid phone number"
	}

	return len(c.Errors) == 0
}

func nums(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return r
		}
		return -1
	}, s)
}

type PostgresContactRepo struct {
	DB *sql.DB
}

func NewPostgresContactRepo(db *sql.DB) *PostgresContactRepo {
	return &PostgresContactRepo{
		DB: db,
	}
}

// Create is a method on a *UserService that takes an email, and password string and returns a *User, and an error. The function changes the email to lowercase, hashes the password, creates the user, then queries the DB to store the email and passwordhash
func (us *PostgresContactRepo) Read(name string) (*[]Contact, error) {
	rows, err := us.DB.Query(`
	SELECT *
	FROM contacts
	WHERE name % $1
	ORDER BY name ASC;`, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []Contact{}
	for rows.Next() {
		user := Contact{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone)
		if err != nil {
			return nil, sql.ErrNoRows
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, sql.ErrNoRows
	}
	return &users, nil
}

func (us *PostgresContactRepo) ReadAll() (*[]Contact, error) {
	rows, err := us.DB.Query(`
	SELECT *
	FROM contacts
	ORDER BY name ASC;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []Contact{}
	for rows.Next() {
		user := Contact{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return &users, nil
}

func (us *PostgresContactRepo) GetContact(id int) (*Contact, error) {
	contact := Contact{}
	row := us.DB.QueryRow(`
		SELECT id, name, email, phone
		FROM contacts
		WHERE id = $1;`, id)
	err := row.Scan(&contact.ID, &contact.Name, &contact.Email, &contact.Phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoResults
		}
		return nil, err
	}
	return &contact, nil
}

func (us *PostgresContactRepo) New(contact *Contact) error {
	_, err := us.DB.Exec(`
		INSERT INTO contacts (name, email, phone)
		VALUES ($1, $2, $3);`, contact.Name, contact.Email, contact.Phone)
	if err != nil {
		// In the future, I might want to see if this is a uniquie email violation.
		// and inform the user that another contact has that email. For now, Im not going
		// to do that becuase most contact apps allow duplicates, but I have my DB setup
		// so that the email must be unique.
		return err
	}
	return nil
}

func (us *PostgresContactRepo) Update(contact *Contact) error {
	_, err := us.DB.Exec(`
		UPDATE contacts
		SET name = $2, phone = $3, email = $4
		WHERE id = $1;`, contact.ID, contact.Name, contact.Phone, contact.Email)
	if err != nil {
		return err
	}
	return nil
}

func (us *PostgresContactRepo) Delete(contact *Contact) error {
	_, err := us.DB.Exec(`
		DELETE FROM contacts
		WHERE id = $1;`, contact.ID)
	if err != nil {
		return err
	}
	return nil
}
