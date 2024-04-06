package controllers_test

import (
	"database/sql"
	"github.com/ethanjmarchand/learnhtmx/internal/controllers"
	"github.com/ethanjmarchand/learnhtmx/internal/errors"
	"github.com/ethanjmarchand/learnhtmx/internal/models"
	"github.com/ethanjmarchand/learnhtmx/internal/models/mocks"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestContacts_Show(t *testing.T) {
	t.Run("gives an internal server error given no contacts in DB and no query string", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTemplate := mocks.NewMockTemplate(ctrl)

		mr := mocks.NewMockContactRepo(ctrl)
		c := models.NewContactService(mr)
		mctrl := controllers.Contacts{
			Templates: struct {
				Contacts controllers.Template
				New      controllers.Template
				View     controllers.Template
				Edit     controllers.Template
			}{
				Contacts: mockTemplate,
			},
			ContactService: c,
		}

		mr.EXPECT().ReadAll().Times(1).Return(nil, sql.ErrNoRows)
		// This validates the template renders with the error
		mockTemplate.EXPECT().
			Execute(gomock.Any(), gomock.Any(), gomock.Any(), errors.Public(sql.ErrNoRows, "No results")).
			AnyTimes()

		req, err := http.NewRequest("GET", "/contacts", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(mctrl.Show)

		// Perform the request
		handler.ServeHTTP(rr, req)

	})

	t.Run("gives an internal server error given no contacts matching query string", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTemplate := mocks.NewMockTemplate(ctrl)

		mr := mocks.NewMockContactRepo(ctrl)
		c := models.NewContactService(mr)
		mctrl := controllers.Contacts{
			Templates: struct {
				Contacts controllers.Template
				New      controllers.Template
				View     controllers.Template
				Edit     controllers.Template
			}{
				Contacts: mockTemplate,
			},
			ContactService: c,
		}

		req, err := http.NewRequest("GET", "/contacts", nil)
		if err != nil {
			t.Fatal(err)
		}

		form := url.Values{}
		form.Add("q", "John Doe")
		req.URL.RawQuery = form.Encode()

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(mctrl.Show)

		mr.EXPECT().Read("John Doe").Times(1).Return(nil, sql.ErrNoRows)
		// This validates the template renders with the error
		mockTemplate.EXPECT().
			Execute(gomock.Any(), gomock.Any(), gomock.Any(), errors.Public(sql.ErrNoRows, "No results")).
			AnyTimes()

		// Perform the request
		handler.ServeHTTP(rr, req)

	})

	t.Run("gives no error given the user can be found in the DB", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTemplate := mocks.NewMockTemplate(ctrl)

		mr := mocks.NewMockContactRepo(ctrl)
		c := models.NewContactService(mr)
		mctrl := controllers.Contacts{
			Templates: struct {
				Contacts controllers.Template
				New      controllers.Template
				View     controllers.Template
				Edit     controllers.Template
			}{
				Contacts: mockTemplate,
			},
			ContactService: c,
		}

		expectedContacts := []models.Contact{{Name: "John Doe", Email: "johndoe@example.com"}}

		req, err := http.NewRequest("GET", "/contacts", nil)
		if err != nil {
			t.Fatal(err)
		}

		form := url.Values{}
		form.Add("q", "John Doe")
		req.URL.RawQuery = form.Encode()

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(mctrl.Show)

		mr.EXPECT().Read("John Doe").Times(1).Return(&expectedContacts, nil)
		// This validates the template renders with the error
		mockTemplate.EXPECT().
			Execute(gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()

		// Perform the request
		handler.ServeHTTP(rr, req)

	})

}
