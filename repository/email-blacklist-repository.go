package repository

import (
	"database/sql"
	"fmt"
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/entity_view_models"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EmailBlacklistRepository interface {
	GetPagination(request commons.Pagination2ndRequest) helper.Response
	Insert(t models.EmailBlacklist) helper.Response
	DeleteById(recordId string) helper.Response
}

type emailBlacklistConnection struct {
	connection        *gorm.DB
	serviceRepository ServiceRepository
	tableName         string
	viewQuery         string
}

func NewEmailBlacklistRepository(db *gorm.DB) EmailBlacklistRepository {
	return &emailBlacklistConnection{
		connection:        db,
		tableName:         models.EmailBlacklist.TableName(models.EmailBlacklist{}),
		viewQuery:         entity_view_models.EntityEmailBlacklistView.ViewModel(entity_view_models.EntityEmailBlacklistView{}),
		serviceRepository: NewServiceRepository(db),
	}
}

func (db *emailBlacklistConnection) GetPagination(request commons.Pagination2ndRequest) helper.Response {
	var response commons.PaginationResponse
	var records []entity_view_models.EntityEmailBlacklistView
	var recordsUnfilter []entity_view_models.EntityEmailBlacklistView

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

	offset := (page - 1) * pageSize

	// #region Ordering
	var orders = "COALESCE(submitted_at, created_at) DESC"
	if len(request.Order) > 0 {
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
	}
	// #endregion

	// #region filter
	var filters = ""
	if len(request.Filter) > 0 {
		total_filter := 0
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

	if err := db.connection.Where(filters).Offset(offset).Order(orders).Limit(pageSize).Find(&records).Error; err != nil {
		return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
	}

	// #region Get Total Data for Pagination
	result := db.connection.Where(filters).Find(&recordsUnfilter)
	if result.Error != nil {
		return helper.ServerResponse(false, fmt.Sprintf("%v,", result.Error), fmt.Sprintf("%v,", result.Error), helper.EmptyObj{})
	}
	response.Total = int(result.RowsAffected)
	// #endregion

	response.Data = records
	return helper.ServerResponse(true, "Ok", "", response)
}

func (db *emailBlacklistConnection) Insert(record models.EmailBlacklist) helper.Response {
	commons.Logger()

	tx := db.connection.Begin()

	record.Id = uuid.New().String()
	record.IsActive = true
	record.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}

	if err := tx.Create(&record).Error; err != nil {
		log.Error(db.serviceRepository.getCurrentFuncName())
		log.Error(fmt.Sprintf("%v,", err))
		tx.Rollback()
		return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
	}

	tx.Commit()
	db.connection.Find(&record)
	return helper.ServerResponse(true, "Ok", "", record)
}

func (db *emailBlacklistConnection) DeleteById(recordId string) helper.Response {
	var record models.EmailBlacklist
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
