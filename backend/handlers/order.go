package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"emoney-603dc/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderHandler struct {
	db *gorm.DB
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{db: db}
}

type CheckoutRequest struct {
	ShippingAddress string `json:"shipping_address" binding:"required"`
	Notes           string `json:"notes"`
	PaymentMethod   string `json:"payment_method" binding:"required"`
}

func (h *OrderHandler) Checkout(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "shipping_address dan payment_method diperlukan",
		})
		return
	}

	var cart models.Cart
	if err := h.db.Where("user_id = ?", userID).
		Preload("Items.Product").
		First(&cart).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Keranjang kosong",
		})
		return
	}

	if len(cart.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Keranjang kosong",
		})
		return
	}

	var totalAmount float64
	var orderItems []models.OrderItem

	for _, item := range cart.Items {
		var product models.Product
		if err := h.db.First(&product, item.ProductID).Error; err != nil {
			continue
		}
		if product.Stock < item.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fmt.Sprintf("Stok produk %s tidak cukup", product.Name),
			})
			return
		}

		orderItems = append(orderItems, models.OrderItem{
			ProductID:   item.ProductID,
			ProductName: product.Name,
			Price:       product.Price,
			Quantity:    item.Quantity,
			Subtotal:    item.Subtotal,
		})

		totalAmount += item.Subtotal
	}

	var vaNumber string
	var gopayDeeplink string

	if req.PaymentMethod == "virtual_account" {
		rand.Seed(time.Now().UnixNano())
		vaNumber = fmt.Sprintf("880880%d", rand.Intn(900000000)+100000000)
	} else if req.PaymentMethod == "gopay" {
		gopayDeeplink = "https://simulator.sandbox.midtrans.com/gopay/landing_page"
	}

	var order models.Order
	err := h.db.Transaction(func(tx *gorm.DB) error {
		order = models.Order{
			UserID:          userID,
			TotalAmount:     totalAmount,
			Status:          "pending",
			ShippingAddress: req.ShippingAddress,
			Notes:           req.Notes,
			PaymentMethod:   req.PaymentMethod,
			VaNumber:        vaNumber,
			GopayDeeplink:   gopayDeeplink,
			Items:           orderItems,
		}

		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		for _, item := range cart.Items {
			if err := tx.Model(&models.Product{}).
				Where("id = ?", item.ProductID).
				UpdateColumn("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membuat pesanan",
		})
		return
	}

	if err := h.db.Preload("Items").First(&order, order.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil pesanan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pesanan berhasil dibuat",
		"data":    order,
	})
}

func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	userID := c.GetUint("user_id")

	var orders []models.Order

	query := h.db.Where("user_id = ?", userID).Preload("Items")

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			query = query.Limit(limit)
		}
	}

	if err := query.Order("created_at desc").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil pesanan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    orders,
	})
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	userID := c.GetUint("user_id")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID pesanan tidak valid",
		})
		return
	}

	var order models.Order
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).
		Preload("Items").
		First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Pesanan tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    order,
	})
}
