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

type ArticleRepository interface {
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetPagination(request commons.Pagination2ndRequest) interface{}
	GetAll(filter map[string]interface{}) []models.Article
	GetArticleCoverById(recordId string) helper.Response
	Insert(t models.Article) helper.Response
	Update(record models.Article) helper.Response
	Submit(recordId string, userId string) helper.Response
	GetById(recordId string) helper.Response
	GetViewById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type articleConnection struct {
	connection        *gorm.DB
	serviceRepository ServiceRepository
	tableName         string
	viewQuery         string
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleConnection{
		connection:        db,
		tableName:         models.Article.TableName(models.Article{}),
		viewQuery:         entity_view_models.EntityArticleView.ViewModel(entity_view_models.EntityArticleView{}),
		serviceRepository: NewServiceRepository(db),
	}
}

func (db *articleConnection) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	var records []entity_view_models.EntityArticleView
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
		res.DataRow = []entity_view_models.EntityArticleView{}
	}
	return res
}

func (db *articleConnection) GetPagination(request commons.Pagination2ndRequest) interface{} {
	var response commons.PaginationResponse
	var records []entity_view_models.EntityArticleView
	var recordsUnfilter []entity_view_models.EntityArticleView

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
	return response
}

func (db *articleConnection) GetAll(filter map[string]interface{}) []models.Article {
	var records []models.Article
	if len(filter) == 0 {
		db.connection.Find(&records)
	} else if len(filter) != 0 {
		db.connection.Where(filter).Find(&records)
	}
	return records
}

func (db *articleConnection) Insert(record models.Article) helper.Response {
	tx := db.connection.Begin()

	record.Id = uuid.New().String()
	record.CreatedAt = sql.NullTime{Time: time.Now().Local().UTC(), Valid: true}

	for i := 0; i < len(record.ArticleTag); i++ {
		record.ArticleTag[i].Id = uuid.New().String()
		record.ArticleTag[i].EntityId = record.EntityId
		record.ArticleTag[i].CreatedBy = record.CreatedBy
		record.ArticleTag[i].CreatedAt = sql.NullTime{Time: time.Now().Local().UTC(), Valid: true}
		record.ArticleTag[i].ArticleId = record.Id
	}

	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
	}

	if len(record.ArticleTag) > 0 {
		if err := tx.Create(&record.ArticleTag).Error; err != nil {
			tx.Rollback()
			return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
		}
	}

	tx.Commit()
	db.connection.Find(&record)
	return helper.ServerResponse(true, "Ok", "", record)
}

func (db *articleConnection) Update(record models.Article) helper.Response {
	tx := db.connection.Begin()
	var oldRecord models.Article
	db.connection.First(&oldRecord, "id = ?", record.Id)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}

	//-- Delete Article Tag
	var tags models.ArticleTag
	if err := tx.Where("article_id = ?", record.Id).Delete(&tags).Error; err != nil {
		tx.Rollback()
		return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
	}

	for i := 0; i < len(record.ArticleTag); i++ {
		fmt.Println(record.ArticleTag[i].TagId)
		record.ArticleTag[i].Id = uuid.New().String()
		record.ArticleTag[i].EntityId = record.EntityId
		record.ArticleTag[i].CreatedBy = record.CreatedBy
		record.ArticleTag[i].CreatedAt = sql.NullTime{Time: time.Now().Local().UTC(), Valid: true}
		record.ArticleTag[i].ArticleId = record.Id
	}
	if len(record.ArticleTag) > 0 {
		if err := tx.Save(&record.ArticleTag).Error; err != nil {
			tx.Rollback()
			return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
		}
	}

	record.IsActive = oldRecord.IsActive
	record.CreatedAt = oldRecord.CreatedAt
	record.CreatedBy = oldRecord.CreatedBy
	record.EntityId = oldRecord.EntityId
	record.UpdatedAt = sql.NullTime{Time: time.Now().Local().UTC(), Valid: true}
	res := tx.Save(&record)
	if res.RowsAffected == 0 {
		return helper.ServerResponse(false, fmt.Sprintf("%v,", res.Error), fmt.Sprintf("%v,", res.Error), helper.EmptyObj{})
	}

	tx.Commit()
	db.connection.Preload(clause.Associations).Find(&record)
	return helper.ServerResponse(true, "Ok", "", record)
}

func (db *articleConnection) Submit(recordId string, userId string) helper.Response {
	tx := db.connection.Begin()
	var existingRecord models.Article
	db.connection.First(&existingRecord, "id = ?", recordId)
	if existingRecord.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}

	existingRecord.SubmittedBy = userId
	existingRecord.SubmittedAt = sql.NullTime{Time: time.Now().Local().UTC(), Valid: true}
	res := tx.Save(&existingRecord)
	if res.RowsAffected == 0 {
		return helper.ServerResponse(false, fmt.Sprintf("%v,", res.Error), fmt.Sprintf("%v,", res.Error), helper.EmptyObj{})
	}

	tx.Commit()
	return helper.ServerResponse(true, "Ok", "", existingRecord)
}

func (db *articleConnection) GetById(recordId string) helper.Response {
	var record models.Article
	db.connection.Preload("ArticleCategory").First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	res := helper.ServerResponse(true, "Ok", "", record)
	return res
}

func (db *articleConnection) GetViewById(recordId string) helper.Response {
	var record entity_view_models.EntityArticleView
	db.connection.First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	res := helper.ServerResponse(true, "Ok", "", record)
	return res
}

func (db *articleConnection) DeleteById(recordId string) helper.Response {
	var record models.Article
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

func (db *articleConnection) GetArticleCoverById(recordId string) helper.Response {
	var record models.Filemaster
	db.connection.First(&record, "record_id = ? AND file_type = ?", recordId, 1)
	if record.Id == "" {
		return helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
	}
	return helper.ServerResponse(true, "Ok", "", record)
}
