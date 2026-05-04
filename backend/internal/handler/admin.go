package handler

import (
	"agricultural-product-mall/internal/model"
	"agricultural-product-mall/pkg/database"
	"agricultural-product-mall/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAdminSalesStats(c *gin.Context) {
	var totalRevenue float64
	var totalOrders int64
	var todayRevenue float64
	var todayOrders int64

	todayStart := time.Now().Format("2006-01-02") + " 00:00:00"
	todayEnd := time.Now().Format("2006-01-02") + " 23:59:59"

	database.DB.Model(&model.Order{}).
		Where("status = ?", "completed").
		Select("COALESCE(SUM(pay_amount), 0)").
		Scan(&totalRevenue)

	database.DB.Model(&model.Order{}).
		Where("status = ?", "completed").
		Count(&totalOrders)

	database.DB.Model(&model.Order{}).
		Where("status = ? AND created_at BETWEEN ? AND ?", "completed", todayStart, todayEnd).
		Select("COALESCE(SUM(pay_amount), 0)").
		Scan(&todayRevenue)

	database.DB.Model(&model.Order{}).
		Where("created_at BETWEEN ? AND ?", todayStart, todayEnd).
		Count(&todayOrders)

	var salesByCategory []struct {
		CategoryID   uint
		CategoryName string
		TotalSales   float64
		TotalCount   int64
	}

	database.DB.Table("order_items").
		Joins("JOIN products ON order_items.product_id = products.id").
		Joins("JOIN categories ON products.category_id = categories.id").
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Where("orders.status = ?", "completed").
		Select("categories.id as category_id, categories.name as category_name, COALESCE(SUM(order_items.sub_total), 0) as total_sales, COUNT(*) as total_count").
		Group("categories.id").
		Scan(&salesByCategory)

	response.Success(c, gin.H{
		"total_revenue":  totalRevenue,
		"total_orders":   totalOrders,
		"today_revenue":  todayRevenue,
		"today_orders":   todayOrders,
		"sales_by_category": salesByCategory,
	})
}

func GetAdminUserStats(c *gin.Context) {
	var totalUsers int64
	var activeUsers int64
	var newTodayUsers int64
	var merchantCount int64

	todayStart := time.Now().Format("2006-01-02") + " 00:00:00"
	todayEnd := time.Now().Format("2006-01-02") + " 23:59:59"

	database.DB.Model(&model.User{}).
		Where("role = ?", "user").
		Count(&totalUsers)

	database.DB.Model(&model.User{}).
		Where("role = ? AND last_login_at >= ?", "user", time.Now().AddDate(0, 0, -7)).
		Count(&activeUsers)

	database.DB.Model(&model.User{}).
		Where("role = ? AND created_at BETWEEN ? AND ?", "user", todayStart, todayEnd).
		Count(&newTodayUsers)

	database.DB.Model(&model.User{}).
		Where("role = ?", "merchant").
		Count(&merchantCount)

	response.Success(c, gin.H{
		"total_users":       totalUsers,
		"active_users":      activeUsers,
		"new_today_users":   newTodayUsers,
		"merchant_count":    merchantCount,
	})
}

func GetCategoryStats(c *gin.Context) {
	var stats []struct {
		ID            uint
		Name          string
		ProductCount  int64
		TotalSales    float64
	}

	database.DB.Table("categories").
		Joins("LEFT JOIN products ON categories.id = products.category_id AND products.status = 1").
		Joins("LEFT JOIN order_items ON products.id = order_items.product_id").
		Joins("LEFT JOIN orders ON order_items.order_id = orders.id AND orders.status = 'completed'").
		Where("categories.status = ? AND categories.parent_id = ?", 1, 0).
		Select("categories.id, categories.name, COUNT(DISTINCT products.id) as product_count, COALESCE(SUM(order_items.sub_total), 0) as total_sales").
		Group("categories.id").
		Order("categories.sort ASC").
		Scan(&stats)

	response.Success(c, stats)
}

func GetProductStats(c *gin.Context) {
	var totalProducts int64
	var activeProducts int64
	var outOfStockProducts int64
	var hotProducts int64

	database.DB.Model(&model.Product{}).Count(&totalProducts)
	database.DB.Model(&model.Product{}).Where("status = ?", 1).Count(&activeProducts)
	database.DB.Model(&model.Product{}).Where("stock <= 0").Count(&outOfStockProducts)
	database.DB.Model(&model.Product{}).Where("is_hot = ?", 1).Count(&hotProducts)

	var topProducts []struct {
		ID          uint
		Name        string
		Sales       int
		TotalRevenue float64
	}

	database.DB.Table("products").
		Joins("JOIN order_items ON products.id = order_items.product_id").
		Joins("JOIN orders ON order_items.order_id = orders.id AND orders.status = 'completed'").
		Where("products.status = ?", 1).
		Select("products.id, products.name, products.sales, COALESCE(SUM(order_items.sub_total), 0) as total_revenue").
		Group("products.id").
		Order("products.sales DESC").
		Limit(10).
		Scan(&topProducts)

	response.Success(c, gin.H{
		"total_products":      totalProducts,
		"active_products":     activeProducts,
		"out_of_stock_products": outOfStockProducts,
		"hot_products":        hotProducts,
		"top_products":        topProducts,
	})
}

func GetAdminProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status := c.Query("status")
	keyword := c.Query("keyword")
	merchantID := c.Query("merchant_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	var products []model.Product

	query := database.DB.Model(&model.Product{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if merchantID != "" {
		query = query.Where("merchant_id = ?", merchantID)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("Category").Preload("Merchant").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&products)

	response.SuccessPage(c, products, total, page, pageSize)
}

func CreateAdminProduct(c *gin.Context) {
	var req struct {
		Name           string  `json:"name" binding:"required"`
		CategoryID     uint    `json:"category_id" binding:"required"`
		MerchantID     uint    `json:"merchant_id" binding:"required"`
		Description    string  `json:"description"`
		MainImage      string  `json:"main_image"`
		Price          float64 `json:"price" binding:"required"`
		Stock          int     `json:"stock" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	product := model.Product{
		Name:        req.Name,
		CategoryID:  req.CategoryID,
		MerchantID:  req.MerchantID,
		Description: req.Description,
		MainImage:   req.MainImage,
		Price:       req.Price,
		Stock:       req.Stock,
		Status:      1,
		Rating:      5.0,
	}

	if err := database.DB.Create(&product).Error; err != nil {
		response.InternalServerError(c, "创建商品失败")
		return
	}

	inventory := model.Inventory{
		ProductID: product.ID,
		Quantity:  req.Stock,
	}
	database.DB.Create(&inventory)

	response.SuccessWithMessage(c, "创建成功", product)
}

func UpdateAdminProduct(c *gin.Context) {
	idStr := c.Param("id")
	productID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var product model.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		response.NotFound(c, "商品不存在")
		return
	}

	var req struct {
		Name           string  `json:"name"`
		CategoryID     uint    `json:"category_id"`
		Description    string  `json:"description"`
		MainImage      string  `json:"main_image"`
		Price          float64 `json:"price"`
		Stock          int     `json:"stock"`
		Status         int     `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.CategoryID > 0 {
		product.CategoryID = req.CategoryID
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.MainImage != "" {
		product.MainImage = req.MainImage
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
		database.DB.Model(&model.Inventory{}).
			Where("product_id = ?", productID).
			Update("quantity", req.Stock)
	}
	if req.Status >= 0 {
		product.Status = req.Status
	}

	database.DB.Save(&product)
	response.SuccessWithMessage(c, "更新成功", product)
}

func DeleteAdminProduct(c *gin.Context) {
	idStr := c.Param("id")
	productID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	result := database.DB.Delete(&model.Product{}, productID)
	if result.RowsAffected == 0 {
		response.NotFound(c, "商品不存在")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

func UpdateAdminProductStatus(c *gin.Context) {
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
		Where("id = ?", productID).
		Update("status", req.Status)

	if result.RowsAffected == 0 {
		response.NotFound(c, "商品不存在")
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

func GetAdminCategories(c *gin.Context) {
	var categories []model.Category
	database.DB.Order("sort ASC, id ASC").Find(&categories)
	response.Success(c, categories)
}

func CreateCategory(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		ParentID uint   `json:"parent_id"`
		Icon     string `json:"icon"`
		Sort     int    `json:"sort"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	category := model.Category{
		Name:     req.Name,
		ParentID: req.ParentID,
		Icon:     req.Icon,
		Sort:     req.Sort,
		Status:   1,
	}

	if err := database.DB.Create(&category).Error; err != nil {
		response.InternalServerError(c, "创建分类失败")
		return
	}

	response.SuccessWithMessage(c, "创建成功", category)
}

func UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	categoryID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var category model.Category
	if err := database.DB.First(&category, categoryID).Error; err != nil {
		response.NotFound(c, "分类不存在")
		return
	}

	var req struct {
		Name   string `json:"name"`
		Icon   string `json:"icon"`
		Sort   int    `json:"sort"`
		Status int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Icon != "" {
		category.Icon = req.Icon
	}
	if req.Sort >= 0 {
		category.Sort = req.Sort
	}
	if req.Status >= 0 {
		category.Status = req.Status
	}

	database.DB.Save(&category)
	response.SuccessWithMessage(c, "更新成功", category)
}

func DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	categoryID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var productCount int64
	database.DB.Model(&model.Product{}).
		Where("category_id = ?", categoryID).
		Count(&productCount)

	if productCount > 0 {
		response.BadRequest(c, "该分类下有商品，不能删除")
		return
	}

	result := database.DB.Delete(&model.Category{}, categoryID)
	if result.RowsAffected == 0 {
		response.NotFound(c, "分类不存在")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

func GetAdminOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status := c.Query("status")
	orderNo := c.Query("order_no")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	var orders []model.Order

	query := database.DB.Model(&model.Order{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if orderNo != "" {
		query = query.Where("order_no LIKE ?", "%"+orderNo+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("User").Preload("Merchant").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&orders)

	response.SuccessPage(c, orders, total, page, pageSize)
}

func GetAdminOrderDetail(c *gin.Context) {
	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var order model.Order
	if err := database.DB.Preload("User").Preload("Merchant").
		First(&order, orderID).Error; err != nil {
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

func UpdateOrderStatus(c *gin.Context) {
	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	result := database.DB.Model(&model.Order{}).
		Where("id = ?", orderID).
		Update("status", req.Status)

	if result.RowsAffected == 0 {
		response.NotFound(c, "订单不存在")
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

func GetAdminCartList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	var carts []model.Cart

	query := database.DB.Model(&model.Cart{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("User").Preload("Product").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&carts)

	response.SuccessPage(c, carts, total, page, pageSize)
}

func GetAdminReviews(c *gin.Context) {
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

	query := database.DB.Model(&model.Review{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("User").Preload("Product").Preload("Order").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&reviews)

	response.SuccessPage(c, reviews, total, page, pageSize)
}

func DeleteReview(c *gin.Context) {
	idStr := c.Param("id")
	reviewID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	result := database.DB.Delete(&model.Review{}, reviewID)
	if result.RowsAffected == 0 {
		response.NotFound(c, "评价不存在")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

func GetAdminLogistics(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	var logistics []model.Logistics

	query := database.DB.Model(&model.Logistics{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("Order").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&logistics)

	response.SuccessPage(c, logistics, total, page, pageSize)
}

func CreateLogistics(c *gin.Context) {
	var req struct {
		OrderID     uint   `json:"order_id" binding:"required"`
		LogisticsNo string `json:"logistics_no" binding:"required"`
		Company     string `json:"company" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	now := time.Now()
	logistics := model.Logistics{
		OrderID:       req.OrderID,
		LogisticsNo:   req.LogisticsNo,
		Company:       req.Company,
		ShipTime:      &now,
		CurrentStatus: "已发货",
	}

	if err := database.DB.Create(&logistics).Error; err != nil {
		response.InternalServerError(c, "创建物流失败")
		return
	}

	response.SuccessWithMessage(c, "创建成功", logistics)
}

func UpdateLogistics(c *gin.Context) {
	idStr := c.Param("id")
	logisticsID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var logistics model.Logistics
	if err := database.DB.First(&logistics, logisticsID).Error; err != nil {
		response.NotFound(c, "物流信息不存在")
		return
	}

	var req struct {
		LogisticsNo   string `json:"logistics_no"`
		Company       string `json:"company"`
		CurrentStatus string `json:"current_status"`
		Trajectory    string `json:"trajectory"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.LogisticsNo != "" {
		logistics.LogisticsNo = req.LogisticsNo
	}
	if req.Company != "" {
		logistics.Company = req.Company
	}
	if req.CurrentStatus != "" {
		logistics.CurrentStatus = req.CurrentStatus
	}
	if req.Trajectory != "" {
		logistics.Trajectory = req.Trajectory
	}

	database.DB.Save(&logistics)
	response.SuccessWithMessage(c, "更新成功", logistics)
}

func DeleteLogistics(c *gin.Context) {
	idStr := c.Param("id")
	logisticsID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	result := database.DB.Delete(&model.Logistics{}, logisticsID)
	if result.RowsAffected == 0 {
		response.NotFound(c, "物流信息不存在")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

func GetAdminInventory(c *gin.Context) {
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

	query := database.DB.Model(&model.Inventory{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("Product").
		Order("updated_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&inventories)

	response.SuccessPage(c, inventories, total, page, pageSize)
}

func GetAdminInventoryRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	var records []model.InventoryRecord

	query := database.DB.Model(&model.InventoryRecord{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("Product").Preload("Operator").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&records)

	response.SuccessPage(c, records, total, page, pageSize)
}

func GetAdminBanners(c *gin.Context) {
	var banners []model.Banner
	database.DB.Order("sort ASC, id ASC").Find(&banners)
	response.Success(c, banners)
}

func CreateBanner(c *gin.Context) {
	var req struct {
		Title     string `json:"title" binding:"required"`
		Image     string `json:"image" binding:"required"`
		LinkType  string `json:"link_type"`
		LinkValue string `json:"link_value"`
		Sort      int    `json:"sort"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	banner := model.Banner{
		Title:     req.Title,
		Image:     req.Image,
		LinkType:  req.LinkType,
		LinkValue: req.LinkValue,
		Sort:      req.Sort,
		Status:    1,
	}

	if err := database.DB.Create(&banner).Error; err != nil {
		response.InternalServerError(c, "创建轮播图失败")
		return
	}

	response.SuccessWithMessage(c, "创建成功", banner)
}

func UpdateBanner(c *gin.Context) {
	idStr := c.Param("id")
	bannerID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var banner model.Banner
	if err := database.DB.First(&banner, bannerID).Error; err != nil {
		response.NotFound(c, "轮播图不存在")
		return
	}

	var req struct {
		Title     string `json:"title"`
		Image     string `json:"image"`
		LinkType  string `json:"link_type"`
		LinkValue string `json:"link_value"`
		Sort      int    `json:"sort"`
		Status    int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.Title != "" {
		banner.Title = req.Title
	}
	if req.Image != "" {
		banner.Image = req.Image
	}
	if req.LinkType != "" {
		banner.LinkType = req.LinkType
	}
	if req.LinkValue != "" {
		banner.LinkValue = req.LinkValue
	}
	if req.Sort >= 0 {
		banner.Sort = req.Sort
	}
	if req.Status >= 0 {
		banner.Status = req.Status
	}

	database.DB.Save(&banner)
	response.SuccessWithMessage(c, "更新成功", banner)
}

func DeleteBanner(c *gin.Context) {
	idStr := c.Param("id")
	bannerID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	result := database.DB.Delete(&model.Banner{}, bannerID)
	if result.RowsAffected == 0 {
		response.NotFound(c, "轮播图不存在")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

func GetAdminUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	role := c.Query("role")
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	var users []model.User

	query := database.DB.Model(&model.User{})

	if role != "" {
		query = query.Where("role = ?", role)
	}
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&users)

	response.SuccessPage(c, users, total, page, pageSize)
}

func CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Role     string `json:"role"`
		Status   int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var existingUser model.User
	database.DB.Where("username = ?", req.Username).First(&existingUser)
	if existingUser.ID > 0 {
		response.BadRequest(c, "用户名已存在")
		return
	}

	database.DB.Where("email = ?", req.Email).First(&existingUser)
	if existingUser.ID > 0 {
		response.BadRequest(c, "邮箱已被注册")
		return
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		response.InternalServerError(c, "密码加密失败")
		return
	}

	user := model.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Role:     req.Role,
		Status:   1,
	}

	if req.Status > 0 {
		user.Status = req.Status
	}

	if err := database.DB.Create(&user).Error; err != nil {
		response.InternalServerError(c, "创建用户失败")
		return
	}

	response.SuccessWithMessage(c, "创建成功", gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})
}

func UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	var req struct {
		Nickname string `json:"nickname"`
		Phone    string `json:"phone"`
		Avatar   string `json:"avatar"`
		Role     string `json:"role"`
		Status   int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Status >= 0 {
		user.Status = req.Status
	}

	database.DB.Save(&user)
	response.SuccessWithMessage(c, "更新成功", user)
}

func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	result := database.DB.Delete(&model.User{}, userID)
	if result.RowsAffected == 0 {
		response.NotFound(c, "用户不存在")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

func UpdateUserStatus(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 64)
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

	result := database.DB.Model(&model.User{}).
		Where("id = ?", userID).
		Update("status", req.Status)

	if result.RowsAffected == 0 {
		response.NotFound(c, "用户不存在")
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

func GetMenuList(c *gin.Context) {
	var menus []model.Menu
	database.DB.Order("sort ASC, id ASC").Find(&menus)
	response.Success(c, menus)
}

func CreateMenu(c *gin.Context) {
	var req struct {
		ParentID   uint   `json:"parent_id"`
		Name       string `json:"name" binding:"required"`
		Path       string `json:"path"`
		Icon       string `json:"icon"`
		Component  string `json:"component"`
		Permission string `json:"permission"`
		Type       int    `json:"type"`
		Sort       int    `json:"sort"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	menu := model.Menu{
		ParentID:   req.ParentID,
		Name:       req.Name,
		Path:       req.Path,
		Icon:       req.Icon,
		Component:  req.Component,
		Permission: req.Permission,
		Type:       req.Type,
		Sort:       req.Sort,
		Status:     1,
	}

	if err := database.DB.Create(&menu).Error; err != nil {
		response.InternalServerError(c, "创建菜单失败")
		return
	}

	response.SuccessWithMessage(c, "创建成功", menu)
}

func UpdateMenu(c *gin.Context) {
	idStr := c.Param("id")
	menuID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var menu model.Menu
	if err := database.DB.First(&menu, menuID).Error; err != nil {
		response.NotFound(c, "菜单不存在")
		return
	}

	var req struct {
		Name       string `json:"name"`
		Path       string `json:"path"`
		Icon       string `json:"icon"`
		Component  string `json:"component"`
		Permission string `json:"permission"`
		Sort       int    `json:"sort"`
		Status     int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.Name != "" {
		menu.Name = req.Name
	}
	if req.Path != "" {
		menu.Path = req.Path
	}
	if req.Icon != "" {
		menu.Icon = req.Icon
	}
	if req.Component != "" {
		menu.Component = req.Component
	}
	if req.Permission != "" {
		menu.Permission = req.Permission
	}
	if req.Sort >= 0 {
		menu.Sort = req.Sort
	}
	if req.Status >= 0 {
		menu.Status = req.Status
	}

	database.DB.Save(&menu)
	response.SuccessWithMessage(c, "更新成功", menu)
}

func DeleteMenu(c *gin.Context) {
	idStr := c.Param("id")
	menuID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	result := database.DB.Delete(&model.Menu{}, menuID)
	if result.RowsAffected == 0 {
		response.NotFound(c, "菜单不存在")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

func GetPermissionList(c *gin.Context) {
	var permissions []model.Permission
	database.DB.Order("id ASC").Find(&permissions)
	response.Success(c, permissions)
}

func CreatePermission(c *gin.Context) {
	var req struct {
		MenuID uint   `json:"menu_id"`
		Name   string `json:"name" binding:"required"`
		Code   string `json:"code" binding:"required"`
		Type   string `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	permission := model.Permission{
		MenuID: req.MenuID,
		Name:   req.Name,
		Code:   req.Code,
		Type:   req.Type,
		Status: 1,
	}

	if err := database.DB.Create(&permission).Error; err != nil {
		response.InternalServerError(c, "创建权限失败")
		return
	}

	response.SuccessWithMessage(c, "创建成功", permission)
}

func UpdatePermission(c *gin.Context) {
	idStr := c.Param("id")
	permissionID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var permission model.Permission
	if err := database.DB.First(&permission, permissionID).Error; err != nil {
		response.NotFound(c, "权限不存在")
		return
	}

	var req struct {
		Name   string `json:"name"`
		Code   string `json:"code"`
		Type   string `json:"type"`
		Status int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.Name != "" {
		permission.Name = req.Name
	}
	if req.Code != "" {
		permission.Code = req.Code
	}
	if req.Type != "" {
		permission.Type = req.Type
	}
	if req.Status >= 0 {
		permission.Status = req.Status
	}

	database.DB.Save(&permission)
	response.SuccessWithMessage(c, "更新成功", permission)
}

func DeletePermission(c *gin.Context) {
	idStr := c.Param("id")
	permissionID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	result := database.DB.Delete(&model.Permission{}, permissionID)
	if result.RowsAffected == 0 {
		response.NotFound(c, "权限不存在")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

func GetRoleList(c *gin.Context) {
	var roles []model.Role
	database.DB.Order("id ASC").Find(&roles)
	response.Success(c, roles)
}

func CreateRole(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	role := model.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      1,
	}

	if err := database.DB.Create(&role).Error; err != nil {
		response.InternalServerError(c, "创建角色失败")
		return
	}

	response.SuccessWithMessage(c, "创建成功", role)
}

func UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	roleID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var role model.Role
	if err := database.DB.First(&role, roleID).Error; err != nil {
		response.NotFound(c, "角色不存在")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.Status >= 0 {
		role.Status = req.Status
	}

	database.DB.Save(&role)
	response.SuccessWithMessage(c, "更新成功", role)
}

func DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	roleID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	result := database.DB.Delete(&model.Role{}, roleID)
	if result.RowsAffected == 0 {
		response.NotFound(c, "角色不存在")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

func GetAnnouncementList(c *gin.Context) {
	var announcements []model.Announcement
	database.DB.Order("created_at DESC").Find(&announcements)
	response.Success(c, announcements)
}

func CreateAnnouncement(c *gin.Context) {
	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Type    string `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	announcement := model.Announcement{
		Title:   req.Title,
		Content: req.Content,
		Type:    req.Type,
		Status:  1,
	}

	if err := database.DB.Create(&announcement).Error; err != nil {
		response.InternalServerError(c, "创建公告失败")
		return
	}

	response.SuccessWithMessage(c, "创建成功", announcement)
}

func UpdateAnnouncement(c *gin.Context) {
	idStr := c.Param("id")
	announcementID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var announcement model.Announcement
	if err := database.DB.First(&announcement, announcementID).Error; err != nil {
		response.NotFound(c, "公告不存在")
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Type    string `json:"type"`
		Status  int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.Title != "" {
		announcement.Title = req.Title
	}
	if req.Content != "" {
		announcement.Content = req.Content
	}
	if req.Type != "" {
		announcement.Type = req.Type
	}
	if req.Status >= 0 {
		announcement.Status = req.Status
	}

	database.DB.Save(&announcement)
	response.SuccessWithMessage(c, "更新成功", announcement)
}

func DeleteAnnouncement(c *gin.Context) {
	idStr := c.Param("id")
	announcementID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	result := database.DB.Delete(&model.Announcement{}, announcementID)
	if result.RowsAffected == 0 {
		response.NotFound(c, "公告不存在")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

func GetAdminNews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	var news []model.News

	query := database.DB.Model(&model.News{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&news)

	response.SuccessPage(c, news, total, page, pageSize)
}

func CreateNews(c *gin.Context) {
	var req struct {
		Title       string `json:"title" binding:"required"`
		Author      string `json:"author"`
		Category    string `json:"category"`
		Content     string `json:"content" binding:"required"`
		Cover       string `json:"cover"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	news := model.News{
		Title:       req.Title,
		Author:      req.Author,
		Category:    req.Category,
		Content:     req.Content,
		Cover:       req.Cover,
		Description: req.Description,
		Status:      1,
	}

	if err := database.DB.Create(&news).Error; err != nil {
		response.InternalServerError(c, "创建资讯失败")
		return
	}

	response.SuccessWithMessage(c, "创建成功", news)
}

func UpdateNews(c *gin.Context) {
	idStr := c.Param("id")
	newsID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var news model.News
	if err := database.DB.First(&news, newsID).Error; err != nil {
		response.NotFound(c, "资讯不存在")
		return
	}

	var req struct {
		Title       string `json:"title"`
		Author      string `json:"author"`
		Category    string `json:"category"`
		Content     string `json:"content"`
		Cover       string `json:"cover"`
		Description string `json:"description"`
		Status      int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.Title != "" {
		news.Title = req.Title
	}
	if req.Author != "" {
		news.Author = req.Author
	}
	if req.Category != "" {
		news.Category = req.Category
	}
	if req.Content != "" {
		news.Content = req.Content
	}
	if req.Cover != "" {
		news.Cover = req.Cover
	}
	if req.Description != "" {
		news.Description = req.Description
	}
	if req.Status >= 0 {
		news.Status = req.Status
	}

	database.DB.Save(&news)
	response.SuccessWithMessage(c, "更新成功", news)
}

func DeleteNews(c *gin.Context) {
	idStr := c.Param("id")
	newsID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	result := database.DB.Delete(&model.News{}, newsID)
	if result.RowsAffected == 0 {
		response.NotFound(c, "资讯不存在")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}
