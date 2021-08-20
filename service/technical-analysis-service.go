package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type TechnicalAnalysisService interface {
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetAll(filter map[string]interface{}) []models.TechnicalAnalysis
	Insert(record models.TechnicalAnalysis) helper.Response
	Update(record models.TechnicalAnalysis) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type technicalAnalysisService struct {
	technicalAnalysisRepository repository.TechnicalAnalysisRepository
	helper.AppSession
}

func NewTechnicalAnalysisService(repo repository.TechnicalAnalysisRepository) TechnicalAnalysisService {
	return &technicalAnalysisService{
		technicalAnalysisRepository: repo,
	}
}

func (service *technicalAnalysisService) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	return service.technicalAnalysisRepository.GetDatatables(request)
}

func (service *technicalAnalysisService) GetAll(filter map[string]interface{}) []models.TechnicalAnalysis {
	return service.technicalAnalysisRepository.GetAll(filter)
}

func (service *technicalAnalysisService) Insert(record models.TechnicalAnalysis) helper.Response {
	return service.technicalAnalysisRepository.Insert(record)
}

func (service *technicalAnalysisService) Update(record models.TechnicalAnalysis) helper.Response {
	return service.technicalAnalysisRepository.Update(record)
}

func (service *technicalAnalysisService) GetById(recordId string) helper.Response {
	return service.technicalAnalysisRepository.GetById(recordId)
}

func (service *technicalAnalysisService) DeleteById(recordId string) helper.Response {
	return service.technicalAnalysisRepository.DeleteById(recordId)
}
