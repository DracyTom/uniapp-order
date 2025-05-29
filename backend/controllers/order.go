package controllers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"uniapp-order-backend/models"
	"uniapp-order-backend/utils"
)

type OrderController struct{}

// 获取订单列表
func (oc *OrderController) GetOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	paymentStatus := c.Query("payment_status")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	offset := (page - 1) * limit
	session := models.DB.NewSession()
	defer session.Close()

	if status != "" {
		session.Where("status = ?", status)
	}
	if paymentStatus != "" {
		session.Where("payment_status = ?", paymentStatus)
	}
	if startDate != "" {
		session.Where("delivery_time >= ?", startDate)
	}
	if endDate != "" {
		session.Where("delivery_time <= ?", endDate)
	}

	var orders []models.Order
	total, err := session.Desc("created_at").Limit(limit, offset).FindAndCount(&orders)
	if err != nil {
		utils.Error(c, 1, "获取订单列表失败")
		return
	}

	// 加载关联数据
	for i := range orders {
		// 加载客户信息
		var customer models.Customer
		has, _ := models.DB.ID(orders[i].CustomerID).Get(&customer)
		if has {
			orders[i].Customer = &customer
		}

		// 加载订单商品
		var orderItems []models.OrderItem
		err = models.DB.Where("order_id = ?", orders[i].ID).Find(&orderItems)
		if err == nil {
			// 加载商品信息
			for j := range orderItems {
				var product models.Product
				has, _ := models.DB.ID(orderItems[j].ProductID).Get(&product)
				if has {
					orderItems[j].Product = &product
				}
			}
			orders[i].OrderItems = orderItems
		}
	}

	utils.Success(c, gin.H{
		"orders": orders,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

// 获取单个订单
func (oc *OrderController) GetOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, 1, "无效的订单ID")
		return
	}

	var order models.Order
	has, err := models.DB.ID(id).Get(&order)
	if err != nil {
		utils.Error(c, 1, "获取订单失败")
		return
	}
	if !has {
		utils.Error(c, 1, "订单不存在")
		return
	}

	// 加载客户信息
	var customer models.Customer
	has, _ = models.DB.ID(order.CustomerID).Get(&customer)
	if has {
		order.Customer = &customer
	}

	// 加载订单商品
	var orderItems []models.OrderItem
	err = models.DB.Where("order_id = ?", order.ID).Find(&orderItems)
	if err == nil {
		// 加载商品信息
		for i := range orderItems {
			var product models.Product
			has, _ := models.DB.ID(orderItems[i].ProductID).Get(&product)
			if has {
				orderItems[i].Product = &product
			}
		}
		order.OrderItems = orderItems
	}

	utils.Success(c, order)
}

// 创建订单
func (oc *OrderController) CreateOrder(c *gin.Context) {
	var req struct {
		models.Order
		Items []models.OrderItem `json:"items"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 1, "参数错误: "+err.Error())
		return
	}

	session := models.DB.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		utils.Error(c, 1, "事务开始失败")
		return
	}

	// 计算总金额
	var totalAmount float64
	for i := range req.Items {
		// 获取商品当前价格
		var product models.Product
		has, err := session.ID(req.Items[i].ProductID).Get(&product)
		if err != nil || !has {
			session.Rollback()
			utils.Error(c, 1, "商品不存在")
			return
		}
		req.Items[i].Price = product.Price
		req.Items[i].Amount = product.Price * float64(req.Items[i].Quantity)
		totalAmount += req.Items[i].Amount
	}
	req.Order.TotalAmount = totalAmount

	// 创建订单
	_, err := session.Insert(&req.Order)
	if err != nil {
		session.Rollback()
		utils.Error(c, 1, "创建订单失败")
		return
	}

	// 创建订单商品
	for i := range req.Items {
		req.Items[i].OrderID = req.Order.ID
	}
	_, err = session.Insert(&req.Items)
	if err != nil {
		session.Rollback()
		utils.Error(c, 1, "创建订单商品失败")
		return
	}

	if err := session.Commit(); err != nil {
		utils.Error(c, 1, "事务提交失败")
		return
	}

	utils.Success(c, gin.H{
		"order":       req.Order,
		"order_items": req.Items,
	})
}

// 更新订单
func (oc *OrderController) UpdateOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, 1, "无效的订单ID")
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.Error(c, 1, "参数错误")
		return
	}

	// 添加更新时间
	updateData["updated_at"] = time.Now()

	_, err = models.DB.Table("order").Where("i_d = ?", id).Update(updateData)
	if err != nil {
		utils.Error(c, 1, "更新订单失败")
		return
	}

	// 获取更新后的订单
	var order models.Order
	has, err := models.DB.ID(id).Get(&order)
	if err != nil || !has {
		utils.Error(c, 1, "获取更新后的订单失败")
		return
	}

	utils.Success(c, order)
}

// 删除订单
func (oc *OrderController) DeleteOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, 1, "无效的订单ID")
		return
	}

	session := models.DB.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		utils.Error(c, 1, "事务开始失败")
		return
	}

	// 删除订单商品
	_, err = session.Where("order_i_d = ?", id).Delete(&models.OrderItem{})
	if err != nil {
		session.Rollback()
		utils.Error(c, 1, "删除订单商品失败")
		return
	}

	// 删除订单
	_, err = session.ID(id).Delete(&models.Order{})
	if err != nil {
		session.Rollback()
		utils.Error(c, 1, "删除订单失败")
		return
	}

	if err := session.Commit(); err != nil {
		utils.Error(c, 1, "事务提交失败")
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}

// 获取日历数据
func (oc *OrderController) GetCalendarData(c *gin.Context) {
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month()))))

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	// 获取该月的订单统计
	var orderStats []struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
		Amount float64 `json:"amount"`
	}

	sql := `
		SELECT 
			DATE(delivery_time) as date,
			COUNT(*) as count,
			SUM(total_amount) as amount
		FROM "order" 
		WHERE delivery_time >= ? AND delivery_time <= ?
		GROUP BY DATE(delivery_time)
		ORDER BY date
	`

	err := models.DB.SQL(sql, startDate, endDate).Find(&orderStats)
	if err != nil {
		utils.Error(c, 1, "获取日历数据失败")
		return
	}

	// 获取节假日信息
	var holidays []models.Holiday
	err = models.DB.Where("date >= ? AND date <= ?", startDate, endDate).Find(&holidays)
	if err != nil {
		utils.Error(c, 1, "获取节假日信息失败")
		return
	}

	// 生成日历数据
	calendarData := make([]gin.H, 0)
	for d := startDate; d.Before(endDate.AddDate(0, 0, 1)); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		
		// 查找当天的订单统计
		var dayStats *struct {
			Date  string `json:"date"`
			Count int    `json:"count"`
			Amount float64 `json:"amount"`
		}
		for _, stat := range orderStats {
			if stat.Date == dateStr {
				dayStats = &stat
				break
			}
		}

		// 查找当天的节假日
		var dayHoliday *models.Holiday
		for _, holiday := range holidays {
			if holiday.Date.Format("2006-01-02") == dateStr {
				dayHoliday = &holiday
				break
			}
		}

		// 农历信息
		lunar := utils.ToLunar(d)

		dayData := gin.H{
			"date":     dateStr,
			"weekday":  int(d.Weekday()),
			"lunar":    lunar,
			"holiday":  dayHoliday,
		}

		if dayStats != nil {
			dayData["orders"] = gin.H{
				"count":  dayStats.Count,
				"amount": dayStats.Amount,
			}
		} else {
			dayData["orders"] = gin.H{
				"count":  0,
				"amount": 0,
			}
		}

		calendarData = append(calendarData, dayData)
	}

	utils.Success(c, gin.H{
		"year":     year,
		"month":    month,
		"calendar": calendarData,
	})
}