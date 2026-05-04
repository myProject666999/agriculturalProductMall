package model

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:50;not null"`
	ParentID  uint           `json:"parent_id" gorm:"default:0"`
	Icon      string         `json:"icon" gorm:"size:255"`
	Sort      int            `json:"sort" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Product struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Name           string         `json:"name" gorm:"size:200;not null"`
	CategoryID     uint           `json:"category_id" gorm:"index;not null"`
	Category       Category       `json:"category" gorm:"foreignKey:CategoryID"`
	MerchantID     uint           `json:"merchant_id" gorm:"index;not null"`
	Merchant       Merchant       `json:"merchant" gorm:"foreignKey:MerchantID"`
	Description    string         `json:"description" gorm:"type:text"`
	MainImage      string         `json:"main_image" gorm:"size:255"`
	SubImages      string         `json:"sub_images" gorm:"type:text"`
	Price          float64        `json:"price" gorm:"type:decimal(10,2);not null"`
	OriginalPrice  float64        `json:"original_price" gorm:"type:decimal(10,2)"`
	Stock          int            `json:"stock" gorm:"default:0"`
	Sales          int            `json:"sales" gorm:"default:0"`
	Unit           string         `json:"unit" gorm:"size:20"`
	Specifications string         `json:"specifications" gorm:"type:text"`
	IsHot          int            `json:"is_hot" gorm:"default:0"`
	IsNew          int            `json:"is_new" gorm:"default:0"`
	IsRecommend    int            `json:"is_recommend" gorm:"default:0"`
	Status         int            `json:"status" gorm:"default:1"`
	Rating         float64        `json:"rating" gorm:"type:decimal(3,2);default:5.00"`
	ReviewCount    int            `json:"review_count" gorm:"default:0"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

type Banner struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"size:200;not null"`
	Image     string         `json:"image" gorm:"size:255;not null"`
	LinkType  string         `json:"link_type" gorm:"size:20;default:'none'"`
	LinkValue string         `json:"link_value" gorm:"size:255"`
	Sort      int            `json:"sort" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Inventory struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	ProductID     uint           `json:"product_id" gorm:"uniqueIndex;not null"`
	Product       Product        `json:"product" gorm:"foreignKey:ProductID"`
	Quantity      int            `json:"quantity" gorm:"default:0"`
	MinStock      int            `json:"min_stock" gorm:"default:10"`
	MaxStock      int            `json:"max_stock" gorm:"default:1000"`
	LastUpdatedAt time.Time      `json:"last_updated_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type InventoryRecord struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ProductID   uint      `json:"product_id" gorm:"index;not null"`
	Product     Product   `json:"product" gorm:"foreignKey:ProductID"`
	Type        string    `json:"type" gorm:"size:20;not null"`
	Quantity    int       `json:"quantity" gorm:"not null"`
	BeforeQty   int       `json:"before_qty" gorm:"not null"`
	AfterQty    int       `json:"after_qty" gorm:"not null"`
	OperatorID  uint      `json:"operator_id" gorm:"not null"`
	Operator    User      `json:"operator" gorm:"foreignKey:OperatorID"`
	Remark      string    `json:"remark" gorm:"size:255"`
	CreatedAt   time.Time `json:"created_at"`
}
