package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/entity_view_models"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WebinarRecordingRepository interface {
	GetPagination(request commons.Pagination2ndRequest) interface{}
	GetAll(filter map[string]interface{}) []models.WebinarRecording
	Lookup(request helper.ReactSelectRequest) []models.WebinarRecording
	Insert(t models.WebinarRecording) helper.Result
	Update(record models.WebinarRecording) helper.Result
	Submit(recordId string, submittedBy string) helper.Result
	GetById(recordId string) helper.Result
	GetViewById(recordId string) helper.Result
	GetByWebinarId(webinarId string) helper.Result
	GetByPathUrl(pathUrl string) helper.Result
	DeleteById(recordId string) helper.Result

	OpenTransaction(trxHandle *gorm.DB) webinarRecordingRepository
}

type webinarRecordingRepository struct {
	DB                *gorm.DB
	serviceRepository ServiceRepository
	tableName         string
	viewQuery         string
}

func NewWebinarRecordingRepository(db *gorm.DB) WebinarRecordingRepository {
	return webinarRecordingRepository{
		DB:                db,
		serviceRepository: NewServiceRepository(db),
	}
}

func (r webinarRecordingRepository) GetPagination(request commons.Pagination2ndRequest) interface{} {
	var response commons.PaginationResponse
	var records []entity_view_models.EntityWebinarRecordingView
	var recordsUnfilter []entity_view_models.EntityWebinarRecordingView

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
	r.DB.Debug().Where(filters).Order(orders).Offset(offset).Limit(pageSize).Find(&records)

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

func (r webinarRecordingRepository) GetAll(filter map[string]interface{}) []models.WebinarRecording {
	var records []models.WebinarRecording
	if len(filter) == 0 {
		r.DB.Find(&records)
	} else if len(filter) != 0 {
		r.DB.Where(filter).Find(&records)
	}
	return records
}

func (r webinarRecordingRepository) Lookup(request helper.ReactSelectRequest) []models.WebinarRecording {
	records := []models.WebinarRecording{}
	r.DB.Order("name asc")

	var orders = "name ASC"
	var filters = ""
	totalFilter := 0
	for _, field := range request.Field {
		if totalFilter == 0 {
			filters += " (LOWER(" + field + ") LIKE " + fmt.Sprint("'%", strings.ToLower(request.Q), "%'")
		} else {
			filters += " OR LOWER(" + field + ") LIKE " + fmt.Sprint("'%", strings.ToLower(request.Q), "%'")
		}
		totalFilter++
	}

	if totalFilter > 0 {
		filters += ")"
	}

	offset := (request.Page - 1) * request.Size
	r.DB.Where(filters).Order(orders).Offset(offset).Limit(request.Size).Find(&records)
	return records
}

func (r webinarRecordingRepository) Insert(record models.WebinarRecording) helper.Result {
	record.Id = uuid.New().String()
	record.CreatedAt = &sql.NullTime{Time: time.Now(), Valid: true}
	record.UpdatedAt = &sql.NullTime{Time: time.Now(), Valid: true}

	if err := r.DB.Create(&record).Error; err != nil {
		return helper.StandartResult(false, fmt.Sprintf("%v,", err.Error()), helper.EmptyObj{})
	}

	return helper.StandartResult(true, "Ok", record)
}

func (r webinarRecordingRepository) Update(record models.WebinarRecording) helper.Result {
	var oldRecord models.WebinarRecording
	r.DB.First(&oldRecord, "id = ?", record.Id)
	if record.Id == "" {
		return helper.StandartResult(false, "Record not found", helper.EmptyObj{})
	}

	mapResult := r.serviceRepository.MapFields(oldRecord, record)
	if !mapResult.Status {
		return mapResult
	}
	json.Unmarshal(mapResult.Data.([]byte), &oldRecord)

	if err := r.DB.Debug().Save(&oldRecord).Error; err != nil {
		return helper.StandartResult(false, fmt.Sprintf("%v,", err.Error()), helper.EmptyObj{})
	}
	return helper.StandartResult(true, "Webinar Recording Data has been updated.", oldRecord)
}

func (r webinarRecordingRepository) Submit(recordId string, submittedBy string) helper.Result {
	tx := r.DB.Begin()

	var currentRecord models.WebinarRecording
	r.DB.First(&currentRecord, "id = ?", recordId)
	if currentRecord.Id == "" {
		return helper.StandartResult(false, "Record not found", helper.EmptyObj{})
	}

	if err := r.DB.Model(&currentRecord).Updates(models.WebinarRecording{SubmittedBy: submittedBy, SubmittedAt: &sql.NullTime{Time: time.Now(), Valid: true}}).Error; err != nil {
		log.Error(r.serviceRepository.getCurrentFuncName())
		log.Error(fmt.Sprintf("%v,", err.Error()))
		tx.Rollback()
		return helper.StandartResult(false, fmt.Sprintf("%v", err.Error()), helper.EmptyObj{})
	}

	tx.Commit()
	return helper.StandartResult(true, "Ok", currentRecord)
}

func (r webinarRecordingRepository) GetById(recordId string) helper.Result {
	var record models.WebinarRecording
	r.DB.First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.StandartResult(false, "Record not found", helper.EmptyObj{})
		return res
	}
	res := helper.StandartResult(true, "Ok", record)
	return res
}

func (r webinarRecordingRepository) GetByWebinarId(webinarId string) helper.Result {
	var record entity_view_models.EntityWebinarRecordingView
	r.DB.First(&record, "webinar_id = ?", webinarId)
	if record.Id == "" {
		res := helper.StandartResult(false, "Record not found", helper.EmptyObj{})
		return res
	}
	res := helper.StandartResult(true, "Ok", record)
	return res
}

func (r webinarRecordingRepository) GetByPathUrl(pathUrl string) helper.Result {
	var record entity_view_models.EntityWebinarRecordingView
	r.DB.First(&record, "path_url = ?", pathUrl)
	if record.Id == "" {
		res := helper.StandartResult(false, "Record not found", helper.EmptyObj{})
		return res
	}
	res := helper.StandartResult(true, "Ok", record)
	return res
}

func (r webinarRecordingRepository) GetViewById(recordId string) helper.Result {
	var record entity_view_models.EntityWebinarRecordingView
	r.DB.First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.StandartResult(false, "Record not found", helper.EmptyObj{})
		return res
	}
	res := helper.StandartResult(true, "Ok", record)
	return res
}

func (r webinarRecordingRepository) DeleteById(recordId string) helper.Result {
	tx := r.DB.Begin()

	var record models.WebinarRecording
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

func (r webinarRecordingRepository) OpenTransaction(trxHandle *gorm.DB) webinarRecordingRepository {
	if trxHandle == nil {
		log.Error("Transaction Database not found")
		log.Error(r.serviceRepository.getCurrentFuncName())
		return r
	}
	r.DB = trxHandle
	return r
}
