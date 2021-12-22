package service

import (
	"joranvest/helper"
	entity_view_models "joranvest/models/entity_view_models"
	"joranvest/repository"
)

type ArticleTagService interface {
	GetById(recordId string) helper.Response
	GetAll(filter map[string]interface{}) []entity_view_models.EntityArticleTagView
}

type articleTagService struct {
	articleTagRepository repository.ArticleTagRepository
	helper.AppSession
}

func NewArticleTagService(repo repository.ArticleTagRepository) ArticleTagService {
	return &articleTagService{
		articleTagRepository: repo,
	}
}

func (service *articleTagService) GetById(recordId string) helper.Response {
	return service.articleTagRepository.GetById(recordId)
}

func (service *articleTagService) GetAll(filter map[string]interface{}) []entity_view_models.EntityArticleTagView {
	return service.articleTagRepository.GetAll(filter)
}
