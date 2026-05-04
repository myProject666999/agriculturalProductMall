package handler

import (
	"agricultural-product-mall/internal/model"
	"agricultural-product-mall/pkg/database"
	"agricultural-product-mall/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AddToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

func GetCartList(c *gin.Context) {
	userID := c.GetUint("user_id")

	var cartItems []model.Cart
	database.DB.Where("user_id = ?", userID).
		Preload("Product").
		Preload("Product.Category").
		Preload("Product.Merchant").
		Order("created_at DESC").
		Find(&cartItems)

	response.Success(c, cartItems)
}

func AddToCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.Quantity < 1 {
		response.BadRequest(c, "数量不能小于1")
		return
	}

	var product model.Product
	if err := database.DB.First(&product, req.ProductID).Error; err != nil {
		response.NotFound(c, "商品不存在")
		return
	}

	if product.Status != 1 {
		response.BadRequest(c, "商品已下架")
		return
	}

	if product.Stock < req.Quantity {
		response.BadRequest(c, "库存不足")
		return
	}

	var existingCart model.Cart
	result := database.DB.Where("user_id = ? AND product_id = ?", userID, req.ProductID).
		First(&existingCart)

	if result.Error == nil {
		existingCart.Quantity += req.Quantity
		if existingCart.Quantity > product.Stock {
			existingCart.Quantity = product.Stock
		}
		database.DB.Save(&existingCart)
		response.SuccessWithMessage(c, "添加成功", existingCart)
		return
	}

	cart := model.Cart{
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Selected:  1,
	}

	if err := database.DB.Create(&cart).Error; err != nil {
		response.InternalServerError(c, "添加购物车失败")
		return
	}

	response.SuccessWithMessage(c, "添加成功", cart)
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity"`
	Selected *int `json:"selected"`
}

func UpdateCartItem(c *gin.Context) {
	userID := c.GetUint("user_id")
	idStr := c.Param("id")
	cartID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var cartItem model.Cart
	if err := database.DB.Where("id = ? AND user_id = ?", cartID, userID).
		First(&cartItem).Error; err != nil {
		response.NotFound(c, "购物车商品不存在")
		return
	}

	if req.Quantity > 0 {
		var product model.Product
		database.DB.First(&product, cartItem.ProductID)
		if req.Quantity > product.Stock {
			response.BadRequest(c, "库存不足")
			return
		}
		cartItem.Quantity = req.Quantity
	}

	if req.Selected != nil {
		cartItem.Selected = *req.Selected
	}

	if err := database.DB.Save(&cartItem).Error; err != nil {
		response.InternalServerError(c, "更新失败")
		return
	}

	response.Success(c, cartItem)
}

func RemoveFromCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	idStr := c.Param("id")
	cartID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	result := database.DB.Where("id = ? AND user_id = ?", cartID, userID).
		Delete(&model.Cart{})

	if result.RowsAffected == 0 {
		response.NotFound(c, "购物车商品不存在")
		return
	}

	response.SuccessWithMessage(c, "移除成功", nil)
}

func ClearCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	database.DB.Where("user_id = ?", userID).Delete(&model.Cart{})

	response.SuccessWithMessage(c, "购物车已清空", nil)
}
