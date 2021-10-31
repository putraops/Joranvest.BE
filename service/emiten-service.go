package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type EmitenService interface {
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetPagination(request commons.PaginationRequest) interface{}
	GetAll(filter map[string]interface{}) []models.Emiten
	Lookup(request helper.ReactSelectRequest) helper.Response
	Insert(record models.Emiten) helper.Response
	Update(record models.Emiten) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type emitenService struct {
	emitenRepository repository.EmitenRepository
	helper.AppSession
}

func NewEmitenService(repo repository.EmitenRepository) EmitenService {
	return &emitenService{
		emitenRepository: repo,
	}
}

func (service *emitenService) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	return service.emitenRepository.GetDatatables(request)
}

func (service *emitenService) GetPagination(request commons.PaginationRequest) interface{} {
	return service.emitenRepository.GetPagination(request)
}

func (service *emitenService) GetAll(filter map[string]interface{}) []models.Emiten {
	return service.emitenRepository.GetAll(filter)
}

func (service *emitenService) Lookup(r helper.ReactSelectRequest) helper.Response {
	var ary helper.ReactSelectResponse

	result := service.emitenRepository.Lookup(r)
	if len(result) > 0 {
		for _, record := range result {
			var p = helper.ReactSelectItem{
				Value:    record.Id,
				Label:    record.EmitenName + " [" + record.EmitenCode + "]",
				ParentId: "",
			}
			ary.Results = append(ary.Results, p)
		}
	}
	ary.Count = len(result)
	return helper.ServerResponse(true, "Ok", "", ary)
}

func (service *emitenService) Insert(record models.Emiten) helper.Response {
	return service.emitenRepository.Insert(record)
}

func (service *emitenService) Update(record models.Emiten) helper.Response {
	return service.emitenRepository.Update(record)
}

func (service *emitenService) GetById(recordId string) helper.Response {
	return service.emitenRepository.GetById(recordId)
}

func (service *emitenService) DeleteById(recordId string) helper.Response {
	return service.emitenRepository.DeleteById(recordId)
}
