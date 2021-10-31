package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type OrganizationService interface {
	Lookup(request helper.Select2Request) helper.Response
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetPagination(request commons.PaginationRequest) interface{}
	GetAll(filter map[string]interface{}) []models.Organization
	Insert(record models.Organization) helper.Response
	Update(record models.Organization) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type organizationService struct {
	organizationRepository repository.OrganizationRepository
	helper.AppSession
}

func NewOrganizationService(repo repository.OrganizationRepository) OrganizationService {
	return &organizationService{
		organizationRepository: repo,
	}
}

func (service *organizationService) Lookup(r helper.Select2Request) helper.Response {
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

	result := service.organizationRepository.Lookup(request)
	if len(result) > 0 {
		for _, record := range result {
			var p = helper.Select2Item{
				Id:          record.Id,
				Value:       record.Id,
				Text:        record.Name,
				Label:       record.Name,
				Description: "",
				Selected:    true,
				Disabled:    false,
			}
			ary.Results = append(ary.Results, p)
		}
	}
	ary.Count = len(result)
	return helper.ServerResponse(true, "Ok", "", ary)
}

func (service *organizationService) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	return service.organizationRepository.GetDatatables(request)
}

func (service *organizationService) GetPagination(request commons.PaginationRequest) interface{} {
	return service.organizationRepository.GetPagination(request)
}

func (service *organizationService) GetAll(filter map[string]interface{}) []models.Organization {
	return service.organizationRepository.GetAll(filter)
}

func (service *organizationService) Insert(record models.Organization) helper.Response {
	return service.organizationRepository.Insert(record)
}

func (service *organizationService) Update(record models.Organization) helper.Response {
	return service.organizationRepository.Update(record)
}

func (service *organizationService) GetById(recordId string) helper.Response {
	return service.organizationRepository.GetById(recordId)
}

func (service *organizationService) DeleteById(recordId string) helper.Response {
	return service.organizationRepository.DeleteById(recordId)
}
