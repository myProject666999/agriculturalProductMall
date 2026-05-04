package handler

import (
	"agricultural-product-mall/internal/model"
	"agricultural-product-mall/internal/service"
	"agricultural-product-mall/pkg/database"
	"agricultural-product-mall/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAddressList(c *gin.Context) {
	userID := c.GetUint("user_id")

	var addresses []model.Address
	database.DB.Where("user_id = ?", userID).
		Order("is_default DESC, created_at DESC").
		Find(&addresses)

	response.Success(c, addresses)
}

type CreateAddressRequest struct {
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Province string `json:"province"`
	City     string `json:"city"`
	District string `json:"district"`
	Detail   string `json:"detail" binding:"required"`
	IsDefault int    `json:"is_default"`
}

func CreateAddress(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.IsDefault == 1 {
		database.DB.Model(&model.Address{}).
			Where("user_id = ?", userID).
			Update("is_default", 0)
	}

	address := model.Address{
		UserID:    userID,
		Name:      req.Name,
		Phone:     req.Phone,
		Province:  req.Province,
		City:      req.City,
		District:  req.District,
		Detail:    req.Detail,
		IsDefault: req.IsDefault,
	}

	if err := database.DB.Create(&address).Error; err != nil {
		response.InternalServerError(c, "添加地址失败")
		return
	}

	response.SuccessWithMessage(c, "添加成功", address)
}

func UpdateAddress(c *gin.Context) {
	userID := c.GetUint("user_id")
	idStr := c.Param("id")
	addressID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var address model.Address
	if err := database.DB.Where("id = ? AND user_id = ?", addressID, userID).
		First(&address).Error; err != nil {
		response.NotFound(c, "地址不存在")
		return
	}

	if req.IsDefault == 1 && address.IsDefault != 1 {
		database.DB.Model(&model.Address{}).
			Where("user_id = ?", userID).
			Update("is_default", 0)
	}

	address.Name = req.Name
	address.Phone = req.Phone
	address.Province = req.Province
	address.City = req.City
	address.District = req.District
	address.Detail = req.Detail
	address.IsDefault = req.IsDefault

	if err := database.DB.Save(&address).Error; err != nil {
		response.InternalServerError(c, "更新地址失败")
		return
	}

	response.SuccessWithMessage(c, "更新成功", address)
}

func DeleteAddress(c *gin.Context) {
	userID := c.GetUint("user_id")
	idStr := c.Param("id")
	addressID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	result := database.DB.Where("id = ? AND user_id = ?", addressID, userID).
		Delete(&model.Address{})

	if result.RowsAffected == 0 {
		response.NotFound(c, "地址不存在")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

func SetDefaultAddress(c *gin.Context) {
	userID := c.GetUint("user_id")
	idStr := c.Param("id")
	addressID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var address model.Address
	if err := database.DB.Where("id = ? AND user_id = ?", addressID, userID).
		First(&address).Error; err != nil {
		response.NotFound(c, "地址不存在")
		return
	}

	database.DB.Model(&model.Address{}).
		Where("user_id = ?", userID).
		Update("is_default", 0)

	address.IsDefault = 1
	database.DB.Save(&address)

	response.SuccessWithMessage(c, "设置成功", address)
}

func GetFavoriteList(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	var favorites []model.Favorite

	query := database.DB.Model(&model.Favorite{}).Where("user_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("Product").Preload("Product.Category").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&favorites)

	response.SuccessPage(c, favorites, total, page, pageSize)
}

type AddFavoriteRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
}

func AddFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req AddFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var product model.Product
	if err := database.DB.First(&product, req.ProductID).Error; err != nil {
		response.NotFound(c, "商品不存在")
		return
	}

	var existingFavorite model.Favorite
	result := database.DB.Where("user_id = ? AND product_id = ?", userID, req.ProductID).
		First(&existingFavorite)

	if result.Error == nil {
		response.SuccessWithMessage(c, "已收藏", existingFavorite)
		return
	}

	favorite := model.Favorite{
		UserID:    userID,
		ProductID: req.ProductID,
	}

	if err := database.DB.Create(&favorite).Error; err != nil {
		response.InternalServerError(c, "收藏失败")
		return
	}

	recService := service.NewRecommendationService()
	recService.RecordBehavior(userID, req.ProductID, "collect", 0)

	response.SuccessWithMessage(c, "收藏成功", favorite)
}

func RemoveFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")
	idStr := c.Param("id")
	favoriteID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	result := database.DB.Where("id = ? AND user_id = ?", favoriteID, userID).
		Delete(&model.Favorite{})

	if result.RowsAffected == 0 {
		response.NotFound(c, "收藏不存在")
		return
	}

	response.SuccessWithMessage(c, "取消收藏成功", nil)
}

func GetReviewList(c *gin.Context) {
	userID := c.GetUint("user_id")
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

	query := database.DB.Model(&model.Review{}).Where("user_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("Product").Preload("User").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&reviews)

	response.SuccessPage(c, reviews, total, page, pageSize)
}
