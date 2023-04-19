package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"loan-api/internal/handler"
	"loan-api/internal/models"
	provider "loan-api/internal/service/accounting"
	"loan-api/internal/service/decision"
	"loan-api/internal/service/user"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	db, err := getDb()

	if err != nil {
		log.Fatal(err)
	}

	userRepo := models.NewUserRepo(db)
	userService := user.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	accRepo := models.NewProviderRepo(db)
	accService := provider.NewAccountingService(accRepo)
	accHandler := handler.NewProviderHandler(accService)

	decService := decision.NewDecisionService()
	decHandler := handler.NewDecisionHandler(decService)

	router := mux.NewRouter()

	router.Use(setAccessControlHeader)

	router.HandleFunc("/login", userHandler.Login).Methods(http.MethodPost)
	router.HandleFunc("/providers", accHandler.GetProviders).Methods(http.MethodGet)
	router.HandleFunc("/balance-sheet", accHandler.GetBalanceSheet).Methods(http.MethodPost)
	router.HandleFunc("/calculate-loan", decHandler.CalculateLoan).Methods(http.MethodPost)

	router.Methods(http.MethodOptions).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	}).Methods(http.MethodOptions)

	log.Println("Application has started. Listening port is 4000")
	http.ListenAndServe(":4000", router)
}

func setAccessControlHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		next.ServeHTTP(w, r)
	})
}

func getDb() (*sqlx.DB, error) {
	connConfig, err := pgx.ParseConfig("postgres://postgres:postgres@db:5432/loan-db?sslmode=disable")
	if err != nil {
		errMsg := err.Error()
		errMsg = regexp.MustCompile(`(://[^:]+:).+(@.+)`).ReplaceAllString(errMsg, "$1*****$2")
		errMsg = regexp.MustCompile(`(password=).+(\s+)`).ReplaceAllString(errMsg, "$1*****$2")
		return nil, fmt.Errorf("parsing DSN failed: %s", errMsg)
	}
	connectionStr := stdlib.RegisterConnConfig(connConfig)
	db, err := sqlx.Open("pgx", connectionStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	instance, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"loan-db",
		instance,
	)
	if err != nil {
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	seed(db)

	return db, nil
}

func seed(db *sqlx.DB) {
	users := []struct {
		Username string
		Password string
	}{
		{Username: "Alice", Password: "password123"},
		{Username: "Bob", Password: "password123"},
		{Username: "Charlie", Password: "password123"},
		{Username: "David", Password: "password123"},
	}
	for _, user := range users {

		db.Exec("INSERT INTO users(username, password) VALUES ($1,$2) ON CONFLICT DO NOTHING", user.Username, user.Password)
	}

	var counter int
	db.QueryRow("SELECT id FROM providers limit 1").Scan(&counter)
	if counter == 0 {
		providers := []struct {
			Name   string
			Slug   string
			Status int64
		}{
			{Name: "Xero", Slug: "xero", Status: 1},
			{Name: "MYOB", Slug: "myob", Status: 1},
			{Name: "MYOB11", Slug: "myob11", Status: 0},
		}

		for _, provider := range providers {

			db.Exec("INSERT INTO providers(name, slug, status) VALUES ($1,$2,$3) ON CONFLICT DO NOTHING", provider.Name, provider.Slug, provider.Status)
			log.Println("Database Seeding line 1 Completed.")

		}
	}
	log.Println("Database Seeding Completed.")
}
