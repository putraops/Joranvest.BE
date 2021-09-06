package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type ApplicationMenuService interface {
	Lookup(request helper.Select2Request) helper.Response
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetAll(filter map[string]interface{}) []models.ApplicationMenu
	GetTree() []commons.JStreeResponse
	GetTreeByRoleId(roleId string) []commons.JStreeResponse
	OrderTree(recordId string, parentId string, orderIndex int) helper.Response
	Insert(record models.ApplicationMenu) helper.Response
	Update(record models.ApplicationMenu) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type applicationMenuService struct {
	applicationMenuRepository repository.ApplicationMenuRepository
	helper.AppSession
}

func NewApplicationMenuService(repo repository.ApplicationMenuRepository) ApplicationMenuService {
	return &applicationMenuService{
		applicationMenuRepository: repo,
	}
}

func (service *applicationMenuService) Lookup(r helper.Select2Request) helper.Response {
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

	result := service.applicationMenuRepository.Lookup(request)
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

func (service *applicationMenuService) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	return service.applicationMenuRepository.GetDatatables(request)
}

func (service *applicationMenuService) GetTree() []commons.JStreeResponse {
	return service.applicationMenuRepository.GetTree()
}

func (service *applicationMenuService) GetTreeByRoleId(roleId string) []commons.JStreeResponse {
	return service.applicationMenuRepository.GetTreeByRoleId(roleId)
}

func (service *applicationMenuService) OrderTree(recordId string, parentId string, orderIndex int) helper.Response {
	return service.applicationMenuRepository.OrderTree(recordId, parentId, orderIndex)
}

func (service *applicationMenuService) GetAll(filter map[string]interface{}) []models.ApplicationMenu {
	return service.applicationMenuRepository.GetAll(filter)
}

func (service *applicationMenuService) Insert(record models.ApplicationMenu) helper.Response {
	return service.applicationMenuRepository.Insert(record)
}

func (service *applicationMenuService) Update(record models.ApplicationMenu) helper.Response {
	return service.applicationMenuRepository.Update(record)
}

func (service *applicationMenuService) GetById(recordId string) helper.Response {
	return service.applicationMenuRepository.GetById(recordId)
}

func (service *applicationMenuService) DeleteById(recordId string) helper.Response {
	return service.applicationMenuRepository.DeleteById(recordId)
}
