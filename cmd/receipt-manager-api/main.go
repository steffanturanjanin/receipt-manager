package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/steffanturanjanin/receipt-manager/internal/controllers"
	"github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/repositories"
	"github.com/steffanturanjanin/receipt-manager/internal/services"
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
	validator := validator.New()

	// Repositories
	userRepository := repositories.NewUserRepository(database.Instance)

	// Services
	authService := services.NewAuthService(userRepository)

	// Controllers
	authController := controllers.NewAuthController(authService, validator)

	mux := mux.NewRouter()
	mux.HandleFunc("/register", authController.RegisterUser).Methods("POST")
	mux.HandleFunc("/login", authController.LoginUser).Methods("POST")
	mux.HandleFunc("/logout", authController.Logout).Methods("POST")

	fmt.Println("Server running at port 8080")
	http.ListenAndServe(":8080", mux)
}
