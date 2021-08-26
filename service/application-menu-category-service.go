package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type ApplicationMenuCategoryService interface {
	Lookup(request helper.Select2Request) helper.Response
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetAll(filter map[string]interface{}) []models.ApplicationMenuCategory
	GetTree() []commons.JStreeResponse
	Insert(record models.ApplicationMenuCategory) helper.Response
	Update(record models.ApplicationMenuCategory) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type applicationMenuCategoryService struct {
	applicationMenuCategoryRepository repository.ApplicationMenuCategoryRepository
	helper.AppSession
}

func NewApplicationMenuCategoryService(repo repository.ApplicationMenuCategoryRepository) ApplicationMenuCategoryService {
	return &applicationMenuCategoryService{
		applicationMenuCategoryRepository: repo,
	}
}

func (service *applicationMenuCategoryService) Lookup(r helper.Select2Request) helper.Response {
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

	result := service.applicationMenuCategoryRepository.Lookup(request)
	if len(result) > 0 {
		for _, record := range result {
			var p = helper.Select2Item{
				Id:          record.Id,
				Text:        record.Name,
				Description: record.Description,
				Selected:    true,
				Disabled:    false,
			}
			ary.Results = append(ary.Results, p)
		}
	} else {
		var p = helper.Select2Item{
			Id:          "",
			Text:        "No result found",
			Description: "",
			Selected:    true,
			Disabled:    true,
		}
		ary.Results = append(ary.Results, p)
	}
	ary.Count = len(result)
	return helper.ServerResponse(true, "Ok", "", ary)
}

func (service *applicationMenuCategoryService) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	return service.applicationMenuCategoryRepository.GetDatatables(request)
}

func (service *applicationMenuCategoryService) GetTree() []commons.JStreeResponse {
	return service.applicationMenuCategoryRepository.GetTree()
}

func (service *applicationMenuCategoryService) GetAll(filter map[string]interface{}) []models.ApplicationMenuCategory {
	return service.applicationMenuCategoryRepository.GetAll(filter)
}

func (service *applicationMenuCategoryService) Insert(record models.ApplicationMenuCategory) helper.Response {
	return service.applicationMenuCategoryRepository.Insert(record)
}

func (service *applicationMenuCategoryService) Update(record models.ApplicationMenuCategory) helper.Response {
	return service.applicationMenuCategoryRepository.Update(record)
}

func (service *applicationMenuCategoryService) GetById(recordId string) helper.Response {
	return service.applicationMenuCategoryRepository.GetById(recordId)
}

func (service *applicationMenuCategoryService) DeleteById(recordId string) helper.Response {
	return service.applicationMenuCategoryRepository.DeleteById(recordId)
}
