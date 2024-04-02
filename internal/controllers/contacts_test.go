package controllers

import (
	"database/sql"
	"github.com/ethanjmarchand/learnhtmx/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Tester struct {
	// here, we can add a few extra data fields to increment when the function runs.
	Calls int
	DB    *sql.DB
}

func NewTesterContactRepo(db *sql.DB) *Tester {
	return &Tester{
		DB: db,
	}
}

func (t Tester) Read(name string) (*[]models.Contact, error) {
	return nil, nil
}

func (t Tester) ReadAll() (*[]models.Contact, error) {
	return nil, nil
}

func (t Tester) GetContact(id int) (*models.Contact, error) {
	t.Calls++
	return nil, nil
}

func (t Tester) New(contact *models.Contact) error {
	return nil
}
func (t Tester) Update(contact *models.Contact) error {
	return nil
}

func (t Tester) Delete(contact *models.Contact) error {
	return nil
}

func TestContacts_Edit(t *testing.T) {
	testCR := NewTesterContactRepo(&sql.DB{})
	testCS := models.NewContactService(testCR)
	contactTester := Contacts{
		ContactService: testCS,
	}
	w := httptest.NewRecorder()
	r := http.Request{}
	contactTester.Edit(w, &r)
}
