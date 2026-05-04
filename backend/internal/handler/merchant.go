package handler

import (
	"agricultural-product-mall/internal/model"
	"agricultural-product-mall/pkg/database"
	"agricultural-product-mall/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func getMerchantIDByUserID(userID uint) uint {
	var merchant model.Merchant
	database.DB.Where("user_id = ?", userID).First(&merchant)
	return merchant.ID
}

func GetMerchantProducts(c *gin.Context) {
	userID := c.GetUint("user_id")
	merchantID := getMerchantIDByUserID(userID)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	var products []model.Product

	query := database.DB.Model(&model.Product{}).Where("merchant_id = ?", merchantID)

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("Category").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&products)

	response.SuccessPage(c, products, total, page, pageSize)
}

type CreateProductRequest struct {
	Name           string  `json:"name" binding:"required"`
	CategoryID     uint    `json:"category_id" binding:"required"`
	Description    string  `json:"description"`
	MainImage      string  `json:"main_image"`
	SubImages      string  `json:"sub_images"`
	Price          float64 `json:"price" binding:"required"`
	OriginalPrice  float64 `json:"original_price"`
	Stock          int     `json:"stock" binding:"required"`
	Unit           string  `json:"unit"`
	Specifications string  `json:"specifications"`
	IsHot          int     `json:"is_hot"`
	IsNew          int     `json:"is_new"`
	IsRecommend    int     `json:"is_recommend"`
}

func CreateProduct(c *gin.Context) {
	userID := c.GetUint("user_id")
	merchantID := getMerchantIDByUserID(userID)

	if merchantID == 0 {
		response.Forbidden(c, "不是商户账号")
		return
	}

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	product := model.Product{
		Name:           req.Name,
		CategoryID:     req.CategoryID,
		MerchantID:     merchantID,
		Description:    req.Description,
		MainImage:      req.MainImage,
		SubImages:      req.SubImages,
		Price:          req.Price,
		OriginalPrice:  req.OriginalPrice,
		Stock:          req.Stock,
		Unit:           req.Unit,
		Specifications: req.Specifications,
		IsHot:          req.IsHot,
		IsNew:          req.IsNew,
		IsRecommend:    req.IsRecommend,
		Status:         1,
		Rating:         5.0,
	}

	if err := database.DB.Create(&product).Error; err != nil {
		response.InternalServerError(c, "创建商品失败")
		return
	}

	inventory := model.Inventory{
		ProductID: product.ID,
		Quantity:  req.Stock,
		MinStock:  10,
		MaxStock:  1000,
	}
	database.DB.Create(&inventory)

	response.SuccessWithMessage(c, "创建成功", product)
}

func UpdateProduct(c *gin.Context) {
	userID := c.GetUint("user_id")
	merchantID := getMerchantIDByUserID(userID)

	idStr := c.Param("id")
	productID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var product model.Product
	if err := database.DB.Where("id = ? AND merchant_id = ?", productID, merchantID).
		First(&product).Error; err != nil {
		response.NotFound(c, "商品不存在")
		return
	}

	oldStock := product.Stock

	product.Name = req.Name
	product.CategoryID = req.CategoryID
	product.Description = req.Description
	product.MainImage = req.MainImage
	product.SubImages = req.SubImages
	product.Price = req.Price
	product.OriginalPrice = req.OriginalPrice
	product.Stock = req.Stock
	product.Unit = req.Unit
	product.Specifications = req.Specifications
	product.IsHot = req.IsHot
	product.IsNew = req.IsNew
	product.IsRecommend = req.IsRecommend

	if err := database.DB.Save(&product).Error; err != nil {
		response.InternalServerError(c, "更新商品失败")
		return
	}

	if oldStock != req.Stock {
		database.DB.Model(&model.Inventory{}).
			Where("product_id = ?", productID).
			Update("quantity", req.Stock)
	}

	response.SuccessWithMessage(c, "更新成功", product)
}

func DeleteProduct(c *gin.Context) {
	userID := c.GetUint("user_id")
	merchantID := getMerchantIDByUserID(userID)

	idStr := c.Param("id")
	productID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	result := database.DB.Where("id = ? AND merchant_id = ?", productID, merchantID).
		Delete(&model.Product{})

	if result.RowsAffected == 0 {
		response.NotFound(c, "商品不存在")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

func UpdateProductStatus(c *gin.Context) {
	userID := c.GetUint("user_id")
	merchantID := getMerchantIDByUserID(userID)

	idStr := c.Param("id")
	productID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	result := database.DB.Model(&model.Product{}).
		Where("id = ? AND merchant_id = ?", productID, merchantID).
		Update("status", req.Status)

	if result.RowsAffected == 0 {
		response.NotFound(c, "商品不存在")
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

func GetInventoryList(c *gin.Context) {
	userID := c.GetUint("user_id")
	merchantID := getMerchantIDByUserID(userID)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	var inventories []model.Inventory

	query := database.DB.Model(&model.Inventory{}).
		Joins("JOIN products ON products.id = inventories.product_id").
		Where("products.merchant_id = ?", merchantID)

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("Product").
		Order("inventories.updated_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&inventories)

	response.SuccessPage(c, inventories, total, page, pageSize)
}

type InventoryOperationRequest struct {
	ProductID uint   `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
	Remark    string `json:"remark"`
}

func InventoryIn(c *gin.Context) {
	userID := c.GetUint("user_id")
	merchantID := getMerchantIDByUserID(userID)

	var req InventoryOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var product model.Product
	if err := database.DB.Where("id = ? AND merchant_id = ?", req.ProductID, merchantID).
		First(&product).Error; err != nil {
		response.NotFound(c, "商品不存在")
		return
	}

	var inventory model.Inventory
	database.DB.Where("product_id = ?", req.ProductID).First(&inventory)

	beforeQty := inventory.Quantity
	afterQty := beforeQty + req.Quantity

	inventory.Quantity = afterQty
	database.DB.Save(&inventory)

	product.Stock = afterQty
	database.DB.Save(&product)

	record := model.InventoryRecord{
		ProductID:  req.ProductID,
		Type:       "in",
		Quantity:   req.Quantity,
		BeforeQty:  beforeQty,
		AfterQty:   afterQty,
		OperatorID: userID,
		Remark:     req.Remark,
	}
	database.DB.Create(&record)

	response.SuccessWithMessage(c, "入库成功", gin.H{
		"before_qty": beforeQty,
		"after_qty":  afterQty,
	})
}

func InventoryOut(c *gin.Context) {
	userID := c.GetUint("user_id")
	merchantID := getMerchantIDByUserID(userID)

	var req InventoryOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var product model.Product
	if err := database.DB.Where("id = ? AND merchant_id = ?", req.ProductID, merchantID).
		First(&product).Error; err != nil {
		response.NotFound(c, "商品不存在")
		return
	}

	var inventory model.Inventory
	database.DB.Where("product_id = ?", req.ProductID).First(&inventory)

	if inventory.Quantity < req.Quantity {
		response.BadRequest(c, "库存不足")
		return
	}

	beforeQty := inventory.Quantity
	afterQty := beforeQty - req.Quantity

	inventory.Quantity = afterQty
	database.DB.Save(&inventory)

	product.Stock = afterQty
	database.DB.Save(&product)

	record := model.InventoryRecord{
		ProductID:  req.ProductID,
		Type:       "out",
		Quantity:   req.Quantity,
		BeforeQty:  beforeQty,
		AfterQty:   afterQty,
		OperatorID: userID,
		Remark:     req.Remark,
	}
	database.DB.Create(&record)

	response.SuccessWithMessage(c, "出库成功", gin.H{
		"before_qty": beforeQty,
		"after_qty":  afterQty,
	})
}

func GetInventoryRecords(c *gin.Context) {
	userID := c.GetUint("user_id")
	merchantID := getMerchantIDByUserID(userID)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	productID := c.Query("product_id")
	recordType := c.Query("type")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	var records []model.InventoryRecord

	query := database.DB.Model(&model.InventoryRecord{}).
		Joins("JOIN products ON products.id = inventory_records.product_id").
		Where("products.merchant_id = ?", merchantID)

	if productID != "" {
		query = query.Where("inventory_records.product_id = ?", productID)
	}
	if recordType != "" {
		query = query.Where("inventory_records.type = ?", recordType)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("Product").Preload("Operator").
		Order("inventory_records.created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&records)

	response.SuccessPage(c, records, total, page, pageSize)
}

func GetMerchantOrders(c *gin.Context) {
	userID := c.GetUint("user_id")
	merchantID := getMerchantIDByUserID(userID)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	var orders []model.Order

	query := database.DB.Model(&model.Order{}).Where("merchant_id = ?", merchantID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("User").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&orders)

	response.SuccessPage(c, orders, total, page, pageSize)
}

func GetMerchantOrderDetail(c *gin.Context) {
	userID := c.GetUint("user_id")
	merchantID := getMerchantIDByUserID(userID)

	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var order model.Order
	if err := database.DB.Preload("User").
		Where("id = ? AND merchant_id = ?", orderID, merchantID).
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

type ShipOrderRequest struct {
	LogisticsNo string `json:"logistics_no" binding:"required"`
	Company     string `json:"company" binding:"required"`
}

func ShipOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	merchantID := getMerchantIDByUserID(userID)

	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req ShipOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var order model.Order
	if err := database.DB.Where("id = ? AND merchant_id = ?", orderID, merchantID).
		First(&order).Error; err != nil {
		response.NotFound(c, "订单不存在")
		return
	}

	if order.Status != "paid" {
		response.BadRequest(c, "该订单不能发货")
		return
	}

	now := time.Now()
	order.Status = "shipped"
	order.ShippingTime = &now
	database.DB.Save(&order)

	logistics := model.Logistics{
		OrderID:      order.ID,
		LogisticsNo:  req.LogisticsNo,
		Company:      req.Company,
		ShipTime:     &now,
		CurrentStatus: "已发货",
	}
	database.DB.Create(&logistics)

	response.SuccessWithMessage(c, "发货成功", order)
}

func CompleteOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	merchantID := getMerchantIDByUserID(userID)

	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var order model.Order
	if err := database.DB.Where("id = ? AND merchant_id = ?", orderID, merchantID).
		First(&order).Error; err != nil {
		response.NotFound(c, "订单不存在")
		return
	}

	if order.Status != "shipped" {
		response.BadRequest(c, "该订单不能完成")
		return
	}

	now := time.Now()
	order.Status = "completed"
	order.CompleteTime = &now
	database.DB.Save(&order)

	response.SuccessWithMessage(c, "订单已完成", order)
}

type HandleRefundRequest struct {
	Agree  bool   `json:"agree" binding:"required"`
	Reason string `json:"reason"`
}

func HandleRefund(c *gin.Context) {
	userID := c.GetUint("user_id")
	merchantID := getMerchantIDByUserID(userID)

	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req HandleRefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var order model.Order
	if err := database.DB.Where("id = ? AND merchant_id = ?", orderID, merchantID).
		First(&order).Error; err != nil {
		response.NotFound(c, "订单不存在")
		return
	}

	if order.RefundStatus != "applying" {
		response.BadRequest(c, "该订单没有退款申请")
		return
	}

	if req.Agree {
		order.RefundStatus = "completed"
	} else {
		order.RefundStatus = "rejected"
	}

	database.DB.Save(&order)

	response.SuccessWithMessage(c, "处理成功", order)
}

func GetMerchantReviews(c *gin.Context) {
	userID := c.GetUint("user_id")
	merchantID := getMerchantIDByUserID(userID)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	var reviews []model.Review

	query := database.DB.Model(&model.Review{}).Where("merchant_id = ?", merchantID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("Product").Preload("User").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&reviews)

	response.SuccessPage(c, reviews, total, page, pageSize)
}

func GetSalesStats(c *gin.Context) {
	userID := c.GetUint("user_id")
	merchantID := getMerchantIDByUserID(userID)

	var totalOrders int64
	var totalRevenue float64

	database.DB.Model(&model.Order{}).
		Where("merchant_id = ? AND status IN ?", merchantID, []string{"completed", "paid", "shipped"}).
		Count(&totalOrders)

	database.DB.Model(&model.Order{}).
		Where("merchant_id = ? AND status = ?", merchantID, "completed").
		Select("COALESCE(SUM(pay_amount), 0)").
		Scan(&totalRevenue)

	var topProducts []struct {
		ProductID   uint
		ProductName string
		TotalSales  int
	}

	database.DB.Table("order_items").
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Where("orders.merchant_id = ? AND orders.status = ?", merchantID, "completed").
		Select("order_items.product_id, order_items.product_name, SUM(order_items.quantity) as total_sales").
		Group("order_items.product_id").
		Order("total_sales DESC").
		Limit(10).
		Scan(&topProducts)

	response.Success(c, gin.H{
		"total_orders":   totalOrders,
		"total_revenue":  totalRevenue,
		"top_products":   topProducts,
	})
}

func GetMerchantInfo(c *gin.Context) {
	userID := c.GetUint("user_id")
	merchantID := getMerchantIDByUserID(userID)

	var merchant model.Merchant
	if err := database.DB.First(&merchant, merchantID).Error; err != nil {
		response.NotFound(c, "商户信息不存在")
		return
	}

	response.Success(c, merchant)
}

func UpdateMerchantInfo(c *gin.Context) {
	userID := c.GetUint("user_id")
	merchantID := getMerchantIDByUserID(userID)

	var req struct {
		StoreName    string `json:"store_name"`
		StoreLogo    string `json:"store_logo"`
		StoreDesc    string `json:"store_desc"`
		ContactName  string `json:"contact_name"`
		ContactPhone string `json:"contact_phone"`
		ContactEmail string `json:"contact_email"`
		Address      string `json:"address"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var merchant model.Merchant
	if err := database.DB.First(&merchant, merchantID).Error; err != nil {
		response.NotFound(c, "商户信息不存在")
		return
	}

	if req.StoreName != "" {
		merchant.StoreName = req.StoreName
	}
	if req.StoreLogo != "" {
		merchant.StoreLogo = req.StoreLogo
	}
	if req.StoreDesc != "" {
		merchant.StoreDesc = req.StoreDesc
	}
	if req.ContactName != "" {
		merchant.ContactName = req.ContactName
	}
	if req.ContactPhone != "" {
		merchant.ContactPhone = req.ContactPhone
	}
	if req.ContactEmail != "" {
		merchant.ContactEmail = req.ContactEmail
	}
	if req.Address != "" {
		merchant.Address = req.Address
	}

	if err := database.DB.Save(&merchant).Error; err != nil {
		response.InternalServerError(c, "更新失败")
		return
	}

	response.SuccessWithMessage(c, "更新成功", merchant)
}

func UpdateMerchantPassword(c *gin.Context) {
	UpdatePassword(c)
}
