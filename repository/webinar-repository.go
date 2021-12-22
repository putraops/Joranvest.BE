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

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WebinarRepository interface {
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetPagination(request commons.Pagination2ndRequest) interface{}
	GetPaginationRegisteredByUser(request commons.Pagination2ndRequest, userId string) interface{}
	GetAll(filter map[string]interface{}) []models.Webinar
	Insert(t models.Webinar) helper.Response
	Submit(recordId string, userId string) helper.Response
	Update(record models.Webinar) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type webinarConnection struct {
	connection        *gorm.DB
	serviceRepository ServiceRepository
	tableName         string
	viewQuery         string
}

func NewWebinarRepository(db *gorm.DB) WebinarRepository {
	return &webinarConnection{
		connection:        db,
		tableName:         models.Webinar.TableName(models.Webinar{}),
		viewQuery:         entity_view_models.EntityWebinarView.ViewModel(entity_view_models.EntityWebinarView{}),
		serviceRepository: NewServiceRepository(db),
	}
}

func (db *webinarConnection) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	var records []entity_view_models.EntityWebinarView
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
	sql.WriteString(strings.Replace(db.viewQuery, "SELECT  r.id", "r.id", -1))
	sql.WriteString(" WHERE 1 = 1 ")
	sql.WriteString(conditions)
	sql.WriteString(") peta_paged ")
	sql.WriteString(fmt.Sprintf("WHERE peta_rn > %s AND peta_rn <= %s ", start, length))
	db.connection.Raw(sql.String()).Scan(&records)

	sqlCount.WriteString(db.serviceRepository.ConvertViewQueryIntoViewCountByPublic(db.viewQuery, db.tableName))
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
		res.DataRow = []entity_view_models.EntityWebinarView{}
	}
	return res
}

func (db *webinarConnection) GetPagination(request commons.Pagination2ndRequest) interface{} {
	var response commons.PaginationResponse
	var records []entity_view_models.EntityWebinarView
	var recordsUnfilter []entity_view_models.EntityWebinarView

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

	// #region Ordering
	var orders = "COALESCE(submitted_at, created_at) DESC"
	isNearest := false
	if len(request.Order) > 0 {
		order_total := 0
		for k, v := range request.Order {
			if k == "nearest" {
				isNearest = true
			} else {
				if order_total == 0 {
					orders = ""
				} else {
					orders += ", "
				}
				orders += fmt.Sprintf("%v %v ", k, v)
				order_total++
			}
		}
	}

	if isNearest {
		t := time.Now()
		year := t.Year()
		month := t.Month()
		day := t.Day()
		currentDate := fmt.Sprintf("'%v-%v-%v 00:00:00'", year, int(month), day)
		filters += fmt.Sprintf("webinar_start_date >= %v ", currentDate)
	}

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
	return response
}

func (db *webinarConnection) GetPaginationRegisteredByUser(request commons.Pagination2ndRequest, userId string) interface{} {
	var response commons.PaginationResponse
	var records []view_models.WebinarUserViewModel
	var totalUnfilter int

	var sql strings.Builder
	sql.WriteString("SELECT")
	sql.WriteString("	r.*,")
	sql.WriteString("	c.name AS webinar_category_name,")
	sql.WriteString("	f.filepath,")
	sql.WriteString("	f.filepath_thumbnail,")
	sql.WriteString("	f.filename,")
	sql.WriteString("	f.extension,")
	sql.WriteString("	CONCAT(u1.first_name, ' ', u1.last_name) AS created_by_fullname,")
	sql.WriteString("	CONCAT(u2.first_name, ' ', u2.last_name) AS updated_by_fullname,")
	sql.WriteString("	CONCAT(u3.first_name, ' ', u3.last_name) AS submitted_by_fullname ")
	sql.WriteString("FROM (")
	sql.WriteString("	SELECT")
	sql.WriteString("	  r.id,")
	sql.WriteString("	  r.title,")
	sql.WriteString("	  r.description,")
	sql.WriteString("	  r.webinar_start_date,")
	sql.WriteString("	  r.webinar_end_date,")
	sql.WriteString("	  r.min_age,")
	sql.WriteString("	  r.webinar_level,")
	sql.WriteString("	  r.price,")
	sql.WriteString("	  r.discount,")
	sql.WriteString("	  r.is_certificate,")
	sql.WriteString("	  r.reward,")
	sql.WriteString("	  r.status,")
	sql.WriteString("	  r.speaker_type,")
	sql.WriteString("	  w.is_active,")
	sql.WriteString("	  w.is_locked,")
	sql.WriteString("	  w.is_default,")
	sql.WriteString("	  w.created_at,")
	sql.WriteString("	  w.created_by,")
	sql.WriteString("	  w.updated_at,")
	sql.WriteString("	  w.updated_by,")
	sql.WriteString("	  w.submitted_at,")
	sql.WriteString("	  w.submitted_by,")
	sql.WriteString("	  w.approved_at,")
	sql.WriteString("	  w.approved_by,")
	sql.WriteString("	  w.owner_id,")
	sql.WriteString("	  w.entity_id,")
	sql.WriteString("	  0 AS webinar_price,")
	sql.WriteString("	  NULL AS payment_date,")
	sql.WriteString("	  NULL AS payment_date_expired,")
	sql.WriteString("	  NULL AS payment_type,")
	sql.WriteString("	  200 AS payment_status,")
	sql.WriteString("	  r.webinar_category_id")
	sql.WriteString("	FROM webinar r")
	sql.WriteString("	INNER JOIN webinar_registration w ON w.webinar_id = r.id ")
	sql.WriteString(fmt.Sprintf("WHERE w.created_by = '%v'", userId))
	sql.WriteString("	AND w.payment_id = ''")
	sql.WriteString("	UNION ALL")
	sql.WriteString("	SELECT")
	sql.WriteString("	  r.id,")
	sql.WriteString("	  r.title,")
	sql.WriteString("	  r.description,")
	sql.WriteString("	  r.webinar_start_date,")
	sql.WriteString("	  r.webinar_end_date,")
	sql.WriteString("	  r.min_age,")
	sql.WriteString("	  r.webinar_level,")
	sql.WriteString("	  r.price,")
	sql.WriteString("	  r.discount,")
	sql.WriteString("	  r.is_certificate,")
	sql.WriteString("	  r.reward,")
	sql.WriteString("	  r.status,")
	sql.WriteString("	  r.speaker_type,")
	sql.WriteString("	  p.is_active,")
	sql.WriteString("	  p.is_locked,")
	sql.WriteString("	  p.is_default,")
	sql.WriteString("	  p.created_at,")
	sql.WriteString("	  p.created_by,")
	sql.WriteString("	  p.updated_at,")
	sql.WriteString("	  p.updated_by,")
	sql.WriteString("	  p.submitted_at,")
	sql.WriteString("	  p.submitted_by,")
	sql.WriteString("	  p.approved_at,")
	sql.WriteString("	  p.approved_by,")
	sql.WriteString("	  p.owner_id,")
	sql.WriteString("	  p.entity_id,")
	sql.WriteString("	  p.price + p.unique_number AS webinar_price,")
	sql.WriteString("	  p.payment_date,")
	sql.WriteString("	  p.payment_date_expired,")
	sql.WriteString("	  p.payment_type,")
	sql.WriteString("	  p.payment_status,")
	sql.WriteString("	  r.webinar_category_id")
	sql.WriteString("	FROM webinar r")
	sql.WriteString("	INNER JOIN payment p ON p.record_id = r.id ")
	sql.WriteString(fmt.Sprintf("WHERE r.created_by = '%v' ", userId))
	sql.WriteString(") AS r ")
	sql.WriteString("LEFT JOIN webinar_category c ON c.id = r.webinar_category_id ")
	sql.WriteString("LEFT JOIN filemaster f ON f.record_id = r.id ")
	sql.WriteString("LEFT JOIN application_user u1 ON u1.id = r.created_by ")
	sql.WriteString("LEFT JOIN application_user u2 ON u2.id = r.updated_by ")
	sql.WriteString("LEFT JOIN application_user u3 ON u3.id = r.submitted_by ")
	sql.WriteString("ORDER BY r.created_by DESC ")

	var sqlTotalUnfilter strings.Builder
	sqlTotalUnfilter.WriteString("SELECT COUNT(r.*) ")
	sqlTotalUnfilter.WriteString("FROM ( ")
	sqlTotalUnfilter.WriteString("	SELECT ")
	sqlTotalUnfilter.WriteString("	r.id ")
	sqlTotalUnfilter.WriteString("	FROM webinar r ")
	sqlTotalUnfilter.WriteString("	INNER JOIN webinar_registration w ON w.webinar_id = r.id ")
	sqlTotalUnfilter.WriteString(fmt.Sprintf("WHERE w.created_by = '%v'", userId))
	sqlTotalUnfilter.WriteString("	AND w.payment_id = '' ")
	sqlTotalUnfilter.WriteString("	UNION ALL ")
	sqlTotalUnfilter.WriteString("	SELECT ")
	sqlTotalUnfilter.WriteString("	r.id ")
	sqlTotalUnfilter.WriteString("	FROM webinar r ")
	sqlTotalUnfilter.WriteString("	INNER JOIN payment p ON p.record_id = r.id  ")
	sqlTotalUnfilter.WriteString(fmt.Sprintf("WHERE r.created_by = '%v' ", userId))
	sqlTotalUnfilter.WriteString(") AS r ")

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

	//#region Get Total Data for Pagination
	result := db.connection.Raw(sqlTotalUnfilter.String()).Find(&totalUnfilter)
	if result.Error != nil {
		return helper.ServerResponse(false, fmt.Sprintf("%v,", result.Error), fmt.Sprintf("%v,", result.Error), helper.EmptyObj{})
	}
	response.Total = totalUnfilter
	//#endregion

	sql.WriteString(fmt.Sprintf("OFFSET %v ROW FETCH NEXT %v ROWS ONLY", offset, pageSize))
	if err := db.connection.Raw(sql.String()).Scan(&records).Error; err != nil {
		return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
	}

	response.Data = records
	return response
}

func (db *webinarConnection) GetAll(filter map[string]interface{}) []models.Webinar {
	var records []models.Webinar
	if len(filter) == 0 {
		db.connection.Find(&records)
	} else if len(filter) != 0 {
		db.connection.Where(filter).Find(&records)
	}
	return records
}

func (db *webinarConnection) Insert(record models.Webinar) helper.Response {
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

func (db *webinarConnection) Update(record models.Webinar) helper.Response {
	tx := db.connection.Begin()
	var oldRecord models.Webinar
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

func (db *webinarConnection) Submit(recordId string, userId string) helper.Response {
	tx := db.connection.Begin()
	var existingRecord models.Webinar
	db.connection.First(&existingRecord, "id = ?", recordId)
	if existingRecord.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}

	existingRecord.SubmittedBy = userId
	existingRecord.SubmittedAt = sql.NullTime{Time: time.Now(), Valid: true}
	res := tx.Save(&existingRecord)
	if res.RowsAffected == 0 {
		return helper.ServerResponse(false, fmt.Sprintf("%v,", res.Error), fmt.Sprintf("%v,", res.Error), helper.EmptyObj{})
	}

	tx.Commit()
	return helper.ServerResponse(true, "Ok", "", existingRecord)
}

func (db *webinarConnection) GetById(recordId string) helper.Response {
	var record entity_view_models.EntityWebinarView
	db.connection.First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	res := helper.ServerResponse(true, "Ok", "", record)
	return res
}

func (db *webinarConnection) DeleteById(recordId string) helper.Response {
	var record models.Webinar
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
