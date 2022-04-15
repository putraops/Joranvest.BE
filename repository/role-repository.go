package repository

import (
	"encoding/json"
	"fmt"
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/entity_view_models"
	"strings"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RoleRepository interface {
	GetPagination(request commons.Pagination2ndRequest) interface{}
	GetAll(filter map[string]interface{}) []models.Role
	Lookup(request helper.ReactSelectRequest) []models.Role
	Insert(t models.Role) helper.Result
	Update(record models.Role) helper.Result
	GetById(recordId string) helper.Result
	GetViewById(recordId string) helper.Result
	DeleteById(recordId string) helper.Result

	OpenTransaction(trxHandle *gorm.DB) roleRepository
}

type roleRepository struct {
	DB                *gorm.DB
	serviceRepository ServiceRepository
	tableName         string
	viewQuery         string
	currentTime       time.Time
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{
		tableName:         models.Role.TableName(models.Role{}),
		viewQuery:         entity_view_models.EntityRoleView.ViewModel(entity_view_models.EntityRoleView{}),
		DB:                db,
		serviceRepository: NewServiceRepository(db),
		currentTime:       time.Now(),
	}
}

func (r roleRepository) GetPagination(request commons.Pagination2ndRequest) interface{} {
	var response commons.PaginationResponse
	var records []entity_view_models.EntityRoleView
	var recordsUnfilter []entity_view_models.EntityRoleView

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

func (r roleRepository) GetAll(filter map[string]interface{}) []models.Role {
	var records []models.Role
	if len(filter) == 0 {
		r.DB.Find(&records)
	} else if len(filter) != 0 {
		r.DB.Where(filter).Find(&records)
	}
	return records
}

func (r roleRepository) Lookup(request helper.ReactSelectRequest) []models.Role {
	records := []models.Role{}
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

func (r roleRepository) Insert(record models.Role) helper.Result {
	newId := uuid.New().String()
	record.Id = &newId
	record.CreatedAt = &r.currentTime

	if err := r.DB.Create(&record).Error; err != nil {
		return helper.StandartResult(false, fmt.Sprintf("%v,", err.Error()), helper.EmptyObj{})
	}

	return helper.StandartResult(true, "Ok", record)
}

func (r roleRepository) Update(record models.Role) helper.Result {
	var oldRecord models.Role
	r.DB.First(&oldRecord, "id = ?", record.Id)
	if record.Id == nil {
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
	return helper.StandartResult(true, "Data has been updated.", oldRecord)
}

func (r roleRepository) GetById(recordId string) helper.Result {
	var record models.Role
	r.DB.First(&record, "id = ?", recordId)
	if record.Id == nil {
		res := helper.StandartResult(false, "Record not found", helper.EmptyObj{})
		return res
	}
	res := helper.StandartResult(true, "Ok", record)
	return res
}

func (r roleRepository) GetViewById(recordId string) helper.Result {
	var record entity_view_models.EntityRoleView
	r.DB.First(&record, "id = ?", recordId)
	if record.Id == nil {
		res := helper.StandartResult(false, "Record not found", helper.EmptyObj{})
		return res
	}
	res := helper.StandartResult(true, "Ok", record)
	return res
}

func (r roleRepository) DeleteById(recordId string) helper.Result {
	tx := r.DB.Begin()

	var record models.Role
	r.DB.First(&record, "id = ?", recordId)

	if record.Id == nil {
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

func (r roleRepository) OpenTransaction(trxHandle *gorm.DB) roleRepository {
	if trxHandle == nil {
		log.Error("Transaction Database not found")
		log.Error(r.serviceRepository.getCurrentFuncName())
		return r
	}
	r.DB = trxHandle
	return r
}
