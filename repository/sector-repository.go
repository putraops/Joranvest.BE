package repository

import (
	"database/sql"
	"fmt"
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/view_models"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SectorRepository interface {
	Lookup(req map[string]interface{}) []models.Sector
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetAll(filter map[string]interface{}) []models.Sector
	Insert(t models.Sector) helper.Response
	Update(record models.Sector) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type sectorConnection struct {
	connection        *gorm.DB
	serviceRepository ServiceRepository
	tableName         string
	viewQuery         string
}

func NewSectorRepository(db *gorm.DB) SectorRepository {
	return &sectorConnection{
		connection:        db,
		tableName:         models.Sector.TableName(models.Sector{}),
		viewQuery:         entity_view_models.EntitySectorView.ViewModel(entity_view_models.EntitySectorView{}),
		serviceRepository: NewServiceRepository(db),
	}
}

func (db *sectorConnection) Lookup(req map[string]interface{}) []models.Sector {
	records := []models.Sector{}
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

func (db *sectorConnection) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	var records []entity_view_models.EntitySectorView
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
		res.DataRow = []entity_view_models.EntitySectorView{}
	}
	return res
}

func (db *sectorConnection) GetAll(filter map[string]interface{}) []models.Sector {
	var records []models.Sector
	if len(filter) == 0 {
		db.connection.Find(&records)
	} else if len(filter) != 0 {
		db.connection.Where(filter).Find(&records)
	}
	return records
}

func (db *sectorConnection) Insert(record models.Sector) helper.Response {
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

func (db *sectorConnection) Update(record models.Sector) helper.Response {
	var oldRecord models.Sector
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

func (db *sectorConnection) GetById(recordId string) helper.Response {
	var record models.Sector
	db.connection.Preload("Emiten").First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	res := helper.ServerResponse(true, "Ok", "", record)
	return res
}

func (db *sectorConnection) DeleteById(recordId string) helper.Response {
	var record models.Sector
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
