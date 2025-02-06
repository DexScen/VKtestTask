package main

import (
	"log"
	"net/http"
	"time"

	"github.com/DexScen/VKtestTask/backend/internal/repository/psql"
	"github.com/DexScen/VKtestTask/backend/pkg/database"
	_ "github.com/lib/pq"
)

func main() {
	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
		Password: "qwerty123",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	containersRepo := psql.NewContainers(db)
	containersService := service.NewContainers(containersRepo)
	handler := rest.NewHandler(containersService)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler.InitRouter(),
	}

	log.Println("Server started at:", time.Now().Format(time.RFC3339))

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
