package service

import (
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type RatingMasterService interface {
	GetAll(filter map[string]interface{}) []models.RatingMaster
	Insert(record models.RatingMaster) helper.Response
	Update(record models.RatingMaster) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type ratingMasterService struct {
	ratingMasterRepository repository.RatingMasterRepository
	helper.AppSession
}

func NewRatingMasterService(repo repository.RatingMasterRepository) RatingMasterService {
	return &ratingMasterService{
		ratingMasterRepository: repo,
	}
}

func (service *ratingMasterService) GetAll(filter map[string]interface{}) []models.RatingMaster {
	return service.ratingMasterRepository.GetAll(filter)
}

func (service *ratingMasterService) Insert(record models.RatingMaster) helper.Response {
	return service.ratingMasterRepository.Insert(record)
}

func (service *ratingMasterService) Update(record models.RatingMaster) helper.Response {
	return service.ratingMasterRepository.Update(record)
}

func (service *ratingMasterService) GetById(recordId string) helper.Response {
	return service.ratingMasterRepository.GetById(recordId)
}

func (service *ratingMasterService) DeleteById(recordId string) helper.Response {
	return service.ratingMasterRepository.DeleteById(recordId)
}
