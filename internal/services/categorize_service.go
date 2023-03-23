package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const (
	CATEGORIZE_PROMPT_SERBIAN = `Ovo je lista kategorija: %s.
Za svaki od sledećih artikala odredi kojoj kategoriji pripada: %s
Svaki artikal pripada samo jednoj kategoriji.
Rezultat vrati u formi uređenih parova, pri čemu je ključ ime artikla, a vrednost ime kategorije.
Pri kreiranju uređenih parova ne menjati originalne vrednosti imena artikla i kategorija.`
)

var (
	ErrUnprocessableCategorization = errors.New("unprocessable categorization")
	ErrCategoryNotFound            = errors.New("category not found")
)

type (
	CategorizeService struct{}

	categorizeCategoryIn struct {
		id   uint
		name string
	}
	categorizeReceiptItemIn struct {
		name string
	}
	categorizeReceiptItemOut struct {
		name       string
		categoryId uint
	}

	CategorizeCategoryIn     []categorizeCategoryIn
	CategorizeReceiptItemIn  []categorizeReceiptItemIn
	CategorizeReceiptItemOut []categorizeReceiptItemOut
)

func (s CategorizeService) Categorize(rilin CategorizeReceiptItemIn, clin CategorizeCategoryIn) (CategorizeReceiptItemOut, error) {
	formatedPrompt := fmt.Sprintf(CATEGORIZE_PROMPT_SERBIAN, rilin.toString(), clin.toString())

	ctx := context.Background()
	client := openai.NewClient(os.Getenv("OPENAI_KEY"))
	response, err := client.CreateCompletion(ctx, openai.CompletionRequest{
		Model:       openai.GPT3TextDavinci003,
		MaxTokens:   2048,
		Temperature: 0.2,
		Prompt:      formatedPrompt,
	})

	if err != nil {
		return nil, ErrUnprocessableCategorization
	}

	categorizationMap := response.Choices[0].Text
	crlo := processCategorizationMap(categorizationMap, clin)

	return crlo, nil
}

func processCategorizationMap(categorizationMap string, clin CategorizeCategoryIn) CategorizeReceiptItemOut {
	crlo := make(CategorizeReceiptItemOut, 0)

	rg := regexp.MustCompile(`(\r\n?|\n){2,}`)
	normalizedMultilines := rg.ReplaceAllString(categorizationMap, "$1")
	rows := strings.Split(normalizedMultilines, "\n")

	for _, row := range rows {
		pairs := strings.Split(row, ":")
		receiptItem := strings.TrimSpace(pairs[0])
		category := strings.TrimSpace(pairs[1])

		categoryId, err := clin.getIdForName(category)
		if err != nil {
			continue
		}

		crlo = append(crlo, categorizeReceiptItemOut{
			name:       receiptItem,
			categoryId: *categoryId,
		})
	}

	return crlo
}

func (cin CategorizeCategoryIn) getIdForName(categoryName string) (*uint, error) {
	for _, c := range cin {
		if c.name == categoryName {
			cid := new(uint)
			*cid = c.id

			return cid, nil
		}
	}

	return nil, ErrCategoryNotFound
}

func (clin CategorizeCategoryIn) toString() string {
	cnames := make([]string, 0)
	for _, clinItem := range clin {
		cnames = append(cnames, clinItem.name)
	}

	return strings.Join(cnames, ", ")
}

func (rilin CategorizeReceiptItemIn) toString() string {
	rnames := make([]string, 0)
	for _, rlinItem := range rilin {
		rnames = append(rnames, rlinItem.name)
	}

	return strings.Join(rnames, ", ")
}
