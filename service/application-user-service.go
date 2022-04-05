package service

import (
	"joranvest/commons"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/entity_view_models"
	"joranvest/repository"

	"log"

	"github.com/mashingan/smapping"
	"gorm.io/gorm"
)

//-- This is user contract
type ApplicationUserService interface {
	GetPagination(request commons.Pagination2ndRequest) interface{}
	Lookup(request helper.ReactSelectRequest) helper.Response
	UserLookup(r helper.ReactSelectRequest) helper.Result
	Update(user dto.ApplicationUserUpdateDto) models.ApplicationUser
	UserProfile(recordId string) models.ApplicationUser
	UpdateProfile(dtoRecord dto.ApplicationUserDescriptionDto) helper.Response
	ChangePhone(recordDto dto.ChangePhoneDto) helper.Response
	ChangePassword(recordDto dto.ChangePasswordDto) helper.Response
	GetById(recordId string) helper.Response
	GetViewById(recordId string) helper.Response
	GetAll() []models.ApplicationUser
	DeleteById(recordId string) helper.Response
	ResetPasswordByEmail(email string) helper.Response
	RecoverPassword(dto dto.RecoverPasswordDto) helper.Response
	EmailVerificationById(recordId string) helper.Response
}

type applicationUserService struct {
	DB                        *gorm.DB
	applicationUserRepository repository.ApplicationUserRepository
	emailLoggingRepository    repository.EmailLoggingRepository
	emailService              EmailService
}

// func NewApplicationUserService(db *gorm.DB, repo repository.ApplicationUserRepository, emailLoggingRepo repository.EmailLoggingRepository, emailService EmailService) ApplicationUserService {
func NewApplicationUserService(db *gorm.DB) ApplicationUserService {
	return &applicationUserService{
		DB:                        db,
		applicationUserRepository: repository.NewApplicationUserRepository(db),
		emailLoggingRepository:    repository.NewEmailLoggingRepository(db),
		emailService:              NewEmailService(db),
	}
}

func (service *applicationUserService) GetPagination(request commons.Pagination2ndRequest) interface{} {
	return service.applicationUserRepository.GetPagination(request)
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

func (service *applicationUserService) UserLookup(r helper.ReactSelectRequest) helper.Result {
	var ary helper.ReactSelectResponse

	result := service.applicationUserRepository.UserLookup(r)
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
	return helper.StandartResult(true, "Ok", ary)
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

func (service *applicationUserService) UpdateProfile(dtoRecord dto.ApplicationUserDescriptionDto) helper.Response {
	return service.applicationUserRepository.UpdateProfile(dtoRecord)
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

func (service *applicationUserService) ChangePhone(recordDto dto.ChangePhoneDto) helper.Response {
	res := service.applicationUserRepository.GetById(recordDto.Id)
	if !res.Status {
		return res
	}

	var userData = (res.Data).(models.ApplicationUser)
	userData.Password = "" //-- Set to empty :: use for restrict change password
	userData.Phone = recordDto.Phone

	_ = service.applicationUserRepository.Update(userData)
	return helper.ServerResponse(true, "Ok", "", helper.EmptyObj{})
}

func (service *applicationUserService) RecoverPassword(dto dto.RecoverPasswordDto) helper.Response {
	user := service.applicationUserRepository.GetById(dto.UserId)

	if user.Data.(models.ApplicationUser).Id == "" {
		return helper.Response{
			Status:  false,
			Message: "User tidak ditemukan. Gagal untuk mengubah Password",
			Data:    helper.EmptyObj{},
			Errors:  helper.EmptyObj{},
		}
	}

	return service.applicationUserRepository.RecoverPassword(user.Data.(models.ApplicationUser).Id, dto.NewPassword)
}

func (service *applicationUserService) EmailVerificationById(recordId string) helper.Response {

	res := service.applicationUserRepository.EmailVerificationById(recordId)
	if res.Status {
		var record = res.Data.(models.ApplicationUser)

		var total = service.emailLoggingRepository.GetLastIntervalLogging(record.Email, commons.MailTypeEmailVerified, commons.MailInterval)
		if total <= commons.MaxSendEmailOneInterval {
			service.emailService.SendEmailVerified(record.Email)
		}
	}
	return res
}

func (service *applicationUserService) ResetPasswordByEmail(email string) helper.Response {
	var response helper.Response
	user := service.applicationUserRepository.GetByEmail(email)
	if user.Id == "" {
		response.Status = false
		response.Message = "Email tidak terdaftar."
		response.Data = helper.EmptyObj{}
		return response
	}
	return service.emailService.ResetPassword(user)
}
