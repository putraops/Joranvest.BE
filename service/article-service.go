package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type ArticleService interface {
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetAll(filter map[string]interface{}) []models.Article
	Insert(record models.Article) helper.Response
	Update(record models.Article) helper.Response
	Submit(recordId string, userId string) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type articleService struct {
	articleRepository repository.ArticleRepository
	helper.AppSession
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{
		articleRepository: repo,
	}
}

func (service *articleService) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	return service.articleRepository.GetDatatables(request)
}

func (service *articleService) GetAll(filter map[string]interface{}) []models.Article {
	return service.articleRepository.GetAll(filter)
}

func (service *articleService) Insert(record models.Article) helper.Response {
	return service.articleRepository.Insert(record)
}

func (service *articleService) Update(record models.Article) helper.Response {
	return service.articleRepository.Update(record)
}

func (service *articleService) Submit(recordId string, userId string) helper.Response {
	return service.articleRepository.Submit(recordId, userId)
}

func (service *articleService) GetById(recordId string) helper.Response {
	return service.articleRepository.GetById(recordId)
}

func (service *articleService) DeleteById(recordId string) helper.Response {
	return service.articleRepository.DeleteById(recordId)
}
