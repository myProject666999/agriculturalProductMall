package model

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"index;not null"`
	ProductID uint          `json:"product_id" gorm:"index;not null"`
	Product   Product        `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int            `json:"quantity" gorm:"default:1"`
	Selected  int            `json:"selected" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Order struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	OrderNo        string         `json:"order_no" gorm:"uniqueIndex;size:50;not null"`
	UserID         uint           `json:"user_id" gorm:"index;not null"`
	User           User           `json:"user" gorm:"foreignKey:UserID"`
	MerchantID     uint           `json:"merchant_id" gorm:"index;not null"`
	Merchant       Merchant       `json:"merchant" gorm:"foreignKey:MerchantID"`
	Status         string         `json:"status" gorm:"size:20;default:'pending'"`
	TotalAmount    float64        `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	DiscountAmount float64        `json:"discount_amount" gorm:"type:decimal(10,2);default:0"`
	PayAmount      float64        `json:"pay_amount" gorm:"type:decimal(10,2);not null"`
	ShippingFee    float64        `json:"shipping_fee" gorm:"type:decimal(10,2);default:0"`
	ReceiverName   string         `json:"receiver_name" gorm:"size:50;not null"`
	ReceiverPhone  string         `json:"receiver_phone" gorm:"size:20;not null"`
	ReceiverAddress string        `json:"receiver_address" gorm:"size:255;not null"`
	PayMethod      string         `json:"pay_method" gorm:"size:20"`
	PayTime        *time.Time     `json:"pay_time"`
	ShippingTime   *time.Time     `json:"shipping_time"`
	CompleteTime   *time.Time     `json:"complete_time"`
	RefundStatus   string         `json:"refund_status" gorm:"size:20;default:'none'"`
	RefundReason   string         `json:"refund_reason" gorm:"type:text"`
	RefundAmount   float64        `json:"refund_amount" gorm:"type:decimal(10,2);default:0"`
	Remark         string         `json:"remark" gorm:"size:255"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

type OrderItem struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	OrderID    uint      `json:"order_id" gorm:"index;not null"`
	ProductID  uint      `json:"product_id" gorm:"index;not null"`
	Product    Product   `json:"product" gorm:"foreignKey:ProductID"`
	ProductName string    `json:"product_name" gorm:"size:200;not null"`
	ProductImage string   `json:"product_image" gorm:"size:255"`
	Price      float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	Quantity   int       `json:"quantity" gorm:"not null"`
	SubTotal   float64   `json:"sub_total" gorm:"type:decimal(10,2);not null"`
	CreatedAt  time.Time `json:"created_at"`
}

type Review struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"index;not null"`
	User       User           `json:"user" gorm:"foreignKey:UserID"`
	ProductID  uint           `json:"product_id" gorm:"index;not null"`
	Product    Product        `json:"product" gorm:"foreignKey:ProductID"`
	OrderID    uint           `json:"order_id" gorm:"index;not null"`
	MerchantID uint           `json:"merchant_id" gorm:"index;not null"`
	Rating     int            `json:"rating" gorm:"default:5"`
	Content    string         `json:"content" gorm:"type:text;not null"`
	Images     string         `json:"images" gorm:"type:text"`
	IsAnonymous int           `json:"is_anonymous" gorm:"default:0"`
	Status     int            `json:"status" gorm:"default:1"`
	Reply      string         `json:"reply" gorm:"type:text"`
	ReplyTime  *time.Time     `json:"reply_time"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

type Logistics struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	OrderID      uint           `json:"order_id" gorm:"uniqueIndex;not null"`
	Order        Order          `json:"order" gorm:"foreignKey:OrderID"`
	LogisticsNo  string         `json:"logistics_no" gorm:"size:50;not null"`
	Company      string         `json:"company" gorm:"size:50;not null"`
	CurrentStatus string        `json:"current_status" gorm:"size:100"`
	Trajectory   string         `json:"trajectory" gorm:"type:text"`
	ShipTime     *time.Time     `json:"ship_time"`
	ReceivedTime *time.Time     `json:"received_time"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
