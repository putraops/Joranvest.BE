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

type EducationService interface {
	GetPagination(request commons.Pagination2ndRequest) interface{}
	GetAll(filter map[string]interface{}) []models.Education
	GetPlaylist(filter map[string]interface{}) []models.EducationPlaylist
	GetPlaylistById(recordId string) helper.Result
	Lookup(request helper.ReactSelectRequest) helper.Result
	Save(dto dto.EducationDto) helper.Result
	AddToPlaylist(dto dto.EducationPlaylistDto) helper.Result
	RemoveFromPlaylist(recordId string) helper.Result
	GetViewById(recordId string) helper.Result
	GetById(recordId string) helper.Result
	DeleteById(recordId string) helper.Result

	GetDirectoryConfig(moduleName string, moduleId string, filetype int) string
	OpenTransaction(*gorm.DB) educationService
}

type educationService struct {
	DB                          *gorm.DB
	jwtService                  JWTService
	educationRepository         repository.EducationRepository
	educationPlaylistRepository repository.EducationPlaylistRepository
}

func NewEducationService(db *gorm.DB, jwtService JWTService) EducationService {
	return educationService{
		DB:                          db,
		jwtService:                  jwtService,
		educationRepository:         repository.NewEducationRepository(db),
		educationPlaylistRepository: repository.NewEducationPlaylistRepository(db),
	}
}

func (r educationService) GetPagination(request commons.Pagination2ndRequest) interface{} {
	return r.educationRepository.GetPagination(request)
}

func (r educationService) GetAll(filter map[string]interface{}) []models.Education {
	return r.educationRepository.GetAll(filter)
}

func (r educationService) Lookup(request helper.ReactSelectRequest) helper.Result {
	var ary helper.ReactSelectResponse

	result := r.educationRepository.Lookup(request)
	if len(result) > 0 {
		for _, record := range result {
			var p = helper.ReactSelectItem{
				Value:    record.Id,
				Label:    record.Title,
				ParentId: "",
			}
			ary.Results = append(ary.Results, p)
		}
	}
	ary.Count = len(result)
	return helper.StandartResult(true, "Ok", ary)
}

func (r educationService) Save(dto dto.EducationDto) helper.Result {
	authHeader := dto.Context.GetHeader("Authorization")
	userIdentity := r.jwtService.GetUserByToken(authHeader)

	//-- Map Dto to Struct
	var newRecord models.Education
	smapping.FillStruct(&newRecord, smapping.MapFields(&dto))
	newRecord.EntityId = userIdentity.UserId

	if dto.Id == "" {
		newRecord.CreatedBy = userIdentity.UserId
		return r.educationRepository.Insert(newRecord)
	} else {
		newRecord.UpdatedBy = userIdentity.UserId
		return r.educationRepository.Update(newRecord)
	}
}

func (r educationService) GetPlaylist(filter map[string]interface{}) []models.EducationPlaylist {
	return r.educationPlaylistRepository.GetAll(filter)
}

func (r educationService) GetPlaylistById(recordId string) helper.Result {
	return r.educationPlaylistRepository.GetById(recordId)
}

func (r educationService) AddToPlaylist(dto dto.EducationPlaylistDto) helper.Result {
	authHeader := dto.Context.GetHeader("Authorization")
	userIdentity := r.jwtService.GetUserByToken(authHeader)

	//-- Map Dto to Struct
	var newRecord models.EducationPlaylist
	smapping.FillStruct(&newRecord, smapping.MapFields(&dto))
	newRecord.EntityId = userIdentity.UserId

	if dto.Id == "" {
		newRecord.CreatedBy = userIdentity.UserId
		return r.educationPlaylistRepository.Insert(newRecord)
	} else {
		newRecord.UpdatedBy = userIdentity.UserId
		return r.educationPlaylistRepository.Update(newRecord)
	}
}

func (r educationService) RemoveFromPlaylist(recordId string) helper.Result {
	return r.educationPlaylistRepository.DeleteById(recordId)
}

func (r educationService) GetById(recordId string) helper.Result {
	return r.educationRepository.GetById(recordId)
}

func (r educationService) GetViewById(recordId string) helper.Result {
	return r.educationRepository.GetViewById(recordId)
}

func (r educationService) DeleteById(recordId string) helper.Result {
	return r.educationRepository.DeleteById(recordId)
}

func (r educationService) GetDirectoryConfig(moduleName string, moduleId string, filetype int) string {
	var res = ""
	if filetype == 1 {
		res = "upload/" + moduleName + "/" + moduleId + "/cover/"
	} else if filetype == 2 {
		res = "upload/" + moduleName + "/" + moduleId + "/files/"
	}
	return res
}

func (r educationService) OpenTransaction(trxHandle *gorm.DB) educationService {
	r.educationRepository = r.educationRepository.OpenTransaction(trxHandle)
	return r
}
