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
	"joranvest/models/request_models"
	"joranvest/service"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"github.com/mashingan/smapping"
	log "github.com/sirupsen/logrus"
)

type EmitenController interface {
	GetDatatables(context *gin.Context)
	GetPagination(context *gin.Context)
	Lookup(context *gin.Context)
	EmitenLookup(context *gin.Context)
	GetById(context *gin.Context)
	DeleteById(context *gin.Context)
	Save(context *gin.Context)
	PatchingEmiten(context *gin.Context)
}

type emitenController struct {
	emitenService service.EmitenService
	jwtService    service.JWTService
}

func NewEmitenController(emitenService service.EmitenService, jwtService service.JWTService) EmitenController {
	return &emitenController{
		emitenService: emitenService,
		jwtService:    jwtService,
	}
}

func (c *emitenController) GetDatatables(context *gin.Context) {
	var dt commons.DataTableRequest
	errDTO := context.Bind(&dt)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.emitenService.GetDatatables(dt)
	context.JSON(http.StatusOK, result)
}

func (c *emitenController) GetPagination(context *gin.Context) {
	commons.Logger()

	var req commons.PaginationRequest
	errDTO := context.Bind(&req)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		log.Error(errDTO.Error())
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.emitenService.GetPagination(req)
	context.JSON(http.StatusOK, result)
}

func (c *emitenController) Lookup(context *gin.Context) {
	var request helper.ReactSelectRequest
	qry := context.Request.URL.Query()

	if _, found := qry["q"]; found {
		request.Q = fmt.Sprint(qry["q"][0])
	}
	request.Field = helper.StringifyToArray(fmt.Sprint(qry["field"]))

	var result = c.emitenService.Lookup(request)
	response := helper.BuildResponse(true, "Ok", result.Data)
	context.JSON(http.StatusOK, response)
}

func (c *emitenController) EmitenLookup(context *gin.Context) {
	var request helper.ReactSelectRequest
	errBind := context.Bind(&request)
	if errBind != nil {
		res := helper.BuildErrorResponse("Failed to process request", errBind.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}

	var result = c.emitenService.Lookup(request)
	response := helper.BuildResponse(true, "Ok", result.Data)
	context.JSON(http.StatusOK, response)
}

func (c *emitenController) Save(context *gin.Context) {
	result := helper.Response{}
	var recordDto dto.EmitenDto

	errDTO := context.Bind(&recordDto)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userIdentity := c.jwtService.GetUserByToken(authHeader)

		var newRecord = models.Emiten{}
		smapping.FillStruct(&newRecord, smapping.MapFields(&recordDto))
		newRecord.EntityId = userIdentity.EntityId

		if recordDto.Id == "" {
			newRecord.CreatedBy = userIdentity.UserId
			result = c.emitenService.Insert(newRecord)
		} else {
			recordDto.UpdatedBy = userIdentity.UserId
			result = c.emitenService.Update(newRecord)
		}

		if result.Status {
			response := helper.BuildResponse(true, result.Message, result.Data)
			context.JSON(http.StatusOK, response)
		} else {
			response := helper.BuildErrorResponse(result.Message, fmt.Sprintf("%v", result.Errors), helper.EmptyObj{})
			context.JSON(http.StatusOK, response)
		}
	}
}

// @Tags Emiten
// @Summary Get Emiten by Id
// @Description Get Emiten By Id
// @Param id path string true "User Id"
// @Success 200 {object} models.Emiten
// @Failure 400,404 {object} object
// @Router /emiten/getById/{id} [get]
func (c *emitenController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.emitenService.GetById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

// @Tags Emiten
// @Summary Delete Emiten By Id
// @Param id path string true "id"
// @Router /emiten/deleteById [delete]
// @Success 200 {object} helper.Response
// @Failure 400,404 {object} object
// @Router /emiten/getById/{id} [get]
func (c *emitenController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	var result = c.emitenService.DeleteById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}

func (c *emitenController) PatchingEmiten(context *gin.Context) {
	result := helper.Response{}
	var record models.Filemaster

	file, errFile := context.FormFile("file")
	if errFile != nil {
		response := helper.BuildResponse(false, fmt.Sprintf("get form err: %s", errFile.Error()), helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
		return
	}

	err := context.Bind(&record)
	if err != nil {
		response := helper.BuildResponse(false, err.Error(), helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
		return
	} else {
		authHeader := context.GetHeader("Authorization")
		userIdentity := c.jwtService.GetUserByToken(authHeader)

		folderUpload := "./upload/patching/emiten/"
		//-- Create folder if not exist
		_, errStat := os.Stat(folderUpload)
		if os.IsNotExist(errStat) {
			errDir := os.MkdirAll(folderUpload, 0755)
			if errDir != nil {
				log.Fatal(errStat)
				response := helper.BuildResponse(false, fmt.Sprintf("Upload File Error: %s", err.Error()), helper.EmptyObj{})
				context.JSON(http.StatusOK, response)
				return
			}
		}

		filename := filepath.Base(file.Filename)
		extension := filepath.Ext(file.Filename)
		if extension != ".xlsx" && extension != ".xls" {
			response := helper.BuildResponse(false, "Only allow file .xls or .xlsx to upload.", helper.EmptyObj{})
			context.JSON(http.StatusOK, response)
			return
		}

		path := folderUpload + filename
		errRemoveDir := os.RemoveAll(path)
		if err != nil {
			log.Fatal(errRemoveDir)
		}

		if err := context.SaveUploadedFile(file, path); err != nil {
			response := helper.BuildResponse(false, fmt.Sprintf("Upload File Error: %s", err.Error()), helper.EmptyObj{})
			context.JSON(http.StatusOK, response)
			return
		}

		fmt.Println("./upload/patching/emiten/" + filename + extension)

		//-- ReadFile
		xlsx, errRead := excelize.OpenFile("./" + path)
		if errRead != nil {
			response := helper.BuildResponse(false, fmt.Sprintf("%v", errRead.Error()), helper.EmptyObj{})
			context.JSON(http.StatusOK, response)
			return
		}
		readSheet1Name := xlsx.GetSheetName(1)

		var PatchingEmiten []request_models.PatchingEmiten
		for index, _ := range xlsx.GetRows(readSheet1Name) {
			if index > 0 {
				var emitenName = xlsx.GetCellValue(readSheet1Name, fmt.Sprintf("A%d", index+1))
				var emitenCode = xlsx.GetCellValue(readSheet1Name, fmt.Sprintf("B%d", index+1))
				var emitenSector = xlsx.GetCellValue(readSheet1Name, fmt.Sprintf("C%d", index+1))
				var emitenCategory = xlsx.GetCellValue(readSheet1Name, fmt.Sprintf("D%d", index+1))

				if emitenName == "" {
					response := helper.BuildResponse(false, fmt.Sprintf("Emiten Name cannot be empty. Please check line %v", index+1), helper.EmptyObj{})
					context.JSON(http.StatusOK, response)
					return
				}
				if emitenCode == "" {
					response := helper.BuildResponse(false, fmt.Sprintf("Emiten Code cannot be empty. Please check line %v", index+1), helper.EmptyObj{})
					context.JSON(http.StatusOK, response)
					return
				}
				if emitenSector == "" {
					response := helper.BuildResponse(false, fmt.Sprintf("Emiten Sector cannot be empty. Please check line %v", index+1), helper.EmptyObj{})
					context.JSON(http.StatusOK, response)
					return
				}
				if emitenCategory == "" {
					response := helper.BuildResponse(false, fmt.Sprintf("Emiten Category cannot be empty. Please check line %v", index+1), helper.EmptyObj{})
					context.JSON(http.StatusOK, response)
					return
				}

				temp := request_models.PatchingEmiten{ // b == Student{"Bob", 0}
					EmitenName:     emitenName,
					EmitenCode:     emitenCode,
					EmitenSector:   emitenSector,
					EmitenCategory: emitenCategory,
				}
				PatchingEmiten = append(PatchingEmiten, temp)
			}
		}
		result = c.emitenService.PatchingEmiten(PatchingEmiten, userIdentity.UserId)

		if result.Status {
			response := helper.BuildResponse(result.Status, "OK", result.Data)
			context.JSON(http.StatusOK, response)
		} else {
			response := helper.BuildErrorResponse(result.Message, fmt.Sprintf("%v", result.Errors), helper.EmptyObj{})
			context.JSON(http.StatusOK, response)
		}
	}
}
