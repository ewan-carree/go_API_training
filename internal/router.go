package router

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"book_keeper/internal/models"
)


func AddRoutes(router *mux.Router, Db *gorm.DB) {
	router.HandleFunc("/people", func(w http.ResponseWriter, r *http.Request) {GetPeople(w, r, Db)}).Methods("GET")
	router.HandleFunc("/people/{id}", func(w http.ResponseWriter, r *http.Request) {GetPerson(w, r, Db)}).Methods("GET")
	router.HandleFunc("/people", func(w http.ResponseWriter, r *http.Request) {CreatePerson(w, r, Db)}).Methods("POST")
	router.HandleFunc("/people/{id}", func(w http.ResponseWriter, r *http.Request) {DeletePerson(w, r, Db)}).Methods("DELETE")
	router.HandleFunc("/people", func(w http.ResponseWriter, r *http.Request) {DeletePeople(w, r, Db)}).Methods("DELETE")
	router.HandleFunc("/people/{id}", func(w http.ResponseWriter, r *http.Request) {UpdatePerson(w, r, Db)}).Methods("PUT")
	router.HandleFunc("/people/{id}", func(w http.ResponseWriter, r *http.Request) {UpdatePerson(w, r, Db)}).Methods("PATCH")
	router.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {GetBooks(w, r, Db)}).Methods("GET")
	router.HandleFunc("/books/{id}", func(w http.ResponseWriter, r *http.Request) {GetBook(w, r, Db)}).Methods("GET")
	router.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {CreateBook(w, r, Db)}).Methods("POST")
	router.HandleFunc("/books/{id}", func(w http.ResponseWriter, r *http.Request) {DeleteBook(w, r, Db)}).Methods("DELETE")
	router.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {DeleteBooks(w, r, Db)}).Methods("DELETE")
	router.HandleFunc("/books/{id}", func(w http.ResponseWriter, r *http.Request) {UpdateBook(w, r, Db)}).Methods("PUT")
	router.HandleFunc("/books/{id}", func(w http.ResponseWriter, r *http.Request) {UpdateBook(w, r, Db)}).Methods("PATCH")
}

// GetPeople returns all people
func GetPeople(w http.ResponseWriter, r *http.Request, Db *gorm.DB) {
	var people []models.Person
	value := Db.Preload("Books").Find(&people)
	if value.Error != nil {
		http.Error(w, "No people found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(people)
}

// GetPerson returns a person
func GetPerson(w http.ResponseWriter, r *http.Request, Db *gorm.DB) {
	params := mux.Vars(r)
	var person models.Person
	id, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}
	value := Db.Preload("Books").First(&person, "id = ?", id)
	if value.Error != nil {
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(person)
}

// CreatePerson creates a person
func CreatePerson(w http.ResponseWriter, r *http.Request, Db *gorm.DB) {
	var person models.Person
	json.NewDecoder(r.Body).Decode(&person)
	value := Db.Create(&person)
	if value.Error != nil {
		http.Error(w, "Person not created", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(person)
}

// DeletePerson deletes a person
func DeletePerson(w http.ResponseWriter, r *http.Request, Db *gorm.DB) {
	params := mux.Vars(r)
	var person models.Person
	id, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}
	Db.Preload("Books").First(&person, "id = ?", id)

	// Update the PersonID field of all the books to NULL
	for _, book := range person.Books {
		Db.Delete(&book)
	}

	Db.Delete(&person)
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(person)
}

// DeletePeople deletes all people
func DeletePeople(w http.ResponseWriter, r *http.Request, Db *gorm.DB) {
	var people []models.Person
	Db.Preload("Books").Find(&people)

	// Update the PersonID field of all the books to NULL
	for _, person := range people {
		for _, book := range person.Books {
			Db.Delete(&book)
		}
	}

	Db.Delete(&people)
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(people)
}

// UpdatePerson updates a person
func UpdatePerson(w http.ResponseWriter, r *http.Request, Db *gorm.DB) {
	params := mux.Vars(r)
	var person models.Person
	id, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}
	Db.First(&person, "id = ?", id)
	json.NewDecoder(r.Body).Decode(&person)
	Db.Save(&person)
	json.NewEncoder(w).Encode(person)
}

// GetBooks returns all books
func GetBooks(w http.ResponseWriter, r *http.Request, Db *gorm.DB) {
	var books []models.Book
	value := Db.Find(&books)
	if value.Error != nil {
		http.Error(w, "No books found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(books)
}

// GetBook returns a book
func GetBook(w http.ResponseWriter, r *http.Request, Db *gorm.DB) {
	params := mux.Vars(r)
	var book models.Book
	id, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}
	value := Db.First(&book, "id = ?", id)
	if value.Error != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(book)
}

// CreateBook creates a book
func CreateBook(w http.ResponseWriter, r *http.Request, Db *gorm.DB) {
	var book models.Book
	json.NewDecoder(r.Body).Decode(&book)
	value := Db.Create(&book)
	if value.Error != nil {
		http.Error(w, "Book not created", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// DeleteBook deletes a book
func DeleteBook(w http.ResponseWriter, r *http.Request, Db *gorm.DB) {
	params := mux.Vars(r)
	var book models.Book
	id, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}
	Db.First(&book, "id = ?", id)
	Db.Delete(&book)
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(book)
}

// DeleteBooks deletes all books
func DeleteBooks(w http.ResponseWriter, r *http.Request, Db *gorm.DB) {
	var books []models.Book
	Db.Find(&books)
	Db.Delete(&books)
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(books)
}

// UpdateBook updates a book
func UpdateBook(w http.ResponseWriter, r *http.Request, Db *gorm.DB) {
	params := mux.Vars(r)
	var book models.Book
	id, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}
	Db.First(&book, "id = ?", id)
	json.NewDecoder(r.Body).Decode(&book)
	Db.Save(&book)
	json.NewEncoder(w).Encode(book)
}
