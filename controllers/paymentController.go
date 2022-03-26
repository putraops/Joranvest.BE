package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/xendit/xendit-go"

	"joranvest/commons"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/midtrans"
	"joranvest/midtrans/coreapi"
	"joranvest/models"
	"joranvest/payment_gateway/xendit/ewallet"
	"joranvest/payment_gateway/xendit/qrcode"
	"joranvest/service"

	"github.com/gin-gonic/gin"
	"github.com/mashingan/smapping"
)

type PaymentController interface {
	GetPagination(context *gin.Context)
	GetById(context *gin.Context)
	GetByProviderRecordId(context *gin.Context)
	GetByProviderReferenceId(context *gin.Context)
	GetEWalletPaymentStatusByReferenceId(context *gin.Context)
	GetUniqueNumber(context *gin.Context)
	UpdateWalletPaymentStatus(context *gin.Context)
	DeleteById(context *gin.Context)
	MembershipPayment(context *gin.Context)
	WebinarPayment(context *gin.Context)
	UpdatePaymentStatus(context *gin.Context)
	CreateTokenIdByCard(context *gin.Context)

	CreateTransferPayment(context *gin.Context)
	CreateEWalletPayment(context *gin.Context)
	CreateQRCode(context *gin.Context)
	Charge(context *gin.Context)
	HookForXendit(context *gin.Context)
}

type paymentController struct {
	paymentService service.PaymentService
	jwtService     service.JWTService
	helper.AppSession
	//ewallet ewallet.Payment
}

func NewPaymentController(paymentService service.PaymentService, jwtService service.JWTService) PaymentController {
	return &paymentController{
		paymentService: paymentService,
		jwtService:     jwtService,
		//ewallet:        ewallet.NewPaymentRequest(&gin.Context{}),
		// orderService: service.NewOrderService(db, jwtService),
	}
}

func (c *paymentController) HookForXendit(context *gin.Context) {
	commons.Logger()
	var callbackBody dto.CallbackBodyDto

	responseToken := context.GetHeader("x-callback-token")
	callbackToken := os.Getenv("XENDIT_CALLBACK_TOKEN")

	if responseToken != callbackToken {
		log.Error("Callback Token and Response Token is not match. [Payment Gateway: Xendit]")
		log.Error(fmt.Sprintf("Response Token: %v", responseToken))
		log.Error(fmt.Sprintf("Callback Token: %v", callbackToken))
		log.Error("Please check the Callback Token.")
		context.JSON(http.StatusUnauthorized, nil)
		return
	}

	errDto := context.Bind(&callbackBody)
	if errDto != nil {
		log.Error("Failed to Bind CallbackDto [Payment Gateway: Xendit]")
		context.JSON(http.StatusBadRequest, nil)
	}

	myMap := make(map[string]interface{})
	myMap = callbackBody.Data

	// convert map to json
	jsonString, errMarshal := json.Marshal(myMap)
	if errMarshal != nil {
		log.Error("Error to Marshal json string from dto.Data")
		context.JSON(http.StatusBadRequest, nil)
	}
	xenditCharge := xendit.EWalletCharge{}
	json.Unmarshal(jsonString, &xenditCharge)

	fmt.Println("------------------------------------------------------")
	fmt.Println("---------------------- Status ------------------------")
	fmt.Println(xenditCharge.Status)
	fmt.Println(xenditCharge.Metadata["record_id"])

	var paymentRecordId string = xenditCharge.Metadata["record_id"].(string)
	fmt.Println("------------------------------------------------------")
	fmt.Println("---------------------- Status ------------------------")
	fmt.Println(paymentRecordId)
	if xenditCharge.Status != string(ewallet.XenditPaymentStatusPending) {
		var paymentStatus int = commons.RejectedPaymentStatus
		if xenditCharge.Status == string(ewallet.XenditPaymentStatusSucceeded) {
			paymentStatus = commons.PaidPaymentStatus
		}

		c.paymentService.UpdateWalletPaymentStatus(dto.UpdatePaymentStatusDto{
			Id:            paymentRecordId,
			PaymentStatus: paymentStatus,
			UpdatedBy:     xenditCharge.Metadata["user_id"].(string),
		})
	}
	fmt.Println("------------------------------------------------------")
	fmt.Println("------------------------------------------------------")
	//get Body
	context.JSON(http.StatusOK, helper.EmptyObj{})
	return
}

func (c *paymentController) GetPagination(context *gin.Context) {
	var req commons.Pagination2ndRequest
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

func (c *paymentController) CreateTransferPayment(context *gin.Context) {
	commons.Logger()
	var recordDto dto.PaymentDto
	errDTO := context.Bind(&recordDto)

	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		log.Error("WebinarPayment: Bind Dto")
		log.Error(fmt.Sprintf("%v,", errDTO.Error()))
		context.JSON(http.StatusBadRequest, res)
		return
	}

	recordDto.Context = context
	result := c.paymentService.CreateTransferPayment(recordDto)

	if !result.Status {
		response := helper.BuildResponse(result.Status, result.Message, result.Data)
		context.JSON(http.StatusOK, response)
		return
	}

	response := helper.BuildResponse(result.Status, result.Message, result.Data)
	context.JSON(http.StatusOK, response)
	return
}

func (c *paymentController) CreateQRCode(context *gin.Context) {
	var dto qrcode.QRCodeDto
	errDto := context.Bind(&dto)
	if errDto != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDto.Error(), helper.EmptyObj{})
		log.Error("CreateQRCode: Bind Dto")
		log.Error(fmt.Sprintf("%v,", errDto.Error()))
		context.JSON(http.StatusBadRequest, res)
		return
	}

	dto.Context = context
	response := c.paymentService.CreateQRCode(dto)
	context.JSON(http.StatusOK, response)
}

func (c *paymentController) CreateEWalletPayment(context *gin.Context) {
	var dto ewallet.PaymentDto
	errDto := context.Bind(&dto)
	if errDto != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDto.Error(), helper.EmptyObj{})
		log.Error("EWalletCharge: Bind Dto")
		log.Error(fmt.Sprintf("%v,", errDto.Error()))
		context.JSON(http.StatusBadRequest, res)
		return
	}

	dto.Context = context
	response := c.paymentService.CreateEWalletPayment(dto)
	context.JSON(http.StatusOK, response)
}

func (c *paymentController) GetEWalletPaymentStatusByReferenceId(context *gin.Context) {
	reference_id := context.Param("reference_id")
	if reference_id == "" {
		response := helper.BuildErrorResponse("Failed to get reference_id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
		return
	}

	response := c.paymentService.GetEWalletPaymentStatus(context, reference_id)
	context.JSON(http.StatusOK, response)
	return
}

func (c *paymentController) UpdateWalletPaymentStatus(context *gin.Context) {
	var dto dto.UpdatePaymentStatusDto
	errDto := context.Bind(&dto)
	if errDto != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDto.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	}

	dto.Context = context
	response := c.paymentService.UpdateWalletPaymentStatus(dto)
	context.JSON(http.StatusOK, response)
	return
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
		return
	}

	result := c.paymentService.GetById(id)
	context.JSON(http.StatusOK, result)
}

func (c *paymentController) GetByProviderRecordId(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
		return
	}

	result := c.paymentService.GetByProviderRecordId(id)
	context.JSON(http.StatusOK, result)
}

func (c *paymentController) GetByProviderReferenceId(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
		return
	}

	result := c.paymentService.GetByProviderReferenceId(id)
	context.JSON(http.StatusOK, result)
}

func (c *paymentController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
		return
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
		return
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
