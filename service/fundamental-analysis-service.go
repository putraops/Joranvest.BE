package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type FundamentalAnalysisService interface {
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetAll(filter map[string]interface{}) []models.FundamentalAnalysis
	Insert(record models.FundamentalAnalysis) helper.Response
	Update(record models.FundamentalAnalysis) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type fundamentalAnalysisService struct {
	fundamentalAnalysisRepository repository.FundamentalAnalysisRepository
	helper.AppSession
}

func NewFundamentalAnalysisService(repo repository.FundamentalAnalysisRepository) FundamentalAnalysisService {
	return &fundamentalAnalysisService{
		fundamentalAnalysisRepository: repo,
	}
}

func (service *fundamentalAnalysisService) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	return service.fundamentalAnalysisRepository.GetDatatables(request)
}

func (service *fundamentalAnalysisService) GetAll(filter map[string]interface{}) []models.FundamentalAnalysis {
	return service.fundamentalAnalysisRepository.GetAll(filter)
}

func (service *fundamentalAnalysisService) Insert(record models.FundamentalAnalysis) helper.Response {
	return service.fundamentalAnalysisRepository.Insert(record)
}

func (service *fundamentalAnalysisService) Update(record models.FundamentalAnalysis) helper.Response {
	return service.fundamentalAnalysisRepository.Update(record)
}

func (service *fundamentalAnalysisService) GetById(recordId string) helper.Response {
	return service.fundamentalAnalysisRepository.GetById(recordId)
}

func (service *fundamentalAnalysisService) DeleteById(recordId string) helper.Response {
	return service.fundamentalAnalysisRepository.DeleteById(recordId)
}
