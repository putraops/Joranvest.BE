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
	log "github.com/sirupsen/logrus"
	"github.com/xendit/xendit-go"
	"gorm.io/gorm"
)

type PaymentService interface {
	GetPagination(request commons.Pagination2ndRequest) interface{}
	GetAll(filter map[string]interface{}) []models.Payment
	GetUniqueNumber() int
	GetEWalletPaymentStatus(ctx *gin.Context, referenceId string) helper.Response
	// UpdateWalletPaymentStatus(dto dto.UpdatePaymentStatusDto) helper.Response
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
	DB                *gorm.DB
	paymentRepository repository.PaymentRepository
	emailService      EmailService
	helper.AppSession
	jwtService JWTService
}

func NewPaymentService(db *gorm.DB) PaymentService {
	return &paymentService{
		DB:                db,
		paymentRepository: repository.NewPaymentRepository(db),
		jwtService:        NewJWTService(),
		emailService:      NewEmailService(db),
	}
}

func (r *paymentService) GetPagination(request commons.Pagination2ndRequest) interface{} {
	return r.paymentRepository.GetPagination(request)
}

func (r *paymentService) GetAll(filter map[string]interface{}) []models.Payment {
	return r.paymentRepository.GetAll(filter)
}

func (r *paymentService) GetUniqueNumber() int {
	return r.paymentRepository.GetUniqueNumber()
}

func (r *paymentService) MembershipPayment(record models.Payment) helper.Response {
	return r.paymentRepository.MembershipPayment(record)
}
func (r *paymentService) WebinarPayment(record models.Payment) helper.Response {
	return r.paymentRepository.WebinarPayment(record)
}

func (r *paymentService) CreateQRCode(dto qrcode.QRCodeDto) helper.Response {
	token := dto.Context.GetHeader("Authorization")

	userIdentity := r.jwtService.GetUserByToken(token)
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
	result = r.paymentRepository.Insert(newRecord)
	if !result.Status {
		return result
	}

	return helper.ServerResponse(true, "Ok", "", res)
}

func (r *paymentService) CreateTransferPayment(dto dto.PaymentDto) helper.Result {
	token := dto.Context.GetHeader("Authorization")
	userIdentity := r.jwtService.GetUserByToken(token)
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
	result = r.paymentRepository.Insert(newRecord)
	if !result.Status {
		return helper.StandartResult(false, result.Message, result.Data)
	}

	// var asds =
	args := []string{"a", "b"}

	_ = r.emailService.PaymentNotificationToTeam(args...)

	return helper.StandartResult(true, "Ok", result.Data)
}

func (r *paymentService) CreateEWalletPayment(dto ewallet.PaymentDto) helper.Response {
	token := dto.Context.GetHeader("Authorization")

	userIdentity := r.jwtService.GetUserByToken(token)
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
	result = r.paymentRepository.Insert(newRecord)
	if !result.Status {
		return result
	}

	return helper.ServerResponse(true, "Ok", "", res)
}

func (r *paymentService) GetEWalletPaymentStatus(ctx *gin.Context, referenceId string) helper.Response {
	xendit := ewallet.NewPaymentRequest(ctx)

	res := r.paymentRepository.GetByProviderReferenceId(referenceId)
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

func (r *paymentService) Update(record models.Payment) helper.Response {
	return r.paymentRepository.Update(record)
}

func (r *paymentService) UpdatePaymentStatus(req dto.UpdatePaymentStatusDto) helper.Response {
	var paymentRecord models.Payment
	res := r.paymentRepository.GetById(req.Id)
	if !res.Status {
		log.Error(res.Message)
		log.Error("Function: UpdatePaymentStatus")
		return res
	}

	currentTime := time.Now()
	paymentRecord.PaymentStatus = req.PaymentStatus
	paymentRecord.UpdatedBy = req.UpdatedBy
	paymentRecord.UpdatedAt = &currentTime
	paymentRecord.PaymentDate = &currentTime

	paymentResult := r.paymentRepository.UpdatePaymentStatus(paymentRecord)
	if paymentResult.Status {
		//-- send email
		//service.emailService.SendWebinarInformationToParticipants()
	}

	return paymentResult
}

// func (r *paymentService) UpdateWalletPaymentStatus(dto dto.UpdatePaymentStatusDto) helper.Response {

// 	return r.paymentRepository.UpdatePaymentStatus(dto)
// }

func (r *paymentService) GetById(recordId string) helper.Response {
	return r.paymentRepository.GetById(recordId)
}

func (r *paymentService) GetByProviderRecordId(id string) helper.Response {
	return r.paymentRepository.GetByProviderRecordId(id)
}

func (r *paymentService) GetByProviderReferenceId(id string) helper.Response {
	return r.paymentRepository.GetByProviderReferenceId(id)
}

func (r *paymentService) DeleteById(recordId string) helper.Response {
	return r.paymentRepository.DeleteById(recordId)
}
