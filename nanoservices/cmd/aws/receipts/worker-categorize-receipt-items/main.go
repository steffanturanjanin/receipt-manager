package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	db "github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
)

const (
	CATEGORIZE_RECEIPT_ITEMS_PROMPT_SERBIAN = `Dat je sledeći spisak kategorija: %s.
Za svaki od sledećih artikala odrediti kojoj od navedenih kategorija pripada: %s.
Svaki artikal pripada samo jednoj kategoriji. Rezultat vratiti u formatu JSON objekta, gde je ključ ime artikla,
a vrednost ime kategorije kojoj artikal pripada. Zadržati imena artikala i kategorija u originalnom obliku.
Ne dodavati nove kategorije i artikle.`
)

func GenerateCategorizationPrompt(categories []string, items []string) string {
	categoriesParam := strings.Join(categories, ", ")
	itemsParam := strings.Join(items, ", ")

	return fmt.Sprintf(CATEGORIZE_RECEIPT_ITEMS_PROMPT_SERBIAN, categoriesParam, itemsParam)
}

var (
	// Database instance
	DB *gorm.DB

	// OpenAI client
	OpenAiClient *openai.Client
)

func init() {
	// Initialize database
	if err := db.InitializeDB(); err != nil {
		os.Exit(1)
	} else {
		DB = db.Instance
	}

	// Initialize OpenAI Client
	OpenAiClient = openai.NewClient(os.Getenv("OpenAiApiKey"))
}

// Process each individual message from the current SQS batch
func processMessage(ctx context.Context, message events.SQSMessage) error {
	// Get Receipt ID from message body
	receiptId, err := strconv.Atoi(message.Body)
	if err != nil {
		return err
	}

	// Retrieve Receipt Items related to Receipt ID
	var dbReceiptItems []models.ReceiptItem
	if dbErr := DB.Where("receipt_id = ?", receiptId).Find(&dbReceiptItems).Error; dbErr != nil {
		return dbErr
	}

	// Retrieve Categories
	var dbCategories []models.Category
	if dbErr := DB.Find(&dbCategories).Error; dbErr != nil {
		return dbErr
	}

	// Receipt receiptItems names
	var receiptItems []string
	for _, receiptItem := range dbReceiptItems {
		receiptItem := receiptItem.Name
		receiptItems = append(receiptItems, receiptItem)
	}

	// Categories names
	var categories []string
	for _, dbCategory := range dbCategories {
		category := dbCategory.Name
		categories = append(categories, category)
	}

	messages := []openai.ChatCompletionMessage{{
		Role:    openai.ChatMessageRoleUser,
		Content: GenerateCategorizationPrompt(categories, receiptItems),
	}}

	response, err := OpenAiClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Messages:    messages,
		Model:       openai.GPT3Dot5Turbo,
		MaxTokens:   2400,
		Temperature: 0.1,
	})

	if err != nil {
		return err
	}

	// Get output
	output := response.Choices[0].Message.Content
	// Normalize output
	exp := regexp.MustCompile(`(\r\n?|\n){2,}`)
	normalizedOutput := exp.ReplaceAllString(output, "$1")

	var itemsCategoriesMap map[string]string
	if err := json.Unmarshal([]byte(normalizedOutput), &itemsCategoriesMap); err != nil {
		return err
	}

	// Update Receipt Items with Categories
	for receiptItemName, categoryName := range itemsCategoriesMap {
		var foundReceiptItem *models.ReceiptItem
		var foundCategory *models.Category

		// Find DB Receipt Item
		for _, dbReceiptItem := range dbReceiptItems {
			if dbReceiptItem.Name == receiptItemName {
				foundReceiptItem = &dbReceiptItem
				break
			}
		}

		// Find DB Category
		for _, dbCategory := range dbCategories {
			if dbCategory.Name == categoryName {
				foundCategory = &dbCategory
				break
			}
		}

		// Update receipt item category
		if foundReceiptItem != nil && foundCategory != nil {
			foundReceiptItem.CategoryID = &foundCategory.ID

			if err := DB.Save(&foundReceiptItem); err != nil {
				break
			}
		}
	}

	return nil
}

func handler(ctx context.Context, event events.SQSEvent) error {
	for _, message := range event.Records {
		if err := processMessage(ctx, message); err != nil {
			log.Printf("Error while trying to categorize receipt items: %+v\n", err)
			continue
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
