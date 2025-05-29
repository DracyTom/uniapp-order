package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"uniapp-order-backend/models"
	"uniapp-order-backend/utils"
)

type DashboardController struct{}

// 获取驾驶仓数据
func (dc *DashboardController) GetDashboardData(c *gin.Context) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	thisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	thisYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())

	// 今日订单统计
	var todayStats struct {
		OrderCount int     `json:"order_count"`
		TotalAmount float64 `json:"total_amount"`
		PaidAmount  float64 `json:"paid_amount"`
	}
	
	models.DB.SQL(`
		SELECT 
			COUNT(*) as order_count,
			COALESCE(SUM(total_amount), 0) as total_amount,
			COALESCE(SUM(CASE WHEN payment_status = 'paid' THEN total_amount ELSE 0 END), 0) as paid_amount
		FROM "order" 
		WHERE delivery_time >= ?
	`, today).Get(&todayStats)

	// 本月订单统计
	var monthStats struct {
		OrderCount int     `json:"order_count"`
		TotalAmount float64 `json:"total_amount"`
		PaidAmount  float64 `json:"paid_amount"`
	}
	
	models.DB.SQL(`
		SELECT 
			COUNT(*) as order_count,
			COALESCE(SUM(total_amount), 0) as total_amount,
			COALESCE(SUM(CASE WHEN payment_status = 'paid' THEN total_amount ELSE 0 END), 0) as paid_amount
		FROM "order" 
		WHERE delivery_time >= ?
	`, thisMonth).Get(&monthStats)

	// 本年订单统计
	var yearStats struct {
		OrderCount int     `json:"order_count"`
		TotalAmount float64 `json:"total_amount"`
		PaidAmount  float64 `json:"paid_amount"`
	}
	
	models.DB.SQL(`
		SELECT 
			COUNT(*) as order_count,
			COALESCE(SUM(total_amount), 0) as total_amount,
			COALESCE(SUM(CASE WHEN payment_status = 'paid' THEN total_amount ELSE 0 END), 0) as paid_amount
		FROM "order" 
		WHERE delivery_time >= ?
	`, thisYear).Get(&yearStats)

	// 订单状态统计
	var statusStats []struct {
		Status string `json:"status"`
		Count  int    `json:"count"`
	}
	models.DB.SQL(`
		SELECT status, COUNT(*) as count 
		FROM "order" 
		GROUP BY status
	`).Find(&statusStats)

	// 热销商品统计（本月）
	var topProducts []struct {
		ProductID   int64   `json:"product_id"`
		ProductName string  `json:"product_name"`
		Quantity    int     `json:"quantity"`
		Amount      float64 `json:"amount"`
	}
	models.DB.SQL(`
		SELECT 
			oi.product_id,
			p.name as product_name,
			SUM(oi.quantity) as quantity,
			SUM(oi.amount) as amount
		FROM order_item oi
		JOIN "order" o ON oi.order_id = o.id
		JOIN product p ON oi.product_id = p.id
		WHERE o.delivery_time >= ?
		GROUP BY oi.product_id, p.name
		ORDER BY quantity DESC
		LIMIT 10
	`, thisMonth).Find(&topProducts)

	// 客户统计
	var customerStats struct {
		TotalCustomers int `json:"total_customers"`
		NewCustomers   int `json:"new_customers"` // 本月新增
	}
	models.DB.SQL(`SELECT COUNT(*) as total_customers FROM customer`).Get(&customerStats.TotalCustomers)
	models.DB.SQL(`SELECT COUNT(*) as new_customers FROM customer WHERE created_at >= ?`, thisMonth).Get(&customerStats.NewCustomers)

	// 商品统计
	var productStats struct {
		TotalProducts  int `json:"total_products"`
		ActiveProducts int `json:"active_products"`
	}
	models.DB.SQL(`SELECT COUNT(*) as total_products FROM product`).Get(&productStats.TotalProducts)
	models.DB.SQL(`SELECT COUNT(*) as active_products FROM product WHERE status = true`).Get(&productStats.ActiveProducts)

	// 近7天订单趋势
	var weeklyTrend []struct {
		Date        string  `json:"date"`
		OrderCount  int     `json:"order_count"`
		TotalAmount float64 `json:"total_amount"`
	}
	
	for i := 6; i >= 0; i-- {
		date := today.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		
		var dayStats struct {
			OrderCount  int     `json:"order_count"`
			TotalAmount float64 `json:"total_amount"`
		}
		
		models.DB.SQL(`
			SELECT 
				COUNT(*) as order_count,
				COALESCE(SUM(total_amount), 0) as total_amount
			FROM "order" 
			WHERE DATE(delivery_time) = ?
		`, dateStr).Get(&dayStats)
		
		weeklyTrend = append(weeklyTrend, struct {
			Date        string  `json:"date"`
			OrderCount  int     `json:"order_count"`
			TotalAmount float64 `json:"total_amount"`
		}{
			Date:        dateStr,
			OrderCount:  dayStats.OrderCount,
			TotalAmount: dayStats.TotalAmount,
		})
	}

	utils.Success(c, gin.H{
		"today": gin.H{
			"orders":       todayStats.OrderCount,
			"total_amount": todayStats.TotalAmount,
			"paid_amount":  todayStats.PaidAmount,
		},
		"month": gin.H{
			"orders":       monthStats.OrderCount,
			"total_amount": monthStats.TotalAmount,
			"paid_amount":  monthStats.PaidAmount,
		},
		"year": gin.H{
			"orders":       yearStats.OrderCount,
			"total_amount": yearStats.TotalAmount,
			"paid_amount":  yearStats.PaidAmount,
		},
		"order_status":    statusStats,
		"top_products":    topProducts,
		"customer_stats":  customerStats,
		"product_stats":   productStats,
		"weekly_trend":    weeklyTrend,
	})
}