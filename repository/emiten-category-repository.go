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

type EmitenCategoryRepository interface {
	Lookup(request helper.ReactSelectRequest) []models.EmitenCategory
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetPagination(request commons.PaginationRequest) interface{}
	GetAll(filter map[string]interface{}) []models.EmitenCategory
	Insert(t models.EmitenCategory) helper.Response
	Update(record models.EmitenCategory) helper.Response
	GetById(recordId string) helper.Response
	GetByName(emitenCategoryName string) helper.Response
	DeleteById(recordId string) helper.Response
}

type emitenCategoryConnection struct {
	connection        *gorm.DB
	serviceRepository ServiceRepository
	tableName         string
	viewQuery         string
}

func NewEmitenCategoryRepository(db *gorm.DB) EmitenCategoryRepository {
	return &emitenCategoryConnection{
		connection:        db,
		tableName:         models.EmitenCategory.TableName(models.EmitenCategory{}),
		viewQuery:         entity_view_models.EntityEmitenCategoryView.ViewModel(entity_view_models.EntityEmitenCategoryView{}),
		serviceRepository: NewServiceRepository(db),
	}
}

func (db *emitenCategoryConnection) Lookup(request helper.ReactSelectRequest) []models.EmitenCategory {
	records := []models.EmitenCategory{}
	db.connection.Order("name asc")

	var orders = "name ASC"
	var filters = ""
	totalFilter := 0
	for _, field := range request.Field {
		if totalFilter == 0 {
			filters += " (LOWER(" + field + ") LIKE " + fmt.Sprint("'%", request.Q, "%'")
		} else {
			filters += " OR LOWER(" + field + ") LIKE " + fmt.Sprint("'%", request.Q, "%'")
		}
		totalFilter++
	}

	if totalFilter > 0 {
		filters += ")"
	}

	offset := (request.Page - 1) * request.Size
	db.connection.Where(filters).Order(orders).Offset(offset).Limit(request.Size).Find(&records)
	return records
}

func (db *emitenCategoryConnection) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	var records []entity_view_models.EntityEmitenCategoryView
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
		res.DataRow = []entity_view_models.EntityEmitenCategoryView{}
	}
	return res
}

func (db *emitenCategoryConnection) GetPagination(request commons.PaginationRequest) interface{} {
	var response commons.PaginationResponse
	var records []entity_view_models.EntityEmitenCategoryView

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
	for k, v := range request.Filter {
		if v != "" {
			if total_filter > 0 {
				filters += "AND "
			}
			filters += fmt.Sprintf("%v = '%v' ", k, v)
			total_filter++
		}
	}
	// #endregion

	offset := (page - 1) * pageSize
	db.connection.Where(filters).Order(orders).Offset(offset).Limit(pageSize).Find(&records)

	var count int64
	db.connection.Model(&entity_view_models.EntityEmitenCategoryView{}).Where(filters).Count(&count)

	response.Data = records
	response.Total = int(count)
	return response
}

func (db *emitenCategoryConnection) GetAll(filter map[string]interface{}) []models.EmitenCategory {
	var records []models.EmitenCategory
	if len(filter) == 0 {
		db.connection.Find(&records)
	} else if len(filter) != 0 {
		db.connection.Where(filter).Find(&records)
	}
	return records
}

func (db *emitenCategoryConnection) Insert(record models.EmitenCategory) helper.Response {
	tx := db.connection.Begin()

	record.Id = uuid.New().String()
	record.IsActive = true
	record.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	record.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
	} else {
		tx.Commit()
		db.connection.Find(&record)
		return helper.ServerResponse(true, "Ok", "", record)
	}
}

func (db *emitenCategoryConnection) Update(record models.EmitenCategory) helper.Response {
	var oldRecord models.EmitenCategory
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
	res := db.connection.Save(&record)
	if res.RowsAffected == 0 {
		return helper.ServerResponse(false, fmt.Sprintf("%v,", res.Error), fmt.Sprintf("%v,", res.Error), helper.EmptyObj{})
	}

	db.connection.Preload(clause.Associations).Find(&record)
	return helper.ServerResponse(true, "Ok", "", record)
}

func (db *emitenCategoryConnection) GetById(recordId string) helper.Response {
	var record models.EmitenCategory
	db.connection.Preload("Emiten").First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	res := helper.ServerResponse(true, "Ok", "", record)
	return res
}

func (db *emitenCategoryConnection) GetByName(emitenCategoryName string) helper.Response {
	var record models.EmitenCategory
	db.connection.First(&record, "name = ?", emitenCategoryName)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	res := helper.ServerResponse(true, "Ok", "", record)
	return res
}

func (db *emitenCategoryConnection) DeleteById(recordId string) helper.Response {
	var record models.EmitenCategory
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
