package repository

import (
	"database/sql"
	"fmt"
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/view_models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoleMenuRepository interface {
	GetAll(filter map[string]interface{}) []models.RoleMenu
	Insert(t models.RoleMenu) helper.Response
	Update(record models.RoleMenu) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
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

func (db *roleMenuConnection) Insert(record models.RoleMenu) helper.Response {
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
