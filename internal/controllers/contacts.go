package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ethanjmarchand/learnhtmx/internal/errors"
	"github.com/ethanjmarchand/learnhtmx/internal/models"
)

type Contacts struct {
	Templates struct {
		Contacts Template
		New      Template
		View     Template
		Edit     Template
	}
	ContactService *models.ContactService
}

// Redirect handler func simply redirects the user to the contacts page.
func (c Contacts) Redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/contacts", http.StatusFound)
}

// Show looks at the parameter query "q" and tries to show you all the contacts
// that match to that. If q == "", it sends you all the contacts.
func (c Contacts) Show(w http.ResponseWriter, r *http.Request) {
	fb, err := GetFlash(w, r, "success")
	if err != nil {
		fmt.Printf("Show() err: %q", err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	var data struct {
		Query   string
		Users   []models.Contact
		Message string
	}
	data.Message = string(fb)
	data.Query = r.FormValue("q")
	if data.Query == "" {
		// fetch and render all the contacts
		users, err := c.ContactService.ReadAll()
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = errors.Public(err, "No results")
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
	users, err := c.ContactService.Read(data.Query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.Public(err, "No results")
			c.Templates.Contacts.Execute(w, r, data, err)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	data.Users = *users
	c.Templates.Contacts.Execute(w, r, data)
}

// New is the handler function that the /contacts/new GET request
func (c Contacts) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Name    string
		Phone   string
		Email   string
		Errors  map[string]string
		Message string
	}
	data.Name = r.FormValue("name")
	data.Phone = r.FormValue("phone")
	data.Email = r.FormValue("email")
	c.Templates.New.Execute(w, r, data)
}

// ProcessNew is the handler function that handles the POST request.
func (c Contacts) ProcessNew(w http.ResponseWriter, r *http.Request) {
	contact := models.Contact{}
	contact.Name = r.FormValue("name")
	contact.Phone = r.FormValue("phone")
	contact.Email = r.FormValue("email")
	if !contact.Valid() {
		c.Templates.New.Execute(w, r, contact)
		return
	}
	// store in the DB.
	err := c.ContactService.New(&contact)
	if err != nil {
		fmt.Printf("ProcessNew() error: %s", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
	// redirect to contacts page.
	SetFlash(w, "success", "Contact was saved successfully")
	http.Redirect(w, r, "/contacts", http.StatusFound)
}

// View is the handler function that handles the GET request for a specific contact
func (c Contacts) View(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("id")
	if sid == "" {
		http.Error(w, "Must provide a valid ID", http.StatusInternalServerError)
	}
	id, err := strconv.Atoi(sid)
	if err != nil {
		fmt.Println("strconv error")
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
	// lookup the contact in the DB
	contact, err := c.ContactService.GetContact(id)
	if err != nil {
		fmt.Printf("View(): GetContact(): err = %s", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	c.Templates.View.Execute(w, r, contact)
}

// TODO: Test the edit function handler by mocking the user function handler.

// Edit is the handler function that handles a GET request to the /contacts/{id}/edit route.
func (c Contacts) Edit(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("id")
	if sid == "" {
		http.Error(w, "Must provide a valid ID", http.StatusInternalServerError)
	}
	id, err := strconv.Atoi(sid)
	if err != nil {
		fmt.Println("strconv error")
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
	contact := models.Contact{
		ID: id,
	}
	// lookup the contact in the DB
	// Use a mock user service. ************
	con, err := c.ContactService.GetContact(contact.ID)
	if err != nil {
		fmt.Printf("Edit(): GetContact(): err = %s", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	c.Templates.Edit.Execute(w, r, con)
}

// ProcessEdit is the handler function that handles the POST request when someone tries to save a contact.
func (c Contacts) ProcessEdit(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(urlParts[len(urlParts)-2])
	if err != nil {
		fmt.Printf("Edit(): err = %s", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	contact := models.Contact{}
	contact.ID = id
	contact.Name = r.FormValue("name")
	contact.Phone = r.FormValue("phone")
	contact.Email = r.FormValue("email")
	if !contact.Valid() {
		c.Templates.Edit.Execute(w, r, contact)
		return
	}
	// Update the db.
	err = c.ContactService.Update(&contact)
	if err != nil {
		fmt.Printf("ProcessEdit() err: %s", err)
	}
	// redirect to contacts page.
	SetFlash(w, "success", "Contact was Updated successfully")
	http.Redirect(w, r, "/contacts", http.StatusFound)
}

// Delete is the handler function that handles a POST request when someone tries to delete a contact.
func (c Contacts) Delete(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(urlParts[len(urlParts)-1])
	if err != nil {
		fmt.Printf("Edit(): err = %s", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	con := models.Contact{
		ID: id,
	}
	err = c.ContactService.Delete(&con)
	if err != nil {
		fmt.Printf("Delete() err: %s", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
	SetFlash(w, "success", "Contact was successfully deleted")
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}
