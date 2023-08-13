package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"

	router "book_keeper/internal"
	"book_keeper/internal/models"

	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	Db     *gorm.DB
}

func (a *App) Initialize() {
	// Wait for the PostgreSQL container to be ready
	cmd := exec.Command("/app/wait-for-it.sh", "postgres:5432", "--timeout=30", "--")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic("Failed to wait for PostgreSQL to be ready: " + err.Error())
	}

	// Construct database connection string
	dbURI := "host=postgres port=5432 user=root password=test dbname=book_keeper sslmode=disable"

	// Connect to the database
	db, err := gorm.Open("postgres", dbURI)
	if err != nil {
		panic(err)
	}
	a.Db = db
	fmt.Println("Connection Established")

	// Migrate the schema
	a.Db.AutoMigrate(&models.Person{}, &models.Book{})

	// Initialize the router
	a.Router = mux.NewRouter()
	router.AddRoutes(a.Router, a.Db)
}

func (a *App) Execute() {
	// Initialize connection to db and router
	a.Initialize()
	defer a.Db.Close()

	// Init cors
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(a.Router)

	// Run server
	http.ListenAndServe(":8080", handler)
}

func main() {
	a := App{}
	a.Execute()
}
