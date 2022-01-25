package service

import (
	"database/sql"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
	"time"

	"github.com/google/uuid"
)

type RoleMenuService interface {
	GetAll(filter map[string]interface{}) []models.RoleMenu
	Insert(record dto.InsertRoleMenuDto) helper.Response
	Update(record models.RoleMenu) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
	DeleteByRoleAndMenuId(roleId string, applicationMenuId string, isParent bool) helper.Response
}

type roleMenuService struct {
	roleMenuRepository repository.RoleMenuRepository
	helper.AppSession
}

func NewRoleMenuService(repo repository.RoleMenuRepository) RoleMenuService {
	return &roleMenuService{
		roleMenuRepository: repo,
	}
}

func (service *roleMenuService) GetAll(filter map[string]interface{}) []models.RoleMenu {
	return service.roleMenuRepository.GetAll(filter)
}

func (service *roleMenuService) Insert(record dto.InsertRoleMenuDto) helper.Response {
	var records []models.RoleMenu

	var newRecord models.RoleMenu
	newRecord.Id = uuid.New().String()
	newRecord.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	newRecord.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	newRecord.ApplicationMenuId = record.ApplicationMenuId
	newRecord.RoleId = record.RoleId
	records = append(records, newRecord)

	if len(record.ChildrenIds) > 0 {
		for _, v := range record.ChildrenIds {
			var newRecord models.RoleMenu
			newRecord.Id = uuid.New().String()
			newRecord.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
			newRecord.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
			newRecord.ApplicationMenuId = v
			newRecord.RoleId = record.RoleId
			records = append(records, newRecord)
		}
	}
	return service.roleMenuRepository.Insert(records)
}

func (service *roleMenuService) Update(record models.RoleMenu) helper.Response {
	return service.roleMenuRepository.Update(record)
}

func (service *roleMenuService) GetById(recordId string) helper.Response {
	return service.roleMenuRepository.GetById(recordId)
}

func (service *roleMenuService) DeleteById(recordId string) helper.Response {
	return service.roleMenuRepository.DeleteById(recordId)
}

func (service *roleMenuService) DeleteByRoleAndMenuId(roleId string, applicationMenuId string, isParent bool) helper.Response {
	return service.roleMenuRepository.DeleteByRoleAndMenuId(roleId, applicationMenuId, isParent)
}
