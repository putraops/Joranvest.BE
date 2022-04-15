package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/entity_view_models"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetPagination(request commons.Pagination2ndRequest) interface{}
	GetAll(filter map[string]interface{}) []models.Product
	Insert(t models.Product) helper.Result
	Update(record models.Product) helper.Result
	GetById(recordId string) helper.Result
	GetProductByRecordId(recordId string) helper.Result
	GetByProductType(product_type string) helper.Result
	GetViewById(recordId string) helper.Result
	DeleteById(recordId string) helper.Result

	OpenTransaction(trxHandle *gorm.DB) productRepository
}

type productRepository struct {
	DB                *gorm.DB
	serviceRepository ServiceRepository
	tableName         string
	viewQuery         string
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return productRepository{
		DB:                db,
		serviceRepository: NewServiceRepository(db),
	}
}

func (r productRepository) GetPagination(request commons.Pagination2ndRequest) interface{} {
	var response commons.PaginationResponse
	var records []entity_view_models.EntityProductView
	var recordsUnfilter []entity_view_models.EntityProductView

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
	if len(request.Filter) > 0 {
		for _, v := range request.Filter {
			if v.Value != "" {
				if total_filter > 0 {
					filters += "AND "
				}

				if v.Operator == "" {
					filters += fmt.Sprintf("%v %v ", v.Field, v.Value)
				} else {
					filters += fmt.Sprintf("%v %v '%v' ", v.Field, v.Operator, v.Value)
				}
				total_filter++
			}
		}
	}
	// #endregion

	offset := (page - 1) * pageSize
	r.DB.Where(filters).Order(orders).Offset(offset).Limit(pageSize).Find(&records)

	// #region Get Total Data for Pagination
	result := r.DB.Where(filters).Find(&recordsUnfilter)
	if result.Error != nil {
		return helper.ServerResponse(false, fmt.Sprintf("%v,", result.Error), fmt.Sprintf("%v,", result.Error), helper.EmptyObj{})
	}
	response.Total = int(result.RowsAffected)
	// #endregion

	response.Data = records
	return response
}

func (r productRepository) GetAll(filter map[string]interface{}) []models.Product {
	var records []models.Product
	if len(filter) == 0 {
		r.DB.Find(&records)
	} else if len(filter) != 0 {
		r.DB.Where(filter).Find(&records)
	}
	return records
}

func (r productRepository) Insert(record models.Product) helper.Result {
	record.Id = uuid.New().String()
	record.CreatedAt = &sql.NullTime{Time: time.Now(), Valid: true}
	record.UpdatedAt = &sql.NullTime{Time: time.Now(), Valid: true}

	if err := r.DB.Create(&record).Error; err != nil {
		return helper.StandartResult(false, fmt.Sprintf("%v,", err.Error()), helper.EmptyObj{})
	}

	return helper.StandartResult(true, "Ok", record)
}

func (r productRepository) Update(record models.Product) helper.Result {
	var oldRecord models.Product
	r.DB.First(&oldRecord, "id = ?", record.Id)
	if record.Id == "" {
		return helper.StandartResult(false, "Record not found", helper.EmptyObj{})
	}

	mapResult := r.serviceRepository.MapFields(oldRecord, record)
	if !mapResult.Status {
		return mapResult
	}
	json.Unmarshal(mapResult.Data.([]byte), &oldRecord)

	if err := r.DB.Save(&oldRecord).Error; err != nil {
		return helper.StandartResult(false, fmt.Sprintf("%v,", err.Error()), helper.EmptyObj{})
	}
	return helper.StandartResult(true, "Product has been updated.", oldRecord)
}

func (r productRepository) GetById(recordId string) helper.Result {
	var record models.Product
	r.DB.First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.StandartResult(false, "Record not found", helper.EmptyObj{})
		return res
	}
	res := helper.StandartResult(true, "Ok", record)
	return res
}

func (r productRepository) GetProductByRecordId(recordId string) helper.Result {
	var webinarRecord models.Webinar
	r.DB.First(&webinarRecord, "id = ?", recordId)
	if webinarRecord.Id != "" {
		return helper.StandartResult(true, "Ok", webinarRecord)
	}

	var membershipRecord models.Membership
	r.DB.First(&membershipRecord, "id = ?", recordId)
	if membershipRecord.Id != "" {
		return helper.StandartResult(true, "Ok", membershipRecord)
	}

	var productRecord models.Product
	r.DB.First(&productRecord, "id = ?", recordId)
	if productRecord.Id != "" {
		return helper.StandartResult(true, "Ok", productRecord)
	}

	return helper.StandartResult(false, "Record not found", helper.EmptyObj{})
}

func (r productRepository) GetByProductType(product_type string) helper.Result {
	var record models.Product
	r.DB.First(&record, "product_type = ?", product_type)
	if record.Id == "" {
		res := helper.StandartResult(false, "Record not found", helper.EmptyObj{})
		return res
	}
	res := helper.StandartResult(true, "Ok", record)
	return res
}

func (r productRepository) GetViewById(recordId string) helper.Result {
	var record entity_view_models.EntityProductView
	r.DB.First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.StandartResult(false, "Record not found", helper.EmptyObj{})
		return res
	}
	res := helper.StandartResult(true, "Ok", record)
	return res
}

func (r productRepository) DeleteById(recordId string) helper.Result {
	tx := r.DB.Begin()

	var record models.Product
	r.DB.First(&record, "id = ?", recordId)

	if record.Id == "" {
		return helper.StandartResult(false, "Record not found", helper.EmptyObj{})
	} else {
		if err := r.DB.Where("id = ?", recordId).Delete(&record).Error; err != nil {
			log.Error(r.serviceRepository.getCurrentFuncName())
			log.Error(fmt.Sprintf("%v,", err.Error()))
			tx.Rollback()
			return helper.StandartResult(false, fmt.Sprintf("%v", err.Error()), helper.EmptyObj{})
		}
	}

	tx.Commit()
	return helper.StandartResult(true, "Ok", helper.EmptyObj{})
}

func (r productRepository) OpenTransaction(trxHandle *gorm.DB) productRepository {
	if trxHandle == nil {
		log.Error("Transaction Database not found")
		log.Error(r.serviceRepository.getCurrentFuncName())
		return r
	}
	r.DB = trxHandle
	return r
}
