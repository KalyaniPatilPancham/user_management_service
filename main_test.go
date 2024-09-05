package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Helper function to create a new HTTP request with JSON body
func createRequest(t *testing.T, method, url string, body interface{}) *http.Request {
	t.Helper()
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Failed to marshal body: %v", err)
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(b))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

// Test adding a user
func TestAddUser(t *testing.T) {
	reqBody := map[string]string{
		"first_name": "Alice",
		"last_name":  "Smith",
		"nickname":   "alice123",
		"password":   "securepassword",
		"email":      "alice@smith.com",
		"country":    "UK",
	}

	req := createRequest(t, http.MethodPost, "/users", reqBody)
	rr := httptest.NewRecorder()

	http.HandlerFunc(AddUser).ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
	}

	var user User
	if err := json.Unmarshal(rr.Body.Bytes(), &user); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if user.FirstName != reqBody["first_name"] || user.LastName != reqBody["last_name"] {
		t.Errorf("Expected user with first_name %s and last_name %s, got %s %s", reqBody["first_name"], reqBody["last_name"], user.FirstName, user.LastName)
	}
}

// Test listing users with pagination and filtering
func TestListUsersPagination(t *testing.T) {
	// Add some test users first
	testUsers := []User{
		{ID: "1", FirstName: "Alice", LastName: "Smith", Country: "UK", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "2", FirstName: "Bob", LastName: "Johnson", Country: "UK", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "3", FirstName: "Charlie", LastName: "Brown", Country: "UK", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "4", FirstName: "Shri", LastName: "Williams", Country: "India", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	// Clear and add test users
	usersMu.Lock()
	users = make(map[string]User)
	for _, u := range testUsers {
		users[u.ID] = u
	}
	usersMu.Unlock()

	// Test with pageSize = 2
	req, err := http.NewRequest(http.MethodGet, "/users?page=1&pageSize=2&country=UK", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()
	http.HandlerFunc(ListUsers).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Validate the total number of users after filtering by country
	if response["total"].(float64) != 3 {
		t.Errorf("Expected total 3, got %v", response["total"])
	}

	// Validate the number of users returned matches the pageSize
	if len(response["users"].([]interface{})) != 2 {
		t.Errorf("Expected 2 users on the first page, got %d", len(response["users"].([]interface{})))
	}

	// Test with pageSize = 2 and page = 2
	req, err = http.NewRequest(http.MethodGet, "/users?page=2&pageSize=2&country=UK", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	rr = httptest.NewRecorder()
	http.HandlerFunc(ListUsers).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Validate the total number of users after filtering by country
	if response["total"].(float64) != 3 {
		t.Errorf("Expected total 3, got %v", response["total"])
	}

	// Validate the number of users returned on the second page
	if len(response["users"].([]interface{})) != 1 {
		t.Errorf("Expected 1 user on the second page with pageSize=2, got %d", len(response["users"].([]interface{})))
	}
}

// Test updating a user
func TestUpdateUser(t *testing.T) {
	testUser := User{
		ID:        "123",
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "jdoe",
		Email:     "johndoe@example.com",
		Country:   "US",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	usersMu.Lock()
	users[testUser.ID] = testUser
	usersMu.Unlock()

	updatedData := map[string]string{
		"first_name": "Jane",
		"last_name":  "Doe",
		"nickname":   "jdoe",
		"email":      "janedoe@example.com",
		"country":    "US",
	}

	req := createRequest(t, http.MethodPut, "/users/"+testUser.ID, updatedData)
	rr := httptest.NewRecorder()

	http.HandlerFunc(UpdateUser).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var updatedUser User
	if err := json.Unmarshal(rr.Body.Bytes(), &updatedUser); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if updatedUser.FirstName != updatedData["first_name"] || updatedUser.LastName != updatedData["last_name"] {
		t.Errorf("Expected updated user with first_name %s and last_name %s, got %s %s", updatedData["first_name"], updatedData["last_name"], updatedUser.FirstName, updatedUser.LastName)
	}
}

// Test deleting a user
func TestDeleteUser(t *testing.T) {
	testUser := User{
		ID:        "456",
		FirstName: "Test",
		LastName:  "User",
		Country:   "Canada",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	usersMu.Lock()
	users[testUser.ID] = testUser
	usersMu.Unlock()

	req, err := http.NewRequest(http.MethodDelete, "/users/"+testUser.ID, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()

	http.HandlerFunc(DeleteUser).ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, rr.Code)
	}

	usersMu.RLock()
	_, exists := users[testUser.ID]
	usersMu.RUnlock()
	if exists {
		t.Errorf("User should have been deleted, but still exists")
	}
}

// Test health check endpoint
func TestHealthCheck(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()

	http.HandlerFunc(HealthCheck).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if strings.TrimSpace(string(body)) != "OK" {
		t.Errorf("Expected response body 'OK', got '%s'", strings.TrimSpace(string(body)))
	}
}
