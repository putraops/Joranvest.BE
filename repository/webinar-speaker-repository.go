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

type WebinarSpeakerRepository interface {
	Insert(records []models.WebinarSpeaker, speakerType int) helper.Response
	GetById(recordId string) helper.Response
	GetAll(filter map[string]interface{}) []entity_view_models.EntityWebinarSpeakerView
}

type webinarSpeakerConnection struct {
	connection        *gorm.DB
	serviceRepository ServiceRepository
	tableName         string
	viewQuery         string
}

func NewWebinarSpeakerRepository(db *gorm.DB) WebinarSpeakerRepository {
	return &webinarSpeakerConnection{
		connection:        db,
		tableName:         models.FundamentalAnalysis.TableName(models.FundamentalAnalysis{}),
		viewQuery:         entity_view_models.EntityWebinarSpeakerView.ViewModel(entity_view_models.EntityWebinarSpeakerView{}),
		serviceRepository: NewServiceRepository(db),
	}
}

func (db *webinarSpeakerConnection) Insert(records []models.WebinarSpeaker, speakerType int) helper.Response {
	tx := db.connection.Begin()

	if len(records) > 0 {
		if err := tx.Where("webinar_id = ?", records[0].WebinarId).Delete(models.WebinarSpeaker{}).Error; err != nil {
			tx.Rollback()
			return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
		}

		for i := 0; i < len(records); i++ {
			records[i].Id = uuid.New().String()
			records[i].CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
		}

		if err := tx.Create(&records).Error; err != nil {
			tx.Rollback()
			return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
		}
	} else {
		return helper.ServerResponse(false, "Speaker cannot be empty", "Speaker cannot be empty", helper.EmptyObj{})
	}

	var webinarRecord models.Webinar
	db.connection.First(&webinarRecord, "id = ?", records[0].WebinarId)
	if webinarRecord.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	webinarRecord.SpeakerType = speakerType
	res := tx.Save(&webinarRecord)
	if res.RowsAffected == 0 {
		return helper.ServerResponse(false, fmt.Sprintf("%v,", res.Error), fmt.Sprintf("%v,", res.Error), helper.EmptyObj{})
	}

	tx.Commit()
	return helper.ServerResponse(true, "Ok", "", records)
}

func (db *webinarSpeakerConnection) GetById(recordId string) helper.Response {
	var record models.FundamentalAnalysis
	db.connection.Preload("Emiten").First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	res := helper.ServerResponse(true, "Ok", "", record)
	return res
}

func (db *webinarSpeakerConnection) GetAll(filter map[string]interface{}) []entity_view_models.EntityWebinarSpeakerView {
	var records []entity_view_models.EntityWebinarSpeakerView
	if len(filter) == 0 {
		db.connection.Find(&records)
	} else if len(filter) != 0 {
		db.connection.Where(filter).Find(&records)
	}
	return records
}
