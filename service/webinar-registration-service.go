package service

import (
	"joranvest/commons"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type WebinarRegistrationService interface {
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetPagination(request commons.PaginationRequest) interface{}
	GetAll(filter map[string]interface{}) []models.WebinarRegistration
	GetById(recordId string) helper.Response
	GetViewById(recordId string) helper.Response
	IsWebinarRegistered(webinarId string, userId string) helper.Response
	Insert(record models.WebinarRegistration) helper.Response
	Update(record models.WebinarRegistration) helper.Response
	UpdatePayment(dto dto.WebinarRegistrationUpdatePaymentDto) helper.Response
	DeleteById(recordId string) helper.Response
}

type webinarRegistrationService struct {
	webinarRegistrationRepository repository.WebinarRegistrationRepository
	helper.AppSession
}

func NewWebinarRegistrationService(repo repository.WebinarRegistrationRepository) WebinarRegistrationService {
	return &webinarRegistrationService{
		webinarRegistrationRepository: repo,
	}
}

func (service *webinarRegistrationService) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	return service.webinarRegistrationRepository.GetDatatables(request)
}

func (service *webinarRegistrationService) GetPagination(request commons.PaginationRequest) interface{} {
	return service.webinarRegistrationRepository.GetPagination(request)
}

func (service *webinarRegistrationService) GetAll(filter map[string]interface{}) []models.WebinarRegistration {
	return service.webinarRegistrationRepository.GetAll(filter)
}

func (service *webinarRegistrationService) Insert(record models.WebinarRegistration) helper.Response {
	return service.webinarRegistrationRepository.Insert(record)
}

func (service *webinarRegistrationService) Update(record models.WebinarRegistration) helper.Response {
	return service.webinarRegistrationRepository.Update(record)
}

func (service *webinarRegistrationService) UpdatePayment(dto dto.WebinarRegistrationUpdatePaymentDto) helper.Response {
	return service.webinarRegistrationRepository.UpdatePayment(dto)
}

func (service *webinarRegistrationService) GetById(recordId string) helper.Response {
	return service.webinarRegistrationRepository.GetById(recordId)
}

func (service *webinarRegistrationService) GetViewById(recordId string) helper.Response {
	return service.webinarRegistrationRepository.GetViewById(recordId)
}

func (service *webinarRegistrationService) IsWebinarRegistered(webinarId string, userId string) helper.Response {
	return service.webinarRegistrationRepository.IsWebinarRegistered(webinarId, userId)
}

func (service *webinarRegistrationService) DeleteById(recordId string) helper.Response {
	return service.webinarRegistrationRepository.DeleteById(recordId)
}
