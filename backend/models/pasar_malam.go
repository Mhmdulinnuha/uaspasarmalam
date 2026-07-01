package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string  `gorm:"type:varchar(255);not null" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	Price       float64 `gorm:"type:decimal(15,2);not null" json:"price"`
	Stock       int     `gorm:"default:0" json:"stock"`
	Category    string  `gorm:"type:varchar(100)" json:"category"`
	ImageUrl    string  `gorm:"type:text" json:"image_url"`
	IsActive    bool    `gorm:"default:true" json:"is_active"`
}

type Cart struct {
	gorm.Model
	UserID uint       `gorm:"uniqueIndex;not null" json:"user_id"`
	User   User       `gorm:"foreignKey:UserID" json:"-"`
	Items  []CartItem `gorm:"foreignKey:CartID" json:"items"`
}

type CartItem struct {
	gorm.Model
	CartID    uint    `gorm:"not null;index" json:"cart_id"`
	ProductID uint    `gorm:"not null;index" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  int     `gorm:"default:1" json:"quantity"`
	Subtotal  float64 `gorm:"type:decimal(15,2)" json:"subtotal"`
}

type Order struct {
	gorm.Model
	UserID          uint        `gorm:"not null;index" json:"user_id"`
	User            User        `gorm:"foreignKey:UserID" json:"-"`
	TotalAmount     float64     `gorm:"type:decimal(15,2);not null" json:"total_amount"`
	Status          string      `gorm:"type:varchar(50);default:'pending'" json:"status"`
	ShippingAddress string      `gorm:"type:text;not null" json:"shipping_address"`
	Notes           string      `gorm:"type:text" json:"notes"`
	PaymentMethod   string      `gorm:"type:varchar(50);not null" json:"payment_method"`
	VaNumber        string      `gorm:"type:varchar(255)" json:"va_number"`
	GopayDeeplink   string      `gorm:"type:text" json:"gopay_deeplink"`
	Items           []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
}

type OrderItem struct {
	gorm.Model
	OrderID     uint    `gorm:"not null;index" json:"order_id"`
	ProductID   uint    `gorm:"not null" json:"product_id"`
	ProductName string  `gorm:"type:varchar(255);not null" json:"product_name"`
	Price       float64 `gorm:"type:decimal(15,2);not null" json:"price"`
	Quantity    int     `gorm:"not null" json:"quantity"`
	Subtotal    float64 `gorm:"type:decimal(15,2);not null" json:"subtotal"`
}
