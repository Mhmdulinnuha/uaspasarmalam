package database

import (
	"fmt"
	"log"

	"emoney-603dc/config"
	"emoney-603dc/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitMySQL(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to MySQL:", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.OTP{},
		&models.Account{},
		&models.Transaction{},
		&models.Product{},
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Tambahkan produk sampel jika belum ada
	seedSampleProducts(db)

	log.Println("MySQL connected and migrated")
	return db
}

func seedSampleProducts(db *gorm.DB) {
	var count int64
	db.Model(&models.Product{}).Count(&count)
	if count > 0 {
		log.Println("Sample products already exist, skipping seeding")
		return
	}

	sampleProducts := []models.Product{
		{
			Name:        "Nasi Goreng Spesial",
			Description: "Nasi goreng dengan telur, ayam suwir, sayuran, dan bumbu khas",
			Price:       25000,
			Stock:       50,
			Category:    "Makanan",
			ImageUrl:    "https://images.unsplash.com/photo-1512058564366-18510be2db19?w=400&h=300&fit=crop",
			IsActive:    true,
		},
		{
			Name:        "Mie Goreng Ayam",
			Description: "Mie goreng dengan potongan ayam, sawi, dan acar timun",
			Price:       22000,
			Stock:       45,
			Category:    "Makanan",
			ImageUrl:    "https://images.unsplash.com/photo-1569718212165-3a8278d5f624?w=400&h=300&fit=crop",
			IsActive:    true,
		},
		{
			Name:        "Es Teh Manis",
			Description: "Es teh manis segar dengan daun teh pilihan",
			Price:       5000,
			Stock:       100,
			Category:    "Minuman",
			ImageUrl:    "https://images.unsplash.com/photo-1461426553385-6f376e4822ed?w=400&h=300&fit=crop",
			IsActive:    true,
		},
		{
			Name:        "Es Jeruk Nipis",
			Description: "Es jeruk nipis segar dengan gula merah",
			Price:       7000,
			Stock:       80,
			Category:    "Minuman",
			ImageUrl:    "https://images.unsplash.com/photo-1544510808-91bcbee1df55?w=400&h=300&fit=crop",
			IsActive:    true,
		},
		{
			Name:        "Bakso Goreng",
			Description: "Bakso goreng renyah dengan saus sambal",
			Price:       15000,
			Stock:       60,
			Category:    "Snack",
			ImageUrl:    "https://images.unsplash.com/photo-1565299624946-b28f40a0ae38?w=400&h=300&fit=crop",
			IsActive:    true,
		},
		{
			Name:        "Klepon",
			Description: "Klepon manis dengan kelapa parut",
			Price:       10000,
			Stock:       40,
			Category:    "Snack",
			ImageUrl:    "https://images.unsplash.com/photo-1571091718767-18b5b1457add?w=400&h=300&fit=crop",
			IsActive:    true,
		},
	}

	if err := db.Create(&sampleProducts).Error; err != nil {
		log.Println("Failed to seed sample products:", err)
	} else {
		log.Println("Successfully seeded sample products")
	}
}
