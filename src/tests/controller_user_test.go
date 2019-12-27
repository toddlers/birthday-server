package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sureshk/birthday-server/src/api/models"
	"gopkg.in/go-playground/assert.v1"
)



func TestCreateUser(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	samples := []struct {
		inputJSON    string
		statusCode   int
		name         string
		email        string
		errorMessage string
	}{
		{
			inputJSON:    `{"name":"sureshk", "email": "suresh@gmail.com", "birthday": "1990-06-04"}`,
			statusCode:   201,
			name:         "sureshk",
			email:        "suresh@gmail.com",
			errorMessage: "",
		},
		{
			inputJSON:    `{"name":"sureshk", "email": "grand@gmail.com", "birthday": "1975-01-08"}`,
			statusCode:   500,
			errorMessage: "Name Already Taken",
		},
		{
			inputJSON:    `{"name":"Kan", "email": "kangmail.com", "birthday": "1985-05-06"}`,
			statusCode:   422,
			errorMessage: "Invalid Email",
		},
		{
			inputJSON:    `{"name": "", "email": "kan@gmail.com", "birthday": "1965-07-04"}`,
			statusCode:   422,
			errorMessage: "Required Name",
		},
		{
			inputJSON:    `{"name": "Kan", "email": "", "birthday": "1932-08-09"}`,
			statusCode:   422,
			errorMessage: "Required Email",
		},
		{
			inputJSON:    `{"name": "Kan", "email": "kan@gmail.com", "birthday": ""}`,
			statusCode:   422,
			errorMessage: "Required Birthday",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateUser)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["name"], v.name)
			assert.Equal(t, responseMap["email"], v.email)
		}
		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetUsers(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	_, err = seedUsers()
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetUsers)
	handler.ServeHTTP(rr, req)

	var users []models.User
	err = json.Unmarshal(rr.Body.Bytes(), &users)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(users), 2)
}

func TestGetUserByID(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}
	userSample := []struct {
		id           string
		statusCode   int
		Name         string
		email        string
		birthday     string
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(user.ID)),
			statusCode: 200,
			Name:       user.Name,
			email:      user.Email,
			birthday:   user.Birthday,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range userSample {

		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetUser)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, user.Name, responseMap["name"])
			assert.Equal(t, user.Email, responseMap["email"])
		}
	}
}

func TestUpdateUser(t *testing.T) {

	//var AuthEmail, AuthPassword string
	var AuthID uint32

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	users, err := seedUsers() //we need atleast two users to properly check the update
	if err != nil {
		log.Fatalf("Error seeding user: %v\n", err)
	}
	// Get only the first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		AuthID = user.ID
	}

	samples := []struct {
		id           string
		updateJSON   string
		statusCode   int
		updateName   string
		updateEmail  string
		tokenGiven   string
		errorMessage string
	}{
		{
			// Convert int32 to int first before converting to string
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"name":"Grand", "email": "grand@gmail.com", "birthday": "1990-02-01"}`,
			statusCode:   200,
			updateName:   "Grand",
			updateEmail:  "grand@gmail.com",
			errorMessage: "",
		},
		{
			// When birthday field is empty
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"name":"Woman", "email": "woman@gmail.com", "birthday": ""}`,
			statusCode:   422,
			errorMessage: "Required Birthday",
		},
		{
			// When birthday field is empty
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"name":"Woman", "email": "woman@gmail.com", "birthday": "2019-12-28"}`,
			statusCode:   422,
			errorMessage: "Birthday should be before today",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"name":"Kan", "email": "kangmail.com", "birthday": "1760-02-08"}`,
			statusCode:   422,
			errorMessage: "Invalid Email",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"name": "", "email": "kan@gmail.com", "birthday": "1991-07-08"}`,
			statusCode:   422,
			errorMessage: "Required Name",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"name": "123su resh", "email": "kan@gmail.com", "birthday": "1991-07-08"}`,
			statusCode:   422,
			errorMessage: "Only alphanumeric usernames accepted",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"name": "Kan", "email": "", "birthday": "1865-02-03"}`,
			statusCode:   422,
			errorMessage: "Required Email",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateUser)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["name"], v.updateName)
			assert.Equal(t, responseMap["email"], v.updateEmail)
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteUser(t *testing.T) {

	//	var AuthEmail, AuthPassword string
	var AuthID uint32

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	users, err := seedUsers() //we need atleast two users to properly check the update
	if err != nil {
		log.Fatalf("Error seeding user: %v\n", err)
	}

	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		AuthID = user.ID
	}

	userSample := []struct {
		id           string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int32 to int first before converting to string
			id:           strconv.Itoa(int(AuthID)),
			statusCode:   204,
			errorMessage: "",
		},
	}

	for _, v := range userSample {

		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteUser)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
