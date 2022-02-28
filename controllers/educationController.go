package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"joranvest/commons"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
	"joranvest/service"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EducationController interface {
	GetPagination(context *gin.Context)
	GetPlaylist(context *gin.Context)
	Lookup(context *gin.Context)
	Save(context *gin.Context)
	AddToPlaylist(c *gin.Context)
	RemoveFromPlaylistById(context *gin.Context)
	GetById(context *gin.Context)
	GetViewById(context *gin.Context)
	GetPlaylistById(context *gin.Context)
	DeleteById(context *gin.Context)

	UploadEducationCover(context *gin.Context)
}

type educationController struct {
	educationService    service.EducationService
	educationRepository repository.EducationRepository
	jwtService          service.JWTService
	db                  *gorm.DB
}

func NewEducationController(db *gorm.DB, jwtService service.JWTService) EducationController {
	return &educationController{
		db:                  db,
		jwtService:          jwtService,
		educationService:    service.NewEducationService(db, jwtService),
		educationRepository: repository.NewEducationRepository(db),
	}
}

// @Tags         Education
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        body body commons.Pagination2ndRequest true "body"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /education/getPagination [post]
func (c educationController) GetPagination(context *gin.Context) {
	var req commons.Pagination2ndRequest
	errDTO := context.Bind(&req)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.educationService.GetPagination(req)
	context.JSON(http.StatusOK, result)
}

// @Tags         Product
// @Security 	 ApiKeyAuth
// @Summary
// @Description
// @Accept       json
// @Produce      json
// @Success      200 {array} models.EducationPlaylist
// @Failure 	 400,404 {object} object
// @Router       /education/getPlaylist [post]
func (c educationController) GetPlaylist(context *gin.Context) {
	qry := context.Request.URL.Query()
	filter := make(map[string]interface{})

	for k, v := range qry {
		filter[k] = v
	}

	var result = c.educationService.GetPlaylist(filter)
	response := helper.StandartResult(true, "Ok", result)
	context.JSON(http.StatusOK, response)
}

// @Tags         Education
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        body body helper.ReactSelectRequest true "body"
// @Param        q query string false "id"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /education/lookup [post]
func (c educationController) Lookup(context *gin.Context) {
	var request helper.ReactSelectRequest

	errDTO := context.Bind(&request)
	if errDTO != nil {
		res := helper.StandartResult(false, errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	}

	var result = c.educationService.Lookup(request)
	response := helper.StandartResult(true, "Ok", result.Data)
	context.JSON(http.StatusOK, response)
}

// @Tags         Education
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.EducationDto true "dto"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /auth/education/save [post]
func (r educationController) Save(c *gin.Context) {
	var result helper.Result
	var dto dto.EducationDto
	dto.Context = c

	errDto := c.Bind(&dto)
	if errDto != nil {
		res := helper.StandartResult(false, errDto.Error(), helper.EmptyObj{})
		c.JSON(http.StatusBadRequest, res)
		return
	}

	result = r.educationService.Save(dto)
	c.JSON(http.StatusOK, helper.StandartResult(result.Status, result.Message, result.Data))
	return
}

// @Tags         Education
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "id"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /education/getById/{id} [get]
func (c educationController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.educationService.GetById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

// @Tags         Education
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "id"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /education/getViewById/{id} [get]
func (c educationController) GetViewById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.educationService.GetViewById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

// @Tags         Education
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "id"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /education/getPlaylistById/{id} [get]
func (c educationController) GetPlaylistById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.educationService.GetPlaylistById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

// @Tags         Education
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "id"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /auth/education/deleteById/{id} [delete]
func (c educationController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.educationService.DeleteById(id)
	context.JSON(http.StatusOK, result)
}

// @Tags         Education
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.EducationDto true "dto"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /auth/education/addToPlaylist [post]
func (r educationController) AddToPlaylist(c *gin.Context) {
	var result helper.Result
	var dto dto.EducationPlaylistDto
	dto.Context = c

	errDto := c.Bind(&dto)
	if errDto != nil {
		res := helper.StandartResult(false, errDto.Error(), helper.EmptyObj{})
		c.JSON(http.StatusBadRequest, res)
		return
	}

	result = r.educationService.AddToPlaylist(dto)
	c.JSON(http.StatusOK, helper.StandartResult(result.Status, result.Message, result.Data))
	return
}

// @Tags         Education
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "id"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /auth/education/removeFromPlaylistById/{id} [delete]
func (c educationController) RemoveFromPlaylistById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	var result = c.educationService.RemoveFromPlaylist(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}

func (c educationController) UploadEducationCover(context *gin.Context) {
	id := context.Param("id")

	file, err1 := context.FormFile("file")
	if err1 != nil {
		context.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err1.Error()))
		return
	}

	authHeader := context.GetHeader("Authorization")
	userIdentity := c.jwtService.GetUserByToken(authHeader)

	folderUpload := c.educationService.GetDirectoryConfig("education", id, 1)
	errRemoveDir := os.RemoveAll(folderUpload)
	if errRemoveDir != nil {
		log.Fatal(errRemoveDir)
	}

	filename := filepath.Base(file.Filename)
	//-- Create folder if not exist
	_, errStat := os.Stat(folderUpload)
	if os.IsNotExist(errStat) {
		errDir := os.MkdirAll(folderUpload, 0755)
		if errDir != nil {
			log.Fatal(errStat)
		}
	}

	path := folderUpload + filename
	if err := context.SaveUploadedFile(file, path); err != nil {
		context.String(http.StatusBadRequest, fmt.Sprintf("Upload File Error: %s", err.Error()))
		return
	}

	var config commons.TConfig
	config.Path = folderUpload
	config.Image.Path = path
	config.Image.Thumbnail.Path = folderUpload
	config.Image.Thumbnail.MaxWidth = 250
	config.Image.Thumbnail.MaxHeight = 250

	path_thumb, errThumb := thumbnailify(config)
	if errThumb != nil {
		log.Fatal(errThumb)
	}

	oldRecordResult := c.educationRepository.GetById(id)
	if !oldRecordResult.Status {
		response := helper.StandartResult(oldRecordResult.Status, oldRecordResult.Message, oldRecordResult.Data)
		context.JSON(http.StatusOK, response)
		return
	}

	var record models.Education
	record = oldRecordResult.Data.(models.Education)
	record.UpdatedBy = userIdentity.UserId
	record.Filepath = path
	record.FilepathThumbnail = path_thumb
	record.Filename = filename
	record.Extension = filepath.Ext(file.Filename)
	record.Size = fmt.Sprint(file.Size)
	result := c.educationRepository.Update(record)

	response := helper.StandartResult(result.Status, result.Message, result.Data)
	context.JSON(http.StatusOK, response)
}
