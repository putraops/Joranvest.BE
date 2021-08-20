package dto

type OrderDto struct {
	Id           string           `json:"id" form:"id"`
	CustomerId   string           `json:"customer_id" form:"customer_id"`
	VoucherId    string           `json:"voucher_id" form:"voucher_id"`
	TotalPrice   float64          `json:"total_price" form:"total_price"`
	TotalPayment float64          `json:"total_payment" form:"total_payment"`
	TotalItem    int              `json:"total_item" form:"total_item"`
	Description  string           `json:"description" form:"description"`
	Detail       string           `json:"-" form:"detail"`
	UpdatedBy    string           `json:"-"`
	OrderDetail  []OrderDetailDto `json:"-"`
}

type OrderDetailDto struct {
	Id        string  `json:"id" form:"id"`
	ProductId string  `json:"product_id" form:"product_id"`
	Price     float64 `json:"price" form:"price"`
	Quantity  float64 `json:"quantity" form:"quantity"`
}

type PaymentDto struct {
	Id           string  `json:"id" form:"id"`
	TotalPayment float64 `json:"total_payment" form:"total_payment"`
	CreatedBy    string  `json:"-"`
	ApprovedBy   string  `json:"-"`
	UpdatedBy    string  `json:"-"`
}

type OrderStatusDto struct {
	Id     string `json:"id" form:"id"`
	Status int    `json:"status" form:"status"`
}

// EstimatedDate time.Time   `gorm:"type:time" json:"estimated_date"`
// PaymentStatus int     `json:"payment_status" form:"customer_id"`
// OrderStatus int `gorm:"type:int" json:"order_status" form:"order_status"`
// OrderDetail   []interface `json:"order_detail"`
