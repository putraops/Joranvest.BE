package service

import (
	"joranvest/commons"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"

	"github.com/mashingan/smapping"
	"gorm.io/gorm"
)

type EducationCategoryService interface {
	GetPagination(request commons.Pagination2ndRequest) interface{}
	GetAll(filter map[string]interface{}) []models.EducationCategory
	Lookup(request helper.ReactSelectRequest) helper.Result
	Save(dto dto.EducationCategoryDto) helper.Result
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response

	OpenTransaction(*gorm.DB) educationCategoryService
}

type educationCategoryService struct {
	DB                          *gorm.DB
	jwtService                  JWTService
	educationCategoryRepository repository.EducationCategoryRepository
}

func NewEducationCategoryService(db *gorm.DB, jwtService JWTService) EducationCategoryService {
	return educationCategoryService{
		DB:                          db,
		jwtService:                  jwtService,
		educationCategoryRepository: repository.NewEducationCategoryRepository(db),
	}
}

func (r educationCategoryService) GetPagination(request commons.Pagination2ndRequest) interface{} {
	return r.educationCategoryRepository.GetPagination(request)
}

func (r educationCategoryService) GetAll(filter map[string]interface{}) []models.EducationCategory {
	return r.educationCategoryRepository.GetAll(filter)
}

func (r educationCategoryService) Lookup(request helper.ReactSelectRequest) helper.Result {
	var ary helper.ReactSelectResponse

	result := r.educationCategoryRepository.Lookup(request)
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

func (r educationCategoryService) Save(dto dto.EducationCategoryDto) helper.Result {
	authHeader := dto.Context.GetHeader("Authorization")
	userIdentity := r.jwtService.GetUserByToken(authHeader)

	//-- Map Dto to Struct
	var newRecord models.EducationCategory
	smapping.FillStruct(&newRecord, smapping.MapFields(&dto))
	newRecord.EntityId = userIdentity.UserId

	if dto.Id == "" {
		newRecord.CreatedBy = userIdentity.UserId
		return r.educationCategoryRepository.Insert(newRecord)
	} else {
		newRecord.UpdatedBy = userIdentity.UserId
		return r.educationCategoryRepository.Update(newRecord)
	}
}

func (r educationCategoryService) GetById(recordId string) helper.Response {
	return r.educationCategoryRepository.GetById(recordId)
}

func (r educationCategoryService) DeleteById(recordId string) helper.Response {
	return r.educationCategoryRepository.DeleteById(recordId)
}

func (r educationCategoryService) OpenTransaction(trxHandle *gorm.DB) educationCategoryService {
	r.educationCategoryRepository = r.educationCategoryRepository.OpenTransaction(trxHandle)
	return r
}
