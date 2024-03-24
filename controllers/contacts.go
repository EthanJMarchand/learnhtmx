package controllers

import (
	"fmt"
	"net/http"

	"github.com/ethanjmarchand/learnhtmx/errors"
	"github.com/ethanjmarchand/learnhtmx/models"
)

type Contacts struct {
	Templates struct {
		Contacts Template
	}
	UserService *models.UserService
}

// this handler func simply redirects the user to the contacts page.
func (c Contacts) Redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/contacts", http.StatusFound)
}

// Show() looks at the perameter query "q" and tries to show you all the contacts
// that match to that. If q == "", it sends you all the contacts.
func (c Contacts) Show(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Query string
		Users []models.User
	}
	data.Query = r.FormValue("q")
	if data.Query == "" {
		// fetch and render all of the contacts
		users, err := c.UserService.ReadAll()
		if err != nil {
			if errors.Is(err, models.ErrNoResults) {
				err = errors.Public(err, "No results")
				fmt.Println(err)
				c.Templates.Contacts.Execute(w, r, data, err)
				return
			}
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		data.Users = *users
		c.Templates.Contacts.Execute(w, r, data)
		return
	}
	// try to find a first name, or last name that matches anything like the q
	// and present them to the screen.
	users, err := c.UserService.Read(data.Query)
	if err != nil {
		if errors.Is(err, models.ErrNoResults) {
			err = errors.Public(err, "No results")
			fmt.Println(err)
			c.Templates.Contacts.Execute(w, r, data, err)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	data.Users = *users
	c.Templates.Contacts.Execute(w, r, data)
}
