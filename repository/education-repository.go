package repository

import (
	"database/sql"
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

type EducationRepository interface {
	GetPagination(request commons.Pagination2ndRequest) interface{}
	GetAll(filter map[string]interface{}) []models.Education
	Lookup(request helper.ReactSelectRequest) []models.Education
	Insert(t models.Education) helper.Result
	Update(record models.Education) helper.Result
	GetById(recordId string) helper.Result
	GetViewById(recordId string) helper.Result
	DeleteById(recordId string) helper.Result

	OpenTransaction(trxHandle *gorm.DB) educationRepository
}

type educationRepository struct {
	DB                *gorm.DB
	serviceRepository ServiceRepository
	tableName         string
	viewQuery         string
}

func NewEducationRepository(db *gorm.DB) EducationRepository {
	return educationRepository{
		DB:                db,
		serviceRepository: NewServiceRepository(db),
	}
}

func (r educationRepository) GetPagination(request commons.Pagination2ndRequest) interface{} {
	var response commons.PaginationResponse
	var records []entity_view_models.EntityEducationView
	var recordsUnfilter []entity_view_models.EntityEducationView

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

func (r educationRepository) GetAll(filter map[string]interface{}) []models.Education {
	var records []models.Education
	if len(filter) == 0 {
		r.DB.Find(&records)
	} else if len(filter) != 0 {
		r.DB.Where(filter).Find(&records)
	}
	return records
}

func (r educationRepository) Lookup(request helper.ReactSelectRequest) []models.Education {
	records := []models.Education{}
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

func (r educationRepository) Insert(record models.Education) helper.Result {
	record.Id = uuid.New().String()
	record.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	record.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	if err := r.DB.Create(&record).Error; err != nil {
		return helper.StandartResult(false, fmt.Sprintf("%v,", err.Error()), helper.EmptyObj{})
	}

	return helper.StandartResult(true, "Ok", record)
}

func (r educationRepository) Update(record models.Education) helper.Result {
	var oldRecord models.Education
	r.DB.First(&oldRecord, "id = ?", record.Id)
	if record.Id == "" {
		return helper.StandartResult(false, "Record not found", helper.EmptyObj{})
	}

	record.IsActive = oldRecord.IsActive
	record.CreatedAt = oldRecord.CreatedAt
	record.CreatedBy = oldRecord.CreatedBy
	record.EntityId = oldRecord.EntityId
	record.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	res := r.DB.Save(&record)
	if res.RowsAffected == 0 {
		return helper.StandartResult(false, fmt.Sprintf("%v,", res.Error), helper.EmptyObj{})
	}

	return helper.StandartResult(true, "Ok", record)
}

func (r educationRepository) GetById(recordId string) helper.Result {
	var record models.Education
	r.DB.First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.StandartResult(false, "Record not found", helper.EmptyObj{})
		return res
	}
	res := helper.StandartResult(true, "Ok", record)
	return res
}

func (r educationRepository) GetViewById(recordId string) helper.Result {
	var record entity_view_models.EntityEducationView
	r.DB.First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.StandartResult(false, "Record not found", helper.EmptyObj{})
		return res
	}
	res := helper.StandartResult(true, "Ok", record)
	return res
}

func (r educationRepository) DeleteById(recordId string) helper.Result {
	tx := r.DB.Begin()

	var record models.Education
	r.DB.First(&record, "id = ?", recordId)

	if record.Id == "" {
		return helper.StandartResult(false, "Record not found", helper.EmptyObj{})
	} else {
		if errEducation := r.DB.Where("id = ?", recordId).Delete(&record).Error; errEducation != nil {
			log.Error(r.serviceRepository.getCurrentFuncName())
			log.Error(fmt.Sprintf("%v,", errEducation.Error()))
			tx.Rollback()
			return helper.StandartResult(false, fmt.Sprintf("%v", errEducation.Error()), helper.EmptyObj{})
		}
	}

	if errEducationPlaylist := r.DB.Delete(models.EducationPlaylist{}, "education_id = ?", recordId).Error; errEducationPlaylist != nil {
		log.Error(r.serviceRepository.getCurrentFuncName())
		log.Error(fmt.Sprintf("%v,", errEducationPlaylist.Error()))
		tx.Rollback()
		return helper.StandartResult(false, fmt.Sprintf("%v,", errEducationPlaylist.Error()), helper.EmptyObj{})
	}

	tx.Commit()
	return helper.StandartResult(true, "Ok", helper.EmptyObj{})
}

func (r educationRepository) OpenTransaction(trxHandle *gorm.DB) educationRepository {
	if trxHandle == nil {
		log.Error("Transaction Database not found")
		log.Error(r.serviceRepository.getCurrentFuncName())
		return r
	}
	r.DB = trxHandle
	return r
}
