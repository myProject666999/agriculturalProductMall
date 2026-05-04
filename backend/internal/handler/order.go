package handler

import (
	"agricultural-product-mall/internal/model"
	"agricultural-product-mall/internal/service"
	"agricultural-product-mall/pkg/database"
	"agricultural-product-mall/pkg/response"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateOrderRequest struct {
	AddressID uint   `json:"address_id" binding:"required"`
	CartIDs   []uint `json:"cart_ids" binding:"required"`
	Remark    string `json:"remark"`
}

type OrderItemInfo struct {
	ProductID uint   `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func generateOrderNo() string {
	now := time.Now()
	timestamp := now.Format("20060102150405")
	nano := now.Nanosecond() / 1000
	return fmt.Sprintf("ORD%s%06d", timestamp, nano%1000000)
}

func GetOrderList(c *gin.Context) {
	userID := c.GetUint("user_id")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	var orders []model.Order

	query := database.DB.Model(&model.Order{}).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("Merchant").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&orders)

	for i := range orders {
		var items []model.OrderItem
		database.DB.Where("order_id = ?", orders[i].ID).Preload("Product").Find(&items)
	}

	response.SuccessPage(c, orders, total, page, pageSize)
}

func CreateOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if len(req.CartIDs) == 0 {
		response.BadRequest(c, "请选择商品")
		return
	}

	var address model.Address
	if err := database.DB.Where("id = ? AND user_id = ?", req.AddressID, userID).
		First(&address).Error; err != nil {
		response.NotFound(c, "收货地址不存在")
		return
	}

	var cartItems []model.Cart
	database.DB.Where("id IN ? AND user_id = ?", req.CartIDs, userID).
		Preload("Product").
		Find(&cartItems)

	if len(cartItems) == 0 {
		response.BadRequest(c, "购物车商品不存在")
		return
	}

	merchantMap := make(map[uint][]model.Cart)
	for _, item := range cartItems {
		merchantMap[item.Product.MerchantID] = append(merchantMap[item.Product.MerchantID], item)
	}

	var orders []model.Order

	for merchantID, items := range merchantMap {
		var totalAmount float64
		var orderItems []model.OrderItem

		for _, item := range items {
			product := item.Product
			if product.Status != 1 {
				response.BadRequest(c, "商品 "+product.Name+" 已下架")
				return
			}
			if product.Stock < item.Quantity {
				response.BadRequest(c, "商品 "+product.Name+" 库存不足")
				return
			}

			subTotal := product.Price * float64(item.Quantity)
			totalAmount += subTotal

			orderItems = append(orderItems, model.OrderItem{
				ProductID:    product.ID,
				ProductName:  product.Name,
				ProductImage: product.MainImage,
				Price:        product.Price,
				Quantity:     item.Quantity,
				SubTotal:     subTotal,
			})
		}

		orderNo := generateOrderNo()
		shippingFee := 0.0
		if totalAmount < 99 {
			shippingFee = 10.0
		}

		order := model.Order{
			OrderNo:         orderNo,
			UserID:          userID,
			MerchantID:      merchantID,
			Status:          "pending",
			TotalAmount:     totalAmount,
			DiscountAmount:  0,
			PayAmount:       totalAmount + shippingFee,
			ShippingFee:     shippingFee,
			ReceiverName:    address.Name,
			ReceiverPhone:   address.Phone,
			ReceiverAddress: address.Province + address.City + address.District + address.Detail,
			Remark:          req.Remark,
		}

		if err := database.DB.Create(&order).Error; err != nil {
			response.InternalServerError(c, "创建订单失败")
			return
		}

		for i := range orderItems {
			orderItems[i].OrderID = order.ID
		}

		if err := database.DB.Create(&orderItems).Error; err != nil {
			response.InternalServerError(c, "创建订单失败")
			return
		}

		orders = append(orders, order)
	}

	database.DB.Where("id IN ? AND user_id = ?", req.CartIDs, userID).Delete(&model.Cart{})

	response.SuccessWithMessage(c, "订单创建成功", orders)
}

func GetOrderDetail(c *gin.Context) {
	userID := c.GetUint("user_id")
	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var order model.Order
	if err := database.DB.Preload("Merchant").
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error; err != nil {
		response.NotFound(c, "订单不存在")
		return
	}

	var orderItems []model.OrderItem
	database.DB.Where("order_id = ?", orderID).Preload("Product").Find(&orderItems)

	response.Success(c, gin.H{
		"order": order,
		"items": orderItems,
	})
}

func CancelOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var order model.Order
	if err := database.DB.Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error; err != nil {
		response.NotFound(c, "订单不存在")
		return
	}

	if order.Status != "pending" {
		response.BadRequest(c, "该订单不能取消")
		return
	}

	order.Status = "cancelled"
	if err := database.DB.Save(&order).Error; err != nil {
		response.InternalServerError(c, "取消订单失败")
		return
	}

	var orderItems []model.OrderItem
	database.DB.Where("order_id = ?", orderID).Find(&orderItems)

	for _, item := range orderItems {
		database.DB.Model(&model.Product{}).
			Where("id = ?", item.ProductID).
			Update("stock", gorm.Expr("stock + ?", item.Quantity))
	}

	response.SuccessWithMessage(c, "订单已取消", nil)
}

func PayOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var order model.Order
	if err := database.DB.Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error; err != nil {
		response.NotFound(c, "订单不存在")
		return
	}

	if order.Status != "pending" {
		response.BadRequest(c, "该订单不能支付")
		return
	}

	var orderItems []model.OrderItem
	database.DB.Where("order_id = ?", orderID).Find(&orderItems)

	for _, item := range orderItems {
		result := database.DB.Model(&model.Product{}).
			Where("id = ? AND stock >= ?", item.ProductID, item.Quantity).
			Update("stock", gorm.Expr("stock - ?", item.Quantity))

		if result.RowsAffected == 0 {
			response.BadRequest(c, "商品库存不足")
			return
		}

		database.DB.Model(&model.Product{}).
			Where("id = ?", item.ProductID).
			Update("sales", gorm.Expr("sales + ?", item.Quantity))

		recService := service.NewRecommendationService()
		recService.RecordBehavior(userID, item.ProductID, "purchase", 0)
	}

	now := time.Now()
	order.Status = "paid"
	order.PayTime = &now
	order.PayMethod = "online"

	if err := database.DB.Save(&order).Error; err != nil {
		response.InternalServerError(c, "支付失败")
		return
	}

	response.SuccessWithMessage(c, "支付成功", order)
}

func ConfirmOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var order model.Order
	if err := database.DB.Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error; err != nil {
		response.NotFound(c, "订单不存在")
		return
	}

	if order.Status != "shipped" {
		response.BadRequest(c, "该订单不能确认收货")
		return
	}

	now := time.Now()
	order.Status = "completed"
	order.CompleteTime = &now

	if err := database.DB.Save(&order).Error; err != nil {
		response.InternalServerError(c, "确认收货失败")
		return
	}

	response.SuccessWithMessage(c, "确认收货成功", order)
}

type ApplyRefundRequest struct {
	Reason string  `json:"reason" binding:"required"`
	Amount float64 `json:"amount"`
}

func ApplyRefund(c *gin.Context) {
	userID := c.GetUint("user_id")
	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req ApplyRefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var order model.Order
	if err := database.DB.Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error; err != nil {
		response.NotFound(c, "订单不存在")
		return
	}

	if order.Status != "paid" && order.Status != "shipped" {
		response.BadRequest(c, "该订单不能申请退款")
		return
	}

	if order.RefundStatus != "none" {
		response.BadRequest(c, "该订单已有退款申请")
		return
	}

	if req.Amount <= 0 || req.Amount > order.PayAmount {
		req.Amount = order.PayAmount
	}

	order.RefundStatus = "applying"
	order.RefundReason = req.Reason
	order.RefundAmount = req.Amount

	if err := database.DB.Save(&order).Error; err != nil {
		response.InternalServerError(c, "申请退款失败")
		return
	}

	response.SuccessWithMessage(c, "退款申请已提交", order)
}

type CreateReviewRequest struct {
	Rating      int    `json:"rating" binding:"required,min=1,max=5"`
	Content     string `json:"content" binding:"required"`
	Images      string `json:"images"`
	IsAnonymous int    `json:"is_anonymous"`
}

func CreateReview(c *gin.Context) {
	userID := c.GetUint("user_id")
	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var order model.Order
	if err := database.DB.Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error; err != nil {
		response.NotFound(c, "订单不存在")
		return
	}

	if order.Status != "completed" {
		response.BadRequest(c, "该订单不能评价")
		return
	}

	var existingReview model.Review
	if database.DB.Where("order_id = ?", orderID).First(&existingReview).Error == nil {
		response.BadRequest(c, "该订单已评价")
		return
	}

	var orderItems []model.OrderItem
	database.DB.Where("order_id = ?", orderID).Find(&orderItems)

	for _, item := range orderItems {
		review := model.Review{
			UserID:      userID,
			ProductID:   item.ProductID,
			OrderID:     order.ID,
			MerchantID:  order.MerchantID,
			Rating:      req.Rating,
			Content:     req.Content,
			Images:      req.Images,
			IsAnonymous: req.IsAnonymous,
			Status:      1,
		}

		if err := database.DB.Create(&review).Error; err != nil {
			response.InternalServerError(c, "评价失败")
			return
		}

		database.DB.Model(&model.Product{}).
			Where("id = ?", item.ProductID).
			Update("review_count", gorm.Expr("review_count + 1"))

		recService := service.NewRecommendationService()
		recService.RecordBehavior(userID, item.ProductID, "rating", float64(req.Rating))
	}

	response.SuccessWithMessage(c, "评价成功", nil)
}

func GetLogisticsInfo(c *gin.Context) {
	orderIDStr := c.Param("order_id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var logistics model.Logistics
	if err := database.DB.Where("order_id = ?", orderID).
		First(&logistics).Error; err != nil {
		response.NotFound(c, "物流信息不存在")
		return
	}

	response.Success(c, logistics)
}
