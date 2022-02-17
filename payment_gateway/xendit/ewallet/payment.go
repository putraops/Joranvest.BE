package ewallet

import (
	"fmt"
	"joranvest/commons"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/ewallet"
)

type Payment interface {
	CreateEWalletCharge(dto PaymentDto) (*xendit.EWalletCharge, *xendit.Error)
	GetEWalletPaymentStatus(providerRecordId string) (*xendit.EWalletCharge, *xendit.Error)
}

type apiRequesterMock struct {
	context   *gin.Context
	secretKey string
	mock.Mock
}

type PaymentDto struct {
	Amount             float64 `json:"amount" form:"amount" binding:"required"`
	PhoneNumber        string  `json:"phone_number" form:"phone_number"`
	PaymentType        string  `json:"payment_type" form:"payment_type" binding:"required"`
	RecordId           string  `json:"record_id" form:"record_id" binding:"required"`
	SuccessRedirectUrl string  `json:"success_redirect_url" form:"success_redirect_url" binding:"required"`
	EntityId           string  `json:"-"`
	UpdatedBy          string
	ApplicationUserId  string
	Context            *gin.Context
}

func NewPaymentRequest(context *gin.Context) Payment {
	return &apiRequesterMock{
		context:   context,
		secretKey: os.Getenv("XENDIT_SECRET_KEY"),
	}
}

func (r apiRequesterMock) CreateEWalletCharge(dto PaymentDto) (*xendit.EWalletCharge, *xendit.Error) {
	commons.Logger()

	xendit.Opt.SecretKey = r.secretKey
	channelProperties := make(map[string]string)
	channelProperties["mobile_number"] = dto.PhoneNumber
	var channelCode = ""
	var checkoutMethod = "ONE_TIME_PAYMENT"
	var referenceId = fmt.Sprintf("%v-%v", strings.ToLower(string(dto.PaymentType)), uuid.New().String())
	var successRedirectUrl = fmt.Sprintf(`%v/%v`, dto.SuccessRedirectUrl, dto.RecordId)

	if dto.PaymentType == string(xendit.EWalletTypeLINKAJA) {
		channelCode = string(EWalletChannelCodeLINKAJA)
		channelProperties["success_redirect_url"] = successRedirectUrl
	} else if dto.PaymentType == string(xendit.EWalletTypeOVO) {
		channelCode = string(EWalletChannelCodeOVO)
	} else {
		log.Error("Payment Type Unsupported")
		return nil, &xendit.Error{Status: 0, Message: "Payment Type Unsupported", ErrorCode: "200"}
	}

	data := ewallet.CreateEWalletChargeParams{
		ReferenceID:       referenceId,
		Currency:          "IDR",
		Amount:            dto.Amount,
		CheckoutMethod:    checkoutMethod,
		ChannelCode:       channelCode,
		ChannelProperties: channelProperties,
		Metadata: map[string]interface{}{
			"merchant":  "joranvest",
			"record_id": dto.RecordId,
			"user_id":   dto.ApplicationUserId,
		},
	}

	//fmt.Println("tess")
	charge, chargeErr := ewallet.CreateEWalletCharge(&data)
	if chargeErr != nil {
		log.Error("Function: CreateEWalletCharge")
		log.Error(data)
		log.Error(chargeErr.Status)
		log.Error(chargeErr.Message)
		log.Error(chargeErr.ErrorCode)
		return nil, chargeErr
	}
	return charge, nil
}

func (r apiRequesterMock) GetEWalletPaymentStatus(providerRecordId string) (*xendit.EWalletCharge, *xendit.Error) {
	xendit.Opt.SecretKey = r.secretKey
	data := ewallet.GetEWalletChargeStatusParams{
		ChargeID: providerRecordId,
	}

	res, resErr := ewallet.GetEWalletChargeStatus(&data)
	if resErr != nil {
		log.Error("Function: GetEWalletPaymentStatus")
		log.Error(resErr.Status)
		log.Error(resErr.Message)
		log.Error(resErr.ErrorCode)
		return nil, resErr
	}
	return res, nil
}

type EWalletChannelCode string
type XenditPaymentStatus string
type XenditPaymentService string

// This consists the values that EWalletChannelCode can take
const (
	EWalletChannelCodeOVO       EWalletChannelCode = "ID_OVO"
	EWalletChannelCodeDANA      EWalletChannelCode = "ID_DANA"
	EWalletChannelCodeLINKAJA   EWalletChannelCode = "ID_LINKAJA"
	EWalletChannelCodeSHOPEEPAY EWalletChannelCode = "ID_SHOPEEPAY"
)

const (
	XenditPaymentStatusPending   XenditPaymentStatus = "PENDING"
	XenditPaymentStatusFailed    XenditPaymentStatus = "FAILED"
	XenditPaymentStatusSucceeded XenditPaymentStatus = "SUCCEEDED"
)

const (
	XenditPaymentServiceOneTimePayment XenditPaymentService = "ONE_TIME_PAYMENT"
)
