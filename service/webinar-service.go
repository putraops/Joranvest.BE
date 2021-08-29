package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type WebinarService interface {
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetAll(filter map[string]interface{}) []models.Webinar
	Insert(record models.Webinar) helper.Response
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

func (service *webinarService) GetAll(filter map[string]interface{}) []models.Webinar {
	return service.webinarRepository.GetAll(filter)
}

func (service *webinarService) Insert(record models.Webinar) helper.Response {
	return service.webinarRepository.Insert(record)
}

func (service *webinarService) Update(record models.Webinar) helper.Response {
	return service.webinarRepository.Update(record)
}

func (service *webinarService) GetById(recordId string) helper.Response {
	return service.webinarRepository.GetById(recordId)
}

func (service *webinarService) DeleteById(recordId string) helper.Response {
	return service.webinarRepository.DeleteById(recordId)
}
