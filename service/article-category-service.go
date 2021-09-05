package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type ArticleCategoryService interface {
	Lookup(request helper.Select2Request) helper.Response
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetAll(filter map[string]interface{}) []models.ArticleCategory
	GetTree() []commons.JStreeResponse
	Insert(record models.ArticleCategory) helper.Response
	Update(record models.ArticleCategory) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type articleCategoryService struct {
	articleCategoryRepository repository.ArticleCategoryRepository
	helper.AppSession
}

func NewArticleCategoryService(repo repository.ArticleCategoryRepository) ArticleCategoryService {
	return &articleCategoryService{
		articleCategoryRepository: repo,
	}
}

func (service *articleCategoryService) Lookup(r helper.Select2Request) helper.Response {
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

	result := service.articleCategoryRepository.Lookup(request)

	if len(result) > 0 {
		for _, record := range result {
			if record.ParentId == "" {
				children := FindArticleCategoryChildren(record.Id, result)
				var hasChildren = false
				if len(children) > 0 {
					hasChildren = true
				}
				var p = helper.Select2Item{
					Id:          record.Id,
					Text:        record.Name,
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

func (service *articleCategoryService) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	return service.articleCategoryRepository.GetDatatables(request)
}

func (service *articleCategoryService) GetTree() []commons.JStreeResponse {
	return service.articleCategoryRepository.GetTree()
}

func (service *articleCategoryService) GetAll(filter map[string]interface{}) []models.ArticleCategory {
	return service.articleCategoryRepository.GetAll(filter)
}

func (service *articleCategoryService) Insert(record models.ArticleCategory) helper.Response {
	return service.articleCategoryRepository.Insert(record)
}

func (service *articleCategoryService) Update(record models.ArticleCategory) helper.Response {
	return service.articleCategoryRepository.Update(record)
}

func (service *articleCategoryService) GetById(recordId string) helper.Response {
	return service.articleCategoryRepository.GetById(recordId)
}

func (service *articleCategoryService) DeleteById(recordId string) helper.Response {
	return service.articleCategoryRepository.DeleteById(recordId)
}

func FindArticleCategoryChildren(parent_id string, records []models.ArticleCategory) []helper.Select2Item {
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
