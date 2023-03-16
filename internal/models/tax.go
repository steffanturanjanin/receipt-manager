package models

import "github.com/steffanturanjanin/receipt-manager/internal/dto"

type Tax struct {
	TaxIdentifier int  `gorm:"primaryKey;uniqueIndex:idx_tax_identifier_receipt_id" json:"tax_identifier"`
	ReceiptID     uint `gorm:"primaryKey;uniqueIndex:idx_tax_identifier_receipt_id" json:"receipt_id"`
}

func (tax Tax) NewTaxDTO() *dto.Tax {
	return dto.TaxIdentifier(tax.TaxIdentifier).Tax()
}
