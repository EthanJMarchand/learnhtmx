package controllers

import (
	"github.com/ethanjmarchand/learnhtmx/internal/models"
	"net/http"
	"testing"
)

func TestContacts_Edit(t *testing.T) {
	type fields struct {
		Templates struct {
			Contacts Template
			New      Template
			View     Template
			Edit     Template
		}
		UserService *models.UserService
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Contacts{
				Templates:   tt.fields.Templates,
				UserService: tt.fields.UserService,
			}
			c.Edit(tt.args.w, tt.args.r)
		})
	}
}
