package repository

import (
	"database/sql"
	"fmt"
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/entity_view_models"
	"joranvest/models/view_models"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EducationPlaylistRepository interface {
	GetPagination(request commons.Pagination2ndRequest) interface{}
	GetAll(filter map[string]interface{}) []models.EducationPlaylist
	GetPlaylistByUserId(educationId string, userId string) helper.Result
	Lookup(request helper.ReactSelectRequest) []models.EducationPlaylist
	Insert(t models.EducationPlaylist) helper.Result
	MarkVideoAsWatched(record models.EducationPlaylistUser) helper.Result
	Update(record models.EducationPlaylist) helper.Result
	GetById(recordId string) helper.Result
	DeleteById(recordId string) helper.Result

	OpenTransaction(trxHandle *gorm.DB) educationPlaylistRepository
}

type educationPlaylistRepository struct {
	DB                *gorm.DB
	serviceRepository ServiceRepository
	tableName         string
	viewQuery         string
}

func NewEducationPlaylistRepository(db *gorm.DB) EducationPlaylistRepository {
	return educationPlaylistRepository{
		DB:                db,
		serviceRepository: NewServiceRepository(db),
	}
}

func (r educationPlaylistRepository) GetPagination(request commons.Pagination2ndRequest) interface{} {
	var response commons.PaginationResponse
	var records []entity_view_models.EntityEducationPlaylistView
	var recordsUnfilter []entity_view_models.EntityEducationPlaylistView

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

func (r educationPlaylistRepository) GetAll(filter map[string]interface{}) []models.EducationPlaylist {
	var records []models.EducationPlaylist
	if len(filter) == 0 {
		r.DB.Find(&records).Order("order_index ASC, created_at ASC")
	} else if len(filter) != 0 {
		r.DB.Where(filter).Find(&records).Order("order_index ASC, created_at ASC")
	}
	return records
}

func (r educationPlaylistRepository) GetViewAll(filter map[string]interface{}) []entity_view_models.EntityEducationPlaylistView {
	var records []entity_view_models.EntityEducationPlaylistView
	if len(filter) == 0 {
		r.DB.Find(&records).Order("order_index ASC, created_at ASC")
	} else if len(filter) != 0 {
		r.DB.Where(filter).Find(&records).Order("order_index ASC, created_at ASC")
	}
	return records
}

func (r educationPlaylistRepository) GetPlaylistByUserId(educationId string, userId string) helper.Result {
	var records []view_models.EducationPlaylistByUserViewModel

	var sql strings.Builder
	sql.WriteString("SELECT r.id,")
	sql.WriteString("  r.is_active,")
	sql.WriteString("  r.education_id,")
	sql.WriteString("  e.title AS education_title,")
	sql.WriteString("  r.title,")
	sql.WriteString("  r.file_url,")
	sql.WriteString("  r.description,")
	sql.WriteString("  r.order_index,")
	sql.WriteString("  pu.id AS education_playlist_user_id,")
	sql.WriteString("  CASE WHEN pu.id IS NULL THEN false ELSE true END AS is_watched ")
	sql.WriteString("FROM education_playlist r ")
	sql.WriteString("LEFT JOIN education e ON e.id::text = r.education_id::text ")
	sql.WriteString(fmt.Sprintf("LEFT JOIN LATERAL (SELECT pu.id FROM education_playlist_user pu WHERE r.id = pu.education_playlist_id AND pu.application_user_id = '%v') AS pu ON true ", userId))
	sql.WriteString(fmt.Sprintf("WHERE e.id = '%v' ", educationId))
	sql.WriteString("ORDER BY COALESCE(r.order_index, 0) ASC, r.created_at ASC ")

	result := r.DB.Raw(sql.String()).Find(&records)
	if result.Error != nil {
		return helper.StandartResult(false, fmt.Sprintf("%v,", result.Error), helper.EmptyObj{})
	}

	return helper.StandartResult(true, "Ok", records)
}

func (r educationPlaylistRepository) Lookup(request helper.ReactSelectRequest) []models.EducationPlaylist {
	records := []models.EducationPlaylist{}
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

func (r educationPlaylistRepository) Insert(record models.EducationPlaylist) helper.Result {
	record.Id = uuid.New().String()
	record.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	record.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	if err := r.DB.Create(&record).Error; err != nil {
		return helper.StandartResult(false, fmt.Sprintf("%v,", err.Error()), helper.EmptyObj{})
	}

	return helper.StandartResult(true, "Ok", record)
}

func (r educationPlaylistRepository) MarkVideoAsWatched(record models.EducationPlaylistUser) helper.Result {
	var currentRecord models.EducationPlaylistUser
	r.DB.First(&currentRecord, "education_playlist_id = ? AND application_user_id = ?", record.EducationPlaylistId, record.ApplicationUserId)
	if currentRecord.Id == "" {
		record.Id = uuid.New().String()
		record.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
		record.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
		if err := r.DB.Create(&record).Error; err != nil {
			return helper.StandartResult(false, fmt.Sprintf("%v,", err.Error()), helper.EmptyObj{})
		}
	}
	return helper.StandartResult(true, "Ok", record)
}

func (r educationPlaylistRepository) Update(record models.EducationPlaylist) helper.Result {
	var oldRecord models.EducationPlaylist
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

func (r educationPlaylistRepository) GetById(recordId string) helper.Result {
	var record models.EducationPlaylist
	r.DB.First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.StandartResult(false, "Record not found", helper.EmptyObj{})
		return res
	}
	res := helper.StandartResult(true, "Ok", record)
	return res
}

func (r educationPlaylistRepository) DeleteById(recordId string) helper.Result {
	var record models.EducationPlaylist
	r.DB.First(&record, "id = ?", recordId)

	if record.Id == "" {
		res := helper.StandartResult(false, "Record not found", helper.EmptyObj{})
		return res
	} else {
		res := r.DB.Where("id = ?", recordId).Delete(&record)
		if res.RowsAffected == 0 {
			return helper.StandartResult(false, fmt.Sprintf("%v", res.Error), helper.EmptyObj{})
		}
		return helper.StandartResult(true, "Ok", helper.EmptyObj{})
	}
}

func (r educationPlaylistRepository) OpenTransaction(trxHandle *gorm.DB) educationPlaylistRepository {
	if trxHandle == nil {
		log.Error("Transaction Database not found")
		log.Error(r.serviceRepository.getCurrentFuncName())
		return r
	}
	r.DB = trxHandle
	return r
}
