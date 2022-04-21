package repository

import (
	"database/sql"
	"fmt"
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/entity_view_models"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PaymentRepository interface {
	Insert(record models.Payment) helper.Response
	GetPagination(request commons.Pagination2ndRequest) interface{}
	GetAll(filter map[string]interface{}) []models.Payment
	GetUniqueNumber() int
	MembershipPayment(t models.Payment) helper.Response
	WebinarPayment(t models.Payment) helper.Response
	Update(record models.Payment) helper.Response
	UpdatePaymentStatus(paymentRecord models.Payment) helper.Response
	GetById(recordId string) helper.Response
	GetViewById(recordId string) helper.Result
	GetByProviderRecordId(id string) helper.Response
	GetByProviderReferenceId(id string) helper.Response
	DeleteById(recordId string) helper.Response
}

type paymentConnection struct {
	connection               *gorm.DB
	serviceRepository        ServiceRepository
	membershipUserRepository MembershipUserRepository
	webinarRegistrationRepo  WebinarRegistrationRepository
	membershipRepository     MembershipRepository
	productRepository        ProductRepository
	filemasterRepository     FilemasterRepository
	tableName                string
	viewQuery                string
	currentTime              time.Time
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentConnection{
		connection:               db,
		tableName:                models.Payment.TableName(models.Payment{}),
		viewQuery:                entity_view_models.EntityPaymentView.ViewModel(entity_view_models.EntityPaymentView{}),
		serviceRepository:        NewServiceRepository(db),
		membershipUserRepository: NewMembershipUserRepository(db),
		webinarRegistrationRepo:  NewWebinarRegistrationRepository(db),
		membershipRepository:     NewMembershipRepository(db),
		productRepository:        NewProductRepository(db),
		filemasterRepository:     NewFilemasterRepository(db),
		currentTime:              time.Now(),
	}
}

func (db *paymentConnection) GetPagination(request commons.Pagination2ndRequest) interface{} {
	var response commons.PaginationResponse
	var records []entity_view_models.EntityPaymentView
	var recordsUnfilter []entity_view_models.EntityPaymentView

	page := request.Page
	if page == 0 {
		page = 1
	}

	pageSize := request.Size
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	// #region order
	var orders = "COALESCE(submitted_at, created_at) DESC"
	order_total := 0
	for k, v := range request.Order {
		if order_total == 0 {
			orders = ""
		} else {
			orders += ", "
		}
		orders += fmt.Sprintf("%v %v ", k, v)
		order_total++
	}
	// #endregion

	// #region filter
	var filters = ""
	total_filter := 0
	if len(request.Filter) > 0 {
		for _, v := range request.Filter {
			if v.Value != "" {
				if total_filter > 0 {
					filters += "AND "
				}

				if v.Operator == "" {
					filters += fmt.Sprintf("%v %v ", v.Field, v.Value)
				} else {
					filters += fmt.Sprintf("%v %v '%v' ", v.Field, v.Operator, v.Value)
				}
				total_filter++
			}
		}
	}
	// #endregion

	offset := (page - 1) * pageSize
	db.connection.Where(filters).Order(orders).Offset(offset).Limit(pageSize).Find(&records)

	// #region Get Total Data for Pagination
	result := db.connection.Where(filters).Find(&recordsUnfilter)
	if result.Error != nil {
		return helper.ServerResponse(false, fmt.Sprintf("%v,", result.Error), fmt.Sprintf("%v,", result.Error), helper.EmptyObj{})
	}
	response.Total = int(result.RowsAffected)
	// #endregion

	response.Data = records
	return response
}

func (db *paymentConnection) GetAll(filter map[string]interface{}) []models.Payment {
	var records []models.Payment
	if len(filter) == 0 {
		db.connection.Find(&records)
	} else if len(filter) != 0 {
		db.connection.Where(filter).Find(&records)
	}
	return records
}

func (db *paymentConnection) GetUniqueNumber() int {
	var record models.Payment
	unique_number := 0
	db.connection.Order("created_at DESC").Where("to_char(created_at, 'YYYY-MM-DD') = to_char(CURRENT_DATE, 'YYYY-MM-DD')").First(&record)
	if record.Id == "" {
		unique_number = 11
	} else {
		unique_number = record.UniqueNumber + 1
	}
	return unique_number
}

func (db *paymentConnection) MembershipPayment(record models.Payment) helper.Response {
	commons.Logger()
	tx := db.connection.Begin()
	// productName := *product

	record.Id = uuid.New().String()
	if record.PaymentStatus == 200 {
		payment_date := time.Now()
		record.PaymentDate = &payment_date
	} else {
		payment_date_expired := time.Now().AddDate(0, 0, 1)
		record.PaymentDateExpired = &payment_date_expired
	}
	record.CreatedAt = &db.currentTime

	if err := tx.Save(&record).Error; err != nil {
		tx.Rollback()
		log.Error(db.serviceRepository.getCurrentFuncName())
		log.Error(fmt.Sprintf("%v,", err))
		return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
	}

	// #region Set After Payment Paid
	if record.PaymentStatus == 200 {
		var membershipRecord models.Membership
		if err := tx.First(&membershipRecord, "id = ?", record.RecordId).Error; err != nil || membershipRecord.Id == "" {
			log.Error(db.serviceRepository.getCurrentFuncName())
			log.Error(fmt.Sprintf("%v,", err))
			return helper.ServerResponse(false, "Membership Record not found", fmt.Sprintf("%v,", err), helper.EmptyObj{})
		}

		var membershipUser models.MembershipUser
		membershipUser.Id = uuid.New().String()
		membershipUser.MembershipId = &record.RecordId
		membershipUser.PaymentId = record.Id
		if record.UpdatedBy != "" {
			membershipUser.CreatedBy = record.UpdatedBy
		} else {
			membershipUser.CreatedBy = record.CreatedBy
		}
		membershipUser.OwnerId = record.OwnerId
		membershipUser.ApplicationUserId = record.CreatedBy
		membershipUser.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
		membershipUser.ExpiredDate = sql.NullTime{
			Time:  record.PaymentDate.AddDate(0, int(membershipRecord.Duration), 0),
			Valid: true,
		}

		if err := tx.Create(&membershipUser).Error; err != nil {
			tx.Rollback()
			log.Error(db.serviceRepository.getCurrentFuncName())
			log.Error(fmt.Sprintf("%v,", err))
			return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
		}
	}
	// #endregion

	tx.Commit()
	db.connection.Find(&record)
	return helper.ServerResponse(true, "Ok", "", record)
}

func (db *paymentConnection) Insert(record models.Payment) helper.Response {
	commons.Logger()
	tx := db.connection.Begin()
	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		log.Error(db.serviceRepository.getCurrentFuncName())
		log.Error(fmt.Sprintf("%v,", err))
		return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
	}

	tx.Commit()
	db.connection.Find(&record)
	return helper.ServerResponse(true, "Ok", "", record)
}

func (db *paymentConnection) WebinarPayment(record models.Payment) helper.Response {
	commons.Logger()
	tx := db.connection.Begin()

	record.Id = uuid.New().String()
	if record.PaymentStatus == 200 {
		payment_date := time.Now()
		record.PaymentDate = &payment_date
	} else {
		payment_date_expired := time.Now().AddDate(0, 0, 1)
		record.PaymentDateExpired = &payment_date_expired
	}
	record.CreatedAt = &db.currentTime

	if err := tx.Save(&record).Error; err != nil {
		tx.Rollback()
		log.Error(db.serviceRepository.getCurrentFuncName())
		log.Error(fmt.Sprintf("%v,", err))
		return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
	}

	//#region Set After Payment Paid
	if record.PaymentStatus == 200 {
		var webinarRecord models.Webinar
		if err := tx.First(&webinarRecord, "id = ?", record.RecordId).Error; err != nil || webinarRecord.Id == "" {
			log.Error(db.serviceRepository.getCurrentFuncName())
			log.Error(fmt.Sprintf("%v,", err))
			return helper.ServerResponse(false, "Webinar Record not found", fmt.Sprintf("%v,", err), helper.EmptyObj{})
		}

		var webinarRegistrationRecord models.WebinarRegistration
		webinarRegistrationRecord.Id = uuid.New().String()
		webinarRegistrationRecord.WebinarId = record.RecordId
		webinarRegistrationRecord.PaymentId = record.Id
		if record.UpdatedBy != "" {
			webinarRegistrationRecord.CreatedBy = record.UpdatedBy
		} else {
			webinarRegistrationRecord.CreatedBy = record.CreatedBy
		}
		webinarRegistrationRecord.OwnerId = record.OwnerId
		webinarRegistrationRecord.ApplicationUserId = record.CreatedBy
		webinarRegistrationRecord.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}

		if err := tx.Create(&webinarRegistrationRecord).Error; err != nil {
			tx.Rollback()
			log.Error(db.serviceRepository.getCurrentFuncName())
			log.Error(fmt.Sprintf("%v,", err))
			return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
		}
	}
	// #endregion

	tx.Commit()
	db.connection.Find(&record)
	return helper.ServerResponse(true, "Ok", "", record)
}

func (db *paymentConnection) Update(record models.Payment) helper.Response {
	tx := db.connection.Begin()
	var oldRecord models.Payment
	db.connection.First(&oldRecord, "id = ?", record.Id)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}

	record.IsActive = oldRecord.IsActive
	record.CreatedAt = oldRecord.CreatedAt
	record.CreatedBy = oldRecord.CreatedBy
	record.EntityId = oldRecord.EntityId
	record.UpdatedAt = &db.currentTime
	res := tx.Save(&record)
	if res.RowsAffected == 0 {
		return helper.ServerResponse(false, fmt.Sprintf("%v,", res.Error), fmt.Sprintf("%v,", res.Error), helper.EmptyObj{})
	}

	tx.Commit()
	db.connection.Preload(clause.Associations).Find(&record)
	return helper.ServerResponse(true, "Ok", "", record)
}

func (db *paymentConnection) UpdatePaymentStatus(paymentRecord models.Payment) helper.Response {
	commons.Logger()

	tx := db.connection.Begin()
	// var paymentRecord models.Payment
	// db.connection.First(&paymentRecord, "id = ?", req.Id)
	// if req.Id == "" {
	// 	log.Error("Record not found")
	// 	log.Error("Function: UpdatePaymentStatus")
	// 	res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
	// 	return res
	// }

	// paymentRecord.PaymentStatus = req.PaymentStatus
	// paymentRecord.UpdatedBy = req.UpdatedBy
	// paymentRecord.UpdatedAt = &db.currentTime
	// paymentRecord.PaymentDate = &db.currentTime

	var viewRecord entity_view_models.EntityPaymentView
	db.connection.First(&viewRecord, "id = ?", paymentRecord.Id)
	if viewRecord.Id == "" {
		log.Error("Payment Record not found")
		res := helper.ServerResponse(false, "Payment Record not found", "Error", helper.EmptyObj{})
		tx.Rollback()
		return res
	}

	if paymentRecord.PaymentStatus == 200 {
		if viewRecord.MembershipName != "" {
			//.. Check Exist
			curentMembershipUserResponse := db.membershipUserRepository.GetExistMembershipByUserId(paymentRecord.ApplicationUserId, true)
			if curentMembershipUserResponse.Status {
				paymentRecord.IsExtendMembership = true

				var membershipUserRecord models.MembershipUser
				membershipUserRecord = curentMembershipUserResponse.Data.(models.MembershipUser)

				//-- Get Membership Record
				membershipResponse := db.membershipRepository.GetById(paymentRecord.RecordId)
				if membershipResponse.Status {
					membershipUserRecord.MembershipId = &paymentRecord.RecordId
				} else {
					log.Error("Membership Record Not Found")
					tx.Rollback()
					return membershipResponse
				}

				membershipUserRecord.ExpiredDate = sql.NullTime{
					Time:  membershipUserRecord.ExpiredDate.Time.AddDate(0, int(membershipResponse.Data.(models.Membership).Duration), 1),
					Valid: true,
				}
				membershipUserRecord.PaymentId = paymentRecord.Id

				//-- Update Membershipuser
				res := tx.Save(&membershipUserRecord)
				if res.RowsAffected == 0 {
					tx.Rollback()
					return helper.ServerResponse(false, fmt.Sprintf("%v,", res.Error), fmt.Sprintf("%v,", res.Error), helper.EmptyObj{})
				}

			} else {
				//.. Insert Membership User
				res := db.membershipUserRepository.SetMembership(viewRecord)
				if !res.Status {
					tx.Rollback()
					return res
				}
			}
		} else if viewRecord.ProductId != "" {
			//.. Joranvest Chart System
			//.. Check Exist
			curentMembershipUserResponse := db.membershipUserRepository.GetExistMembershipByUserId(paymentRecord.ApplicationUserId, false)
			if curentMembershipUserResponse.Status {

				fmt.Println("----------------------- paymentRecord -----------------------")
				fmt.Println(viewRecord)
				fmt.Println(viewRecord.PaymentDate)
				fmt.Println(viewRecord.PaymentDate)
				fmt.Println(viewRecord.PaymentDateExpired)
				fmt.Println(viewRecord.PaymentDateExpired)
				fmt.Println("----------------------- paymentRecord -----------------------")

				paymentRecord.IsExtendMembership = true

				var membershipUserRecord models.MembershipUser
				membershipUserRecord = curentMembershipUserResponse.Data.(models.MembershipUser)

				//-- Get Product Record
				productResponse := db.productRepository.GetById(paymentRecord.RecordId)
				if productResponse.Status {
					membershipUserRecord.ProductId = &paymentRecord.RecordId
				} else {
					log.Error("Product Record Not Found")
					tx.Rollback()
					return helper.ServerResponse(false, productResponse.Message, productResponse.Message, productResponse.Data)
				}

				duration := productResponse.Data.(models.Product).Duration
				membershipUserRecord.ExpiredDate = sql.NullTime{
					Time:  membershipUserRecord.ExpiredDate.Time.AddDate(0, *duration, 1),
					Valid: true,
				}
				membershipUserRecord.PaymentId = paymentRecord.Id

				//-- Update Membershipuser
				res := tx.Save(&membershipUserRecord)
				if res.RowsAffected == 0 {
					tx.Rollback()
					return helper.ServerResponse(false, fmt.Sprintf("%v,", res.Error), fmt.Sprintf("%v,", res.Error), helper.EmptyObj{})
				}

			} else {
				//.. Insert Membership User
				res := db.membershipUserRepository.SetMembership(viewRecord)
				if !res.Status {
					tx.Rollback()
					return res
				}
			}
		} else if viewRecord.WebinarTitle != "" {
			//.. Insert Webinar Registration
			var webinarRegistrationRecord models.WebinarRegistration
			webinarRegistrationRecord.CreatedBy = paymentRecord.UpdatedBy
			webinarRegistrationRecord.PaymentId = paymentRecord.Id
			webinarRegistrationRecord.ApplicationUserId = paymentRecord.CreatedBy
			webinarRegistrationRecord.WebinarId = viewRecord.RecordId

			res := db.webinarRegistrationRepo.Insert(webinarRegistrationRecord)
			if !res.Status {
				tx.Rollback()
				return res
			}
		}
	}

	res := tx.Save(&paymentRecord)
	if res.RowsAffected == 0 {
		tx.Rollback()
		return helper.ServerResponse(false, fmt.Sprintf("%v,", res.Error), fmt.Sprintf("%v,", res.Error), helper.EmptyObj{})
	}

	tx.Commit()
	return helper.ServerResponse(true, "Ok", "", paymentRecord)
}

func (db *paymentConnection) GetById(recordId string) helper.Response {
	var record models.Payment
	db.connection.First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	res := helper.ServerResponse(true, "Ok", "", record)
	return res
}

func (db *paymentConnection) GetViewById(recordId string) helper.Result {
	var result helper.Result

	var record entity_view_models.EntityPaymentView
	db.connection.First(&record, "id = ?", recordId)
	if record.Id == "" {
		result = helper.StandartResult(false, "Record not found", nil)
	}
	result = helper.StandartResult(true, "Ok", record)

	return result
}

func (db *paymentConnection) GetByProviderRecordId(id string) helper.Response {
	var record models.Payment
	db.connection.First(&record, "provider_record_id = ?", id)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	res := helper.ServerResponse(true, "Ok", "", record)
	return res
}

func (db *paymentConnection) GetByProviderReferenceId(id string) helper.Response {
	var record models.Payment
	db.connection.First(&record, "provider_reference_id = ?", id)
	if record.Id == "" {
		fmt.Println("Record not found")
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	fmt.Println(record.Id)
	res := helper.ServerResponse(true, "Ok", "", record)
	return res
}

func (db *paymentConnection) DeleteById(recordId string) helper.Response {
	var record models.Payment
	db.connection.First(&record, "id = ?", recordId)

	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	} else {
		res := db.connection.Where("id = ?", recordId).Delete(&record)
		if res.RowsAffected == 0 {
			return helper.ServerResponse(false, "Error", fmt.Sprintf("%v", res.Error), helper.EmptyObj{})
		}
		return helper.ServerResponse(true, "Ok", "", helper.EmptyObj{})
	}
}
