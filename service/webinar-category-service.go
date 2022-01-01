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
	GetTreeParent() []commons.JStreeResponse
	GetTree() []commons.JStreeResponse
	OrderTree(recordId string, parentId string, orderIndex int) helper.Response
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
			if record.ParentId == "" {
				children := FindChildren(record.Id, result)
				var hasChildren = false
				if len(children) > 0 {
					hasChildren = true
				}
				var p = helper.Select2Item{
					Id:          record.Id,
					Value:       record.Id,
					Text:        record.Name,
					Label:       record.Name,
					Description: record.Description,
					ParentId:    record.ParentId,
					Selected:    false,
					Disabled:    false,
					HasChildren: hasChildren,
					Children:    children,
				}
				ary.Results = append(ary.Results, p)
			}
		}
	}
	ary.Count = len(result)
	return helper.ServerResponse(true, "Ok", "", ary)
}

func (service *webinarCategoryService) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	return service.webinarCategoryRepository.GetDatatables(request)
}

func (service *webinarCategoryService) GetTreeParent() []commons.JStreeResponse {
	return service.webinarCategoryRepository.GetTreeParent()
}

func (service *webinarCategoryService) GetTree() []commons.JStreeResponse {
	return service.webinarCategoryRepository.GetTree()
}

func (service *webinarCategoryService) OrderTree(recordId string, parentId string, orderIndex int) helper.Response {
	return service.webinarCategoryRepository.OrderTree(recordId, parentId, orderIndex)
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

func FindChildren(parent_id string, records []models.WebinarCategory) []helper.Select2Item {
	res := []helper.Select2Item{}

	for _, v := range records {
		if v.ParentId == parent_id {
			var p = helper.Select2Item{
				Id:          v.Id,
				Text:        v.Name,
				Description: v.Description,
				ParentId:    parent_id,
				Selected:    true,
				Disabled:    false,
				Children:    nil,
			}
			res = append(res, p)
		}
	}
	return res
}
