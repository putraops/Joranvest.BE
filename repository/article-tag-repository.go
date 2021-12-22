package repository

import (
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/entity_view_models"

	"gorm.io/gorm"
)

type ArticleTagRepository interface {
	GetById(recordId string) helper.Response
	GetAll(filter map[string]interface{}) []entity_view_models.EntityArticleTagView
}

type articleTagConnection struct {
	connection        *gorm.DB
	serviceRepository ServiceRepository
	tableName         string
	viewQuery         string
}

func NewArticleTagRepository(db *gorm.DB) ArticleTagRepository {
	return &articleTagConnection{
		connection:        db,
		tableName:         models.ArticleTag.TableName(models.ArticleTag{}),
		viewQuery:         entity_view_models.EntityArticleTagView.ViewModel(entity_view_models.EntityArticleTagView{}),
		serviceRepository: NewServiceRepository(db),
	}
}

func (db *articleTagConnection) GetById(recordId string) helper.Response {
	var record models.ArticleTag
	db.connection.Preload("Emiten").First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	res := helper.ServerResponse(true, "Ok", "", record)
	return res
}

func (db *articleTagConnection) GetAll(filter map[string]interface{}) []entity_view_models.EntityArticleTagView {
	var records []entity_view_models.EntityArticleTagView
	if len(filter) == 0 {
		db.connection.Find(&records)
	} else if len(filter) != 0 {
		db.connection.Where(filter).Find(&records)
	}

	return records
}
