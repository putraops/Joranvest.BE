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

type ApplicationMenuRepository interface {
	Lookup(req map[string]interface{}) []models.ApplicationMenu
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetAll(filter map[string]interface{}) []models.ApplicationMenu
	GetTree() []commons.JStreeResponse
	GetTreeByRoleId(roleId string) []commons.JStreeResponse
	Insert(t models.ApplicationMenu) helper.Response
	Update(record models.ApplicationMenu) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type applicationMenuConnection struct {
	connection        *gorm.DB
	serviceRepository ServiceRepository
	tableName         string
	viewQuery         string
}

func NewApplicationMenuRepository(db *gorm.DB) ApplicationMenuRepository {
	return &applicationMenuConnection{
		connection:        db,
		tableName:         models.ApplicationMenu.TableName(models.ApplicationMenu{}),
		viewQuery:         entity_view_models.EntityApplicationMenuView.ViewModel(entity_view_models.EntityApplicationMenuView{}),
		serviceRepository: NewServiceRepository(db),
	}
}

func (db *applicationMenuConnection) Lookup(req map[string]interface{}) []models.ApplicationMenu {
	records := []models.ApplicationMenu{}
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

func (db *applicationMenuConnection) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	var records []entity_view_models.EntityWebinarCategoryView
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
		res.DataRow = []entity_view_models.EntityApplicationMenuView{}
	}
	return res
}

func (db *applicationMenuConnection) GetTree() []commons.JStreeResponse {
	var res []commons.JStreeResponse
	var records []models.ApplicationMenu
	db.connection.Where("parent_id = ?", "").Find(&records)

	if len(records) > 0 {
		for _, s := range records {
			var item commons.JStreeResponse
			item.Id = s.Id
			item.Text = s.Name
			item.Description = s.Description
			item.JStreeState.Opened = true
			item.JStreeState.Disabled = false
			item.JStreeState.Selected = true
			item.JStreeState.Checked = false

			var childrenModel []models.ApplicationMenu
			db.connection.Where("parent_id = ?", item.Id).Find(&childrenModel)
			if len(childrenModel) > 0 {
				var children []commons.JStreeResponse
				for _, c := range childrenModel {
					var child commons.JStreeResponse
					child.Id = c.Id
					child.Text = c.Name
					child.Description = c.Description
					child.JStreeState.Opened = true
					child.JStreeState.Disabled = false
					child.JStreeState.Selected = true
					child.JStreeState.Checked = false
					children = append(children, child)
				}
				item.Children = children
			}
			res = append(res, item)
		}
	}
	return res
}

func (db *applicationMenuConnection) GetTreeByRoleId(roleId string) []commons.JStreeResponse {
	var res []commons.JStreeResponse
	var records []entity_view_models.EntityApplicationMenuView

	var sqlQuery strings.Builder
	sqlQuery.WriteString("SELECT ")
	sqlQuery.WriteString("* ")
	sqlQuery.WriteString("FROM ( ")
	sqlQuery.WriteString("	SELECT r.*, true AS is_checked ")
	sqlQuery.WriteString("	FROM application_menu r ")
	sqlQuery.WriteString("	WHERE r.parent_id = '' AND r.id IN (SELECT application_menu_id FROM role_menu WHERE role_id = '" + roleId + "') ")
	sqlQuery.WriteString("	UNION ")
	sqlQuery.WriteString("	SELECT r.*, false AS is_checked ")
	sqlQuery.WriteString("	FROM application_menu r ")
	sqlQuery.WriteString("	WHERE r.parent_id = '' AND r.id NOT IN (SELECT application_menu_id FROM role_menu WHERE role_id = '" + roleId + "') ")
	sqlQuery.WriteString(") AS r ")
	sqlQuery.WriteString("ORDER BY r.order_index ASC")
	db.connection.Raw(sqlQuery.String()).Scan(&records)
	if len(records) > 0 {
		for _, s := range records {
			var item commons.JStreeResponse
			item.Id = s.Id
			item.Text = s.Name
			item.Description = s.Description
			item.JStreeState.Opened = true
			item.JStreeState.Disabled = false
			item.JStreeState.Selected = true
			item.JStreeState.Checked = s.IsChecked
			item.JStreeState.HasSibling = false

			var childrenModel []entity_view_models.EntityApplicationMenuView
			var sqlChildrenQuery strings.Builder
			sqlChildrenQuery.WriteString("SELECT ")
			sqlChildrenQuery.WriteString("* ")
			sqlChildrenQuery.WriteString("FROM ( ")
			sqlChildrenQuery.WriteString("	SELECT r.*, true AS is_checked ")
			sqlChildrenQuery.WriteString("	FROM application_menu r ")
			sqlChildrenQuery.WriteString("	WHERE r.parent_id = '" + s.Id + "' AND r.id IN (SELECT application_menu_id FROM role_menu WHERE role_id = '" + roleId + "') ")
			sqlChildrenQuery.WriteString("	UNION ")
			sqlChildrenQuery.WriteString("	SELECT r.*, false AS is_checked ")
			sqlChildrenQuery.WriteString("	FROM application_menu r ")
			sqlChildrenQuery.WriteString("	WHERE r.parent_id = '" + s.Id + "' AND r.id NOT IN (SELECT application_menu_id FROM role_menu WHERE role_id = '" + roleId + "') ")
			sqlChildrenQuery.WriteString(") AS r ")
			sqlChildrenQuery.WriteString("ORDER BY r.order_index ASC")
			db.connection.Raw(sqlChildrenQuery.String()).Scan(&childrenModel)

			if len(childrenModel) > 0 {
				item.JStreeState.Checked = false
				var children []commons.JStreeResponse
				for _, c := range childrenModel {
					var child commons.JStreeResponse
					child.Id = c.Id
					child.Text = c.Name
					child.Description = c.Description
					child.JStreeState.Opened = true
					child.JStreeState.Disabled = false
					child.JStreeState.Selected = true
					child.JStreeState.Checked = c.IsChecked
					child.JStreeState.HasSibling = false
					if len(childrenModel) > 1 {
						child.JStreeState.HasSibling = true
					}
					children = append(children, child)
				}
				item.Children = children
			}

			res = append(res, item)
		}

	}
	return res
}

func (db *applicationMenuConnection) GetAll(filter map[string]interface{}) []models.ApplicationMenu {
	var records []models.ApplicationMenu
	if len(filter) == 0 {
		db.connection.Find(&records)
	} else if len(filter) != 0 {
		db.connection.Where(filter).Find(&records)
	}
	return records
}

func (db *applicationMenuConnection) Insert(record models.ApplicationMenu) helper.Response {
	tx := db.connection.Begin()

	fmt.Println(record.ParentId)
	fmt.Println(record.ParentId)
	fmt.Println(record.ParentId)
	fmt.Println(record.ParentId)
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

func (db *applicationMenuConnection) Update(record models.ApplicationMenu) helper.Response {
	var oldRecord models.ApplicationMenu
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

func (db *applicationMenuConnection) GetById(recordId string) helper.Response {
	var record entity_view_models.EntityApplicationMenuView
	db.connection.First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	res := helper.ServerResponse(true, "Ok", "", record)
	return res
}

func (db *applicationMenuConnection) DeleteById(recordId string) helper.Response {
	var record models.ApplicationMenu
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
