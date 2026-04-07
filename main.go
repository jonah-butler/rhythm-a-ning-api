package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rhythmapi/db"
	"rhythmapi/route"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("APP_ENV")
	var envFile string
	if env == "prod" {
		envFile = "prod.env"
	} else {
		envFile = "local.env"
	}

	fmt.Println("Loading environment variables from " + envFile)

	err := godotenv.Load(envFile)
	if err != nil {
		log.Println("Error loading .env file: " + err.Error())
	}

	db, err := db.InitDB()
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
	}

	r := route.SetupRoutes(db)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to listen on port 8080: %v", err)
		}
	}()

	fmt.Println("--------------------------")
	fmt.Println("|                        |")
	fmt.Println("| Listening on PORT 8080 |")
	fmt.Println("|                        |")
	fmt.Println("--------------------------")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Forced shutdown:", err)
	}

	log.Println("Server exited")
}
