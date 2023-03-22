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
	Session    *session.Session
	SqsService *sqs.SQS
)

func init() {
	Session = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	SqsService = sqs.New(Session)
}

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
	mux.HandleFunc("/receipts", receiptController.List).Methods("GET")
	mux.HandleFunc("/receipts/{id}", receiptController.Delete).Methods("DELETE")

	mux.HandleFunc("/test-function", func(w http.ResponseWriter, r *http.Request) {
		type url struct {
			Url string `json:"url"`
		}

		u := url{}
		controllers.ParseBody(&u, r)

		queueName := queue.RECEIPT_URL_QUEUE
		urlResult, _ := SqsService.GetQueueUrl(&sqs.GetQueueUrlInput{
			QueueName: &queueName,
		})

		producer := queue.NewReceiptUrlProducer(SqsService, *urlResult.QueueUrl)
		if err := producer.Produce(u.Url); err != nil {
			w.Write([]byte("{'error': 'Error occured while writing url to queue.'}"))
			return
		}
	}).Methods("POST")

	queueName := queue.RECEIPT_URL_QUEUE
	urlResult, _ := SqsService.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})

	consumer, err := queue.NewReceiptUrlConsumer(*urlResult.QueueUrl)
	if err != nil {
		log.Fatal(err)
	}

	consumer.StartWorker()

	fmt.Println("Server running at port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println(err)
	}
}
