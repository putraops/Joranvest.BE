package repository

import (
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/view_models"

	"gorm.io/gorm"
)

type WebinarSpeakerRepository interface {
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
