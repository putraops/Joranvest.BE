package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type WebinarRegistrationService interface {
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetPagination(request commons.Pagination2ndRequest) interface{}
	GetAll(filter map[string]interface{}) []models.WebinarRegistration
	GetById(recordId string) helper.Response
	GetViewById(recordId string) helper.Response
	Insert(record models.WebinarRegistration) helper.Response
	Update(record models.WebinarRegistration) helper.Response
	IsWebinarRegistered(webinarId string, userId string) helper.Response
	DeleteById(recordId string) helper.Response

	SendWebinarInformationViaEmail(webinarId string) helper.Response
}

type webinarRegistrationService struct {
	webinarRegistrationRepository repository.WebinarRegistrationRepository
	emailService                  EmailService
	helper.AppSession
}

func NewWebinarRegistrationService(repo repository.WebinarRegistrationRepository, emailService EmailService) WebinarRegistrationService {
	return &webinarRegistrationService{
		webinarRegistrationRepository: repo,
		emailService:                  emailService,
	}
}

func (service *webinarRegistrationService) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	return service.webinarRegistrationRepository.GetDatatables(request)
}

func (service *webinarRegistrationService) GetPagination(request commons.Pagination2ndRequest) interface{} {
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

func (service *webinarRegistrationService) SendWebinarInformationViaEmail(webinarId string) helper.Response {
	var result helper.Response
	var participants = service.webinarRegistrationRepository.GetParticipantsByWebinarId(webinarId)
	if len(participants) > 0 {
		var temp []string
		for _, item := range participants {
			temp = append(temp, item.UserEmail)
		}
		result = service.emailService.SendWebinarInformationToParticipants(temp)
	} else {
		result.Data = helper.EmptyObj{}
		result.Errors = helper.EmptyObj{}
		result.Message = "There is no Participants"
		result.Status = false
	}
	return result
}
