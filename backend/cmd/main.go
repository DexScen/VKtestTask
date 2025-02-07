package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DexScen/VKtestTask/backend/internal/repository/psql"
	"github.com/DexScen/VKtestTask/backend/internal/service"
	"github.com/DexScen/VKtestTask/backend/internal/transport/rest"
	"github.com/DexScen/VKtestTask/backend/pkg/database"

	_ "github.com/lib/pq"
)

func main() {
	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     os.Getenv("DB_HOST"),
		Port:     5432,
		Username: os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  "disable",
		Password: os.Getenv("DB_PASSWORD"),
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
