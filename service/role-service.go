package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type RoleService interface {
	Lookup(request helper.Select2Request) helper.Response
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetAll(filter map[string]interface{}) []models.Role
	Insert(record models.Role) helper.Response
	Update(record models.Role) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type roleService struct {
	roleRepository repository.RoleRepository
	helper.AppSession
}

func NewRoleService(repo repository.RoleRepository) RoleService {
	return &roleService{
		roleRepository: repo,
	}
}

func (service *roleService) Lookup(r helper.Select2Request) helper.Response {
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

	result := service.roleRepository.Lookup(request)
	if len(result) > 0 {
		for _, record := range result {
			var p = helper.Select2Item{
				Id:          record.Id,
				Text:        record.Name,
				Description: "",
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

func (service *roleService) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	return service.roleRepository.GetDatatables(request)
}

func (service *roleService) GetAll(filter map[string]interface{}) []models.Role {
	return service.roleRepository.GetAll(filter)
}

func (service *roleService) Insert(record models.Role) helper.Response {
	return service.roleRepository.Insert(record)
}

func (service *roleService) Update(record models.Role) helper.Response {
	return service.roleRepository.Update(record)
}

func (service *roleService) GetById(recordId string) helper.Response {
	return service.roleRepository.GetById(recordId)
}

func (service *roleService) DeleteById(recordId string) helper.Response {
	return service.roleRepository.DeleteById(recordId)
}
