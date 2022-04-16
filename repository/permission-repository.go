package repository

import (
	"joranvest/models"
	"time"

	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type PermissionRepository interface {
	GetAll(filter map[string]interface{}) []models.Team

	OpenTransaction(trxHandle *gorm.DB) permissionRepository
}

type permissionRepository struct {
	DB                *gorm.DB
	currentTime       time.Time
	serviceRepository ServiceRepository
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return permissionRepository{
		DB:                db,
		serviceRepository: NewServiceRepository(db),
		currentTime:       time.Now(),
	}
}

func (r permissionRepository) GetAll(filter map[string]interface{}) []models.Team {
	var records []models.Team
	if len(filter) == 0 {
		r.DB.Find(&records)
	} else if len(filter) != 0 {
		r.DB.Where(filter).Find(&records)
	}
	return records
}

func (r permissionRepository) OpenTransaction(trxHandle *gorm.DB) permissionRepository {
	if trxHandle == nil {
		log.Error("Transaction Database not found")
		log.Error(r.serviceRepository.getCurrentFuncName())
		return r
	}
	r.DB = trxHandle
	return r
}
