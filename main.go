package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "kasir-api/docs"
)

type Config struct {
	Port          string `mapstructure:"PORT"`
	DBConn        string `mapstructure:"DB_CONN"`
	EnableSwagger bool   `mapstructure:"ENABLE_SWAGGER"`
}

func main() {
	mux := http.NewServeMux()

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:          viper.GetString("PORT"),
		DBConn:        viper.GetString("DB_CONN"),
		EnableSwagger: viper.GetBool("ENABLE_SWAGGER"),
	}

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Repositories
	productRepo := repositories.NewProductRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	// Services
	productService := services.NewProductService(productRepo, categoryRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	transactionService := services.NewTransactionService(transactionRepo)

	// Handlers
	productHandler := handlers.NewProductHandler(productService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Product routes
	mux.HandleFunc("/api/product", productHandler.HandleProducts)
	mux.HandleFunc("/api/product/", productHandler.HandleProductByID)

	// Category routes
	mux.HandleFunc("/api/categories", categoryHandler.HandleCategorys)
	mux.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)

	// Transaction routes
	mux.HandleFunc("/api/checkout", transactionHandler.HandleCheckout)

	// Health check route
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]string{
			"status":  strconv.Itoa(http.StatusOK),
			"message": "API running",
		})
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}

	})

	if config.EnableSwagger {
		mux.Handle("/swagger/", httpSwagger.WrapHandler)
		fmt.Println("swagger di enable")
	}

	fmt.Println("server running di localhost:8080")

	err = http.ListenAndServe(":"+config.Port, mux)

	if err != nil {
		fmt.Println("gagal running server")
	}
}
