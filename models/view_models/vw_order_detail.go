package entity_view_models

import (
	"database/sql"
	"strings"
	"time"
)

type EntityOrderDetailView struct {
	Id           string       `json:"id"`
	IsActive     bool         `json:"is_active"`
	IsLocked     bool         `json:"is_locked"`
	IsDefault    bool         `json:"is_default"`
	CreatedAt    time.Time    `json:"created_at"`
	CreatedBy    string       `json:"created_by"`
	UpdatedAt    time.Time    `json:"updated_at"`
	UpdatedBy    string       `json:"updated_by"`
	ApprovedAt   sql.NullTime `json:"approved_at"`
	ApprovedBy   string       `json:"approved_by"`
	EntityId     string       `json:"entity_id"`
	OrderId      string       `json:"order_id"`
	ProductId    string       `json:"product_id"`
	CategoryName string       `json:"category_name"`
	ProductName  string       `json:"product_name"`
	IsUnit       bool         `json:"is_unit"`
	Price        float64      `json:"price"`
	Quantity     float64      `json:"quantity"`
	Description  string       `json:"description"`
}

func (EntityOrderDetailView) TableName() string {
	return "vw_order_detail"
}

func (EntityOrderDetailView) ViewModel() string {
	var sql strings.Builder
	sql.WriteString("SELECT")
	sql.WriteString("  r.id,")
	sql.WriteString("  r.is_active,")
	sql.WriteString("  r.is_locked,")
	sql.WriteString("  r.is_default,")
	sql.WriteString("  r.created_at,")
	sql.WriteString("  r.created_by,")
	sql.WriteString("  r.updated_at,")
	sql.WriteString("  r.updated_by,")
	sql.WriteString("  r.approved_at,")
	sql.WriteString("  r.approved_by,")
	sql.WriteString("  r.entity_id,")
	sql.WriteString("  r.order_id,")
	sql.WriteString("  r.product_id,")
	sql.WriteString("  c.name AS category_name,")
	sql.WriteString("  p.name AS product_name,")
	sql.WriteString("  p.is_unit,")
	sql.WriteString("  r.price,")
	sql.WriteString("  r.quantity,")
	sql.WriteString("  r.description ")
	sql.WriteString("FROM order_detail r")
	sql.WriteString("  LEFT JOIN product p ON p.id = r.product_id")
	sql.WriteString("  LEFT JOIN category c ON c.id = p.category_id")
	return sql.String()
}

func (EntityOrderDetailView) Migration() map[string]string {
	var view = EntityOrderDetailView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
