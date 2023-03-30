package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/steffanturanjanin/receipt-manager/internal/controllers"
	"github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/queue"
	"github.com/steffanturanjanin/receipt-manager/internal/repositories"
	"github.com/steffanturanjanin/receipt-manager/internal/services"
	"github.com/steffanturanjanin/receipt-manager/internal/validator"
)

var (
	Session *session.Session = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	SqsService *sqs.SQS = sqs.New(Session)
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
	categoryRepository := repositories.NewCategoryRepository(database.Instance)
	receiptRepository := repositories.NewReceiptRepository(database.Instance)
	receiptItemRepository := repositories.NewReceiptItemRepository(database.Instance)
	statisticRepository := repositories.NewStatisticRepository(database.Instance)

	// Services
	authService := services.NewAuthService(userRepository)
	categoryService := services.NewCategoryService(categoryRepository)
	receiptService := services.NewReceiptService(receiptRepository)
	receiptItemService := services.NewReceiptItemService(receiptItemRepository, categoryService)
	statisticService := services.NewStatisticService(statisticRepository, categoryService)
	queueService := queue.NewQueueService(SqsService)

	// Controllers
	authController := controllers.NewAuthController(authService, validator)
	receiptController := controllers.NewReceiptController(receiptService, queueService, validator)
	receiptItemController := controllers.NewReceiptItemController(receiptItemService)
	statisticController := controllers.NewStatisticController(statisticService)

	// Routes
	mux := mux.NewRouter()
	// Auth routes
	mux.HandleFunc("/register", authController.RegisterUser).Methods("POST")
	mux.HandleFunc("/login", authController.LoginUser).Methods("POST")
	mux.HandleFunc("/logout", authController.Logout).Methods("POST")
	// Receipt routes
	mux.HandleFunc("/receipts", receiptController.CreateFromUrl).Methods("POST")
	mux.HandleFunc("/receipts", receiptController.List).Methods("GET")
	mux.HandleFunc("/receipts/{id}", receiptController.Show).Methods("GET")
	mux.HandleFunc("/receipts/{id}", receiptController.Delete).Methods("DELETE")
	// Receipt Item routes
	mux.HandleFunc("/receipt-items/{id}", receiptItemController.UpdateCategory).Methods("PUT")
	// Statistics routes
	mux.HandleFunc("/statistics/categories", statisticController.ListCategoriesStatistics).Methods("GET")
	mux.HandleFunc("/statistics/categories/{id}/stores", statisticController.ListStoreStatisticsForCategory).Methods("GET")

	// Workers
	receiptUrlQueueWorker, err := queue.NewReceiptUrlQueueWorker(queueService)
	if err != nil {
		log.Printf("Error trying to instantiate worker: %+v\n", err)
	}
	queueService.SpawnWorker(receiptUrlQueueWorker)

	receiptQueueWorker, err := queue.NewParsedReceiptQueueWorker(queueService, receiptService, categoryService)
	if err != nil {
		log.Printf("Error trying to instantiate worker: %+v\n", err)
	}
	queueService.SpawnWorker(receiptQueueWorker)

	fmt.Println("Server running at port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println(err)
	}
}
