package main_test

import (
	"bytes"
	"encoding/json"
	"goAPI/cmd/book_keeper"
	"goAPI/internal/models"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
)

var randomUUID = uuid.New()

var a book_keeper.App

func TestMain(m *testing.M) {
	a.Initialize()

	code := m.Run()
	clearTable()
	os.Exit(code)
}

func clearTable() {
	a.Db.Exec("DELETE FROM people")
	a.Db.Exec("DELETE FROM books")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestEmptyPeopleTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/people", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]\n" {
		t.Errorf("Expected []. Got %s", body)
	}
}

func TestGetNonExistentPerson(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/people/"+randomUUID.String(), nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	if body:= response.Body.String(); body != "Person not found\n" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Person not found'. Got '%s'", body)
	}
}

func TestGetPersonWithWrongUUIDFormat(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/people/123", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

	if body:= response.Body.String(); body != "Invalid UUID\n" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Invalid UUID'. Got '%s'", body)
	}
}

func TestCreatePerson(t *testing.T) {
	clearTable()

	p := models.Person{ID: randomUUID, Name: "p1", EMail: "p1@gmail.com"}
	payload, _ := json.Marshal(p)

	req, _ := http.NewRequest("POST", "/people", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["Name"] != "p1" {
		t.Errorf("Expected person Name to be 'p1'. Got '%v'", m["Name"])
	}

	if m["EMail"] != "p1@gmail.com" {
		t.Errorf("Expected person EMail to be 'p1@gmail.com'. Got '%v'", m["EMail"])
	}

	if m["ID"] != p.ID.String() {
		t.Errorf("Expected person ID to be '%v'. Got '%v'", p.ID, m["ID"])
	}
}

func TestGetPerson(t *testing.T) {
	TestCreatePerson(t)

	req, _ := http.NewRequest("GET", "/people/"+randomUUID.String(), nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["Name"] != "p1" {
		t.Errorf("Expected person Name to be 'p1'. Got '%v'", m["Name"])
	}

	if m["EMail"] != "p1@gmail.com" {
		t.Errorf("Expected person EMail to be 'p1@gmail.com'. Got '%v'", m["EMail"])
	}

	if m["ID"] != randomUUID.String() {
		t.Errorf("Expected person ID to be '%v'. Got '%v'", randomUUID.String(), m["ID"])
	}
}

func TestUpdatePerson(t *testing.T) {
	TestCreatePerson(t)

	p := models.Person{ID: randomUUID, Name: "p1-updated", EMail: "p1@gmail.com"}
	payload, _ := json.Marshal(p)

	req, _ := http.NewRequest("PUT", "/people/"+randomUUID.String(), bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["Name"] != "p1-updated" {
		t.Errorf("Expected person Name to be 'p1-updated'. Got '%v'", m["Name"])
	}
	if m["EMail"] != "p1@gmail.com" {
		t.Errorf("Expected person EMail to be 'p1@gmail.com'. Got '%v'", m["EMail"])
	}
	if m["ID"] != randomUUID.String() {
		t.Errorf("Expected person ID to be '%v'. Got '%v'", randomUUID.String(), m["ID"])
	}
}

func TestDeletePerson(t *testing.T) {
	TestCreatePerson(t)

	req, _ := http.NewRequest("DELETE", "/people/"+randomUUID.String(), nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNoContent, response.Code)

	req, _ = http.NewRequest("GET", "/people/"+randomUUID.String(), nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestDeletePeople(t *testing.T) {
	TestCreatePerson(t)

	req, _ := http.NewRequest("DELETE", "/people", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNoContent, response.Code)

	req, _ = http.NewRequest("GET", "/people", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]\n" {
		t.Errorf("Expected []. Got %s", body)
	}
}

func TestEmptyBooksTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/books", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]\n" {
		t.Errorf("Expected []. Got %s", body)
	}
}

func TestGetNonExistentBook(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/books/"+randomUUID.String(), nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	if body:= response.Body.String(); body != "Book not found\n" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Book not found'. Got '%s'", body)
	}
}

func TestGetBookWithWrongUUIDFormat(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/books/123", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

	if body:= response.Body.String(); body != "Invalid UUID\n" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Invalid UUID'. Got '%s'", body)
	}
}

func TestCreateBook(t *testing.T) {
	clearTable()

	b := models.Book{ID: randomUUID, Title: "b1", Author: "a1"}
	payload, _ := json.Marshal(b)

	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["Title"] != "b1" {
		t.Errorf("Expected book Title to be 'b1'. Got '%v'", m["Title"])
	}

	if m["Author"] != "a1" {
		t.Errorf("Expected book Author to be 'a1'. Got '%v'", m["Author"])
	}

	if m["ID"] != b.ID.String() {
		t.Errorf("Expected book ID to be '%v'. Got '%v'", b.ID, m["ID"])
	}
}

func TestGetBook(t *testing.T) {
	TestCreateBook(t)

	req, _ := http.NewRequest("GET", "/books/"+randomUUID.String(), nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["Title"] != "b1" {
		t.Errorf("Expected book Title to be 'b1'. Got '%v'", m["Title"])
	}

	if m["Author"] != "a1" {
		t.Errorf("Expected book Author to be 'a1'. Got '%v'", m["Author"])
	}

	if m["ID"] != randomUUID.String() {
		t.Errorf("Expected book ID to be '%v'. Got '%v'", randomUUID.String(), m["ID"])
	}
}

func TestUpdateBook(t *testing.T) {
	TestCreateBook(t)

	b := models.Book{ID: randomUUID, Title: "b1-updated", Author: "a1"}
	payload, _ := json.Marshal(b)

	req, _ := http.NewRequest("PUT", "/books/"+randomUUID.String(), bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["Title"] != "b1-updated" {
		t.Errorf("Expected book Title to be 'b1-updated'. Got '%v'", m["Title"])
	}
	if m["Author"] != "a1" {
		t.Errorf("Expected book Author to be 'a1'. Got '%v'", m["Author"])
	}
	if m["ID"] != randomUUID.String() {
		t.Errorf("Expected book ID to be '%v'. Got '%v'", randomUUID.String(), m["ID"])
	}
}

func TestDeleteBook(t *testing.T) {
	TestCreateBook(t)

	req, _ := http.NewRequest("DELETE", "/books/"+randomUUID.String(), nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNoContent, response.Code)

	req, _ = http.NewRequest("GET", "/books/"+randomUUID.String(), nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestDeleteBooks(t *testing.T) {
	TestCreateBook(t)

	req, _ := http.NewRequest("DELETE", "/books", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNoContent, response.Code)

	req, _ = http.NewRequest("GET", "/books", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]\n" {
		t.Errorf("Expected []. Got %s", body)
	}
}