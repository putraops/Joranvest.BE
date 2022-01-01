package controllers

import (
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/service"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
)

type FilemasterController interface {
	GetAll(context *gin.Context)
	GetAllByRecordIds(context *gin.Context)
	SingleUpload(context *gin.Context)
	UploadByType(context *gin.Context)
	UploadPDFDocuments(context *gin.Context)
	UploadProfilePicture(context *gin.Context)
	UploadWebinarCover(context *gin.Context)
	Insert(context *gin.Context)
	DeleteById(context *gin.Context)
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

func (c *filemasterController) GetAllByRecordIds(context *gin.Context) {
	var request []string
	errBind := context.Bind(&request)
	if errBind != nil {
		res := helper.BuildErrorResponse("Failed to process request", errBind.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}

	var result = c.filemasterService.GetAllByRecordIds(request)
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

		var config commons.TConfig
		config.Path = folderUpload
		config.Image.Path = path
		config.Image.Thumbnail.Path = folderUpload
		config.Image.Thumbnail.MaxWidth = 250
		config.Image.Thumbnail.MaxHeight = 250

		path_thumb, errThumb := thumbnailify(config)
		if err != nil {
			log.Fatal(errThumb)
		}

		record.RecordId = id
		record.EntityId = userIdentity.EntityId
		record.CreatedBy = userIdentity.UserId
		record.Filepath = path
		record.FilepathThumbnail = path_thumb
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

func (c *filemasterController) UploadWebinarCover(context *gin.Context) {
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

		folderUpload := c.filemasterService.GetDirectoryConfig("webinar", id, 1)
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

		var config commons.TConfig
		config.Path = folderUpload
		config.Image.Path = path
		config.Image.Thumbnail.Path = folderUpload
		config.Image.Thumbnail.MaxWidth = 250
		config.Image.Thumbnail.MaxHeight = 250

		path_thumb, errThumb := thumbnailify(config)
		if err != nil {
			log.Fatal(errThumb)
		}

		record.RecordId = id

		var webinarRecord models.Webinar
		webinarRecord.Id = id
		webinarRecord.UpdatedBy = userIdentity.UserId
		webinarRecord.Filepath = path
		webinarRecord.FilepathThumbnail = path_thumb
		webinarRecord.Filename = filename
		webinarRecord.Extension = filepath.Ext(file.Filename)
		webinarRecord.Size = fmt.Sprint(file.Size)
		result = c.filemasterService.UpdateWebinarCoverImage(webinarRecord)

		if result.Status {
			response := helper.BuildResponse(true, "OK", result.Data)
			context.JSON(http.StatusOK, response)
		} else {
			response := helper.BuildErrorResponse(result.Message, fmt.Sprintf("%v", result.Errors), helper.EmptyObj{})
			context.JSON(http.StatusOK, response)
		}
	}
}

func (c *filemasterController) UploadPDFDocuments(context *gin.Context) {
	id := context.Param("id")
	module := context.Param("module")

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

		folderUpload := "upload/" + module + "/" + id + "/supported_files/"
		//-- Create folder if not exist
		_, errStat := os.Stat(folderUpload)
		if os.IsNotExist(errStat) {
			errDir := os.MkdirAll(folderUpload, 0755)
			if errDir != nil {
				log.Fatal(errStat)
			}
		}
		filename := filepath.Base(file.Filename)
		extension := filepath.Ext(file.Filename)
		if extension != ".pdf" {
			response := helper.BuildResponse(false, "Only allow file .pdf to upload.", helper.EmptyObj{})
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
		record.RecordId = id
		record.EntityId = userIdentity.EntityId
		record.CreatedBy = userIdentity.UserId
		record.Filepath = path
		// record.FilepathThumbnail = path_thumb
		record.Filename = filename
		record.Extension = filepath.Ext(file.Filename)
		record.Size = fmt.Sprint(file.Size)
		record.FileType = 3 //-- Documents
		result = c.filemasterService.Insert(record)

		if result.Status {
			response := helper.BuildResponse(result.Status, "OK", result.Data)
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

func (c *filemasterController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	var result = c.filemasterService.DeleteById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		errRemoveDir := os.RemoveAll(result.Data.(string))
		if errRemoveDir != nil {
			log.Fatal(errRemoveDir)
		}

		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
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

func (c *filemasterController) UploadProfilePicture(context *gin.Context) {
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

		folderUpload := "upload/users/profile_picture/" + id
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

		path := folderUpload + "/" + filename
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
		if err != nil {
			log.Fatal(errThumb)
		}

		record.RecordId = id
		record.EntityId = userIdentity.EntityId
		record.CreatedBy = userIdentity.UserId
		record.Filepath = path
		record.FilepathThumbnail = path_thumb
		record.Filename = filename
		record.Extension = filepath.Ext(file.Filename)
		record.Size = fmt.Sprint(file.Size)
		record.FileType = 1
		result = c.filemasterService.UploadProfilePicture(record)

		if result.Status {
			response := helper.BuildResponse(true, "OK", result.Data)
			context.JSON(http.StatusOK, response)
		} else {
			response := helper.BuildErrorResponse(result.Message, fmt.Sprintf("%v", result.Errors), helper.EmptyObj{})
			context.JSON(http.StatusOK, response)
		}
	}
}

func thumbnailify(config commons.TConfig) (outputPath string, err error) {
	var (
		file     *os.File
		img      image.Image
		filename = "thumb_" + path.Base(config.Image.Path)
	)

	extname := strings.ToLower(path.Ext(config.Image.Path))

	outputPath = path.Join(config.Image.Thumbnail.Path, filename)
	println("outputPath")
	println(outputPath)

	//-- Baca File
	if file, err = os.Open(config.Image.Path); err != nil {
		log.Fatal(err)
		return
	}

	defer file.Close()

	// decode jpeg into image.Image
	switch extname {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
		break
	case ".png":
		img, err = png.Decode(file)
		break
	case ".gif":
		img, err = gif.Decode(file)
		break
	default:
		err = errors.New("Unsupport file type" + extname)
		return
	}

	if img == nil {
		err = errors.New("Generate thumbnail fail...")
		return

	}

	m := resize.Thumbnail(uint(config.Image.Thumbnail.MaxWidth), uint(config.Image.Thumbnail.MaxHeight), img, resize.Lanczos3)

	out, err := os.Create(outputPath)
	if err != nil {
		return
	}
	defer out.Close()

	// write new image to file
	//decode jpeg/png/gif into image.Image
	switch extname {
	case ".jpg", ".jpeg":
		jpeg.Encode(out, m, nil)
		break
	case ".png":
		png.Encode(out, m)
		break
	case ".gif":
		gif.Encode(out, m, nil)
		break
	default:
		err = errors.New("Unsupport file type" + extname)
		return
	}

	return
}
