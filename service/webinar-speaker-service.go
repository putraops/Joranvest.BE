package service

import (
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/entity_view_models"
	"joranvest/repository"
)

type WebinarSpeakerService interface {
	Insert(records []models.WebinarSpeaker, speakerType int) helper.Response
	GetById(recordId string) helper.Response
	GetAll(filter map[string]interface{}) []entity_view_models.EntityWebinarSpeakerView
	GetSpeakersRatingByWebinarId(webinarId string) helper.Response
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

func (service *webinarSpeakerService) Insert(records []models.WebinarSpeaker, speakerType int) helper.Response {
	return service.webinarSpeakerRepository.Insert(records, speakerType)
}

func (service *webinarSpeakerService) GetById(recordId string) helper.Response {
	return service.webinarSpeakerRepository.GetById(recordId)
}

func (service *webinarSpeakerService) GetAll(filter map[string]interface{}) []entity_view_models.EntityWebinarSpeakerView {
	return service.webinarSpeakerRepository.GetAll(filter)
}

func (service *webinarSpeakerService) GetSpeakersRatingByWebinarId(webinarId string) helper.Response {
	return service.webinarSpeakerRepository.GetSpeakersRatingByWebinarId(webinarId)
}
