package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sashabaranov/go-openai"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/steffanturanjanin/receipt-manager/internal/controllers"
	"github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/repositories"
	"github.com/steffanturanjanin/receipt-manager/internal/services"
	"github.com/steffanturanjanin/receipt-manager/internal/validator"
)

var (
	Session    *session.Session
	SqsService *sqs.SQS
)

const (
	PROMPT_SERBIAN_2 = `Ovo je lista kategorija: %s.
Za svaki od sledećih artikala odredi kojoj kategoriji pripada: %s
Svaki artikal pripada samo jednoj kategoriji.
Rezultat vrati u formi uređenih parova, pri čemu je ključ ime artikla, a vrednost ime kategorije.
Pri kreiranju uređenih parova ne menjati originalne vrednosti imena artikla i kategorija.`
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
		// type url struct {
		// 	Url string `json:"url"`
		// }

		// u := url{}
		// controllers.ParseBody(&u, r)

		// queueName := queue.RECEIPT_URL_QUEUE
		// urlResult, _ := SqsService.GetQueueUrl(&sqs.GetQueueUrlInput{
		// 	QueueName: &queueName,
		// })

		// producer := queue.NewReceiptUrlProducer(SqsService, *urlResult.QueueUrl)
		// if err := producer.Produce(u.Url); err != nil {
		// 	w.Write([]byte("{'error': 'Error occured while writing url to queue.'}"))
		// 	return
		// }

		receiptItems := []string{
			"CARNEX STISNJENA SUNKA",
			"MH DUGOTRAJNO MLEKO 2,8%MM 1L",
			"TIS JASTREBACKI KACKAVALJ 49%MM",
			"ZLATIBORAC LJUTA DELIKATES KOBASICA 250G ZAS.A",
			"DLC SAMPINJONI REZANI 700G",
			"TUNJEVINA KOMADI REGINA DELMARE U ULJU",
			"BANANA PILOS 3KOM",
			"STOLICA KOZNA HERMAN MILLER",
			"JASTUK JYSK 30x20",
			"FARMERKE LEVIS 75",
			"KECAP POLIMARK 250ml",
			"NEXT CLASSIC POMORANDZA",
			"MUSKI DUKS BELI",
			"SIJALICA LED 2KOM",
			"RANAC ADDIDAS CRNI",
		}

		categories := []string{
			"Hrana",
			"Voće i povrće",
			"Delikates",
			"Meso",
			"Riba",
			"Mlečni proizvodi i jaja",
			"Hleb i pekarski proizvodi",
			"Namirnice za pripremu jela",
			"Slatkiši i grickalice",
			"Gotova jela",
			"Piće",
			"Alkoholna pića",
			"Sokovi",
			"Voda",
			"Čaj i kafa",
			"Kućna hemija",
			"Lična higijena i nega",
			"Tehnika",
			"Odeća i obuća",
			"Zdravlje",
			"Domaćinstvo",
		}

		categoryList := strings.Join(categories, ", ")
		receiptItemsList := strings.Join(receiptItems, ", ")

		formatedPrompt := fmt.Sprintf(PROMPT_SERBIAN_2, categoryList, receiptItemsList)

		//fmt.Println(formatedPrompt)

		ctx := context.Background()
		client := openai.NewClient("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")
		response, err := client.CreateCompletion(ctx, openai.CompletionRequest{
			Model:       openai.GPT3TextDavinci003,
			MaxTokens:   2048,
			Temperature: 0.2,
			Prompt:      formatedPrompt,
		})

		if err != nil {
			log.Fatalf("ERROR OPENAI: %+v\n", err)
		}

		fmt.Println(response.Choices[0].Text)

	}).Methods("POST")

	// queueName := queue.RECEIPT_URL_QUEUE
	// urlResult, _ := SqsService.GetQueueUrl(&sqs.GetQueueUrlInput{
	// 	QueueName: &queueName,
	// })

	// consumer, err := queue.NewReceiptUrlConsumer(*urlResult.QueueUrl)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// consumer.StartWorker()

	fmt.Println("Server running at port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println(err)
	}
}
