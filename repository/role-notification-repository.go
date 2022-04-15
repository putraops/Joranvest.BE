package repository

import (
	"fmt"
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/entity_view_models"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RoleNotificationRepository interface {
	SetNotification(record models.RoleNotification) helper.Result
	GetRoleById(roleId string) helper.Result
	DeleteByRoleId(roleId string) helper.Result

	OpenTransaction(trxHandle *gorm.DB) roleNotificationRepository
}

type roleNotificationRepository struct {
	DB                *gorm.DB
	serviceRepository ServiceRepository
	tableName         string
	viewQuery         string
	currentTime       time.Time
}

func NewRoleNotificationRepository(db *gorm.DB) RoleNotificationRepository {
	return &roleNotificationRepository{
		tableName:         models.RoleNotification.TableName(models.RoleNotification{}),
		viewQuery:         entity_view_models.EntityRoleNotificationView.ViewModel(entity_view_models.EntityRoleNotificationView{}),
		DB:                db,
		serviceRepository: NewServiceRepository(db),
		currentTime:       time.Now(),
	}
}

func (r roleNotificationRepository) SetNotification(record models.RoleNotification) helper.Result {
	tx := r.DB.Begin()

	if record.Id == nil {
		newId := uuid.New().String()
		record.Id = &newId
		record.CreatedAt = &r.currentTime
		if err := r.DB.Create(&record).Error; err != nil {
			log.Error(r.serviceRepository.getCurrentFuncName())
			log.Error(fmt.Sprintf("%v,", err.Error()))
			tx.Rollback()
			return helper.StandartResult(false, fmt.Sprintf("%v,", err.Error()), nil)
		}
	} else {
		record.UpdatedAt = &r.currentTime
		if err := r.DB.Save(&record).Error; err != nil {
			log.Error(r.serviceRepository.getCurrentFuncName())
			log.Error(fmt.Sprintf("%v,", err.Error()))
			tx.Rollback()
			return helper.StandartResult(false, fmt.Sprintf("%v,", err.Error()), nil)
		}
	}
	tx.Commit()
	return helper.StandartResult(true, "Ok", record)
}

func (r roleNotificationRepository) GetRoleById(roleId string) helper.Result {
	var record models.RoleNotification
	r.DB.First(&record, "role_id = ?", roleId)
	if record.Id == nil {
		res := helper.StandartResult(false, "Record not found", nil)
		return res
	}
	res := helper.StandartResult(true, "Ok", record)
	return res
}

func (r roleNotificationRepository) DeleteByRoleId(roleId string) helper.Result {
	var record models.RoleNotification
	tx := r.DB.Begin()

	if err := r.DB.Where("role_id = ?", roleId).Delete(&record).Error; err != nil {
		log.Error(r.serviceRepository.getCurrentFuncName())
		log.Error(fmt.Sprintf("%v,", err.Error()))
		tx.Rollback()
		return helper.StandartResult(false, fmt.Sprintf("%v", err.Error()), nil)
	}

	tx.Commit()
	return helper.StandartResult(true, "Ok", nil)
}

func (r roleNotificationRepository) OpenTransaction(trxHandle *gorm.DB) roleNotificationRepository {
	if trxHandle == nil {
		log.Error("Transaction Database not found")
		log.Error(r.serviceRepository.getCurrentFuncName())
		return r
	}
	r.DB = trxHandle
	return r
}
