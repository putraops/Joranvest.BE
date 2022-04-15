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

type EmailLoggingRepository interface {
	GetPagination(request commons.Pagination2ndRequest) helper.Response
	GetLastIntervalLogging(email string, mailType string, intervalMinutes int) int64
	Insert(t models.EmailLogging) helper.Response
	DeleteById(recordId string) helper.Response
}

type emailLoggingConnection struct {
	connection        *gorm.DB
	serviceRepository ServiceRepository
	tableName         string
	viewQuery         string
}

func NewEmailLoggingRepository(db *gorm.DB) EmailLoggingRepository {
	return &emailLoggingConnection{
		connection:        db,
		tableName:         models.EmailLogging.TableName(models.EmailLogging{}),
		viewQuery:         entity_view_models.EntityEmailLoggingView.ViewModel(entity_view_models.EntityEmailLoggingView{}),
		serviceRepository: NewServiceRepository(db),
	}
}

func (db *emailLoggingConnection) GetPagination(request commons.Pagination2ndRequest) helper.Response {
	var response commons.PaginationResponse
	var records []entity_view_models.EntityEmailLoggingView
	var recordsUnfilter []entity_view_models.EntityEmailLoggingView

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

func (db *emailLoggingConnection) GetLastIntervalLogging(email string, mailType string, intervalMinutes int) int64 {
	nowTime := sql.NullTime{Time: time.Now().Add(time.Minute), Valid: true}
	lastInterval := sql.NullTime{Time: time.Now().Add(time.Minute * (-1 * time.Duration(intervalMinutes))), Valid: true}

	var total int64
	db.connection.Model(&models.EmailLogging{}).
		Where("email = ? AND (last_sent BETWEEN ? AND ?)", email, lastInterval, nowTime).
		Count(&total)
	return total
}

func (db *emailLoggingConnection) Insert(record models.EmailLogging) helper.Response {
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
	return helper.ServerResponse(true, "Ok", "", record)
}

func (db *emailLoggingConnection) DeleteById(recordId string) helper.Response {
	var record models.EmailLogging
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
