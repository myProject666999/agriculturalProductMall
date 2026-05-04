package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	Username        string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password        string         `json:"-" gorm:"size:255;not null"`
	Email           string         `json:"email" gorm:"uniqueIndex;size:100;not null"`
	Phone           string         `json:"phone" gorm:"size:20"`
	Avatar          string         `json:"avatar" gorm:"size:255"`
	Nickname        string         `json:"nickname" gorm:"size:50"`
	Gender          int            `json:"gender" gorm:"default:0"`
	Role            string         `json:"role" gorm:"size:20;default:'user'"`
	Status          int            `json:"status" gorm:"default:1"`
	LastLoginAt     *time.Time     `json:"last_login_at"`
	LastLoginIP     string         `json:"last_login_ip" gorm:"size:50"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type Address struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"index;not null"`
	Name      string         `json:"name" gorm:"size:50;not null"`
	Phone     string         `json:"phone" gorm:"size:20;not null"`
	Province  string         `json:"province" gorm:"size:50"`
	City      string         `json:"city" gorm:"size:50"`
	District  string         `json:"district" gorm:"size:50"`
	Detail    string         `json:"detail" gorm:"size:255;not null"`
	IsDefault int            `json:"is_default" gorm:"default:0"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Merchant struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	UserID          uint           `json:"user_id" gorm:"uniqueIndex;not null"`
	StoreName       string         `json:"store_name" gorm:"size:100;not null"`
	StoreLogo       string         `json:"store_logo" gorm:"size:255"`
	StoreDesc       string         `json:"store_desc" gorm:"type:text"`
	ContactName     string         `json:"contact_name" gorm:"size:50"`
	ContactPhone    string         `json:"contact_phone" gorm:"size:20"`
	ContactEmail    string         `json:"contact_email" gorm:"size:100"`
	Address         string         `json:"address" gorm:"size:255"`
	BusinessLicense string         `json:"business_license" gorm:"size:255"`
	Status          int            `json:"status" gorm:"default:1"`
	Rating          float64        `json:"rating" gorm:"type:decimal(3,2);default:5.00"`
	TotalSales      int64          `json:"total_sales" gorm:"default:0"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type Favorite struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"uniqueIndex:idx_user_product;not null"`
	ProductID uint           `json:"product_id" gorm:"uniqueIndex:idx_user_product;not null"`
	Product   Product        `json:"product" gorm:"foreignKey:ProductID"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
