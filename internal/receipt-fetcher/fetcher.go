package receipt_fetcher

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func Get(url string) (*Receipt, error) {
	pageContent, err := getUrlContent(url)
	if err != nil {
		return nil, err
	}

	receiptContent, err := getReceiptContent(pageContent)
	if err != nil {
		return nil, err
	}

	receipt, err := parseReceipt(receiptContent)
	if err != nil {
		return nil, err
	}

	return receipt, nil
}

func getUrlContent(url string) (io.ReadCloser, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

func getReceiptContent(rc io.ReadCloser) (string, error) {
	doc, err := goquery.NewDocumentFromReader(rc)
	if err != nil {
		return "", err
	}

	defer rc.Close()

	html, err := doc.Find("pre").Html()
	if err != nil {
		return "", err
	}

	return html, nil
}

func parseReceipt(receiptContent string) (*Receipt, error) {
	receiptContent = receiptContent[strings.Index(receiptContent, "\n")+1:]
	receiptContent = receiptContent[:strings.LastIndex(receiptContent, "\n")]

	startDelimiter := "========================================"
	endDelimiter := "----------------------------------------"

	recSections := strings.Split(receiptContent, startDelimiter)
	for index := range recSections {
		recSections[index] = normalizeMultilineStrings(recSections[index])
	}

	header, body, taxData, fiscalizationData, qrCodeData := recSections[0], recSections[1], recSections[2], recSections[3], recSections[4]

	meta := extractKeyValuePairs(header)
	store := extractStoreData(header)

	taxLines := strings.Split(taxData, "\n")
	taxes := []TaxItem{}
	taxTypes := make(map[string]Tax)
	totalTaxAmount := RsdAmount{}

	for i := 1; i < len(taxLines); i++ {
		if taxLines[i] == endDelimiter {
			totals := extractKeyValuePairs(taxLines[i+1])
			totalTaxAmount, _ = NewRsdAmountFromString(totals["Укупан износ пореза"])

			break
		}
		parts := strings.Fields(taxLines[i])

		identifier := parts[0]
		extractRate := strings.Split(parts[2], ",")[0]
		rate, _ := strconv.Atoi(extractRate)

		tax := Tax{
			Name:       parts[1],
			Identifier: identifier,
			Rate:       rate,
		}

		taxTypes[identifier] = tax
		taxItemRsdAmount, _ := NewRsdAmountFromString(parts[3])

		taxItem := TaxItem{
			Tax:    tax,
			Amount: taxItemRsdAmount,
		}

		taxes = append(taxes, taxItem)
	}

	bodySplit := strings.Split(body, endDelimiter)
	itemsData, paymentData := normalizeMultilineStrings(bodySplit[0]), normalizeMultilineStrings(bodySplit[1])

	paymentPairs := extractKeyValuePairs(paymentData)
	paymentSummary := make(map[string]RsdAmount)
	for index, payment := range paymentPairs {
		paymentSummary[index], _ = NewRsdAmountFromString(payment)
	}

	receiptItems := extractItemData(itemsData, taxTypes)

	fiscalization := extractKeyValuePairs(fiscalizationData)
	dsLayout := "02.01.2006. 15:04:05"
	fiscDateTime, _ := time.Parse(dsLayout, fiscalization["ПФР време"])

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(qrCodeData))
	qrCode, _ := doc.Find("img").Attr("src")
	qrCodeString := qrCode[len("data:image/gif;base64,"):]

	receipt := &Receipt{
		Store:               store,
		Number:              fiscalization["ПФР број рачуна"],
		Counter:             fiscalization["Бројач рачуна"],
		Meta:                meta,
		Items:               receiptItems,
		Taxes:               taxes,
		PaymentSummary:      paymentSummary,
		TotalPurchaseAmount: paymentSummary["Укупан износ"],
		TotalTaxAmount:      totalTaxAmount,
		Date:                fiscDateTime,
		QrCod:               qrCodeString,
	}

	return receipt, nil
}

func normalizeMultilineStrings(str string) string {
	out := []rune(str)
	if out[0] == '\n' {
		out = out[1:]
	}

	if last := len(out) - 1; out[last] == '\n' {
		out = out[:last]
	}

	return string(out)
}

func extractKeyValuePairs(str string) map[string]string {
	pairs := make(map[string]string)

	lines := strings.Split(str, "\n")
	var prevLineIndex string = ""

	for _, line := range lines {
		pair := strings.SplitN(line, ":", 2)

		if len(pair) == 1 {
			if prevLineIndex == "" {
				continue
			}

			if _, err := strconv.Atoi(strings.TrimSpace(pair[0])); err != nil {
				continue
			}

			repairedPairValue := pairs[prevLineIndex] + strings.TrimSpace(pair[0])
			pairs[prevLineIndex] = repairedPairValue
			prevLineIndex = ""
		} else {
			pairs[pair[0]] = strings.TrimSpace(pair[1])
			prevLineIndex = pair[0]
		}
	}

	return pairs
}

func extractStoreData(content string) Store {
	lines := strings.Split(content, "\n")
	locationData := strings.Split(lines[2], "-")
	locationId, locationName := locationData[0], locationData[1]

	return Store{
		Name:         strings.TrimSpace(lines[1]),
		Tin:          strings.TrimSpace(lines[0]),
		LocationId:   strings.TrimSpace(locationId),
		LocationName: strings.TrimSpace(locationName),
		Address:      strings.TrimSpace(lines[3]),
		City:         strings.TrimSpace(lines[4]),
	}
}

func extractItemData(itemData string, taxes map[string]Tax) []ReceiptItem {
	lines := strings.Split(itemData, "\n")

	optionalPrefixItemCode := `(?:[0-9]{3,}(?: |,|\-))?`
	optionalSuffixItemCode := `(?:(?: |,|\-)[0-9]{3,})?`
	itemName := `(?P<name>.*)`
	unit := `(?:\/|\/ | )(?P<unit>kom|kg|l|lit|kut|m|pce|fl)`
	taxIdentifier := `\((?P<taxIdentifier>е|ђ|a)\)`

	items := []ReceiptItem{}
	var itemLine string

	for i := 1; i < len(lines); i++ {
		amountReg := regexp.MustCompile(`([0-9,.]+)\s+([0-9,.]+)\s+([0-9,.]+)`)
		amountMatches := amountReg.FindStringSubmatch(lines[i])

		if len(amountMatches) == 0 {
			itemLine += lines[i]
			continue
		}

		lineVariants := []string{
			fmt.Sprintf("%s%s%s %s", optionalPrefixItemCode, itemName, unit, taxIdentifier),
			fmt.Sprintf("%s%s%s %s", itemName, optionalSuffixItemCode, unit, taxIdentifier),
			fmt.Sprintf("%s%s %s", optionalPrefixItemCode, itemName, taxIdentifier),
			fmt.Sprintf("%s%s %s", itemName, optionalSuffixItemCode, taxIdentifier),
		}

		lineMatches := make(map[string]string)
		for _, variant := range lineVariants {
			matchExp := regexp.MustCompile("(?i)" + variant)
			matches := matchExp.FindStringSubmatch(itemLine)

			if len(matches) > 0 {
				for i, name := range matchExp.SubexpNames() {
					if i != 0 && name != "" {
						lineMatches[name] = matches[i]
					}
				}

				break
			}
		}

		name := strings.TrimSpace(lineMatches["name"])
		quantity, _ := strconv.ParseFloat(strings.Replace(amountMatches[2], ",", ".", -1), 64)
		unit := strings.TrimSpace(lineMatches["unit"])
		tax := taxes[lineMatches["taxIdentifier"]]
		singleAmount, _ := NewRsdAmountFromString(amountMatches[1])
		totalAmount, _ := NewRsdAmountFromString(amountMatches[3])

		receiptItem := ReceiptItem{}
		receiptItem.Name = name
		receiptItem.Quantity = quantity
		receiptItem.Unit = unit
		receiptItem.Tax = tax
		receiptItem.SingleAmount = singleAmount
		receiptItem.TotalAmount = totalAmount

		itemLine = ""
		items = append(items, receiptItem)
	}

	return items
}
