package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ethanjmarchand/learnhtmx/internal/controllers"
	"github.com/ethanjmarchand/learnhtmx/internal/migrations"
	"github.com/ethanjmarchand/learnhtmx/internal/models"
	"github.com/ethanjmarchand/learnhtmx/internal/templates"
	"github.com/ethanjmarchand/learnhtmx/internal/views"
	"github.com/joho/godotenv"
)

type config struct {
	PSQL   models.PostgresConfig
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return cfg, err
	}
	// Setup psql
	cfg.PSQL = models.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		Database: os.Getenv("PSQL_DATABASE"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
	}
	if cfg.PSQL.Host == "" && cfg.PSQL.Port == "" {
		return cfg, fmt.Errorf("no psql config provided")
	}

	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")

	return cfg, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}
	err = run(cfg)
	if err != nil {
		panic(err)
	}

}

func run(cfg config) error {

	db, err := models.Open(cfg.PSQL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		return err
	}
	// setup services
	userService := &models.UserService{
		DB: db,
	}

	// Setup Controllers, and pass the services in.

	contactsC := controllers.Contacts{
		UserService: userService,
	}

	contactsC.Templates.Contacts = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "contacts.gohtml"))
	contactsC.Templates.New = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "new.gohtml"))
	contactsC.Templates.View = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "view.gohtml"))
	contactsC.Templates.Edit = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "edit.gohtml"))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", contactsC.Redirect)
	mux.HandleFunc("GET /contacts", contactsC.Show)
	mux.HandleFunc("GET /contacts/new", contactsC.New)
	mux.HandleFunc("POST /contacts/new", contactsC.ProcessNew)
	mux.HandleFunc("GET /contacts/{id}", contactsC.View)
	mux.HandleFunc("GET /contacts/{id}/edit", contactsC.Edit)
	mux.HandleFunc("POST /contacts/{id}/edit", contactsC.ProcessEdit)
	mux.HandleFunc("DELETE /contacts/{id}", contactsC.Delete)

	assetsHandler := http.FileServer(http.Dir("assets"))
	mux.HandleFunc("GET /assets/*", http.StripPrefix("/assets", assetsHandler).ServeHTTP)

	fmt.Printf("Starting server on port %s", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, mux)
	if err != nil {
		return err
	}
	return nil
}
