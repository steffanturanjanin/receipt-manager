package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/steffanturanjanin/receipt-manager/internal/controllers"
	"github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/repositories"
	"github.com/steffanturanjanin/receipt-manager/internal/services"
	"github.com/steffanturanjanin/receipt-manager/internal/validator"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed loading env variables with error: %s", err.Error())
	}

	err := database.InitializeDB()
	if err != nil {
		log.Fatalf("Failed trying to initialize the database with error: %s", err.Error())
	}

	// Validator
	translator := validator.NewEnglishTranslator()
	validator := validator.NewValidator(translator)
	validator.ConfigureValidator()

	// Repositories
	userRepository := repositories.NewUserRepository(database.Instance)
	receiptRepository := repositories.NewReceiptRepository(database.Instance)

	// Services
	authService := services.NewAuthService(userRepository)
	receiptService := services.NewReceiptService(receiptRepository)

	// Controllers
	authController := controllers.NewAuthController(authService, validator)
	receiptController := controllers.NewReceiptController(receiptService)

	mux := mux.NewRouter()
	mux.HandleFunc("/register", authController.RegisterUser).Methods("POST")
	mux.HandleFunc("/login", authController.LoginUser).Methods("POST")
	mux.HandleFunc("/logout", authController.Logout).Methods("POST")

	mux.HandleFunc("/receipts", receiptController.CreateFromUrl).Methods("POST")
	mux.HandleFunc("/receipts/{id}", receiptController.Delete).Methods("DELETE")

	fmt.Println("Server running at port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println(err)
	}
}
