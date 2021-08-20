package service

import (
	"joranvest/commons"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
	"log"

	"github.com/mashingan/smapping"
)

//-- This is user contract
type ApplicationUserService interface {
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	Update(user dto.ApplicationUserUpdateDto) models.ApplicationUser
	UserProfile(recordId string) models.ApplicationUser
	GetById(recordId string) helper.Response
	GetAll() []models.ApplicationUser
	DeleteById(recordId string) helper.Response
}

type applicationUserService struct {
	applicationUserRepository repository.ApplicationUserRepository
}

func NewApplicationUserService(repo repository.ApplicationUserRepository) ApplicationUserService {
	return &applicationUserService{
		applicationUserRepository: repo,
	}
}

func (service *applicationUserService) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	return service.applicationUserRepository.GetDatatables(request)
}

func (service *applicationUserService) Update(record dto.ApplicationUserUpdateDto) models.ApplicationUser {
	recordToUpdate := models.ApplicationUser{}
	err := smapping.FillStruct(&recordToUpdate, smapping.MapFields(record))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}
	updatedRecord := service.applicationUserRepository.Update(recordToUpdate)
	return updatedRecord
}

func (service *applicationUserService) UserProfile(recordId string) models.ApplicationUser {
	return service.applicationUserRepository.UserProfile(recordId)
}

func (service *applicationUserService) GetById(recordId string) helper.Response {
	return service.applicationUserRepository.GetById(recordId)
}

func (service *applicationUserService) GetAll() []models.ApplicationUser {
	result := service.applicationUserRepository.GetAll()
	return result
}

func (service *applicationUserService) DeleteById(userId string) helper.Response {
	return service.applicationUserRepository.DeleteById(userId)
}
