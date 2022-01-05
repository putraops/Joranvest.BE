package repository

import (
	"database/sql"
	"fmt"
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/entity_view_models"
	"joranvest/models/request_models"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EmitenRepository interface {
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetPagination(request commons.PaginationRequest) interface{}
	GetAll(filter map[string]interface{}) []models.Emiten
	Lookup(request helper.ReactSelectRequest) []models.Emiten
	Insert(t models.Emiten) helper.Response
	Update(record models.Emiten) helper.Response
	GetById(recordId string) helper.Response
	GetByEmitenCode(emiten_code string) helper.Response
	DeleteById(recordId string) helper.Response
	PatchingEmiten(data []request_models.PatchingEmiten, userId string) helper.Response
}

type emitenConnection struct {
	connection               *gorm.DB
	serviceRepository        ServiceRepository
	emitenCategoryRepository EmitenCategoryRepository
	sectorRepository         SectorRepository
	tableName                string
	viewQuery                string
}

func NewEmitenRepository(db *gorm.DB) EmitenRepository {
	return &emitenConnection{
		connection:               db,
		tableName:                models.Emiten.TableName(models.Emiten{}),
		viewQuery:                entity_view_models.EntityEmitenView.ViewModel(entity_view_models.EntityEmitenView{}),
		serviceRepository:        NewServiceRepository(db),
		sectorRepository:         NewSectorRepository(db),
		emitenCategoryRepository: NewEmitenCategoryRepository(db),
	}
}

func (db *emitenConnection) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	var records []entity_view_models.EntityEmitenView
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
		res.DataRow = []entity_view_models.EntityEmitenView{}
	}
	return res
}

func (db *emitenConnection) GetPagination(request commons.PaginationRequest) interface{} {
	var response commons.PaginationResponse
	var records []entity_view_models.EntityEmitenView

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
	db.connection.Model(&entity_view_models.EntityEmitenView{}).Where(filters).Count(&count)

	response.Data = records
	response.Total = int(count)
	return response
}

func (db *emitenConnection) GetAll(filter map[string]interface{}) []models.Emiten {
	var records []models.Emiten
	if len(filter) == 0 {
		db.connection.Find(&records)
	} else if len(filter) != 0 {
		db.connection.Where(filter).Find(&records)
	}
	return records
}

func (db *emitenConnection) Lookup(request helper.ReactSelectRequest) []models.Emiten {
	records := []models.Emiten{}
	db.connection.Order("emiten_name asc")

	var orders = "emiten_code ASC"
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
	db.connection.Where(filters).Order(orders).Offset(offset).Limit(request.Size).Find(&records)
	return records
}

func (db *emitenConnection) Insert(record models.Emiten) helper.Response {
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

func (db *emitenConnection) Update(record models.Emiten) helper.Response {
	var oldRecord models.Emiten
	db.connection.First(&oldRecord, "id = ?", record.Id)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}

	println("Updating...")
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

func (db *emitenConnection) GetById(recordId string) helper.Response {
	var record models.Emiten
	db.connection.Preload("Sector").Preload("EmitenCategory").First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	res := helper.ServerResponse(true, "Ok", "", record)
	return res
}

func (db *emitenConnection) GetByEmitenCode(emiten_code string) helper.Response {
	var record models.Emiten
	db.connection.First(&record, "emiten_code = ?", emiten_code)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	res := helper.ServerResponse(true, "Ok", "", record)
	return res
}

func (db *emitenConnection) DeleteById(recordId string) helper.Response {
	var record models.Emiten
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

func (db *emitenConnection) PatchingEmiten(data []request_models.PatchingEmiten, userId string) helper.Response {
	tx := db.connection.Begin()

	if len(data) > 0 {
		var emitenData []models.Emiten
		for _, datum := range data {
			//-- Check Sector and getId
			var emitenResult = db.GetByEmitenCode(datum.EmitenCode)
			if emitenResult.Status {
				// Exist
				res := (emitenResult.Data).(models.Emiten)
				emitenData = append(emitenData, res)
			} else {
				var emitenRecord models.Emiten
				emitenRecord.Id = uuid.New().String()
				emitenRecord.EmitenName = datum.EmitenName
				emitenRecord.EmitenCode = datum.EmitenCode
				emitenRecord.IsActive = true
				emitenRecord.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
				emitenRecord.CreatedBy = userId
				emitenRecord.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}

				//-- Check Sector and getId
				var sectorResult = db.sectorRepository.GetByName(datum.EmitenSector)
				if sectorResult.Status {
					// Exist
					res := (sectorResult.Data).(models.Sector)
					emitenRecord.SectorId = res.Id
				} else {
					var new models.Sector
					new.Name = datum.EmitenSector
					new.Description = ""
					new.IsActive = true
					new.CreatedBy = userId
					var resNew = db.sectorRepository.Insert(new)
					if resNew.Status {
						res := (resNew.Data).(models.Sector)
						emitenRecord.SectorId = res.Id
					}
				}

				//-- Check Category and GetId
				var emitenCategoryResult = db.emitenCategoryRepository.GetByName(datum.EmitenCategory)
				if emitenCategoryResult.Status {
					// Exist
					res := (emitenCategoryResult.Data).(models.EmitenCategory)
					emitenRecord.EmitenCategoryId = res.Id
				} else {
					var new models.EmitenCategory
					new.Name = datum.EmitenCategory
					new.Description = ""
					new.IsActive = true
					new.CreatedBy = userId
					var resNew = db.emitenCategoryRepository.Insert(new)
					if resNew.Status {
						res := (resNew.Data).(models.EmitenCategory)
						emitenRecord.EmitenCategoryId = res.Id
					}
				}
				emitenData = append(emitenData, emitenRecord)
			}
		}

		if err := tx.Save(&emitenData).Error; err != nil {
			tx.Rollback()
			return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
		} else {
			tx.Commit()
			return helper.ServerResponse(true, "Ok", "", emitenData)
		}
	} else {
		return helper.ServerResponse(false, "File is empty. Please check your file.", "File is empty. Please check your file.", helper.EmptyObj{})
	}
}
