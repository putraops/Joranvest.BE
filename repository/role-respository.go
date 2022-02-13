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

type RoleRepository interface {
	Lookup(req map[string]interface{}) []models.Role
	GetPagination(request commons.Pagination2ndRequest) interface{}
	GetAll(filter map[string]interface{}) []models.Role
	Insert(t models.Role) helper.Response
	Update(record models.Role) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type roleConnection struct {
	connection        *gorm.DB
	serviceRepository ServiceRepository
	tableName         string
	viewQuery         string
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleConnection{
		connection:        db,
		tableName:         models.Role.TableName(models.Role{}),
		viewQuery:         entity_view_models.EntityRoleView.ViewModel(entity_view_models.EntityRoleView{}),
		serviceRepository: NewServiceRepository(db),
	}
}

func (db *roleConnection) GetPagination(request commons.Pagination2ndRequest) interface{} {
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
	db.connection.Where(filters).Order(orders).Offset(offset).Limit(pageSize).Find(&records)

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

func (db *roleConnection) GetAll(filter map[string]interface{}) []models.Role {
	var records []models.Role
	if len(filter) == 0 {
		db.connection.Find(&records)
	} else if len(filter) != 0 {
		db.connection.Where(filter).Find(&records)
	}
	return records
}

func (db *roleConnection) Lookup(req map[string]interface{}) []models.Role {
	records := []models.Role{}
	db.connection.Order("name asc")

	var sqlQuery strings.Builder
	sqlQuery.WriteString("SELECT * FROM " + db.tableName)

	v, found := req["condition"]
	if found {
		sqlQuery.WriteString(" WHERE 1 = 1")
		requests := v.(helper.DataFilter).Request
		for _, v := range requests {
			totalFilter := 0
			if v.Operator == "like" {
				for _, valueDetail := range v.Field {
					if totalFilter == 0 {
						sqlQuery.WriteString(" AND (LOWER(" + valueDetail + ") LIKE " + fmt.Sprint("'%", v.Value, "%'"))
					} else {
						sqlQuery.WriteString(" OR LOWER(" + valueDetail + ") LIKE " + fmt.Sprint("'%", v.Value, "%'"))
					}
					totalFilter++
				}
			}

			if totalFilter > 0 {
				sqlQuery.WriteString(")")
			}
		}
	}

	fmt.Println("Query: ", sqlQuery.String())

	db.connection.Raw(sqlQuery.String()).Scan(&records)
	return records
}

func (db *roleConnection) Insert(record models.Role) helper.Response {
	tx := db.connection.Begin()

	record.Id = uuid.New().String()
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

func (db *roleConnection) Update(record models.Role) helper.Response {
	var oldRecord models.Role
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

func (db *roleConnection) GetById(recordId string) helper.Response {
	var record models.Role
	db.connection.First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	res := helper.ServerResponse(true, "Ok", "", record)
	return res
}

func (db *roleConnection) DeleteById(recordId string) helper.Response {
	var record models.Role
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
