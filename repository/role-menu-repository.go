package repository

import (
	"database/sql"
	"fmt"
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/view_models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoleMenuRepository interface {
	GetAll(filter map[string]interface{}) []models.RoleMenu
	Insert(t []models.RoleMenu) helper.Response
	Update(record models.RoleMenu) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
	DeleteByRoleAndMenuId(roleId string, applicationMenuId string, isParent bool) helper.Response
}

type roleMenuConnection struct {
	connection        *gorm.DB
	serviceRepository ServiceRepository
	tableName         string
	viewQuery         string
}

func NewRoleMenuRepository(db *gorm.DB) RoleMenuRepository {
	return &roleMenuConnection{
		connection:        db,
		tableName:         models.RoleMenu.TableName(models.RoleMenu{}),
		viewQuery:         entity_view_models.EntityRoleMenuView.ViewModel(entity_view_models.EntityRoleMenuView{}),
		serviceRepository: NewServiceRepository(db),
	}
}

func (db *roleMenuConnection) GetAll(filter map[string]interface{}) []models.RoleMenu {
	var records []models.RoleMenu
	if len(filter) == 0 {
		db.connection.Find(&records)
	} else if len(filter) != 0 {
		db.connection.Where(filter).Find(&records)
	}
	return records
}

func (db *roleMenuConnection) Insert(record []models.RoleMenu) helper.Response {
	tx := db.connection.Begin()

	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
	} else {
		tx.Commit()
		db.connection.Find(&record)
		return helper.ServerResponse(true, "Ok", "", record)
	}
}

func (db *roleMenuConnection) Update(record models.RoleMenu) helper.Response {
	var oldRecord models.RoleMenu
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

func (db *roleMenuConnection) GetById(recordId string) helper.Response {
	var record models.RoleMenu
	db.connection.Preload("Emiten").First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	res := helper.ServerResponse(true, "Ok", "", record)
	return res
}

func (db *roleMenuConnection) DeleteById(recordId string) helper.Response {
	var record models.RoleMenu
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

func (db *roleMenuConnection) DeleteByRoleAndMenuId(roleId string, applicationMenuId string, isParent bool) helper.Response {
	tx := db.connection.Begin()

	var menuRecord models.ApplicationMenu
	db.connection.First(&menuRecord, "id = ?", applicationMenuId)
	if menuRecord.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	} else {
		if !isParent {
			//-- Delete Child
			if err := tx.Where("application_menu_id = ? AND role_id = ?", applicationMenuId, roleId).Delete(&models.RoleMenu{}).Error; err != nil {
				tx.Rollback()
				return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
			}
		}
		if isParent && menuRecord.ParentId == "" {
			var children []models.ApplicationMenu
			db.connection.Find(&children, "parent_id = ?", applicationMenuId)

			//-- Delete Parent
			if err := tx.Where("application_menu_id = ? AND role_id = ?", menuRecord.Id, roleId).Delete(&models.RoleMenu{}).Error; err != nil {
				tx.Rollback()
				return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
			}

			//-- Delete Children
			if len(children) > 0 {
				fmt.Println("has children")
				fmt.Println(len(children))

				for _, v := range children {
					fmt.Println(v.Id)
					if err := tx.Where("application_menu_id = ? AND role_id = ?", v.Id, roleId).Delete(&models.RoleMenu{}).Error; err != nil {
						tx.Rollback()
						return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
					}
				}
			}
		}
	}
	tx.Commit()
	return helper.ServerResponse(true, "Ok", "", helper.EmptyObj{})
}
