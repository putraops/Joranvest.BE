package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type EmitenCategoryService interface {
	Lookup(request helper.Select2Request) helper.Response
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetPagination(request commons.PaginationRequest) interface{}
	GetAll(filter map[string]interface{}) []models.EmitenCategory
	Insert(record models.EmitenCategory) helper.Response
	Update(record models.EmitenCategory) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type emitenCategoryService struct {
	emitenCategoryRepository repository.EmitenCategoryRepository
	helper.AppSession
}

func NewEmitenCategoryService(repo repository.EmitenCategoryRepository) EmitenCategoryService {
	return &emitenCategoryService{
		emitenCategoryRepository: repo,
	}
}

func (service *emitenCategoryService) Lookup(r helper.Select2Request) helper.Response {
	var ary helper.Select2Response

	request := make(map[string]interface{})
	request["qry"] = r.Q
	request["condition"] = helper.DataFilter{
		Request: []helper.Operator{
			{
				Operator: "like",
				Field:    r.Field,
				Value:    r.Q,
			},
		},
	}

	result := service.emitenCategoryRepository.Lookup(request)
	if len(result) > 0 {
		for _, record := range result {
			var p = helper.Select2Item{
				Id:          record.Id,
				Value:       record.Id,
				Text:        record.Name,
				Label:       record.Name,
				Description: record.Description,
				Selected:    true,
				Disabled:    false,
			}
			ary.Results = append(ary.Results, p)
		}
	}
	ary.Count = len(result)
	return helper.ServerResponse(true, "Ok", "", ary)
}

func (service *emitenCategoryService) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	return service.emitenCategoryRepository.GetDatatables(request)
}

func (service *emitenCategoryService) GetPagination(request commons.PaginationRequest) interface{} {
	return service.emitenCategoryRepository.GetPagination(request)
}

func (service *emitenCategoryService) GetAll(filter map[string]interface{}) []models.EmitenCategory {
	return service.emitenCategoryRepository.GetAll(filter)
}

func (service *emitenCategoryService) Insert(record models.EmitenCategory) helper.Response {
	return service.emitenCategoryRepository.Insert(record)
}

func (service *emitenCategoryService) Update(record models.EmitenCategory) helper.Response {
	return service.emitenCategoryRepository.Update(record)
}

func (service *emitenCategoryService) GetById(recordId string) helper.Response {
	return service.emitenCategoryRepository.GetById(recordId)
}

func (service *emitenCategoryService) DeleteById(recordId string) helper.Response {
	return service.emitenCategoryRepository.DeleteById(recordId)
}
