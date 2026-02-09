package main

import (
	"encoding/json"
	"fmt"
	"kasir-go-api/database"
	"kasir-go-api/handlers"
	"kasir-go-api/repositories"
	"kasir-go-api/services"
	"log"
	"net/http"
	"os"

	// "strconv"
	"strings"

	"github.com/spf13/viper"
)

// ubah Config
type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

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

	log.Printf("Server akan berjalan di port: %s", config.Port)
	if config.DBConn != "" {
		log.Printf("Database connection string: %s", config.DBConn)
	}

	// ✅ Setup product repository, service, dan handler
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// ✅ Setup routes untuk produk
	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)

	// ✅ Setup category repository, service, dan handler
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// ✅ Setup routes untuk category
	http.HandleFunc("/api/categories", categoryHandler.HandleCategory)
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)

	// ✅ Setup transaction repository, service, dan handler
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	http.HandleFunc("/api/checkout", transactionHandler.HandleCheckout) // POST

	// ✅ Setup report repository, service, dan handler
	reportRepo := repositories.NewReportRepository(db)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)

	// ✅ Setup routes untuk report
	http.HandleFunc("/api/report/hari-ini", reportHandler.HandleDailySalesReport) // GET
	http.HandleFunc("/api/report", reportHandler.HandleSalesReport)                // GET with query params

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running di", addr)
	fmt.Println("\n📋 Available endpoints:")
	fmt.Println("  Health:")
	fmt.Println("    - GET  /health")
	fmt.Println("\n  Products:")
	fmt.Println("    - GET    /api/produk")
	fmt.Println("    - POST   /api/produk")
	fmt.Println("    - GET    /api/produk/{id}")
	fmt.Println("    - PUT    /api/produk/{id}")
	fmt.Println("    - DELETE /api/produk/{id}")
	fmt.Println("\n  Categories:")
	fmt.Println("    - GET    /api/categories")
	fmt.Println("    - POST   /api/categories")
	fmt.Println("    - GET    /api/categories/{id}")
	fmt.Println("    - PUT    /api/categories/{id}")
	fmt.Println("    - DELETE /api/categories/{id}")
	fmt.Println("\n  Transaction:")
	fmt.Println("    - POST   /api/checkout")
	fmt.Println("\n  Reports:")
	fmt.Println("    - GET    /api/report/hari-ini")
	fmt.Println("    - GET    /api/report?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD")

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("Gagal running server:", err)
	}
}