package handlers

import (
	"net/http"
	"strconv"

	"emoney-603dc/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CartHandler struct {
	db *gorm.DB
}

func NewCartHandler(db *gorm.DB) *CartHandler {
	return &CartHandler{db: db}
}

type AddToCartRequest struct {
	ProductID int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required,min=1"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

func (h *CartHandler) getOrCreateCart(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := h.db.Where("user_id = ?", userID).
		Preload("Items.Product").
		First(&cart).Error

	if err == gorm.ErrRecordNotFound {
		cart = models.Cart{UserID: userID}
		if err := h.db.Create(&cart).Error; err != nil {
			return nil, err
		}
		return &cart, nil
	}

	return &cart, err
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	cart, err := h.getOrCreateCart(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil keranjang",
		})
		return
	}

	var total float64
	itemCount := 0
	for _, item := range cart.Items {
		total += item.Subtotal
		itemCount += item.Quantity
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"items":      cart.Items,
			"total":      total,
			"item_count": itemCount,
		},
	})
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "product_id dan quantity diperlukan",
		})
		return
	}

	var product models.Product
	if err := h.db.First(&product, req.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Produk tidak ditemukan",
		})
		return
	}

	if product.Stock < req.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Stok produk tidak cukup",
		})
		return
	}

	cart, err := h.getOrCreateCart(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengakses keranjang",
		})
		return
	}

	var existingItem models.CartItem
	err = h.db.Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&existingItem).Error

	if err == nil {
		existingItem.Quantity += req.Quantity
		existingItem.Subtotal = float64(existingItem.Quantity) * product.Price
		if err := h.db.Save(&existingItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Gagal update keranjang",
			})
			return
		}
	} else {
		newItem := models.CartItem{
			CartID:    cart.ID,
			ProductID: uint(req.ProductID),
			Quantity:  req.Quantity,
			Subtotal:  float64(req.Quantity) * product.Price,
		}
		if err := h.db.Create(&newItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Gagal menambahkan ke keranjang",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Produk ditambahkan ke keranjang",
	})
}

func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	userID := c.GetUint("user_id")

	itemIDStr := c.Param("id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID item tidak valid",
		})
		return
	}

	var req UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "quantity diperlukan",
		})
		return
	}

	var cartItem models.CartItem
	if err := h.db.Joins("JOIN carts ON carts.id = cart_items.cart_id").
		Where("cart_items.id = ? AND carts.user_id = ?", itemID, userID).
		First(&cartItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Item tidak ditemukan",
		})
		return
	}

	var product models.Product
	if err := h.db.First(&product, cartItem.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Produk tidak ditemukan",
		})
		return
	}

	if product.Stock < req.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Stok produk tidak cukup",
		})
		return
	}

	cartItem.Quantity = req.Quantity
	cartItem.Subtotal = float64(req.Quantity) * product.Price
	if err := h.db.Save(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal update item",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Item keranjang diperbarui",
	})
}

func (h *CartHandler) RemoveCartItem(c *gin.Context) {
	userID := c.GetUint("user_id")

	itemIDStr := c.Param("id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID item tidak valid",
		})
		return
	}

	result := h.db.Joins("JOIN carts ON carts.id = cart_items.cart_id").
		Where("cart_items.id = ? AND carts.user_id = ?", itemID, userID).
		Delete(&models.CartItem{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghapus item",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Item tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Item dihapus dari keranjang",
	})
}

func (h *CartHandler) ClearCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	cart, err := h.getOrCreateCart(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengakses keranjang",
		})
		return
	}

	if err := h.db.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengosongkan keranjang",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Keranjang dikosongkan",
	})
}
