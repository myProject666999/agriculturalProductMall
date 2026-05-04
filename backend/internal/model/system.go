package model

import (
	"time"

	"gorm.io/gorm"
)

type News struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"size:200;not null"`
	Author      string         `json:"author" gorm:"size:50"`
	Category    string         `json:"category" gorm:"size:50"`
	Content     string         `json:"content" gorm:"type:longtext;not null"`
	Cover       string         `json:"cover" gorm:"size:255"`
	Description string        `json:"description" gorm:"size:500"`
	Views       int            `json:"views" gorm:"default:0"`
	Status      int            `json:"status" gorm:"default:1"`
	Sort        int            `json:"sort" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Announcement struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"size:200;not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	Type      string         `json:"type" gorm:"size:20;default:'general'"`
	Status    int            `json:"status" gorm:"default:1"`
	StartAt   *time.Time     `json:"start_at"`
	EndAt     *time.Time     `json:"end_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Menu struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ParentID  uint           `json:"parent_id" gorm:"default:0"`
	Name      string         `json:"name" gorm:"size:50;not null"`
	Path      string         `json:"path" gorm:"size:255"`
	Icon      string         `json:"icon" gorm:"size:100"`
	Component string         `json:"component" gorm:"size:255"`
	Permission string        `json:"permission" gorm:"size:100"`
	Type      int            `json:"type" gorm:"default:1"`
	Sort      int            `json:"sort" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Permission struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	MenuID    uint           `json:"menu_id" gorm:"index"`
	Name      string         `json:"name" gorm:"size:50;not null"`
	Code      string         `json:"code" gorm:"size:100;not null"`
	Type      string         `json:"type" gorm:"size:20;default:'button'"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Role struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:50;not null"`
	Code        string         `json:"code" gorm:"uniqueIndex;size:50;not null"`
	Description string         `json:"description" gorm:"size:255"`
	Status      int            `json:"status" gorm:"default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type RolePermission struct {
	ID           uint `json:"id" gorm:"primaryKey"`
	RoleID       uint `json:"role_id" gorm:"index:idx_role_permission,unique"`
	PermissionID uint `json:"permission_id" gorm:"index:idx_role_permission,unique"`
}

type UserRole struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	UserID uint `json:"user_id" gorm:"index:idx_user_role,unique"`
	RoleID uint `json:"role_id" gorm:"index:idx_user_role,unique"`
}

type VerificationCode struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"index;size:100;not null"`
	Code      string    `json:"code" gorm:"size:10;not null"`
	Type      string    `json:"type" gorm:"size:20;not null"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      int       `json:"used" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
}

type RecommendRecord struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"index;not null"`
	ProductID   uint      `json:"product_id" gorm:"index;not null"`
	Product     Product   `json:"product" gorm:"foreignKey:ProductID"`
	Score       float64   `json:"score" gorm:"type:decimal(5,2);not null"`
	RecommendType string  `json:"recommend_type" gorm:"size:20;not null"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserBehavior struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	ProductID uint      `json:"product_id" gorm:"index;not null"`
	Behavior  string    `json:"behavior" gorm:"size:20;not null"`
	Score     float64   `json:"score" gorm:"type:decimal(5,2);not null"`
	CreatedAt time.Time `json:"created_at"`
}
