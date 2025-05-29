package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"uniapp-order-backend/controllers"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// CORS配置
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	r.Use(cors.New(config))

	// 控制器实例
	productController := &controllers.ProductController{}
	customerController := &controllers.CustomerController{}
	orderController := &controllers.OrderController{}
	dashboardController := &controllers.DashboardController{}

	// API路由组
	api := r.Group("/api")
	{
		// 商品管理
		products := api.Group("/products")
		{
			products.GET("", productController.GetProducts)
			products.GET("/:id", productController.GetProduct)
			products.POST("", productController.CreateProduct)
			products.PUT("/:id", productController.UpdateProduct)
			products.DELETE("/:id", productController.DeleteProduct)
			products.GET("/categories", productController.GetCategories)
		}

		// 客户管理
		customers := api.Group("/customers")
		{
			customers.GET("", customerController.GetCustomers)
			customers.GET("/:id", customerController.GetCustomer)
			customers.POST("", customerController.CreateCustomer)
			customers.PUT("/:id", customerController.UpdateCustomer)
			customers.DELETE("/:id", customerController.DeleteCustomer)
			customers.GET("/:id/addresses", customerController.GetCustomerAddresses)
			customers.POST("/:id/addresses", customerController.AddCustomerAddress)
		}

		// 订单管理
		orders := api.Group("/orders")
		{
			orders.GET("", orderController.GetOrders)
			orders.GET("/:id", orderController.GetOrder)
			orders.POST("", orderController.CreateOrder)
			orders.PUT("/:id", orderController.UpdateOrder)
			orders.DELETE("/:id", orderController.DeleteOrder)
			orders.GET("/calendar", orderController.GetCalendarData)
		}

		// 驾驶仓
		api.GET("/dashboard", dashboardController.GetDashboardData)
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}