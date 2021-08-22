package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type WebinarCategoryService interface {
	Lookup(request helper.Select2Request) helper.Response
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetAll(filter map[string]interface{}) []models.WebinarCategory
	Insert(record models.WebinarCategory) helper.Response
	Update(record models.WebinarCategory) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type webinarCategoryService struct {
	webinarCategoryRepository repository.WebinarCategoryRepository
	helper.AppSession
}

func NewWebinarCategoryService(repo repository.WebinarCategoryRepository) WebinarCategoryService {
	return &webinarCategoryService{
		webinarCategoryRepository: repo,
	}
}

func (service *webinarCategoryService) Lookup(r helper.Select2Request) helper.Response {
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

	result := service.webinarCategoryRepository.Lookup(request)
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

func (service *webinarCategoryService) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	return service.webinarCategoryRepository.GetDatatables(request)
}

func (service *webinarCategoryService) GetAll(filter map[string]interface{}) []models.WebinarCategory {
	return service.webinarCategoryRepository.GetAll(filter)
}

func (service *webinarCategoryService) Insert(record models.WebinarCategory) helper.Response {
	return service.webinarCategoryRepository.Insert(record)
}

func (service *webinarCategoryService) Update(record models.WebinarCategory) helper.Response {
	return service.webinarCategoryRepository.Update(record)
}

func (service *webinarCategoryService) GetById(recordId string) helper.Response {
	return service.webinarCategoryRepository.GetById(recordId)
}

func (service *webinarCategoryService) DeleteById(recordId string) helper.Response {
	return service.webinarCategoryRepository.DeleteById(recordId)
}
