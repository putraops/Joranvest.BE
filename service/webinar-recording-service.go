package service

import (
	"joranvest/commons"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"

	"github.com/gin-gonic/gin"
	"github.com/mashingan/smapping"
	"gorm.io/gorm"
)

type WebinarRecordingService interface {
	GetPagination(request commons.Pagination2ndRequest) interface{}
	GetAll(filter map[string]interface{}) []models.WebinarRecording
	// GetPlaylist(filter map[string]interface{}) []models.EducationPlaylist
	// GetPlaylistById(recordId string) helper.Result
	// GetPlaylistByUserId(educationId string, userId string) helper.Result
	// Lookup(request helper.ReactSelectRequest) helper.Result
	Save(dto dto.WebinarRecordingDto) helper.Result
	Submit(recordId string, context *gin.Context) helper.Result
	// AddToPlaylist(dto dto.EducationPlaylistDto) helper.Result
	// MarkVideoAsWatched(dto dto.EducationPlaylistUserDto) helper.Result
	// RemoveFromPlaylist(recordId string) helper.Result
	GetViewById(recordId string) helper.Result
	GetById(recordId string) helper.Result
	GetByWebinarId(webinarId string) helper.Result
	DeleteById(recordId string) helper.Result

	GetDirectoryConfig(moduleName string, moduleId string, filetype int) string
	OpenTransaction(*gorm.DB) webinarRecordingService
}

type webinarRecordingService struct {
	DB                         *gorm.DB
	jwtService                 JWTService
	webinarRepository          repository.WebinarRepository
	webinarRecordingRepository repository.WebinarRecordingRepository
}

func NewWebinarRecordingService(db *gorm.DB, jwtService JWTService) WebinarRecordingService {
	return webinarRecordingService{
		DB:                         db,
		jwtService:                 jwtService,
		webinarRepository:          repository.NewWebinarRepository(db),
		webinarRecordingRepository: repository.NewWebinarRecordingRepository(db),
	}
}

func (r webinarRecordingService) GetPagination(request commons.Pagination2ndRequest) interface{} {
	return r.webinarRecordingRepository.GetPagination(request)
}

func (r webinarRecordingService) GetAll(filter map[string]interface{}) []models.WebinarRecording {
	return r.webinarRecordingRepository.GetAll(filter)
}

func (r webinarRecordingService) Save(dto dto.WebinarRecordingDto) helper.Result {
	authHeader := dto.Context.GetHeader("Authorization")
	userIdentity := r.jwtService.GetUserByToken(authHeader)

	result := r.webinarRepository.GetById(*dto.WebinarId)
	if !result.Status {
		return helper.StandartResult(result.Status, result.Message, result.Data)
	}

	//-- Map Dto to Struct
	var newRecord models.WebinarRecording
	smapping.FillStruct(&newRecord, smapping.MapFields(&dto))
	newRecord.EntityId = &userIdentity.EntityId
	newRecord.WebinarId = *dto.WebinarId
	newRecord.PathUrl = helper.StringToPathUrl(result.Data.(models.Webinar).Title)

	if dto.Id == "" {
		newRecord.CreatedBy = userIdentity.UserId
		return r.webinarRecordingRepository.Insert(newRecord)
	} else {
		newRecord.UpdatedBy = userIdentity.UserId
		return r.webinarRecordingRepository.Update(newRecord)
	}
}

func (r webinarRecordingService) Submit(recordId string, context *gin.Context) helper.Result {
	authHeader := context.GetHeader("Authorization")
	userIdentity := r.jwtService.GetUserByToken(authHeader)

	return r.webinarRecordingRepository.Submit(recordId, userIdentity.UserId)
}

// func (r webinarRecordingService) AddToPlaylist(dto dto.EducationPlaylistDto) helper.Result {
// 	authHeader := dto.Context.GetHeader("Authorization")
// 	userIdentity := r.jwtService.GetUserByToken(authHeader)

// 	//-- Map Dto to Struct
// 	var newRecord models.EducationPlaylist
// 	smapping.FillStruct(&newRecord, smapping.MapFields(&dto))
// 	newRecord.EntityId = userIdentity.UserId

// 	if dto.Id == "" {
// 		newRecord.CreatedBy = userIdentity.UserId
// 		return r.educationPlaylistRepository.Insert(newRecord)
// 	} else {
// 		newRecord.UpdatedBy = userIdentity.UserId
// 		return r.educationPlaylistRepository.Update(newRecord)
// 	}
// }

func (r webinarRecordingService) GetById(recordId string) helper.Result {
	return r.webinarRecordingRepository.GetById(recordId)
}

func (r webinarRecordingService) GetViewById(recordId string) helper.Result {
	return r.webinarRecordingRepository.GetViewById(recordId)
}

func (r webinarRecordingService) GetByWebinarId(webinarId string) helper.Result {
	return r.webinarRecordingRepository.GetByWebinarId(webinarId)
}

func (r webinarRecordingService) DeleteById(recordId string) helper.Result {
	return r.webinarRecordingRepository.DeleteById(recordId)
}

func (r webinarRecordingService) GetDirectoryConfig(moduleName string, moduleId string, filetype int) string {
	var res = ""
	if filetype == 1 {
		res = "upload/" + moduleName + "/" + moduleId + "/cover/"
	} else if filetype == 2 {
		res = "upload/" + moduleName + "/" + moduleId + "/files/"
	}
	return res
}

func (r webinarRecordingService) OpenTransaction(trxHandle *gorm.DB) webinarRecordingService {
	r.webinarRecordingRepository = r.webinarRecordingRepository.OpenTransaction(trxHandle)
	return r
}
