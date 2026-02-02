package main

import (
	"crud-kategori/database"
	"crud-kategori/handlers"
	"crud-kategori/repositories"
	"crud-kategori/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Category represents a product in the cashier system
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ubah Config
type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

// In-memory storage (sementara, nanti ganti database)
var categories = []Category{
	{ID: 1, Name: "Minuman", Description: "Minuman segar dan hangat"},
	{ID: 2, Name: "Makanan", Description: "Makanan berat dan ringan"},
	{ID: 3, Name: "Pakaian", Description: "Pakaian untuk segala usia"},
	{ID: 4, Name: "Aksesoris", Description: "Aksesoris untuk segala usia"},
	{ID: 5, Name: "Elektronik", Description: "Perangkat elektronik dan gadget"},
	{ID: 6, Name: "Olahraga", Description: "Peralatan fitness dan apparel gym"},
}

func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	// Parse ID dari URL path
	// URL: /api/kategori/123 -> ID = 123
	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Kategori ID", http.StatusBadRequest)
		return
	}

	// Cari produk dengan ID tersebut
	for _, p := range categories {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	// Kalau tidak found
	http.Error(w, "Kategori belum ada", http.StatusNotFound)
}

// PUT localhost:8080/api/kategori/{id}
func updateKategori(w http.ResponseWriter, r *http.Request) {
	// get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")

	// ganti int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	// get data dari request
	var updateProduk Category
	err = json.NewDecoder(r.Body).Decode(&updateProduk)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// loop produk, cari id, ganti sesuai data dari request
	for i := range categories {
		if categories[i].ID == id {
			updateProduk.ID = id
			categories[i] = updateProduk

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateProduk)
			return
		}
	}

	http.Error(w, "Kategori belum ada", http.StatusNotFound)
}

func deleteKategori(w http.ResponseWriter, r *http.Request) {
	// get id
	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")

	// ganti id int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	// loop produk cari ID, dapet index yang mau dihapus
	for i, p := range categories {
		if p.ID == id {
			// bikin slice baru dengan data sebelum dan sesudah index
			categories = append(categories[:i], categories[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "sukses delete",
			})
			return
		}
	}

	http.Error(w, "Kategori belum ada", http.StatusNotFound)
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}
	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	// GET localhost:8080/api/kategori
	// POST localhost:8080/api/kategori
	http.HandleFunc("/api/kategori", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories)
		case "POST":
			// baca data dari request
			var kategoriBaru Category
			err := json.NewDecoder(r.Body).Decode(&kategoriBaru)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			// masukkin data ke dalam variable produk
			kategoriBaru.ID = len(categories) + 1
			categories = append(categories, kategoriBaru)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated) // 201
			json.NewEncoder(w).Encode(kategoriBaru)
		case "PUT":
			updateKategori(w, r)
		case "DELETE":
			deleteKategori(w, r)
		}
	})

	// Handler untuk /categories (GET: list all, POST: create)
	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories)
		case "POST":
			var kategoriBaru Category
			err := json.NewDecoder(r.Body).Decode(&kategoriBaru)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}
			kategoriBaru.ID = len(categories) + 1
			categories = append(categories, kategoriBaru)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(kategoriBaru)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Handler untuk /categories/{id} (GET: get by id, PUT: update, DELETE: delete)
	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid Category ID", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case "GET":
			for _, cat := range categories {
				if cat.ID == id {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(cat)
					return
				}
			}
			http.Error(w, "Category not found", http.StatusNotFound)
		case "PUT":
			var updateCat Category
			err := json.NewDecoder(r.Body).Decode(&updateCat)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}
			for i := range categories {
				if categories[i].ID == id {
					updateCat.ID = id
					categories[i] = updateCat
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(updateCat)
					return
				}
			}
			http.Error(w, "Category not found", http.StatusNotFound)
		case "DELETE":
			for i, cat := range categories {
				if cat.ID == id {
					categories = append(categories[:i], categories[i+1:]...)
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(map[string]string{"message": "Category deleted"})
					return
				}
			}
			http.Error(w, "Category not found", http.StatusNotFound)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Setup routes
	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// Setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running di", addr)

}
