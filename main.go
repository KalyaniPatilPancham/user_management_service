package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system.
type User struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Nickname  string    `json:"nickname"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var (
	users   = make(map[string]User)
	usersMu sync.RWMutex
)

// AddUser handles the addition of a new user.
func AddUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	usersMu.Lock()
	users[user.ID] = user
	usersMu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUser handles retrieving a user by ID.
func GetUser(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")

	usersMu.RLock()
	user, exists := users[id]
	usersMu.RUnlock()

	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// ListUsers handles listing all users with optional pagination and filtering.
func ListUsers(w http.ResponseWriter, r *http.Request) {
	// Get query parameters for pagination and filtering
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	countryFilter := r.URL.Query().Get("country")

	// Set default values for pagination if not provided
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// Filter users by country
	usersMu.RLock()
	defer usersMu.RUnlock()

	var filteredUsers []User
	for _, user := range users {
		if countryFilter == "" || strings.EqualFold(user.Country, countryFilter) {
			filteredUsers = append(filteredUsers, user)
		}
	}

	// Calculate total count after filtering
	total := len(filteredUsers)

	// Apply pagination logic
	start := (page - 1) * pageSize
	if start >= total {
		start = total // If the start index is out of range, start from the total (empty page)
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	paginatedUsers := filteredUsers[start:end]

	// Create response
	response := map[string]interface{}{
		"total": total,
		"users": paginatedUsers,
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateUser handles updating an existing user.
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	var updatedUser User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	usersMu.Lock()
	user, exists := users[id]
	if !exists {
		usersMu.Unlock()
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	updatedUser.ID = id
	updatedUser.CreatedAt = user.CreatedAt
	updatedUser.UpdatedAt = time.Now()
	users[id] = updatedUser
	usersMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

// DeleteUser handles deleting a user by User_ID.
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")

	usersMu.Lock()
	_, exists := users[id]
	if !exists {
		usersMu.Unlock()
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	delete(users, id)
	usersMu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

// HealthCheck handles a simple health check.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

func main() {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ListUsers(w, r)
		case http.MethodPost:
			AddUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetUser(w, r)
		case http.MethodPut:
			UpdateUser(w, r)
		case http.MethodDelete:
			DeleteUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/health", HealthCheck)

	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
