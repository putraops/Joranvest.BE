package repository

import (
	"database/sql"
	"fmt"
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/view_models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PaymentRepository interface {
	GetPagination(request commons.PaginationRequest) interface{}
	GetAll(filter map[string]interface{}) []models.Payment
	GetUniqueNumber() int
	Insert(t models.Payment) helper.Response
	Update(record models.Payment) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type paymentConnection struct {
	connection           *gorm.DB
	serviceRepository    ServiceRepository
	filemasterRepository FilemasterRepository
	tableName            string
	viewQuery            string
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentConnection{
		connection:           db,
		tableName:            models.Payment.TableName(models.Payment{}),
		viewQuery:            entity_view_models.EntityPaymentView.ViewModel(entity_view_models.EntityPaymentView{}),
		serviceRepository:    NewServiceRepository(db),
		filemasterRepository: NewFilemasterRepository(db),
	}
}

func (db *paymentConnection) GetPagination(request commons.PaginationRequest) interface{} {
	var response commons.PaginationResponse
	var records []entity_view_models.EntityPaymentView

	page := request.Page
	if page == 0 {
		page = 1
	}

	pageSize := request.Size
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	// #region order
	var orders = "COALESCE(submitted_at, created_at) DESC"
	order_total := 0
	for k, v := range request.Order {
		if order_total == 0 {
			orders = ""
		} else {
			orders += ", "
		}
		orders += fmt.Sprintf("%v %v ", k, v)
		order_total++
	}
	// #endregion

	// #region filter
	var filters = ""
	total_filter := 0
	for k, v := range request.Filter {
		if v != "" {
			if total_filter > 0 {
				filters += "AND "
			}
			filters += fmt.Sprintf("%v = '%v' ", k, v)
			total_filter++
		}
	}
	// #endregion

	offset := (page - 1) * pageSize
	db.connection.Where(filters).Order(orders).Offset(offset).Limit(pageSize).Find(&records)

	var count int64
	db.connection.Model(&entity_view_models.EntityPaymentView{}).Where(filters).Count(&count)

	response.Data = records
	response.Total = int(count)
	return response
}

func (db *paymentConnection) GetAll(filter map[string]interface{}) []models.Payment {
	var records []models.Payment
	if len(filter) == 0 {
		db.connection.Find(&records)
	} else if len(filter) != 0 {
		db.connection.Where(filter).Find(&records)
	}
	return records
}

func (db *paymentConnection) GetUniqueNumber() int {
	var record models.Payment
	unique_number := 0
	db.connection.Order("created_at DESC").Where("to_char(created_at, 'YYYY-MM-DD') = to_char(CURRENT_DATE, 'YYYY-MM-DD')").First(&record)
	if record.Id == "" {
		unique_number = 11
	} else {
		unique_number = record.UniqueNumber + 1
	}
	return unique_number
}

func (db *paymentConnection) Insert(record models.Payment) helper.Response {
	tx := db.connection.Begin()

	record.Id = uuid.New().String()
	if record.PaymentStatus == 200 {
		record.PaymentDate = sql.NullTime{Time: time.Now(), Valid: true}
	} else {
		record.PaymentDateExpired = sql.NullTime{Time: time.Now().AddDate(0, 0, 1), Valid: true}
	}
	record.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}

	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
	}

	tx.Commit()
	db.connection.Find(&record)
	return helper.ServerResponse(true, "Ok", "", record)
}

func (db *paymentConnection) Update(record models.Payment) helper.Response {
	tx := db.connection.Begin()
	var oldRecord models.Payment
	db.connection.First(&oldRecord, "id = ?", record.Id)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}

	record.IsActive = oldRecord.IsActive
	record.CreatedAt = oldRecord.CreatedAt
	record.CreatedBy = oldRecord.CreatedBy
	record.EntityId = oldRecord.EntityId
	record.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	res := tx.Save(&record)
	if res.RowsAffected == 0 {
		return helper.ServerResponse(false, fmt.Sprintf("%v,", res.Error), fmt.Sprintf("%v,", res.Error), helper.EmptyObj{})
	}

	tx.Commit()
	db.connection.Preload(clause.Associations).Find(&record)
	return helper.ServerResponse(true, "Ok", "", record)
}

func (db *paymentConnection) GetById(recordId string) helper.Response {
	var record models.Payment
	db.connection.First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	res := helper.ServerResponse(true, "Ok", "", record)
	return res
}

func (db *paymentConnection) DeleteById(recordId string) helper.Response {
	var record models.Payment
	db.connection.First(&record, "id = ?", recordId)

	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	} else {
		res := db.connection.Where("id = ?", recordId).Delete(&record)
		if res.RowsAffected == 0 {
			return helper.ServerResponse(false, "Error", fmt.Sprintf("%v", res.Error), helper.EmptyObj{})
		}
		return helper.ServerResponse(true, "Ok", "", helper.EmptyObj{})
	}
}
