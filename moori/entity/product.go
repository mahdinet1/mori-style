package entity

import "gorm.io/gorm"

type ProductID = uint64
type CurrentPrice = float32
type OffPercent = uint
type Color int

type Product struct {
	gorm.Model   `json:"-"`
	ID           ProductID    `json:"id" gorm:"primaryKey;not null"`
	Name         string       `json:"name" gorm:"not null;size:255;"`
	Description  string       `json:"description" gorm:"type:text;not null"`
	Material     string       `json:"material" gorm:"size:255;not null"`
	ShopName     string       `json:"shop_name" gorm:"size:255;not null"`
	Link         string       `json:"link" gorm:"type:text;not null"`
	Region       string       `json:"region" gorm:"size:255;not null"`
	CategoryName string       `json:"category_name" gorm:"size:255;"`
	OffPercent   OffPercent   `json:"off_percent"`
	CurrentPrice CurrentPrice `json:"current_price"`
	Images       []Image      `json:"images" gorm:"foreignKey:ProductId;constraint:OnDelete:CASCADE;"`
}

type Image struct {
	gorm.Model `json:"-"`
	ID         uint64    `json:"id" gorm:"primaryKey;not null"`
	Address    string    `json:"address" gorm:"not null;size:255;unique"`
	ProductId  ProductID `json:"-" gorm:"not null;index"`
}
