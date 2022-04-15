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

type RoleMemberRepository interface {
	GetPagination(request commons.Pagination2ndRequest) interface{}
	GetAll(filter map[string]interface{}) []models.RoleMember
	GetViewAll(filter map[string]interface{}) []entity_view_models.EntityRoleMemberView
	Lookup(request helper.ReactSelectRequest) []models.RoleMember
	Insert(t models.RoleMember) helper.Result
	Update(record models.RoleMember) helper.Result
	GetById(recordId string) helper.Result
	GetViewById(recordId string) helper.Result
	DeleteById(recordId string) helper.Result
	GetUsersInRole(roleId string) []entity_view_models.EntityRoleMemberView
	GetUsersNotInRole(roleId string, search string) []entity_view_models.EntityApplicationUserView

	OpenTransaction(trxHandle *gorm.DB) roleMemberRepository
}

type roleMemberRepository struct {
	DB                *gorm.DB
	serviceRepository ServiceRepository
	tableName         string
	viewQuery         string
	currentTime       time.Time
}

func NewRoleMemberRepository(db *gorm.DB) RoleMemberRepository {
	return &roleMemberRepository{
		tableName:         models.RoleMember.TableName(models.RoleMember{}),
		viewQuery:         entity_view_models.EntityRoleMemberView.ViewModel(entity_view_models.EntityRoleMemberView{}),
		DB:                db,
		serviceRepository: NewServiceRepository(db),
		currentTime:       time.Now(),
	}
}

func (r roleMemberRepository) GetPagination(request commons.Pagination2ndRequest) interface{} {
	var response commons.PaginationResponse
	var records []entity_view_models.EntityRoleMemberView
	var recordsUnfilter []entity_view_models.EntityRoleMemberView

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

func (r roleMemberRepository) GetAll(filter map[string]interface{}) []models.RoleMember {
	var records []models.RoleMember
	if len(filter) == 0 {
		r.DB.Find(&records)
	} else if len(filter) != 0 {
		r.DB.Where(filter).Find(&records)
	}
	return records
}

func (r roleMemberRepository) GetViewAll(filter map[string]interface{}) []entity_view_models.EntityRoleMemberView {
	var records []entity_view_models.EntityRoleMemberView
	if len(filter) == 0 {
		r.DB.Find(&records)
	} else if len(filter) != 0 {
		r.DB.Where(filter).Find(&records)
	}
	return records
}

func (r roleMemberRepository) Lookup(request helper.ReactSelectRequest) []models.RoleMember {
	records := []models.RoleMember{}
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

func (r roleMemberRepository) Insert(record models.RoleMember) helper.Result {
	newId := uuid.New().String()
	record.Id = &newId
	record.CreatedAt = &r.currentTime

	if err := r.DB.Create(&record).Error; err != nil {
		message := err.Error()
		if strings.Contains(message, "duplicate key value violates unique constraint") {
			message = "User already in Role"
		}
		return helper.StandartResult(false, message, helper.EmptyObj{})
	}

	return helper.StandartResult(true, "Ok", record)
}

func (r roleMemberRepository) Update(record models.RoleMember) helper.Result {
	var oldRecord models.RoleMember
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

func (r roleMemberRepository) GetById(recordId string) helper.Result {
	var record models.RoleMember
	r.DB.First(&record, "id = ?", recordId)
	if record.Id == nil {
		res := helper.StandartResult(false, "Record not found", helper.EmptyObj{})
		return res
	}
	res := helper.StandartResult(true, "Ok", record)
	return res
}

func (r roleMemberRepository) GetViewById(recordId string) helper.Result {
	var record entity_view_models.EntityRoleMemberView
	r.DB.First(&record, "id = ?", recordId)
	if record.Id == nil {
		res := helper.StandartResult(false, "Record not found", helper.EmptyObj{})
		return res
	}
	res := helper.StandartResult(true, "Ok", record)
	return res
}

func (r roleMemberRepository) DeleteById(recordId string) helper.Result {
	tx := r.DB.Begin()

	var record models.RoleMember
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

func (r roleMemberRepository) GetUsersInRole(roleId string) []entity_view_models.EntityRoleMemberView {
	records := []entity_view_models.EntityRoleMemberView{}
	r.DB.Where("role_id = ?", roleId).Find(&records)
	return records
}

func (r roleMemberRepository) GetUsersNotInRole(roleId string, search string) []entity_view_models.EntityApplicationUserView {
	records := []entity_view_models.EntityApplicationUserView{}
	r.DB.
		Where("is_admin <> true AND (first_name LIKE ? OR last_name LIKE ?) AND id NOT IN (?)", search, search, r.DB.Where("role_id = ? ", roleId).Table("role_member").Select("application_user_id")).
		Find(&records)
	r.DB.
		Where("is_admin <> true AND (LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ?) AND id NOT IN (?)",
			"%"+strings.ToLower(search)+"%",
			"%"+strings.ToLower(search)+"%",
			r.DB.Where("role_id = ? ", roleId).Table("role_member").Select("application_user_id")).
		Find(&records)
	return records
}

func (r roleMemberRepository) OpenTransaction(trxHandle *gorm.DB) roleMemberRepository {
	if trxHandle == nil {
		log.Error("Transaction Database not found")
		log.Error(r.serviceRepository.getCurrentFuncName())
		return r
	}
	r.DB = trxHandle
	return r
}
