package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/mvrilo/go-redoc"
	"github.com/onlytunesradio/go-api-template/src/api"
	config "github.com/onlytunesradio/go-api-template/src/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Redoc Environment
	// ================
	doc := redoc.Redoc{
		DocsPath: "/docs",
		// Change SpecPath && SpecFile to ./static/swagger.json when developing locally!
		SpecPath:    "/srv/static/swagger.json",
		SpecFile:    "/srv/static/swagger.json",
		Title:       "OnlyTunes API Template",
		Description: "API Documentation for OnlyTunes API Template",
	}
	// =================
	// Declaring Environment variables
	// =================
	var DBHost, DBUser, DBPass, DBName, DBPort string
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file, Using container variables: ERR: %v", err)
	}
	DBHost = os.Getenv("DB_HOST")
	DBUser = os.Getenv("DB_USER")
	DBPass = os.Getenv("DB_PASS")
	DBName = os.Getenv("DB_NAME")
	DBPort = os.Getenv("DB_PORT")
	// =================
	// Declaring Database Connection
	// =================
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", DBHost, DBPort, DBUser, DBPass, DBName)
	config.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Uncomment the line below to disable the logging Gorm does by default
		//Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Printf("Error connecting to the database: ERR: %v", err)
	}
	// =================
	// Initialize Router and WebServer
	// =================
	r := chi.NewRouter()
	// Uncomment these during development / debugging
	//r.Use(middleware.Logger)
	//r.Use(middleware.RealIP)
	//r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)

	// =================
	// Set Custom status messages for 404 && 405
	// =================
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("404 - Page not found - Check the API Docs for more info"))
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		w.Write([]byte("405 - Method not allowed - Check the API Docs for more info"))
	})
	// =================
	// Initialize API Middleware
	// =================
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.NoCache)
	// =================
	// Initialize API Routes
	// =================
	r.Mount("/test", api.TestRouter())
	// =================
	// Initialize API Documentation
	// =================
	// Change the http.Dir to ./static for local development!
	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("/srv/static"))))
	r.Handle("/docs", doc.Handler())
	// =================
	// Start WebServer
	// =================
	fmt.Printf("Starting Server on port 4000\n")
	err = http.ListenAndServe(":4000", r)
	if err != nil {
		log.Fatal(err)
	}
}
