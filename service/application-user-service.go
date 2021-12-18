package service

import (
	"joranvest/commons"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/view_models"
	"joranvest/repository"

	"log"

	"github.com/mashingan/smapping"
)

//-- This is user contract
type ApplicationUserService interface {
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	Lookup(request helper.ReactSelectRequest) helper.Response
	Update(user dto.ApplicationUserUpdateDto) models.ApplicationUser
	UserProfile(recordId string) models.ApplicationUser
	ChangePassword(recordDto dto.ChangePasswordDto) helper.Response
	GetById(recordId string) helper.Response
	GetViewById(recordId string) helper.Response
	GetAll() []models.ApplicationUser
	DeleteById(recordId string) helper.Response
	RecoverPassword(recordId string, oldPassword string) helper.Response
	EmailVerificationById(recordId string) helper.Response
}

type applicationUserService struct {
	applicationUserRepository repository.ApplicationUserRepository
	emailService              EmailService
}

func NewApplicationUserService(repo repository.ApplicationUserRepository, emailService EmailService) ApplicationUserService {
	return &applicationUserService{
		applicationUserRepository: repo,
		emailService:              emailService,
	}
}

func (service *applicationUserService) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	return service.applicationUserRepository.GetDatatables(request)
}

func (service *applicationUserService) Lookup(r helper.ReactSelectRequest) helper.Response {
	var ary helper.ReactSelectResponse

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

	result := service.applicationUserRepository.Lookup(request)
	if len(result) > 0 {
		for _, record := range result {
			var p = helper.ReactSelectItem{
				Value:    record.Id,
				Label:    record.FirstName + " " + record.LastName,
				ParentId: "",
			}
			ary.Results = append(ary.Results, p)
		}
	}
	ary.Count = len(result)
	return helper.ServerResponse(true, "Ok", "", ary)
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

func (service *applicationUserService) GetViewById(recordId string) helper.Response {
	return service.applicationUserRepository.GetViewById(recordId)
}

func (service *applicationUserService) GetAll() []models.ApplicationUser {
	result := service.applicationUserRepository.GetAll()
	return result
}

func (service *applicationUserService) DeleteById(userId string) helper.Response {
	return service.applicationUserRepository.DeleteById(userId)
}

func (service *applicationUserService) ChangePassword(recordDto dto.ChangePasswordDto) helper.Response {
	res := service.applicationUserRepository.GetViewUserByUsernameOrEmail(recordDto.Username, recordDto.Email)
	if res == nil {
		return helper.ServerResponse(false, "Record not found", "NotFound", helper.EmptyObj{})
	}
	if v, ok := res.(entity_view_models.EntityApplicationUserView); ok {
		comparedPassword := comparePassword(v.Password, []byte(recordDto.OldPassword))
		if (v.Email == recordDto.Email || v.Username == recordDto.Username) && comparedPassword {
			user := (service.applicationUserRepository.GetById(v.Id).Data).(models.ApplicationUser)
			user.Password = recordDto.NewPassword

			newUserData := service.applicationUserRepository.Update(user)
			return helper.ServerResponse(true, "Ok", "", newUserData)

		} else {
			return helper.ServerResponse(false, "Password is not match", "Error", helper.EmptyObj{})
		}
	}
	return helper.ServerResponse(true, "Ok", "", helper.EmptyObj{})
}

func (service *applicationUserService) RecoverPassword(recordId string, oldPassword string) helper.Response {
	return service.applicationUserRepository.RecoverPassword(recordId, oldPassword)
}

func (service *applicationUserService) EmailVerificationById(recordId string) helper.Response {

	res := service.applicationUserRepository.EmailVerificationById(recordId)
	if res.Status {
		var record = res.Data.(models.ApplicationUser)
		to := []string{record.Email}
		service.emailService.SendEmailVerified(to)
	}
	return res
}
