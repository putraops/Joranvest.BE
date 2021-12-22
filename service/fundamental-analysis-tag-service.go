package service

import (
	"joranvest/helper"
	entity_view_models "joranvest/models/entity_view_models"
	"joranvest/repository"
)

type FundamentalAnalysisTagService interface {
	GetById(recordId string) helper.Response
	GetAll(filter map[string]interface{}) []entity_view_models.EntityFundamentalAnalysisTagView
}

type fundamentalAnalysisTagService struct {
	fundamentalAnalysisTagRepository repository.FundamentalAnalysisTagRepository
	helper.AppSession
}

func NewFundamentalAnalysisTagService(repo repository.FundamentalAnalysisTagRepository) FundamentalAnalysisTagService {
	return &fundamentalAnalysisTagService{
		fundamentalAnalysisTagRepository: repo,
	}
}

func (service *fundamentalAnalysisTagService) GetById(recordId string) helper.Response {
	return service.fundamentalAnalysisTagRepository.GetById(recordId)
}

func (service *fundamentalAnalysisTagService) GetAll(filter map[string]interface{}) []entity_view_models.EntityFundamentalAnalysisTagView {
	return service.fundamentalAnalysisTagRepository.GetAll(filter)
}
