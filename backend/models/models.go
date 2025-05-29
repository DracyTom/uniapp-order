package models

import (
	"time"
)

// 商品模型
type Product struct {
	ID          int64             `json:"id" xorm:"pk autoincr"`
	Name        string            `json:"name" xorm:"varchar(255) not null"`
	Category    string            `json:"category" xorm:"varchar(100)"`
	Price       float64           `json:"price" xorm:"decimal(10,2)"`
	Status      bool              `json:"status" xorm:"default(true)"` // 上架状态
	Image       string            `json:"image" xorm:"text"`
	CustomAttrs map[string]string `json:"custom_attrs" xorm:"json"` // 自定义属性
	CreatedAt   time.Time         `json:"created_at" xorm:"created"`
	UpdatedAt   time.Time         `json:"updated_at" xorm:"updated"`
}

// 客户模型
type Customer struct {
	ID        int64     `json:"id" xorm:"pk autoincr"`
	Name      string    `json:"name" xorm:"varchar(255) not null"`
	Phone     string    `json:"phone" xorm:"varchar(50)"`
	Email     string    `json:"email" xorm:"varchar(255)"`
	CreatedAt time.Time `json:"created_at" xorm:"created"`
	UpdatedAt time.Time `json:"updated_at" xorm:"updated"`
}

// 客户地址模型
type CustomerAddress struct {
	ID         int64  `json:"id" xorm:"pk autoincr"`
	CustomerID int64  `json:"customer_id" xorm:"index"`
	Address    string `json:"address" xorm:"text not null"`
	IsDefault  bool   `json:"is_default" xorm:"default(false)"`
	Label      string `json:"label" xorm:"varchar(50)"` // 地址标签，如"家"、"公司"等
}

// 订单模型
type Order struct {
	ID             int64       `json:"id" xorm:"pk autoincr"`
	CustomerID     int64       `json:"customer_id" xorm:"index"`
	CustomerName   string      `json:"customer_name" xorm:"varchar(255)"`
	DeliveryTime   time.Time   `json:"delivery_time"`
	DeliveryAddr   string      `json:"delivery_addr" xorm:"text"`
	Recipient      string      `json:"recipient" xorm:"varchar(255)"`
	TotalAmount    float64     `json:"total_amount" xorm:"decimal(10,2)"`
	PaymentStatus  string      `json:"payment_status" xorm:"varchar(20) default('unpaid')"` // paid, unpaid
	SpecialReq     string      `json:"special_req" xorm:"text"`
	Remarks        string      `json:"remarks" xorm:"text"`
	Status         string      `json:"status" xorm:"varchar(20) default('pending')"` // pending, confirmed, delivered, cancelled
	CreatedAt      time.Time   `json:"created_at" xorm:"created"`
	UpdatedAt      time.Time   `json:"updated_at" xorm:"updated"`
	
	// 关联数据
	Customer   *Customer    `json:"customer,omitempty" xorm:"-"`
	OrderItems []OrderItem  `json:"order_items,omitempty" xorm:"-"`
}

// 订单商品模型
type OrderItem struct {
	ID        int64   `json:"id" xorm:"pk autoincr"`
	OrderID   int64   `json:"order_id" xorm:"index"`
	ProductID int64   `json:"product_id" xorm:"index"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price" xorm:"decimal(10,2)"` // 下单时的价格
	Amount    float64 `json:"amount" xorm:"decimal(10,2)"` // 小计金额
	
	// 关联数据
	Product *Product `json:"product,omitempty" xorm:"-"`
}

// 节假日模型
type Holiday struct {
	ID       int64     `json:"id" xorm:"pk autoincr"`
	Date     time.Time `json:"date" xorm:"date unique"`
	Name     string    `json:"name" xorm:"varchar(255)"`
	Type     string    `json:"type" xorm:"varchar(50)"` // legal, traditional
	IsWorkday bool     `json:"is_workday" xorm:"default(false)"` // 是否为调休工作日
}