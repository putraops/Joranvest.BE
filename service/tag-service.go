package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type TagService interface {
	Lookup(request helper.ReactSelectRequest) helper.Response
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetAll(filter map[string]interface{}) []models.Tag
	Insert(record models.Tag) helper.Response
	Update(record models.Tag) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type tagService struct {
	tagRepository repository.TagRepository
	helper.AppSession
}

func NewTagService(repo repository.TagRepository) TagService {
	return &tagService{
		tagRepository: repo,
	}
}

func (service *tagService) Lookup(r helper.ReactSelectRequest) helper.Response {
	var ary helper.ReactSelectResponse

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

	result := service.tagRepository.Lookup(request)
	if len(result) > 0 {
		for _, record := range result {
			var p = helper.ReactSelectItem{
				Value: record.Id,
				Label: record.Name,
			}
			ary.Results = append(ary.Results, p)
		}
	}
	ary.Count = len(result)
	return helper.ServerResponse(true, "Ok", "", ary)
}

func (service *tagService) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	return service.tagRepository.GetDatatables(request)
}

func (service *tagService) GetAll(filter map[string]interface{}) []models.Tag {
	return service.tagRepository.GetAll(filter)
}

func (service *tagService) Insert(record models.Tag) helper.Response {
	return service.tagRepository.Insert(record)
}

func (service *tagService) Update(record models.Tag) helper.Response {
	return service.tagRepository.Update(record)
}

func (service *tagService) GetById(recordId string) helper.Response {
	return service.tagRepository.GetById(recordId)
}

func (service *tagService) DeleteById(recordId string) helper.Response {
	return service.tagRepository.DeleteById(recordId)
}
