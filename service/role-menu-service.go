package service

import (
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type RoleMenuService interface {
	GetAll(filter map[string]interface{}) []models.RoleMenu
	Insert(record models.RoleMenu) helper.Response
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

func (service *roleMenuService) Insert(record models.RoleMenu) helper.Response {
	return service.roleMenuRepository.Insert(record)
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
