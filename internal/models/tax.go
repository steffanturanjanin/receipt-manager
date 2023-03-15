package models

type Tax struct {
	TaxIdentifier int  `gorm:"primaryKey;uniqueIndex:idx_tax_identifier_receipt_id" json:"tax_identifier"`
	ReceiptID     uint `gorm:"primaryKey;uniqueIndex:idx_tax_identifier_receipt_id" json:"receipt_id"`
}
