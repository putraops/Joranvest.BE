package service

import (
	"encoding/json"
	"fmt"
	"joranvest/commons"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/models"
	payment_gateway_providers "joranvest/payment_gateway"
	"joranvest/payment_gateway/xendit/ewallet"
	"joranvest/payment_gateway/xendit/qrcode"
	"joranvest/repository"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mashingan/smapping"
	"github.com/xendit/xendit-go"
)

type PaymentService interface {
	GetPagination(request commons.Pagination2ndRequest) interface{}
	GetAll(filter map[string]interface{}) []models.Payment
	GetUniqueNumber() int
	GetEWalletPaymentStatus(ctx *gin.Context, referenceId string) helper.Response
	UpdateWalletPaymentStatus(dto dto.UpdatePaymentStatusDto) helper.Response
	MembershipPayment(record models.Payment) helper.Response
	WebinarPayment(record models.Payment) helper.Response

	CreateTransferPayment(dto dto.PaymentDto) helper.Result
	CreateEWalletPayment(dto ewallet.PaymentDto) helper.Response
	CreateQRCode(dto qrcode.QRCodeDto) helper.Response

	Update(record models.Payment) helper.Response
	UpdatePaymentStatus(req dto.UpdatePaymentStatusDto) helper.Response
	GetById(recordId string) helper.Response
	GetByProviderRecordId(id string) helper.Response
	GetByProviderReferenceId(id string) helper.Response
	DeleteById(recordId string) helper.Response
}

type paymentService struct {
	paymentRepository repository.PaymentRepository
	helper.AppSession
	jwtService JWTService
}

func NewPaymentService(repo repository.PaymentRepository) PaymentService {
	return &paymentService{
		paymentRepository: repo,
		jwtService:        NewJWTService(),
	}
}

func (service *paymentService) GetPagination(request commons.Pagination2ndRequest) interface{} {
	return service.paymentRepository.GetPagination(request)
}

func (service *paymentService) GetAll(filter map[string]interface{}) []models.Payment {
	return service.paymentRepository.GetAll(filter)
}

func (service *paymentService) GetUniqueNumber() int {
	return service.paymentRepository.GetUniqueNumber()
}

func (service *paymentService) MembershipPayment(record models.Payment) helper.Response {
	return service.paymentRepository.MembershipPayment(record)
}
func (service *paymentService) WebinarPayment(record models.Payment) helper.Response {
	return service.paymentRepository.WebinarPayment(record)
}

func (service *paymentService) CreateQRCode(dto qrcode.QRCodeDto) helper.Response {
	token := dto.Context.GetHeader("Authorization")

	userIdentity := service.jwtService.GetUserByToken(token)
	xenditService := qrcode.NewQRCode(dto.Context)

	var newRecord = models.Payment{}
	newRecord.Id = uuid.New().String()
	newRecord.RecordId = dto.RecordId //-- WebinarId or MembershipId

	dto.RecordId = newRecord.Id //-- RecordId replace by NewPaymentId
	dto.ApplicationUserId = userIdentity.UserId
	res, err := xenditService.CreateQRCode(dto)
	if err != nil {
		return helper.ServerResponse(false, fmt.Sprintf("%v", err.Message), "", helper.EmptyObj{})
	}

	currentTime := time.Now()

	newRecord.CreatedAt = &currentTime
	newRecord.IsActive = true
	newRecord.EntityId = userIdentity.EntityId
	newRecord.CreatedBy = userIdentity.UserId
	newRecord.OwnerId = userIdentity.UserId
	newRecord.ApplicationUserId = userIdentity.UserId
	newRecord.Currency = "IDR"
	newRecord.OrderNumber = fmt.Sprintf("%v/QR/%v/%v/%v", "JORAN", strconv.Itoa(time.Now().Year()), helper.NumberMonthToRoman(int(time.Now().Month())), strings.ToUpper((strconv.Itoa(time.Now().Nanosecond()))[0:5]))

	newRecord.PaymentType = string(qrcode.QRChannelCodeQRIS)
	newRecord.PaymentStatus = 2
	newRecord.Price = int(dto.Amount)
	newRecord.UniqueNumber = 0

	payment_date_expired := time.Now().Add(time.Minute * 5)
	newRecord.PaymentDateExpired = &payment_date_expired
	newRecord.ProviderName = string(payment_gateway_providers.Xendit)
	newRecord.ProviderRecordId = res.ID

	var result helper.Response
	result = service.paymentRepository.Insert(newRecord)
	if !result.Status {
		return result
	}

	return helper.ServerResponse(true, "Ok", "", res)
}

func (service *paymentService) CreateTransferPayment(dto dto.PaymentDto) helper.Result {
	token := dto.Context.GetHeader("Authorization")
	userIdentity := service.jwtService.GetUserByToken(token)
	// xenditService := ewallet.NewPaymentRequest(dto.Context)

	var newRecord = models.Payment{}
	smapping.FillStruct(&newRecord, smapping.MapFields(&dto))

	//dto.RecordId = newRecord.Id //-- RecordId replace by NewPaymentId
	//dto.ApplicationUserId = userIdentity.UserId

	// res, err := xenditService.CreateEWalletCharge(dto)
	// if err != nil {
	// 	return helper.ServerResponse(false, fmt.Sprintf("%v", err.Message), "", helper.EmptyObj{})
	// }

	currentTime := time.Now()
	newRecord.Id = uuid.New().String()
	// newRecord.RecordId = dto.RecordId //-- WebinarId or MembershipId
	newRecord.CreatedAt = &currentTime
	newRecord.EntityId = userIdentity.EntityId
	newRecord.CreatedBy = userIdentity.UserId
	newRecord.OwnerId = userIdentity.UserId
	newRecord.ApplicationUserId = userIdentity.UserId
	newRecord.Currency = "IDR"
	newRecord.OrderNumber = fmt.Sprintf("%v/TRF/%v/%v/%v", "JORAN", strconv.Itoa(time.Now().Year()), helper.NumberMonthToRoman(int(time.Now().Month())), strings.ToUpper((strconv.Itoa(time.Now().Nanosecond()))[0:5]))

	// newRecord.PaymentType = dto.PaymentType
	// newRecord.PaymentStatus = 2
	// newRecord.Price = dto.Price
	// newRecord.UniqueNumber = dto.UniqueNumber

	payment_date_expired := time.Now().AddDate(0, 0, 1)
	newRecord.PaymentDate = &currentTime
	newRecord.PaymentDateExpired = &payment_date_expired

	var result helper.Response
	result = service.paymentRepository.Insert(newRecord)
	if !result.Status {
		return helper.StandartResult(false, result.Message, result.Data)
	}

	return helper.StandartResult(true, "Ok", result.Data)
}

func (service *paymentService) CreateEWalletPayment(dto ewallet.PaymentDto) helper.Response {
	token := dto.Context.GetHeader("Authorization")

	userIdentity := service.jwtService.GetUserByToken(token)
	xenditService := ewallet.NewPaymentRequest(dto.Context)

	var newRecord = models.Payment{}
	newRecord.Id = uuid.New().String()
	newRecord.RecordId = dto.RecordId //-- WebinarId or MembershipId

	dto.RecordId = newRecord.Id //-- RecordId replace by NewPaymentId
	dto.ApplicationUserId = userIdentity.UserId
	res, err := xenditService.CreateEWalletCharge(dto)
	if err != nil {
		return helper.ServerResponse(false, fmt.Sprintf("%v", err.Message), "", helper.EmptyObj{})
	}

	currentTime := time.Now()

	newRecord.CreatedAt = &currentTime
	newRecord.IsActive = true
	newRecord.EntityId = userIdentity.EntityId
	newRecord.CreatedBy = userIdentity.UserId
	newRecord.OwnerId = userIdentity.UserId
	newRecord.ApplicationUserId = userIdentity.UserId
	newRecord.Currency = "IDR"
	newRecord.OrderNumber = fmt.Sprintf("%v/EW/%v/%v/%v", "JORAN", strconv.Itoa(time.Now().Year()), helper.NumberMonthToRoman(int(time.Now().Month())), strings.ToUpper((strconv.Itoa(time.Now().Nanosecond()))[0:5]))

	newRecord.PaymentType = dto.PaymentType
	newRecord.PaymentStatus = 2
	newRecord.Price = int(dto.Amount)
	newRecord.UniqueNumber = 0

	newRecord.PaymentDate = &currentTime
	if dto.PaymentType == string(xendit.EWalletTypeLINKAJA) {
		payment_date_expired := time.Now().Add(time.Minute * 5)
		newRecord.PaymentDateExpired = &payment_date_expired
	} else if dto.PaymentType == string(xendit.EWalletTypeOVO) {
		payment_date_expired := time.Now().Add(time.Second * 55)
		newRecord.PaymentDateExpired = &payment_date_expired
	}
	newRecord.ProviderName = string(payment_gateway_providers.Xendit)
	newRecord.ProviderRecordId = res.ID
	newRecord.ProviderBusinessId = res.BusinessID
	newRecord.ProviderReferenceId = res.ReferenceID
	// res.BusinessID

	var result helper.Response
	result = service.paymentRepository.Insert(newRecord)
	if !result.Status {
		return result
	}

	return helper.ServerResponse(true, "Ok", "", res)
}

func (service *paymentService) GetEWalletPaymentStatus(ctx *gin.Context, referenceId string) helper.Response {
	xendit := ewallet.NewPaymentRequest(ctx)

	res := service.paymentRepository.GetByProviderReferenceId(referenceId)
	if !res.Status {
		return res
	}

	xenditPaymentResult, err := xendit.GetEWalletPaymentStatus(res.Data.(models.Payment).ProviderRecordId)
	if err != nil {
		return helper.ServerResponse(false, fmt.Sprintf("%v", err.Message), "", helper.EmptyObj{})
	}

	if xenditPaymentResult.Status != "PENDING" {
		//-- Update Payment
	}

	//------------ Payment Record -------------
	var mapResult map[string]interface{}
	temp, _ := json.Marshal(res.Data)
	json.Unmarshal(temp, &mapResult)

	//-------- Payment Status E-Wallet --------
	var xenditPaymentStatus map[string]interface{}
	tempXenditResult, _ := json.Marshal(xenditPaymentResult)
	json.Unmarshal(tempXenditResult, &xenditPaymentStatus)

	mapResult["payment_status_ewallet"] = xenditPaymentStatus

	return helper.Response{Status: true, Message: "Ok", Data: mapResult}
}

func (service *paymentService) Update(record models.Payment) helper.Response {
	return service.paymentRepository.Update(record)
}

func (service *paymentService) UpdatePaymentStatus(req dto.UpdatePaymentStatusDto) helper.Response {
	return service.paymentRepository.UpdatePaymentStatus(req)
}

func (service *paymentService) UpdateWalletPaymentStatus(dto dto.UpdatePaymentStatusDto) helper.Response {
	return service.paymentRepository.UpdatePaymentStatus(dto)
}

func (service *paymentService) GetById(recordId string) helper.Response {
	return service.paymentRepository.GetById(recordId)
}

func (service *paymentService) GetByProviderRecordId(id string) helper.Response {
	return service.paymentRepository.GetByProviderRecordId(id)
}

func (service *paymentService) GetByProviderReferenceId(id string) helper.Response {
	return service.paymentRepository.GetByProviderReferenceId(id)
}

func (service *paymentService) DeleteById(recordId string) helper.Response {
	return service.paymentRepository.DeleteById(recordId)
}
