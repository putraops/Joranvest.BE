package repository

import (
	"database/sql"
	"fmt"
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/entity_view_models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FilemasterRepository interface {
	GetAll(filter map[string]interface{}) []models.Filemaster
	GetAllByRecordIds(ids []string) []models.Filemaster
	SingleUpload(t models.Filemaster) helper.Response
	UploadByType(t models.Filemaster) helper.Response
	UploadProfilePicture(t models.Filemaster) helper.Response
	Insert(t models.Filemaster) helper.Response
	DeleteById(recordId string) helper.Response
	DeleteByRecordId(recordId string) helper.Response
}

type filemasterConnection struct {
	connection          *gorm.DB
	serviceRepository   ServiceRepository
	applicationUserRepo ApplicationUserRepository
	tableName           string
	viewQuery           string
}

func NewFilemasterRepository(db *gorm.DB) FilemasterRepository {
	return &filemasterConnection{
		connection:          db,
		tableName:           models.Filemaster.TableName(models.Filemaster{}),
		viewQuery:           entity_view_models.EntityFilemasterView.ViewModel(entity_view_models.EntityFilemasterView{}),
		serviceRepository:   NewServiceRepository(db),
		applicationUserRepo: NewApplicationUserRepository(db),
	}
}

func (db *filemasterConnection) GetAll(filter map[string]interface{}) []models.Filemaster {
	var records []models.Filemaster
	if len(filter) == 0 {
		db.connection.Find(&records)
	} else if len(filter) != 0 {
		db.connection.Where(filter).Find(&records)
	}

	return records
}

func (db *filemasterConnection) GetAllByRecordIds(ids []string) []models.Filemaster {
	var records []models.Filemaster
	db.connection.Where("record_id IN ?", ids).Find(&records)
	return records
}

func (db *filemasterConnection) SingleUpload(record models.Filemaster) helper.Response {
	var filemasterRecord models.Filemaster
	tx := db.connection.Begin()

	if err := tx.Where("record_id = ?", record.RecordId).Delete(&filemasterRecord).Error; err != nil {
		tx.Rollback()
		return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
	}

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

func (db *filemasterConnection) UploadByType(record models.Filemaster) helper.Response {
	var filemasterRecord models.Filemaster
	tx := db.connection.Begin()

	if err := tx.Where("record_id = ? AND file_type = ?", record.RecordId, record.FileType).Delete(&filemasterRecord).Error; err != nil {
		tx.Rollback()
		return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
	}

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

func (db *filemasterConnection) UploadProfilePicture(record models.Filemaster) helper.Response {
	var filemasterRecord models.Filemaster
	tx := db.connection.Begin()

	if err := tx.Where("record_id = ? AND file_type = 1", record.RecordId, record.FileType).Delete(&filemasterRecord).Error; err != nil {
		tx.Rollback()
		return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
	}

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

func (db *filemasterConnection) Insert(record models.Filemaster) helper.Response {
	tx := db.connection.Begin()

	record.Id = uuid.New().String()
	record.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	record.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
	} else {
		tx.Commit()
		return helper.ServerResponse(true, "Ok", "", record)
	}
}
func (db *filemasterConnection) DeleteById(id string) helper.Response {
	var filepath string
	var record models.Filemaster
	db.connection.First(&record, "id = ?", id)

	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	} else {
		filepath = record.Filepath
		res := db.connection.Delete(&record)
		if res.RowsAffected == 0 {
			return helper.ServerResponse(false, "Error", fmt.Sprintf("%v", res.Error), helper.EmptyObj{})
		}
		return helper.ServerResponse(true, "Ok", "", filepath)
	}
}

func (db *filemasterConnection) DeleteByRecordId(recordId string) helper.Response {
	var record models.Filemaster
	db.connection.First(&record, "record_id = ?", recordId)

	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	} else {
		res := db.connection.Where("record_id = ?", recordId).Delete(&record)
		if res.RowsAffected == 0 {
			return helper.ServerResponse(false, "Error", fmt.Sprintf("%v", res.Error), helper.EmptyObj{})
		}
		return helper.ServerResponse(true, "Ok", "", helper.EmptyObj{})
	}
}
