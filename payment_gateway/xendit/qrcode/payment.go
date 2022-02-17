package qrcode

import (
	"joranvest/commons"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/qrcode"
)

type Payment interface {
	//CreateEWalletCharge(dto QRCodeDto) (*xendit.EWalletCharge, *xendit.Error)
	CreateQRCode(dto QRCodeDto) (*xendit.QRCode, *xendit.Error)
}

type apiRequesterMock struct {
	context   *gin.Context
	secretKey string
	mock.Mock
}

type QRCodeDto struct {
	Amount            float64 `json:"amount" form:"amount" binding:"required"`
	RecordId          string  `json:"record_id" form:"record_id" binding:"required"`
	ProductName       string  `json:"product_name" form:"product_name" binding:"required"`
	EntityId          string  `json:"-"`
	UpdatedBy         string
	ApplicationUserId string
	Context           *gin.Context
}

func NewQRCode(context *gin.Context) Payment {
	return &apiRequesterMock{
		context:   context,
		secretKey: os.Getenv("XENDIT_SECRET_KEY"),
	}
}

func (r apiRequesterMock) CreateQRCode(dto QRCodeDto) (*xendit.QRCode, *xendit.Error) {
	commons.Logger()

	xendit.Opt.SecretKey = r.secretKey

	data := &qrcode.CreateQRCodeParams{
		ExternalID:  dto.RecordId,
		CallbackURL: "https://api.joranvest.com",
		Type:        xendit.DynamicQRCode,
		Amount:      dto.Amount,
	}

	response, err := qrcode.CreateQRCode(data)
	if err != nil {
		log.Error("Function: CreateQRCode")
		log.Error(data)
		log.Error(err.Status)
		log.Error(err.Message)
		log.Error(err.ErrorCode)
		return nil, err
	}
	return response, nil
}

// func (r apiRequesterMock) GetQRPaymentStatus(providerRecordId string) (*xendit.EWalletCharge, *xendit.Error) {
// 	xendit.Opt.SecretKey = r.secretKey
// 	data := ewallet.GetEWalletChargeStatusParams{
// 		ChargeID: providerRecordId,
// 	}

// 	res, resErr := qrcode.ge
// 	if resErr != nil {
// 		log.Error("Function: GetEWalletPaymentStatus")
// 		log.Error(resErr.Status)
// 		log.Error(resErr.Message)
// 		log.Error(resErr.ErrorCode)
// 		return nil, resErr
// 	}
// 	return res, nil
// }

type QRChannelCode string

const (
	QRChannelCodeQRIS QRChannelCode = "QRIS"
)
