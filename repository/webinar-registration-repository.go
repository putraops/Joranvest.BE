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

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WebinarRegistrationRepository interface {
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetPagination(request commons.Pagination2ndRequest) interface{}
	GetAll(filter map[string]interface{}) []models.WebinarRegistration
	GetParticipantsByIds(ids []string) []entity_view_models.EntityWebinarRegistrationView
	GetParticipantsByWebinarId(webinarId string) []entity_view_models.EntityWebinarRegistrationView
	GetById(recordId string) helper.Response
	GetViewById(recordId string) helper.Response
	Insert(record models.WebinarRegistration) helper.Response
	Update(record models.WebinarRegistration) helper.Response
	UpdateInvitationStatusById(id string)
	IsWebinarRegistered(webinarId string, userId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type webinarRegistrationConnection struct {
	connection           *gorm.DB
	serviceRepository    ServiceRepository
	filemasterRepository FilemasterRepository
	tableName            string
	viewQuery            string
}

func NewWebinarRegistrationRepository(db *gorm.DB) WebinarRegistrationRepository {
	return &webinarRegistrationConnection{
		connection:           db,
		tableName:            models.WebinarRegistration.TableName(models.WebinarRegistration{}),
		viewQuery:            entity_view_models.EntityWebinarRegistrationView.ViewModel(entity_view_models.EntityWebinarRegistrationView{}),
		serviceRepository:    NewServiceRepository(db),
		filemasterRepository: NewFilemasterRepository(db),
	}
}

func (db *webinarRegistrationConnection) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	var records []entity_view_models.EntityWebinarRegistrationView
	var res commons.DataTableResponse

	var conditions = ""
	var orderpart = ""
	if request.Draw == 1 && request.DataTableDefaultOrder.Column != "" {
		var column = request.DataTableDefaultOrder.Column
		orderpart = column + " " + request.DataTableDefaultOrder.Dir
	} else {
		var column = request.DataTableColumn[request.DataTableOrder[0].Column].Name
		orderpart = column + " " + request.DataTableOrder[0].Dir
	}
	start := fmt.Sprintf("%v", request.Start)
	length := fmt.Sprintf("%v", (request.Start + request.Length))

	if len(request.Filter) > 0 {
		for _, s := range request.Filter {
			conditions += " AND (" + s.Column + " = '" + s.Value + "') "
		}
	}

	if request.Search.Value != "" {
		conditions += " AND ("
		var totalFilter int = 0
		for _, s := range request.DataTableColumn {
			if s.Searchable {
				if totalFilter > 0 {
					conditions += " OR "
				}
				conditions += fmt.Sprintf("LOWER(CAST (%v AS varchar))", s.Name) + " LIKE '%" + request.Search.Value + "%' "
				totalFilter++
			}
		}
		conditions += ")"
	}

	var sql strings.Builder
	var sqlCount strings.Builder
	sql.WriteString(fmt.Sprintf("SELECT * FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s) peta_rn, ", orderpart))
	sql.WriteString(strings.Replace(db.viewQuery, "SELECT", "", -1))
	sql.WriteString(" WHERE 1 = 1 ")
	sql.WriteString(conditions)
	sql.WriteString(") peta_paged ")
	sql.WriteString(fmt.Sprintf("WHERE peta_rn > %s AND peta_rn <= %s ", start, length))
	db.connection.Raw(sql.String()).Scan(&records)

	sqlCount.WriteString(db.serviceRepository.ConvertViewQueryIntoViewCount(db.viewQuery))
	sqlCount.WriteString("WHERE 1=1")
	sqlCount.WriteString(conditions)
	db.connection.Raw(sqlCount.String()).Scan(&res.RecordsFiltered)

	res.Draw = request.Draw
	if len(records) > 0 {
		res.RecordsTotal = res.RecordsFiltered
		res.DataRow = records
	} else {
		res.RecordsTotal = 0
		res.RecordsFiltered = 0
		res.DataRow = []entity_view_models.EntityWebinarRegistrationView{}
	}
	return res
}

func (db *webinarRegistrationConnection) GetPagination(request commons.Pagination2ndRequest) interface{} {
	var response commons.PaginationResponse
	var records []entity_view_models.EntityWebinarRegistrationView

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
	db.connection.Where(filters).Order(orders).Offset(offset).Limit(pageSize).Find(&records)
	db.connection.Debug().Where(filters).Order(orders).Offset(offset).Limit(pageSize).Find(&records)

	var count int64
	db.connection.Model(&entity_view_models.EntityWebinarRegistrationView{}).Where(filters).Count(&count)

	response.Data = records
	response.Total = int(count)
	return response
}

func (db *webinarRegistrationConnection) GetAll(filter map[string]interface{}) []models.WebinarRegistration {
	var records []models.WebinarRegistration
	if len(filter) == 0 {
		db.connection.Find(&records)
	} else if len(filter) != 0 {
		db.connection.Where(filter).Find(&records)
	}
	return records
}

func (db *webinarRegistrationConnection) GetParticipantsByWebinarId(webinarId string) []entity_view_models.EntityWebinarRegistrationView {
	var records []entity_view_models.EntityWebinarRegistrationView
	db.connection.Where("webinar_id = ?", webinarId).Find(&records)
	return records
}

func (db *webinarRegistrationConnection) GetParticipantsByIds(ids []string) []entity_view_models.EntityWebinarRegistrationView {
	var records []entity_view_models.EntityWebinarRegistrationView
	db.connection.Debug().Where("id IN ?", ids).Find(&records)
	return records
}

func (db *webinarRegistrationConnection) Insert(record models.WebinarRegistration) helper.Response {
	tx := db.connection.Begin()

	record.Id = uuid.New().String()
	record.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}

	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
	}

	tx.Commit()
	db.connection.Find(&record)
	return helper.ServerResponse(true, "Ok", "", record)
}

func (db *webinarRegistrationConnection) Update(record models.WebinarRegistration) helper.Response {
	tx := db.connection.Begin()
	var oldRecord models.WebinarRegistration
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

func (db *webinarRegistrationConnection) UpdateInvitationStatusById(id string) {
	db.connection.Exec("UPDATE webinar_registration SET is_invitation_sent = true, invitation_sent_at = ? WHERE id = ?", sql.NullTime{Time: time.Now(), Valid: true}, id)
}

func (db *webinarRegistrationConnection) GetById(recordId string) helper.Response {
	var record models.WebinarRegistration
	db.connection.First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	res := helper.ServerResponse(true, "Ok", "", record)
	return res
}

func (db *webinarRegistrationConnection) GetViewById(recordId string) helper.Response {
	var record entity_view_models.EntityWebinarRegistrationView
	db.connection.First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	} else {
		res := helper.ServerResponse(true, "Ok", "", record)
		return res
	}
}

func (db *webinarRegistrationConnection) IsWebinarRegistered(webinarId string, userId string) helper.Response {
	var record models.WebinarRegistration
	db.connection.First(&record, "webinar_id = ? AND application_user_id = ?", webinarId, userId)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	res := helper.ServerResponse(true, "Ok", "", record)
	return res
}

func (db *webinarRegistrationConnection) DeleteById(recordId string) helper.Response {
	var record models.WebinarRegistration
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
