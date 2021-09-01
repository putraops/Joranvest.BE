package service

import (
	"joranvest/helper"
	entity_view_models "joranvest/models/view_models"
	"joranvest/repository"
)

type WebinarSpeakerService interface {
	GetById(recordId string) helper.Response
	GetAll(filter map[string]interface{}) []entity_view_models.EntityWebinarSpeakerView
}

type webinarSpeakerService struct {
	webinarSpeakerRepository repository.WebinarSpeakerRepository
	helper.AppSession
}

func NewWebinarSpeakerService(repo repository.WebinarSpeakerRepository) WebinarSpeakerService {
	return &webinarSpeakerService{
		webinarSpeakerRepository: repo,
	}
}

func (service *webinarSpeakerService) GetById(recordId string) helper.Response {
	return service.webinarSpeakerRepository.GetById(recordId)
}

func (service *webinarSpeakerService) GetAll(filter map[string]interface{}) []entity_view_models.EntityWebinarSpeakerView {
	return service.webinarSpeakerRepository.GetAll(filter)
}
