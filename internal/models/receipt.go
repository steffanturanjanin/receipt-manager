package models

import (
	"time"
)

type Receipt struct {
	ID                          uint      `gorm:"primaryKey; autoIncrement"`
	PfrNumber                   string    `gorm:"unique; not null"`
	Counter                     string    `gorm:"unique; not null"`
	TotalPurchaseAmountValue    uint      `gorm:"not null; default:0"`
	TotalPurchaseAmountFraction uint      `gorm:"not null; default:0"`
	TotalTaxAmountValue         uint      `gorm:"not null; default:0"`
	TotalTaxMountFraction       uint      `gorm:"not null; default:0"`
	Date                        time.Time `gorm:"not null"`
	QrCode                      string    `gorm:"not null"`
	CreatedAt                   time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
}
