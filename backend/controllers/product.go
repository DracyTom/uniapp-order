package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"uniapp-order-backend/models"
	"uniapp-order-backend/utils"
)

type ProductController struct{}

// 获取商品列表
func (pc *ProductController) GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	category := c.Query("category")
	status := c.Query("status")

	offset := (page - 1) * limit
	session := models.DB.NewSession()
	defer session.Close()

	if category != "" {
		session.Where("category = ?", category)
	}
	if status != "" {
		statusBool, _ := strconv.ParseBool(status)
		session.Where("status = ?", statusBool)
	}

	var products []models.Product
	total, err := session.Limit(limit, offset).FindAndCount(&products)
	if err != nil {
		utils.Error(c, 1, "获取商品列表失败")
		return
	}

	utils.Success(c, gin.H{
		"products": products,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

// 获取单个商品
func (pc *ProductController) GetProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, 1, "无效的商品ID")
		return
	}

	var product models.Product
	has, err := models.DB.ID(id).Get(&product)
	if err != nil {
		utils.Error(c, 1, "获取商品失败")
		return
	}
	if !has {
		utils.Error(c, 1, "商品不存在")
		return
	}

	utils.Success(c, product)
}

// 创建商品
func (pc *ProductController) CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		utils.Error(c, 1, "参数错误")
		return
	}

	_, err := models.DB.Insert(&product)
	if err != nil {
		utils.Error(c, 1, "创建商品失败")
		return
	}

	utils.Success(c, product)
}

// 更新商品
func (pc *ProductController) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, 1, "无效的商品ID")
		return
	}

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		utils.Error(c, 1, "参数错误")
		return
	}

	product.ID = id
	_, err = models.DB.ID(id).Update(&product)
	if err != nil {
		utils.Error(c, 1, "更新商品失败")
		return
	}

	utils.Success(c, product)
}

// 删除商品
func (pc *ProductController) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, 1, "无效的商品ID")
		return
	}

	_, err = models.DB.ID(id).Delete(&models.Product{})
	if err != nil {
		utils.Error(c, 1, "删除商品失败")
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}

// 获取商品分类列表
func (pc *ProductController) GetCategories(c *gin.Context) {
	var categories []string
	err := models.DB.Table("product").Distinct("category").Find(&categories)
	if err != nil {
		utils.Error(c, 1, "获取分类失败")
		return
	}

	utils.Success(c, categories)
}