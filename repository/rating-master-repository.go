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
	"gorm.io/gorm/clause"
)

type RatingMasterRepository interface {
	GetAll(filter map[string]interface{}) []models.RatingMaster
	Insert(t models.RatingMaster) helper.Response
	Update(record models.RatingMaster) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type ratingMasterConnection struct {
	connection        *gorm.DB
	serviceRepository ServiceRepository
	tableName         string
	viewQuery         string
}

func NewRatingMasterRepository(db *gorm.DB) RatingMasterRepository {
	return &ratingMasterConnection{
		connection:        db,
		tableName:         models.RatingMaster.TableName(models.RatingMaster{}),
		viewQuery:         entity_view_models.EntityRatingMasterView.ViewModel(entity_view_models.EntityRatingMasterView{}),
		serviceRepository: NewServiceRepository(db),
	}
}

func (db *ratingMasterConnection) GetAll(filter map[string]interface{}) []models.RatingMaster {
	var records []models.RatingMaster
	if len(filter) == 0 {
		db.connection.Find(&records)
	} else if len(filter) != 0 {
		db.connection.Where(filter).Find(&records)
	}
	return records
}

func (db *ratingMasterConnection) Insert(record models.RatingMaster) helper.Response {
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

func (db *ratingMasterConnection) Update(record models.RatingMaster) helper.Response {
	var oldRecord models.RatingMaster
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

func (db *ratingMasterConnection) GetById(recordId string) helper.Response {
	var record models.RatingMaster
	db.connection.Preload("Emiten").First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	res := helper.ServerResponse(true, "Ok", "", record)
	return res
}

func (db *ratingMasterConnection) DeleteById(recordId string) helper.Response {
	var record models.RatingMaster
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
