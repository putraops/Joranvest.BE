package service

import (
	"joranvest/commons"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"

	"gorm.io/gorm"
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

	SendInvitation(request dto.SendWebinarInformationDto) helper.Response
	SendWebinarInformationViaEmail(request dto.SendWebinarInformationDto) helper.Response
}

type webinarRegistrationService struct {
	webinarRegistrationRepository repository.WebinarRegistrationRepository
	emailService                  EmailService
	helper.AppSession
	DB *gorm.DB
}

func NewWebinarRegistrationService(db *gorm.DB) WebinarRegistrationService {
	return &webinarRegistrationService{
		DB:                            db,
		webinarRegistrationRepository: repository.NewWebinarRegistrationRepository(db),
		emailService:                  NewEmailService(db),
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

func (service *webinarRegistrationService) SendInvitation(request dto.SendWebinarInformationDto) helper.Response {
	var result helper.Response
	var participants = service.webinarRegistrationRepository.GetParticipantsByIds(request.WebinarRegistrationIds)
	if len(participants) > 0 {
		for _, item := range participants {
			service.emailService.SendWebinarInformationToParticipants(request, item)
			service.webinarRegistrationRepository.UpdateInvitationStatusById(item.Id)
		}
		result.Data = helper.EmptyObj{}
		result.Errors = helper.EmptyObj{}
		result.Message = "Email Sent"
		result.Status = true
	} else {
		result.Data = helper.EmptyObj{}
		result.Errors = helper.EmptyObj{}
		result.Message = "There is no Participants"
		result.Status = false
	}
	return result
}

func (service *webinarRegistrationService) SendWebinarInformationViaEmail(request dto.SendWebinarInformationDto) helper.Response {
	var result helper.Response
	var participants = service.webinarRegistrationRepository.GetParticipantsByWebinarId(request.WebinarId)
	if len(participants) > 0 {
		for _, item := range participants {
			service.emailService.SendWebinarInformationToParticipants(request, item)
		}
		result.Data = helper.EmptyObj{}
		result.Errors = helper.EmptyObj{}
		result.Message = "Email Sent"
		result.Status = true
	} else {
		result.Data = helper.EmptyObj{}
		result.Errors = helper.EmptyObj{}
		result.Message = "There is no Participants"
		result.Status = false
	}
	return result
}
