package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"joranvest/helper"
	"joranvest/models"
	"joranvest/service"

	"github.com/gin-gonic/gin"
)

type FilemasterController interface {
	GetAll(context *gin.Context)
	SingleUpload(context *gin.Context)
	UploadByType(context *gin.Context)
	Insert(context *gin.Context)
	DeleteByRecordId(context *gin.Context)
}

type filemasterController struct {
	filemasterService service.FilemasterService
	jwtService        service.JWTService
}

func NewFilemasterController(_service service.FilemasterService, jwtService service.JWTService) FilemasterController {
	return &filemasterController{
		filemasterService: _service,
		jwtService:        jwtService,
	}
}

func (c *filemasterController) GetAll(context *gin.Context) {
	qry := context.Request.URL.Query()
	filter := make(map[string]interface{})

	for k, v := range qry {
		filter[k] = v
	}

	var result = c.filemasterService.GetAll(filter)
	response := helper.BuildResponse(true, "Ok", result)
	context.JSON(http.StatusOK, response)
}

func (c *filemasterController) SingleUpload(context *gin.Context) {
	id := context.Param("id")

	result := helper.Response{}
	var record models.Filemaster

	file, err1 := context.FormFile("file")
	if err1 != nil {
		context.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err1.Error()))
		return
	}

	err := context.Bind(&record)
	if err != nil {
		res := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userIdentity := c.jwtService.GetUserByToken(authHeader)

		//folderDir := "upload/" + id
		folderUpload := "upload/" + id + "/"

		errRemoveDir := os.RemoveAll(folderUpload)
		if err != nil {
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

		record.RecordId = id
		record.EntityId = userIdentity.EntityId
		record.CreatedBy = userIdentity.UserId
		record.Filepath = path
		record.Filename = filename
		record.Extension = filepath.Ext(file.Filename)
		record.Size = fmt.Sprint(file.Size)
		result = c.filemasterService.SingleUpload(record)

		if result.Status {
			response := helper.BuildResponse(true, "OK", result.Data)
			context.JSON(http.StatusOK, response)
		} else {
			response := helper.BuildErrorResponse(result.Message, fmt.Sprintf("%v", result.Errors), helper.EmptyObj{})
			context.JSON(http.StatusOK, response)
		}
	}
}

func (c *filemasterController) UploadByType(context *gin.Context) {
	id := context.Param("id")
	module := context.Param("module")
	filetype := context.Param("filetype")

	result := helper.Response{}
	var record models.Filemaster

	file, err1 := context.FormFile("file")
	if err1 != nil {
		context.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err1.Error()))
		return
	}

	err := context.Bind(&record)
	if err != nil {
		res := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userIdentity := c.jwtService.GetUserByToken(authHeader)

		//folderDir := "upload/" + id
		_filetype, errConvert := strconv.Atoi(filetype)
		if errConvert != nil {
			log.Fatal(errConvert)
		}
		folderUpload := c.filemasterService.GetDirectoryConfig(module, id, _filetype)

		errRemoveDir := os.RemoveAll(folderUpload)
		if err != nil {
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

		record.RecordId = id
		record.EntityId = userIdentity.EntityId
		record.CreatedBy = userIdentity.UserId
		record.Filepath = path
		record.Filename = filename
		record.Extension = filepath.Ext(file.Filename)
		record.Size = fmt.Sprint(file.Size)
		record.FileType = _filetype
		result = c.filemasterService.UploadByType(record)

		if result.Status {
			response := helper.BuildResponse(true, "OK", result.Data)
			context.JSON(http.StatusOK, response)
		} else {
			response := helper.BuildErrorResponse(result.Message, fmt.Sprintf("%v", result.Errors), helper.EmptyObj{})
			context.JSON(http.StatusOK, response)
		}
	}
}

func (c *filemasterController) Insert(context *gin.Context) {
	id := context.Param("id")
	is_multiple := context.Param("is_multiple")

	fmt.Println(is_multiple)

	result := helper.Response{}
	var record models.Filemaster

	file, err1 := context.FormFile("file")
	if err1 != nil {
		context.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err1.Error()))
		return
	}

	err := context.Bind(&record)
	if err != nil {
		res := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userIdentity := c.jwtService.GetUserByToken(authHeader)

		folderUpload := "upload/" + id + "/"
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

		record.RecordId = id
		record.EntityId = userIdentity.EntityId
		record.CreatedBy = userIdentity.UserId
		record.Filepath = path
		record.Filename = filename
		record.Extension = filepath.Ext(file.Filename)
		record.Size = fmt.Sprint(file.Size)
		result = c.filemasterService.Insert(record)

		if result.Status {
			response := helper.BuildResponse(true, "OK", result.Data)
			context.JSON(http.StatusOK, response)
		} else {
			response := helper.BuildErrorResponse(result.Message, fmt.Sprintf("%v", result.Errors), helper.EmptyObj{})
			context.JSON(http.StatusOK, response)
		}
	}
}

func (c *filemasterController) DeleteByRecordId(context *gin.Context) {
	recordId := context.Param("recordId")
	if recordId == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	var result = c.filemasterService.DeleteByRecordId(recordId)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}
