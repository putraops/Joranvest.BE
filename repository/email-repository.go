package repository

import (
	"joranvest/models"
	entity_view_models "joranvest/models/entity_view_models"

	"gorm.io/gorm"
)

type EmailRepository interface {
}

type emailConnection struct {
	connection                *gorm.DB
	serviceRepository         ServiceRepository
	applicationUserRepository ApplicationUserRepository
	tableName                 string
	viewQuery                 string
}

func NewEmailRepository(db *gorm.DB) EmailRepository {
	return &emailConnection{
		connection:                db,
		tableName:                 models.Emiten.TableName(models.Emiten{}),
		viewQuery:                 entity_view_models.EntityEmitenView.ViewModel(entity_view_models.EntityEmitenView{}),
		serviceRepository:         NewServiceRepository(db),
		applicationUserRepository: NewApplicationUserRepository(db),
	}
}
