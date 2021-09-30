package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type MembershipUserService interface {
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetPagination(request commons.PaginationRequest) interface{}
	GetAll(filter map[string]interface{}) []models.MembershipUser
	Insert(record models.MembershipUser, payment models.MembershipPayment) helper.Response
	Update(record models.MembershipUser) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type membershipUserService struct {
	membershipUserRepository repository.MembershipUserRepository
	helper.AppSession
}

func NewMembershipUserService(repo repository.MembershipUserRepository) MembershipUserService {
	return &membershipUserService{
		membershipUserRepository: repo,
	}
}

func (service *membershipUserService) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	return service.membershipUserRepository.GetDatatables(request)
}

func (service *membershipUserService) GetPagination(request commons.PaginationRequest) interface{} {
	return service.membershipUserRepository.GetPagination(request)
}

func (service *membershipUserService) GetAll(filter map[string]interface{}) []models.MembershipUser {
	return service.membershipUserRepository.GetAll(filter)
}

func (service *membershipUserService) Insert(record models.MembershipUser, payment models.MembershipPayment) helper.Response {
	return service.membershipUserRepository.Insert(record, payment)
}

func (service *membershipUserService) Update(record models.MembershipUser) helper.Response {
	return service.membershipUserRepository.Update(record)
}

func (service *membershipUserService) GetById(recordId string) helper.Response {
	return service.membershipUserRepository.GetById(recordId)
}

func (service *membershipUserService) DeleteById(recordId string) helper.Response {
	return service.membershipUserRepository.DeleteById(recordId)
}
