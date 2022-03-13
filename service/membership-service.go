package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type MembershipService interface {
	Lookup(request helper.ReactSelectRequest) helper.Result
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetPagination(request commons.Pagination2ndRequest) interface{}
	GetAll(filter map[string]interface{}) []models.Membership
	Insert(record models.Membership) helper.Response
	Update(record models.Membership) helper.Response
	SetRecommendationById(recordId string, isChecked bool) helper.Response
	GetById(recordId string) helper.Response
	GetViewById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type membershipService struct {
	membershipRepository repository.MembershipRepository
	helper.AppSession
}

func NewMembershipService(repo repository.MembershipRepository) MembershipService {
	return &membershipService{
		membershipRepository: repo,
	}
}

func (service *membershipService) Lookup(request helper.ReactSelectRequest) helper.Result {
	var ary helper.ReactSelectResponse

	result := service.membershipRepository.Lookup(request)
	if len(result) > 0 {
		for _, record := range result {
			var p = helper.ReactSelectItem{
				Value:    record.Id,
				Label:    record.Name,
				ParentId: "",
			}
			ary.Results = append(ary.Results, p)
		}
	}
	ary.Count = len(result)
	return helper.StandartResult(true, "Ok", ary)
}

func (service *membershipService) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	return service.membershipRepository.GetDatatables(request)
}

func (service *membershipService) GetPagination(request commons.Pagination2ndRequest) interface{} {
	return service.membershipRepository.GetPagination(request)
}

func (service *membershipService) GetAll(filter map[string]interface{}) []models.Membership {
	return service.membershipRepository.GetAll(filter)
}

func (service *membershipService) Insert(record models.Membership) helper.Response {
	return service.membershipRepository.Insert(record)
}

func (service *membershipService) Update(record models.Membership) helper.Response {
	return service.membershipRepository.Update(record)
}

func (service *membershipService) SetRecommendationById(recordId string, isChecked bool) helper.Response {
	return service.membershipRepository.SetRecomendationById(recordId, isChecked)
}

func (service *membershipService) GetById(recordId string) helper.Response {
	return service.membershipRepository.GetById(recordId)
}

func (service *membershipService) GetViewById(recordId string) helper.Response {
	return service.membershipRepository.GetViewById(recordId)
}

func (service *membershipService) DeleteById(recordId string) helper.Response {
	return service.membershipRepository.DeleteById(recordId)
}
