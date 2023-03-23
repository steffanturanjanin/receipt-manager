package models

type Category struct {
	ID          uint        `gorm:"primaryKey; autoIncrement" json:"id"`
	Name        string      `gorm:"unique; not null" json:"name"`
	ReceiptItem ReceiptItem `gorm:"foreignKey:ReceiptID;references:ID;constraint:OnDelete:SET NULL" json:"receipt_item"`
}
