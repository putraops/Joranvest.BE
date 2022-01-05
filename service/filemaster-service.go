package service

import (
	"joranvest/helper"
	"joranvest/models"
	"joranvest/models/request_models"
	"joranvest/repository"
)

type FilemasterService interface {
	GetAll(filter map[string]interface{}) []models.Filemaster
	GetAllByRecordIds(ids []string) []models.Filemaster
	SingleUpload(record models.Filemaster) helper.Response
	UploadByType(record models.Filemaster) helper.Response
	UploadProfilePicture(request request_models.FileRequestDto) helper.Response
	Insert(record models.Filemaster) helper.Response
	UpdateWebinarCoverImage(record models.Webinar) helper.Response
	DeleteById(recordId string) helper.Response
	DeleteByRecordId(recordId string) helper.Response
	GetDirectoryConfig(moduleName string, moduleId string, filetype int) string
}

type filemasterService struct {
	repo                      repository.FilemasterRepository
	webinarRepository         repository.WebinarRepository
	applicationUserRepository repository.ApplicationUserRepository
	organizationRepository    repository.OrganizationRepository
	helper.AppSession
}

func NewFilemasterService(
	repo repository.FilemasterRepository,
	webinarRepo repository.WebinarRepository,
	applicationUserRepo repository.ApplicationUserRepository,
	organizationRepo repository.OrganizationRepository) FilemasterService {
	return &filemasterService{
		repo:                      repo,
		webinarRepository:         webinarRepo,
		applicationUserRepository: applicationUserRepo,
		organizationRepository:    organizationRepo,
	}
}

func (service *filemasterService) GetAll(filter map[string]interface{}) []models.Filemaster {
	return service.repo.GetAll(filter)
}

func (service *filemasterService) GetAllByRecordIds(ids []string) []models.Filemaster {
	return service.repo.GetAllByRecordIds(ids)
}

func (service *filemasterService) SingleUpload(record models.Filemaster) helper.Response {
	return service.repo.SingleUpload(record)
}

func (service *filemasterService) UploadByType(record models.Filemaster) helper.Response {
	return service.repo.UploadByType(record)
}

func (service *filemasterService) UploadProfilePicture(request request_models.FileRequestDto) helper.Response {
	return service.applicationUserRepository.UpdateProfilePicture(request)
}

func (service *filemasterService) Insert(record models.Filemaster) helper.Response {
	return service.repo.Insert(record)
}

func (service *filemasterService) UpdateWebinarCoverImage(record models.Webinar) helper.Response {
	return service.webinarRepository.UpdateCoverImage(record)
}

func (service *filemasterService) DeleteById(id string) helper.Response {
	return service.repo.DeleteById(id)
}

func (service *filemasterService) DeleteByRecordId(recordId string) helper.Response {
	return service.repo.DeleteByRecordId(recordId)
}

func (service *filemasterService) GetDirectoryConfig(moduleName string, moduleId string, filetype int) string {
	var res = ""
	if filetype == 1 {
		res = "upload/" + moduleName + "/" + moduleId + "/cover/"
	} else if filetype == 2 {
		res = "upload/" + moduleName + "/" + moduleId + "/files/"
	}
	return res
}
