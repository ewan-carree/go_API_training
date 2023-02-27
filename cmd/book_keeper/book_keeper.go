package book_keeper

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"goAPI/internal"
	"goAPI/internal/models"
	"goAPI/pkg"
)

type App struct {
    Router *mux.Router
    Db     *gorm.DB
}

func (a *App) Initialize() {
	// load env variables
	dialect, host, port, user, password, dbname := pkg.GetEnvVariables()

	// db connection string
	dbURI := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", host, port, user, password)

	// open connection to default database template1
	db, err := gorm.Open(dialect, dbURI+" dbname=template1")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// create database if it doesn't exist
	err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname)).Error
	if err != nil {
		// assume error is related to existing database, ignore and proceed
		fmt.Printf("Skipping creating database %s, assuming it already exists", dbname)
	} else {
		fmt.Printf("Database %s created successfully\n", dbname)
	}

	// update dbURI to use the created database
	dbURI += " dbname=" + dbname

	// connect to the created database
	a.Db, err = gorm.Open(dialect, dbURI)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("\nConnection Established")

	a.initTables(uuid.New().String())
	
	// Migrate the schema
	a.Db.AutoMigrate(&models.Person{}, &models.Book{})
	
	// API routes
	a.Router = mux.NewRouter()
	router.AddRoutes(a.Router, a.Db)
}

const tablePeopleCreationQuery = 
`	
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE IF NOT EXISTS people (
		id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
		name VARCHAR(100),
		email VARCHAR(100)
	);
`

const tableBooksCreationQuery =
`
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE IF NOT EXISTS books (
		id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
		title VARCHAR(100),
		author VARCHAR(100),
		person_id UUID
	);
`

func (a *App) initTables(name string) {
	if value := a.Db.Exec(tablePeopleCreationQuery); value.Error != nil {
		log.Fatal(value.Error)
	}
	if value := a.Db.Exec(tableBooksCreationQuery); value.Error != nil {
		log.Fatal(value.Error)
	}
}

func (a *App) Execute() {
	// open connection to db
	a.Initialize()
	defer a.Db.Close()

	// Run server
	http.ListenAndServe(":8080", a.Router)
}