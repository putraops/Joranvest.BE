package repository

import (
	"database/sql"
	"fmt"
	"joranvest/commons"
	"joranvest/dto"
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
	GetPagination(request commons.Pagination2ndRequest) interface{}
	GetAll(filter map[string]interface{}) []models.Payment
	GetUniqueNumber() int
	MembershipPayment(t models.Payment) helper.Response
	WebinarPayment(t models.Payment) helper.Response
	Update(record models.Payment) helper.Response
	UpdatePaymentStatus(req dto.UpdatePaymentStatusDto) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type paymentConnection struct {
	connection               *gorm.DB
	serviceRepository        ServiceRepository
	membershipUserRepository MembershipUserRepository
	webinarRegistrationRepo  WebinarRegistrationRepository
	filemasterRepository     FilemasterRepository
	tableName                string
	viewQuery                string
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentConnection{
		connection:               db,
		tableName:                models.Payment.TableName(models.Payment{}),
		viewQuery:                entity_view_models.EntityPaymentView.ViewModel(entity_view_models.EntityPaymentView{}),
		serviceRepository:        NewServiceRepository(db),
		membershipUserRepository: NewMembershipUserRepository(db),
		webinarRegistrationRepo:  NewWebinarRegistrationRepository(db),
		filemasterRepository:     NewFilemasterRepository(db),
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
		record.PaymentDate = sql.NullTime{Time: time.Now().Local().UTC(), Valid: true}
	} else {
		record.PaymentDateExpired = sql.NullTime{Time: time.Now().Local().UTC().AddDate(0, 0, 1), Valid: true}
	}
	record.CreatedAt = sql.NullTime{Time: time.Now().Local().UTC(), Valid: true}

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
		membershipUser.MembershipId = record.RecordId
		membershipUser.PaymentId = record.Id
		if record.UpdatedBy != "" {
			membershipUser.CreatedBy = record.UpdatedBy
		} else {
			membershipUser.CreatedBy = record.CreatedBy
		}
		membershipUser.OwnerId = record.OwnerId
		membershipUser.ApplicationUserId = record.CreatedBy
		membershipUser.CreatedAt = sql.NullTime{Time: time.Now().Local().UTC(), Valid: true}
		membershipUser.ExpiredDate = sql.NullTime{
			Time:  record.PaymentDate.Time.AddDate(0, int(membershipRecord.Duration), 0),
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

func (db *paymentConnection) WebinarPayment(record models.Payment) helper.Response {
	commons.Logger()
	tx := db.connection.Begin()

	record.Id = uuid.New().String()
	if record.PaymentStatus == 200 {
		record.PaymentDate = sql.NullTime{Time: time.Now().Local().UTC(), Valid: true}
	} else {
		record.PaymentDateExpired = sql.NullTime{Time: time.Now().Local().UTC().AddDate(0, 0, 1), Valid: true}
	}
	record.CreatedAt = sql.NullTime{Time: time.Now().Local().UTC(), Valid: true}

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
		webinarRegistrationRecord.CreatedAt = sql.NullTime{Time: time.Now().Local().UTC(), Valid: true}

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
	record.UpdatedAt = sql.NullTime{Time: time.Now().Local().UTC(), Valid: true}
	res := tx.Save(&record)
	if res.RowsAffected == 0 {
		return helper.ServerResponse(false, fmt.Sprintf("%v,", res.Error), fmt.Sprintf("%v,", res.Error), helper.EmptyObj{})
	}

	tx.Commit()
	db.connection.Preload(clause.Associations).Find(&record)
	return helper.ServerResponse(true, "Ok", "", record)
}

func (db *paymentConnection) UpdatePaymentStatus(req dto.UpdatePaymentStatusDto) helper.Response {
	tx := db.connection.Begin()
	var paymentRecord models.Payment
	db.connection.First(&paymentRecord, "id = ?", req.Id)
	if req.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}

	paymentRecord.PaymentStatus = req.PaymentStatus
	paymentRecord.UpdatedBy = req.UpdatedBy
	paymentRecord.PaymentDate = sql.NullTime{Time: time.Now().Local().UTC(), Valid: true}
	paymentRecord.UpdatedAt = paymentRecord.PaymentDate

	var viewRecord entity_view_models.EntityPaymentView
	db.connection.First(&viewRecord, "id = ?", req.Id)
	if req.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}

	if paymentRecord.PaymentStatus == 200 {
		if viewRecord.MembershipName != "" {
			//.. Insert Membership User
			res := db.membershipUserRepository.SetMembership(viewRecord.RecordId, paymentRecord)
			if !res.Status {
				return res
			}
		} else if viewRecord.WebinarTitle != "" {
			//.. Insert Webinar Registration
			var webinarRegistrationRecord models.WebinarRegistration
			webinarRegistrationRecord.CreatedBy = req.UpdatedBy
			webinarRegistrationRecord.PaymentId = req.Id
			webinarRegistrationRecord.ApplicationUserId = paymentRecord.CreatedBy
			webinarRegistrationRecord.WebinarId = viewRecord.RecordId

			res := db.webinarRegistrationRepo.Insert(webinarRegistrationRecord)
			if !res.Status {
				return res
			}
		}
	}

	res := tx.Save(&paymentRecord)
	if res.RowsAffected == 0 {
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
