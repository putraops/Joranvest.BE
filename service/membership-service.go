package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type MembershipService interface {
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetAll(filter map[string]interface{}) []models.Membership
	Insert(record models.Membership) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type membershipService struct {
	membershipRepository repository.MembershipRepository
	// orderDetailRepository repository.OrderDetailRepository
	// paymentRepository     repository.PaymentRepository
	// entityRepository      repository.EntityRepository
	helper.AppSession
}

func NewMembershipService(repo repository.MembershipRepository) MembershipService {
	return &membershipService{
		membershipRepository: repo,
	}
}

func (service *membershipService) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	return service.membershipRepository.GetDatatables(request)
}

func (service *membershipService) GetAll(filter map[string]interface{}) []models.Membership {
	return service.membershipRepository.GetAll(filter)
}

func (service *membershipService) Insert(record models.Membership) helper.Response {
	return service.membershipRepository.Insert(record)
}

func (service *membershipService) GetById(recordId string) helper.Response {
	return service.membershipRepository.GetById(recordId)
}

func (service *membershipService) DeleteById(recordId string) helper.Response {
	return service.membershipRepository.DeleteById(recordId)
}