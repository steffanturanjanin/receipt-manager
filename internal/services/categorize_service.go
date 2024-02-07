package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const (
	CATEGORIZE_PROMPT_SERBIAN = `Dat je spisak kategorija: %s.
Za svaki od sledećih artikala odredi kojoj od navedenih kategorija pripada: %s
Svaki artikal pripada samo jednoj kategoriji. Rezultat vratiti u JSON formatu gde je ključ ime artikla,
a vrednost ime kategorije kojoj artikal pripada. Zadržati imena artikala i kategorija u originalnom obliku.`
)

var (
	ErrUnprocessableCategorization = errors.New("unprocessable categorization")
	ErrCategoryNotFound            = errors.New("category not found")
)

type (
	CategorizeService struct{}

	CategoryIn struct {
		Id   uint
		Name string
	}
	ReceiptItemIn struct {
		Name string
	}
	ReceiptItemOut struct {
		Name       string
		CategoryId uint
	}

	CategoryInList     []CategoryIn
	ReceiptItemInList  []ReceiptItemIn
	ReceiptItemOutList []ReceiptItemOut

	CategorizedReceiptItemsMap map[string]*uint
)

func (s CategorizeService) Categorize(rilin ReceiptItemInList, clin CategoryInList) (CategorizedReceiptItemsMap, error) {
	formattedPrompt := fmt.Sprintf(CATEGORIZE_PROMPT_SERBIAN, rilin.toString(), clin.toString())

	ctx := context.Background()
	client := openai.NewClient(os.Getenv("OPEN_AI_API_KEY"))
	response, err := client.CreateCompletion(ctx, openai.CompletionRequest{
		Model:       openai.GPT3TextDavinci003,
		MaxTokens:   2000,
		Temperature: 0.2,
		Prompt:      formattedPrompt,
	})

	if err != nil {
		return nil, err
	}

	categorizationMap := response.Choices[0].Text

	cat := processCategorizationMap(categorizationMap, clin)

	return cat, nil
}

func processCategorizationMap(categorizationMap string, clin CategoryInList) CategorizedReceiptItemsMap {
	cat := make(CategorizedReceiptItemsMap, 0)

	rg := regexp.MustCompile(`(\r\n?|\n){2,}`)
	normalizedMultilines := rg.ReplaceAllString(categorizationMap, "$1")
	//rows := strings.Split(normalizedMultilines, "\n")

	var receiptItemNameCategoryNameMap map[string]string
	json.Unmarshal([]byte(normalizedMultilines), &receiptItemNameCategoryNameMap)

	for receiptItemName, categoryName := range receiptItemNameCategoryNameMap {
		var categoryId *uint
		if id, err := clin.getIdForName(categoryName); err == nil {
			categoryId = id
		}

		cat[receiptItemName] = categoryId
	}

	// for _, row := range rows {
	// 	pairs := strings.Split(row, ":")
	// 	if len(pairs) != 2 {
	// 		continue
	// 	}
	// 	receiptItem := strings.TrimSpace(pairs[0])
	// 	category := strings.TrimSpace(pairs[1])

	// 	categoryId, err := clin.getIdForName(category)
	// 	if err != nil {
	// 		continue
	// 	}

	// 	cat[receiptItem] = int(*categoryId)
	// }

	return cat
}

func (cin CategoryInList) getIdForName(categoryName string) (*uint, error) {
	for _, c := range cin {
		if c.Name == categoryName {
			cid := new(uint)
			*cid = c.Id

			return cid, nil
		}
	}

	return nil, ErrCategoryNotFound
}

func (clin CategoryInList) toString() string {
	cnames := make([]string, 0)
	for _, clinItem := range clin {
		cnames = append(cnames, clinItem.Name)
	}

	return strings.Join(cnames, ", ")
}

func (rilin ReceiptItemInList) toString() string {
	rnames := make([]string, 0)
	for _, rlinItem := range rilin {
		rnames = append(rnames, rlinItem.Name)
	}

	return strings.Join(rnames, ", ")
}
