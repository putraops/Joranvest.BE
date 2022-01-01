package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type SectorService interface {
	Lookup(request helper.ReactSelectRequest) helper.Response
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetPagination(request commons.PaginationRequest) interface{}
	GetAll(filter map[string]interface{}) []models.Sector
	Insert(record models.Sector) helper.Response
	Update(record models.Sector) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type sectorService struct {
	sectorRepository repository.SectorRepository
	helper.AppSession
}

func NewSectorService(repo repository.SectorRepository) SectorService {
	return &sectorService{
		sectorRepository: repo,
	}
}

func (service *sectorService) Lookup(r helper.ReactSelectRequest) helper.Response {
	var ary helper.ReactSelectResponse

	result := service.sectorRepository.Lookup(r)
	if len(result) > 0 {
		for _, record := range result {
			var p = helper.ReactSelectItem{
				Value:    record.Id,
				Label:    record.Name,
				ParentId: "",
			}
			ary.Results = append(ary.Results, p)
		}
	}
	ary.Count = len(result)
	return helper.ServerResponse(true, "Ok", "", ary)
}

func (service *sectorService) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	return service.sectorRepository.GetDatatables(request)
}

func (service *sectorService) GetPagination(request commons.PaginationRequest) interface{} {
	return service.sectorRepository.GetPagination(request)
}

func (service *sectorService) GetAll(filter map[string]interface{}) []models.Sector {
	return service.sectorRepository.GetAll(filter)
}

func (service *sectorService) Insert(record models.Sector) helper.Response {
	return service.sectorRepository.Insert(record)
}

func (service *sectorService) Update(record models.Sector) helper.Response {
	return service.sectorRepository.Update(record)
}

func (service *sectorService) GetById(recordId string) helper.Response {
	return service.sectorRepository.GetById(recordId)
}

func (service *sectorService) DeleteById(recordId string) helper.Response {
	return service.sectorRepository.DeleteById(recordId)
}
