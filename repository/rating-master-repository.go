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
	"gorm.io/gorm/clause"
)

type RatingMasterRepository interface {
	GetPagination(request commons.Pagination2ndRequest) helper.Response
	GetAll(filter map[string]interface{}) []models.RatingMaster
	Insert(t models.RatingMaster) helper.Response
	Update(record models.RatingMaster) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type ratingMasterConnection struct {
	connection        *gorm.DB
	serviceRepository ServiceRepository
	tableName         string
	viewQuery         string
}

func NewRatingMasterRepository(db *gorm.DB) RatingMasterRepository {
	return &ratingMasterConnection{
		connection:        db,
		tableName:         models.RatingMaster.TableName(models.RatingMaster{}),
		viewQuery:         entity_view_models.EntityRatingMasterView.ViewModel(entity_view_models.EntityRatingMasterView{}),
		serviceRepository: NewServiceRepository(db),
	}
}

func (db *ratingMasterConnection) GetPagination(request commons.Pagination2ndRequest) helper.Response {
	var response commons.PaginationResponse
	var records []entity_view_models.EntityRatingMasterView
	var recordsUnfilter []entity_view_models.EntityRatingMasterView

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

func (db *ratingMasterConnection) GetAll(filter map[string]interface{}) []models.RatingMaster {
	var records []models.RatingMaster
	if len(filter) == 0 {
		db.connection.Find(&records)
	} else if len(filter) != 0 {
		db.connection.Where(filter).Find(&records)
	}
	return records
}

func (db *ratingMasterConnection) Insert(record models.RatingMaster) helper.Response {
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

func (db *ratingMasterConnection) Update(record models.RatingMaster) helper.Response {
	var oldRecord models.RatingMaster
	db.connection.First(&oldRecord, "id = ?", record.Id)
	if record.Id == "" {
		log.Error(db.serviceRepository.getCurrentFuncName())
		log.Error("Record not found")
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}

	record.IsActive = oldRecord.IsActive
	record.CreatedAt = oldRecord.CreatedAt
	record.CreatedBy = oldRecord.CreatedBy
	record.EntityId = oldRecord.EntityId
	record.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	res := db.connection.Save(&record)
	if res.RowsAffected == 0 {
		log.Error(db.serviceRepository.getCurrentFuncName())
		log.Error(fmt.Sprintf("%v,", res.Error))
		return helper.ServerResponse(false, fmt.Sprintf("%v,", res.Error), fmt.Sprintf("%v,", res.Error), helper.EmptyObj{})
	}

	db.connection.Preload(clause.Associations).Find(&record)
	return helper.ServerResponse(true, "Ok", "", record)
}

func (db *ratingMasterConnection) GetById(recordId string) helper.Response {
	var record models.RatingMaster
	db.connection.Preload("Emiten").First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	res := helper.ServerResponse(true, "Ok", "", record)
	return res
}

func (db *ratingMasterConnection) DeleteById(recordId string) helper.Response {
	var record models.RatingMaster
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
