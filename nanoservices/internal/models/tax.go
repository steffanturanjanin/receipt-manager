package models

// type Tax struct {
// 	TaxIdentifier int  `gorm:"primaryKey;uniqueIndex:idx_tax_identifier_receipt_id" json:"tax_identifier"`
// 	ReceiptID     uint `gorm:"primaryKey;uniqueIndex:idx_tax_identifier_receipt_id" json:"receipt_id"`
// }

type Tax struct {
	ID         uint   `gorm:"primaryKey; autoIncrement" json:"id"`
	Identifier string `gorm:"not null" json:"identifier"`
	Name       string `gorm:"not null" json:"name"`
	Rate       int    `gorm:"not null" json:"rate"`
}

// func (tax Tax) NewTaxDTO() *dto.Tax {
// 	return dto.TaxIdentifier(tax.Identifier).Tax()
// }
