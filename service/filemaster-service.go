package service

import (
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type FilemasterService interface {
	GetAll(filter map[string]interface{}) []models.Filemaster
	Insert(record models.Filemaster) helper.Response
	DeleteByRecordId(recordId string) helper.Response
}

type filemasterService struct {
	repo repository.FilemasterRepository
	helper.AppSession
}

func NewFilemasterService(repo repository.FilemasterRepository) FilemasterService {
	return &filemasterService{
		repo: repo,
	}
}

func (service *filemasterService) GetAll(filter map[string]interface{}) []models.Filemaster {
	return service.repo.GetAll(filter)
}

func (service *filemasterService) Insert(record models.Filemaster) helper.Response {
	return service.repo.Insert(record)
}

func (service *filemasterService) DeleteByRecordId(recordId string) helper.Response {
	return service.repo.DeleteByRecordId(recordId)
}
