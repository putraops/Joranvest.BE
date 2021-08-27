package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type RoleMemberService interface {
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetAll(filter map[string]interface{}) []models.RoleMember
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
