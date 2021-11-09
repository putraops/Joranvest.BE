package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"joranvest/helper"
	"joranvest/midtrans"
	"joranvest/midtrans/coreapi"
	"joranvest/service"

	"github.com/gin-gonic/gin"
)

type PaymentController interface {
	CreateTokenIdByCard(context *gin.Context)
	Charge(context *gin.Context)
}

type paymentController struct {
	jwtService service.JWTService
}

func NewPaymentController(jwtService service.JWTService) PaymentController {
	return &paymentController{
		jwtService: jwtService,
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
