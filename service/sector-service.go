package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type SectorService interface {
	Lookup(request helper.Select2Request) helper.Response
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
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

func (service *sectorService) Lookup(r helper.Select2Request) helper.Response {
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

	result := service.sectorRepository.Lookup(request)
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

func (service *sectorService) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	return service.sectorRepository.GetDatatables(request)
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
