package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type WebinarService interface {
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetPagination(request commons.Pagination2ndRequest) interface{}
	GetPaginationRegisteredByUser(request commons.Pagination2ndRequest, userId string) interface{}
	GetAll(filter map[string]interface{}) []models.Webinar
	Insert(record models.Webinar) helper.Response
	Submit(recordId string, userId string) helper.Response
	Update(record models.Webinar) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type webinarService struct {
	webinarRepository repository.WebinarRepository
	helper.AppSession
}

func NewWebinarService(repo repository.WebinarRepository) WebinarService {
	return &webinarService{
		webinarRepository: repo,
	}
}

func (service *webinarService) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	return service.webinarRepository.GetDatatables(request)
}

func (service *webinarService) GetPagination(request commons.Pagination2ndRequest) interface{} {
	return service.webinarRepository.GetPagination(request)
}

func (service *webinarService) GetPaginationRegisteredByUser(request commons.Pagination2ndRequest, userId string) interface{} {
	return service.webinarRepository.GetPaginationRegisteredByUser(request, userId)
}

func (service *webinarService) GetAll(filter map[string]interface{}) []models.Webinar {
	return service.webinarRepository.GetAll(filter)
}

func (service *webinarService) Insert(record models.Webinar) helper.Response {
	return service.webinarRepository.Insert(record)
}

func (service *webinarService) Update(record models.Webinar) helper.Response {
	return service.webinarRepository.Update(record)
}

func (service *webinarService) Submit(recordId string, userId string) helper.Response {
	return service.webinarRepository.Submit(recordId, userId)
}

func (service *webinarService) GetById(recordId string) helper.Response {
	return service.webinarRepository.GetById(recordId)
}

func (service *webinarService) DeleteById(recordId string) helper.Response {
	return service.webinarRepository.DeleteById(recordId)
}
