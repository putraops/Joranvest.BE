package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type EmitenService interface {
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetAll(filter map[string]interface{}) []models.Emiten
	Lookup(request helper.Select2Request) helper.Response
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

func (service *emitenService) GetAll(filter map[string]interface{}) []models.Emiten {
	return service.emitenRepository.GetAll(filter)
}

func (service *emitenService) Lookup(r helper.Select2Request) helper.Response {
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

	result := service.emitenRepository.Lookup(request)
	if len(result) > 0 {
		for _, record := range result {
			var p = helper.Select2Item{
				Id:          record.Id,
				Text:        record.EmitenName + " [" + record.EmitenCode + "]",
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
