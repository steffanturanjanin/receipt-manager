package models

type ReceiptItem struct {
	ID           uint    `gorm:"primaryKey; autoIncrement" json:"id"`
	Name         string  `gorm:"size:255,not null" json:"name"`
	Unit         string  `gorm:"not null" json:"unit"`
	Quantity     float64 `gorm:"not null" json:"quantity"`
	SingleAmount int     `gorm:"not null; default:0" json:"single_amount"`
	TotalAmount  int     `gorm:"not null; default:0" json:"total_amount"`
	Tax          int     `gorm:"nullable" json:"tax"`
	ReceiptID    uint    `json:"receipt_id"`
}
