package handler

import (
	"agricultural-product-mall/internal/model"
	"agricultural-product-mall/internal/service"
	"agricultural-product-mall/pkg/database"
	"agricultural-product-mall/pkg/response"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetProductList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	categoryID := c.Query("category_id")
	sort := c.DefaultQuery("sort", "default")
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	var products []model.Product

	query := database.DB.Model(&model.Product{}).Where("status = ?", 1)

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize

	switch sort {
	case "price_asc":
		query = query.Order("price ASC")
	case "price_desc":
		query = query.Order("price DESC")
	case "sales":
		query = query.Order("sales DESC")
	case "rating":
		query = query.Order("rating DESC")
	default:
		query = query.Order("is_recommend DESC, sales DESC")
	}

	query.Preload("Category").Preload("Merchant").
		Offset(offset).Limit(pageSize).
		Find(&products)

	response.SuccessPage(c, products, total, page, pageSize)
}

func GetProductDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var product model.Product
	if err := database.DB.Preload("Category").Preload("Merchant").
		First(&product, id).Error; err != nil {
		response.NotFound(c, "商品不存在")
		return
	}

	if product.Status != 1 {
		response.NotFound(c, "商品不存在或已下架")
		return
	}

	userID, exists := c.Get("user_id")
	if exists {
		recService := service.NewRecommendationService()
		recService.RecordBehavior(userID.(uint), product.ID, "view", 0)
	}

	response.Success(c, product)
}

func GetCategoryList(c *gin.Context) {
	var categories []model.Category
	database.DB.Where("status = ?", 1).
		Order("sort ASC, id ASC").
		Find(&categories)

	response.Success(c, categories)
}

func GetProductsByCategory(c *gin.Context) {
	idStr := c.Param("id")
	categoryID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sort := c.DefaultQuery("sort", "default")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	var products []model.Product

	query := database.DB.Model(&model.Product{}).
		Where("status = ? AND category_id = ?", 1, categoryID)

	query.Count(&total)

	offset := (page - 1) * pageSize

	switch sort {
	case "price_asc":
		query = query.Order("price ASC")
	case "price_desc":
		query = query.Order("price DESC")
	case "sales":
		query = query.Order("sales DESC")
	default:
		query = query.Order("is_recommend DESC, sales DESC")
	}

	query.Preload("Category").Preload("Merchant").
		Offset(offset).Limit(pageSize).
		Find(&products)

	response.SuccessPage(c, products, total, page, pageSize)
}

func GetBannerList(c *gin.Context) {
	var banners []model.Banner
	database.DB.Where("status = ?", 1).
		Order("sort ASC, id ASC").
		Find(&banners)

	response.Success(c, banners)
}

func GetNewsList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	category := c.Query("category")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	var news []model.News

	query := database.DB.Model(&model.News{}).Where("status = ?", 1)

	if category != "" {
		query = query.Where("category = ?", category)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("sort ASC, created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&news)

	response.SuccessPage(c, news, total, page, pageSize)
}

func GetNewsDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var news model.News
	if err := database.DB.First(&news, id).Error; err != nil {
		response.NotFound(c, "资讯不存在")
		return
	}

	if news.Status != 1 {
		response.NotFound(c, "资讯不存在")
		return
	}

	news.Views += 1
	database.DB.Save(&news)

	response.Success(c, news)
}

func SearchProducts(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.BadRequest(c, "请输入搜索关键词")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sort := c.DefaultQuery("sort", "default")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	keywords := strings.Fields(keyword)
	var total int64
	var products []model.Product

	query := database.DB.Model(&model.Product{}).Where("status = ?", 1)

	if len(keywords) > 0 {
		conditions := make([]string, 0)
		params := make([]interface{}, 0)
		for _, kw := range keywords {
			conditions = append(conditions, "(name LIKE ? OR description LIKE ?)")
			params = append(params, "%"+kw+"%", "%"+kw+"%")
		}
		query = query.Where(strings.Join(conditions, " AND "), params...)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize

	switch sort {
	case "price_asc":
		query = query.Order("price ASC")
	case "price_desc":
		query = query.Order("price DESC")
	case "sales":
		query = query.Order("sales DESC")
	default:
		query = query.Order("is_recommend DESC, sales DESC")
	}

	query.Preload("Category").Preload("Merchant").
		Offset(offset).Limit(pageSize).
		Find(&products)

	response.SuccessPage(c, products, total, page, pageSize)
}

func GetRecommendProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	var userID uint
	userIDVal, exists := c.Get("user_id")
	if exists {
		userID = userIDVal.(uint)
	}

	recService := service.NewRecommendationService()
	products, err := recService.GetRecommendProducts(userID, limit)
	if err != nil {
		var hotProducts []model.Product
		database.DB.Where("status = ?", 1).
			Order("sales DESC, rating DESC").
			Limit(limit).
			Find(&hotProducts)
		response.Success(c, hotProducts)
		return
	}

	response.Success(c, products)
}
