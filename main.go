package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Category represents a category in the cashier system
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// In-memory storage (sementara, nanti ganti database)
var category = []Category{
	{ID: 1, Name: "Electron", Description: "Electron products (Smartphone, Laptop, etc)"},
	{ID: 2, Name: "Health", Description: "Health products (Vitamins, Medicine, etc)"},
	{ID: 3, Name: "Hobbie", Description: "Hobbie products (Toys, Games, etc)"},
}

func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	// Parse ID dari URL path
	// URL: /api/category/123 -> ID = 123
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	// Cari category dengan ID tersebut
	for _, p := range category {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	// Kalau tidak found
	http.Error(w, "Category not found", http.StatusNotFound)
}

// PUT localhost:8080/api/category/{id}
func updateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	var updateData Category
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	for i := range category {
		if category[i].ID == id {
			// Update field satu-satu, jangan ganti ID
			category[i].Name = updateData.Name
			category[i].Description = updateData.Description

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(category[i])
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	// get id
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")

	// ganti id int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	// loop category cari ID, dapet index yang mau dihapus
	for i, p := range category {
		if p.ID == id {
			// bikin slice baru dengan data sebelum dan sesudah index
			category = append(category[:i], category[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "success delete",
			})
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

func main() {
	// GET localhost:8080/api/category/{id}
	// PUT localhost:8080/api/category/{id}
	// DELETE localhost:8080/api/category/{id}
	http.HandleFunc("/api/category/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getCategoryByID(w, r)
		case "PUT":
			updateCategory(w, r)
		case "DELETE":
			deleteCategory(w, r)
		default:
			// TAMBAHKAN INI
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// GET localhost:8080/api/category
	// POST localhost:8080/api/category
	http.HandleFunc("/api/category", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(category)
		case "POST":
			var categoryNew Category
			err := json.NewDecoder(r.Body).Decode(&categoryNew)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			categoryNew.ID = len(category) + 1
			category = append(category, categoryNew)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(categoryNew)
		default:
			// TAMBAHKAN INI
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	fmt.Println("Server running di localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}
