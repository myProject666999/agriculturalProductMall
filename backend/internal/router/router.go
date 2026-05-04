package router

import (
	"agricultural-product-mall/internal/handler"
	"agricultural-product-mall/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.LoggerMiddleware())
	r.Use(gin.Recovery())
	r.Use(middleware.CorsMiddleware())

	api := r.Group("/api")
	{
		public := api.Group("")
		{
			public.POST("/register", handler.Register)
			public.POST("/login", handler.Login)
			public.POST("/send-code", handler.SendVerificationCode)
			public.GET("/products", handler.GetProductList)
			public.GET("/products/:id", handler.GetProductDetail)
			public.GET("/categories", handler.GetCategoryList)
			public.GET("/categories/:id/products", handler.GetProductsByCategory)
			public.GET("/banners", handler.GetBannerList)
			public.GET("/news", handler.GetNewsList)
			public.GET("/news/:id", handler.GetNewsDetail)
			public.GET("/search", handler.SearchProducts)
			public.GET("/recommend", handler.GetRecommendProducts)
		}

		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/user/info", handler.GetUserInfo)
			auth.PUT("/user/info", handler.UpdateUserInfo)
			auth.PUT("/user/password", handler.UpdatePassword)

			auth.GET("/addresses", handler.GetAddressList)
			auth.POST("/addresses", handler.CreateAddress)
			auth.PUT("/addresses/:id", handler.UpdateAddress)
			auth.DELETE("/addresses/:id", handler.DeleteAddress)
			auth.PUT("/addresses/:id/default", handler.SetDefaultAddress)

			auth.GET("/favorites", handler.GetFavoriteList)
			auth.POST("/favorites", handler.AddFavorite)
			auth.DELETE("/favorites/:id", handler.RemoveFavorite)

			auth.GET("/cart", handler.GetCartList)
			auth.POST("/cart", handler.AddToCart)
			auth.PUT("/cart/:id", handler.UpdateCartItem)
			auth.DELETE("/cart/:id", handler.RemoveFromCart)
			auth.POST("/cart/clear", handler.ClearCart)

			auth.GET("/orders", handler.GetOrderList)
			auth.POST("/orders", handler.CreateOrder)
			auth.GET("/orders/:id", handler.GetOrderDetail)
			auth.PUT("/orders/:id/cancel", handler.CancelOrder)
			auth.PUT("/orders/:id/pay", handler.PayOrder)
			auth.PUT("/orders/:id/confirm", handler.ConfirmOrder)
			auth.POST("/orders/:id/refund", handler.ApplyRefund)
			auth.POST("/orders/:id/review", handler.CreateReview)
			auth.GET("/logistics/:order_id", handler.GetLogisticsInfo)

			auth.GET("/reviews", handler.GetReviewList)
		}

		merchant := api.Group("/merchant")
		merchant.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("merchant"))
		{
			merchant.GET("/products", handler.GetMerchantProducts)
			merchant.POST("/products", handler.CreateProduct)
			merchant.PUT("/products/:id", handler.UpdateProduct)
			merchant.DELETE("/products/:id", handler.DeleteProduct)
			merchant.PUT("/products/:id/status", handler.UpdateProductStatus)

			merchant.GET("/inventory", handler.GetInventoryList)
			merchant.POST("/inventory/in", handler.InventoryIn)
			merchant.POST("/inventory/out", handler.InventoryOut)
			merchant.GET("/inventory/records", handler.GetInventoryRecords)

			merchant.GET("/orders", handler.GetMerchantOrders)
			merchant.GET("/orders/:id", handler.GetMerchantOrderDetail)
			merchant.PUT("/orders/:id/ship", handler.ShipOrder)
			merchant.PUT("/orders/:id/complete", handler.CompleteOrder)
			merchant.PUT("/orders/:id/refund", handler.HandleRefund)

			merchant.GET("/reviews", handler.GetMerchantReviews)
			merchant.GET("/sales/stats", handler.GetSalesStats)

			merchant.GET("/info", handler.GetMerchantInfo)
			merchant.PUT("/info", handler.UpdateMerchantInfo)
			merchant.PUT("/password", handler.UpdateMerchantPassword)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
		{
			admin.GET("/stats/sales", handler.GetAdminSalesStats)
			admin.GET("/stats/users", handler.GetAdminUserStats)
			admin.GET("/stats/categories", handler.GetCategoryStats)
			admin.GET("/stats/products", handler.GetProductStats)

			admin.GET("/products", handler.GetAdminProducts)
			admin.POST("/products", handler.CreateAdminProduct)
			admin.PUT("/products/:id", handler.UpdateAdminProduct)
			admin.DELETE("/products/:id", handler.DeleteAdminProduct)
			admin.PUT("/products/:id/status", handler.UpdateAdminProductStatus)

			admin.GET("/categories", handler.GetAdminCategories)
			admin.POST("/categories", handler.CreateCategory)
			admin.PUT("/categories/:id", handler.UpdateCategory)
			admin.DELETE("/categories/:id", handler.DeleteCategory)

			admin.GET("/orders", handler.GetAdminOrders)
			admin.GET("/orders/:id", handler.GetAdminOrderDetail)
			admin.PUT("/orders/:id/status", handler.UpdateOrderStatus)

			admin.GET("/cart", handler.GetAdminCartList)

			admin.GET("/reviews", handler.GetAdminReviews)
			admin.DELETE("/reviews/:id", handler.DeleteReview)

			admin.GET("/logistics", handler.GetAdminLogistics)
			admin.POST("/logistics", handler.CreateLogistics)
			admin.PUT("/logistics/:id", handler.UpdateLogistics)
			admin.DELETE("/logistics/:id", handler.DeleteLogistics)

			admin.GET("/inventory", handler.GetAdminInventory)
			admin.GET("/inventory/records", handler.GetAdminInventoryRecords)

			admin.GET("/banners", handler.GetAdminBanners)
			admin.POST("/banners", handler.CreateBanner)
			admin.PUT("/banners/:id", handler.UpdateBanner)
			admin.DELETE("/banners/:id", handler.DeleteBanner)

			admin.GET("/users", handler.GetAdminUsers)
			admin.POST("/users", handler.CreateUser)
			admin.PUT("/users/:id", handler.UpdateUser)
			admin.DELETE("/users/:id", handler.DeleteUser)
			admin.PUT("/users/:id/status", handler.UpdateUserStatus)

			admin.GET("/menus", handler.GetMenuList)
			admin.POST("/menus", handler.CreateMenu)
			admin.PUT("/menus/:id", handler.UpdateMenu)
			admin.DELETE("/menus/:id", handler.DeleteMenu)

			admin.GET("/permissions", handler.GetPermissionList)
			admin.POST("/permissions", handler.CreatePermission)
			admin.PUT("/permissions/:id", handler.UpdatePermission)
			admin.DELETE("/permissions/:id", handler.DeletePermission)
			admin.GET("/roles", handler.GetRoleList)
			admin.POST("/roles", handler.CreateRole)
			admin.PUT("/roles/:id", handler.UpdateRole)
			admin.DELETE("/roles/:id", handler.DeleteRole)

			admin.GET("/announcements", handler.GetAnnouncementList)
			admin.POST("/announcements", handler.CreateAnnouncement)
			admin.PUT("/announcements/:id", handler.UpdateAnnouncement)
			admin.DELETE("/announcements/:id", handler.DeleteAnnouncement)

			admin.GET("/news", handler.GetAdminNews)
			admin.POST("/news", handler.CreateNews)
			admin.PUT("/news/:id", handler.UpdateNews)
			admin.DELETE("/news/:id", handler.DeleteNews)
		}
	}

	return r
}
