package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"uniapp-order-backend/models"
	"uniapp-order-backend/utils"
)

type CustomerController struct{}

// 获取客户列表
func (cc *CustomerController) GetCustomers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	keyword := c.Query("keyword")

	offset := (page - 1) * limit
	session := models.DB.NewSession()
	defer session.Close()

	if keyword != "" {
		session.Where("name LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var customers []models.Customer
	total, err := session.Limit(limit, offset).FindAndCount(&customers)
	if err != nil {
		utils.Error(c, 1, "获取客户列表失败")
		return
	}

	utils.Success(c, gin.H{
		"customers": customers,
		"total":     total,
		"page":      page,
		"limit":     limit,
	})
}

// 获取单个客户
func (cc *CustomerController) GetCustomer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, 1, "无效的客户ID")
		return
	}

	var customer models.Customer
	has, err := models.DB.ID(id).Get(&customer)
	if err != nil {
		utils.Error(c, 1, "获取客户失败")
		return
	}
	if !has {
		utils.Error(c, 1, "客户不存在")
		return
	}

	// 获取客户地址
	var addresses []models.CustomerAddress
	err = models.DB.Where("customer_id = ?", id).Find(&addresses)
	if err != nil {
		utils.Error(c, 1, "获取客户地址失败")
		return
	}

	utils.Success(c, gin.H{
		"customer":  customer,
		"addresses": addresses,
	})
}

// 创建客户
func (cc *CustomerController) CreateCustomer(c *gin.Context) {
	var req struct {
		Customer  models.Customer          `json:"customer"`
		Addresses []models.CustomerAddress `json:"addresses"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 1, "参数错误")
		return
	}

	session := models.DB.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		utils.Error(c, 1, "事务开始失败")
		return
	}

	// 创建客户
	_, err := session.Insert(&req.Customer)
	if err != nil {
		session.Rollback()
		utils.Error(c, 1, "创建客户失败")
		return
	}

	// 创建地址
	for i := range req.Addresses {
		req.Addresses[i].CustomerID = req.Customer.ID
	}
	if len(req.Addresses) > 0 {
		_, err = session.Insert(&req.Addresses)
		if err != nil {
			session.Rollback()
			utils.Error(c, 1, "创建客户地址失败")
			return
		}
	}

	if err := session.Commit(); err != nil {
		utils.Error(c, 1, "事务提交失败")
		return
	}

	utils.Success(c, gin.H{
		"customer":  req.Customer,
		"addresses": req.Addresses,
	})
}

// 更新客户
func (cc *CustomerController) UpdateCustomer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, 1, "无效的客户ID")
		return
	}

	var customer models.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		utils.Error(c, 1, "参数错误")
		return
	}

	customer.ID = id
	_, err = models.DB.ID(id).Update(&customer)
	if err != nil {
		utils.Error(c, 1, "更新客户失败")
		return
	}

	utils.Success(c, customer)
}

// 删除客户
func (cc *CustomerController) DeleteCustomer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, 1, "无效的客户ID")
		return
	}

	session := models.DB.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		utils.Error(c, 1, "事务开始失败")
		return
	}

	// 删除客户地址
	_, err = session.Where("customer_id = ?", id).Delete(&models.CustomerAddress{})
	if err != nil {
		session.Rollback()
		utils.Error(c, 1, "删除客户地址失败")
		return
	}

	// 删除客户
	_, err = session.ID(id).Delete(&models.Customer{})
	if err != nil {
		session.Rollback()
		utils.Error(c, 1, "删除客户失败")
		return
	}

	if err := session.Commit(); err != nil {
		utils.Error(c, 1, "事务提交失败")
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}

// 获取客户地址
func (cc *CustomerController) GetCustomerAddresses(c *gin.Context) {
	customerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, 1, "无效的客户ID")
		return
	}

	var addresses []models.CustomerAddress
	err = models.DB.Where("customer_id = ?", customerID).Find(&addresses)
	if err != nil {
		utils.Error(c, 1, "获取客户地址失败")
		return
	}

	utils.Success(c, addresses)
}

// 添加客户地址
func (cc *CustomerController) AddCustomerAddress(c *gin.Context) {
	customerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, 1, "无效的客户ID")
		return
	}

	var address models.CustomerAddress
	if err := c.ShouldBindJSON(&address); err != nil {
		utils.Error(c, 1, "参数错误")
		return
	}

	address.CustomerID = customerID
	_, err = models.DB.Insert(&address)
	if err != nil {
		utils.Error(c, 1, "添加地址失败")
		return
	}

	utils.Success(c, address)
}