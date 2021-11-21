package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"joranvest/commons"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/midtrans"
	"joranvest/midtrans/coreapi"
	"joranvest/models"
	"joranvest/service"

	"github.com/gin-gonic/gin"
	"github.com/mashingan/smapping"
)

type PaymentController interface {
	GetPagination(context *gin.Context)
	GetById(context *gin.Context)
	GetUniqueNumber(context *gin.Context)
	DeleteById(context *gin.Context)
	MembershipPayment(context *gin.Context)
	WebinarPayment(context *gin.Context)
	UpdatePaymentStatus(context *gin.Context)
	CreateTokenIdByCard(context *gin.Context)
	Charge(context *gin.Context)
}

type paymentController struct {
	paymentService service.PaymentService
	jwtService     service.JWTService
	helper.AppSession
}

func NewPaymentController(paymentService service.PaymentService, jwtService service.JWTService) PaymentController {
	return &paymentController{
		paymentService: paymentService,
		jwtService:     jwtService,
	}
}

func (c *paymentController) GetPagination(context *gin.Context) {
	var req commons.PaginationRequest
	errDTO := context.Bind(&req)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.paymentService.GetPagination(req)
	context.JSON(http.StatusOK, result)
}

func (c *paymentController) GetUniqueNumber(context *gin.Context) {
	var result = c.paymentService.GetUniqueNumber()
	response := helper.BuildResponse(true, "Ok", result)
	context.JSON(http.StatusOK, response)
}

func (c *paymentController) MembershipPayment(context *gin.Context) {
	result := helper.Response{}
	var recordDto dto.PaymentDto
	fmt.Println(recordDto)
	errDTO := context.Bind(&recordDto)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		fmt.Println("not error")
		authHeader := context.GetHeader("Authorization")
		userIdentity := c.jwtService.GetUserByToken(authHeader)

		var newRecord = models.Payment{}
		smapping.FillStruct(&newRecord, smapping.MapFields(&recordDto))
		newRecord.EntityId = userIdentity.EntityId
		newRecord.CreatedBy = userIdentity.UserId
		newRecord.OwnerId = userIdentity.UserId
		result = c.paymentService.MembershipPayment(newRecord)

		if result.Status {
			response := helper.BuildResponse(true, "OK", result.Data)
			context.JSON(http.StatusOK, response)
		} else {
			response := helper.BuildErrorResponse(result.Message, fmt.Sprintf("%v", result.Errors), helper.EmptyObj{})
			context.JSON(http.StatusOK, response)
		}
	}
}

func (c *paymentController) WebinarPayment(context *gin.Context) {
	commons.Logger()
	result := helper.Response{}
	var recordDto dto.PaymentDto
	fmt.Println(recordDto)
	errDTO := context.Bind(&recordDto)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		log.Error("WebinarPayment: Bind Dto")
		log.Error(fmt.Sprintf("%v,", errDTO.Error()))
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userIdentity := c.jwtService.GetUserByToken(authHeader)

		var newRecord = models.Payment{}
		smapping.FillStruct(&newRecord, smapping.MapFields(&recordDto))
		newRecord.EntityId = userIdentity.EntityId
		newRecord.CreatedBy = userIdentity.UserId
		newRecord.OwnerId = userIdentity.UserId
		result = c.paymentService.WebinarPayment(newRecord)

		if result.Status {
			response := helper.BuildResponse(true, "OK", result.Data)
			context.JSON(http.StatusOK, response)
		} else {
			response := helper.BuildErrorResponse(result.Message, fmt.Sprintf("%v", result.Errors), helper.EmptyObj{})
			context.JSON(http.StatusOK, response)
		}
	}
}

func (c *paymentController) UpdatePaymentStatus(context *gin.Context) {
	result := helper.Response{}
	var recordDto dto.UpdatePaymentStatusDto
	fmt.Println(recordDto)
	errDTO := context.Bind(&recordDto)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		fmt.Println("not error")
		authHeader := context.GetHeader("Authorization")
		userIdentity := c.jwtService.GetUserByToken(authHeader)

		recordDto.UpdatedBy = userIdentity.UserId
		result = c.paymentService.UpdatePaymentStatus(recordDto)

		if result.Status {
			response := helper.BuildResponse(true, "OK", result.Data)
			context.JSON(http.StatusOK, response)
		} else {
			response := helper.BuildErrorResponse(result.Message, fmt.Sprintf("%v", result.Errors), helper.EmptyObj{})
			context.JSON(http.StatusOK, response)
		}
	}
}

func (c *paymentController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.paymentService.GetById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

func (c *paymentController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	var result = c.paymentService.DeleteById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}

func (c *paymentController) CreateTokenIdByCard(context *gin.Context) {
	var r midtrans.CardDetails

	errDTO := context.Bind(&r)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}

	fmt.Println("=================================================")
	fmt.Println("---------------- Configuration ------------------")
	fmt.Println("=================================================")

	res, err := coreapi.CardToken(r)

	if err != nil {
		fmt.Println("Error get card token", err.GetMessage())
	}
	fmt.Println("response card token", res)
	fmt.Println("Request")
	fmt.Println("Request: ", r)
	fmt.Println("CardNumber: ", r.CardNumber)
	fmt.Println("ExpMonth: ", r.ExpMonth)
	fmt.Println("ExpYear: ", r.ExpYear)
	fmt.Println("CVV: ", r.CVV)
	fmt.Println("=================================================")
	fmt.Println("--------------------- End -----------------------")
	fmt.Println("=================================================")
	status_code, _ := strconv.Atoi(res.StatusCode)

	if status_code == 200 {
		response := helper.BuildResponse(true, "OK", res)
		context.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildErrorResponse(res.StatusMessage, fmt.Sprintf("%v", res.ValidationMessage), res.ValidationMessage)
		context.JSON(http.StatusOK, response)
	}
}

func (c *paymentController) Charge(context *gin.Context) {
	var r *coreapi.ChargeReq

	errDTO := context.Bind(&r)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	//result := helper.Response{}

	fmt.Println("=================================================")
	fmt.Println("---------------- Configuration ------------------")
	fmt.Println("=================================================")

	res, err := coreapi.ChargeTransaction(r)

	if err != nil {
		fmt.Println("Error get card token", err.GetMessage())
	}
	fmt.Println("response card token", res)
	fmt.Println("Request")
	fmt.Println("Request: ", r)
	fmt.Println("=================================================")
	fmt.Println("--------------------- End -----------------------")
	fmt.Println("=================================================")
	status_code, _ := strconv.Atoi(res.StatusCode)

	if status_code == 200 {
		response := helper.BuildResponse(true, "OK", res)
		context.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildErrorResponse(res.StatusMessage, fmt.Sprintf("%v", helper.EmptyObj{}), res)
		context.JSON(http.StatusOK, response)
	}
}
