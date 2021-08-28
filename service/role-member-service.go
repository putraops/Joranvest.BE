package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/view_models"
	"joranvest/repository"
)

type RoleMemberService interface {
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetAll(filter map[string]interface{}) []models.RoleMember
	GetUsersInRole(roleId string) []entity_view_models.EntityRoleMemberView
	GetUsersNotInRole(roleId string) []models.ApplicationUser
	Insert(record models.RoleMember) helper.Response
	Update(record models.RoleMember) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type roleMemberService struct {
	roleMemberRepository repository.RoleMemberRepository
	helper.AppSession
}

func NewRoleMemberService(repo repository.RoleMemberRepository) RoleMemberService {
	return &roleMemberService{
		roleMemberRepository: repo,
	}
}

func (service *roleMemberService) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	return service.roleMemberRepository.GetDatatables(request)
}

func (service *roleMemberService) GetAll(filter map[string]interface{}) []models.RoleMember {
	return service.roleMemberRepository.GetAll(filter)
}

func (service *roleMemberService) GetUsersInRole(roleId string) []entity_view_models.EntityRoleMemberView {
	result := service.roleMemberRepository.GetUsersInRole(roleId)
	return result
}

func (service *roleMemberService) GetUsersNotInRole(roleId string) []models.ApplicationUser {
	result := service.roleMemberRepository.GetUsersNotInRole(roleId)
	return result
}

func (service *roleMemberService) Insert(record models.RoleMember) helper.Response {
	return service.roleMemberRepository.Insert(record)
}

func (service *roleMemberService) Update(record models.RoleMember) helper.Response {
	return service.roleMemberRepository.Update(record)
}

func (service *roleMemberService) GetById(recordId string) helper.Response {
	return service.roleMemberRepository.GetById(recordId)
}

func (service *roleMemberService) DeleteById(recordId string) helper.Response {
	return service.roleMemberRepository.DeleteById(recordId)
}
