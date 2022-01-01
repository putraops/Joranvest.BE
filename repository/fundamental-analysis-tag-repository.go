package repository

import (
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/entity_view_models"

	"gorm.io/gorm"
)

type FundamentalAnalysisTagRepository interface {
	GetById(recordId string) helper.Response
	GetAll(filter map[string]interface{}) []entity_view_models.EntityFundamentalAnalysisTagView
}

type fundamentalAnalysisTagConnection struct {
	connection        *gorm.DB
	serviceRepository ServiceRepository
	tableName         string
	viewQuery         string
}

func NewFundamentalAnalysisTagRepository(db *gorm.DB) FundamentalAnalysisTagRepository {
	return &fundamentalAnalysisTagConnection{
		connection:        db,
		tableName:         models.FundamentalAnalysis.TableName(models.FundamentalAnalysis{}),
		viewQuery:         entity_view_models.EntityFundamentalAnalysisTagView.ViewModel(entity_view_models.EntityFundamentalAnalysisTagView{}),
		serviceRepository: NewServiceRepository(db),
	}
}

func (db *fundamentalAnalysisTagConnection) GetById(recordId string) helper.Response {
	var record models.FundamentalAnalysis
	db.connection.Preload("Emiten").First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	res := helper.ServerResponse(true, "Ok", "", record)
	return res
}

func (db *fundamentalAnalysisTagConnection) GetAll(filter map[string]interface{}) []entity_view_models.EntityFundamentalAnalysisTagView {
	var records []entity_view_models.EntityFundamentalAnalysisTagView
	if len(filter) == 0 {
		db.connection.Find(&records)
	} else if len(filter) != 0 {
		db.connection.Where(filter).Find(&records)
	}

	return records
}
