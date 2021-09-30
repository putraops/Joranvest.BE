package repository

import (
	"database/sql"
	"fmt"
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/view_models"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MembershipUserRepository interface {
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	GetPagination(request commons.PaginationRequest) interface{}
	GetAll(filter map[string]interface{}) []models.MembershipUser
	Insert(membershipUser models.MembershipUser, payment models.MembershipPayment) helper.Response
	Update(record models.MembershipUser) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type membershipUserConnection struct {
	connection           *gorm.DB
	serviceRepository    ServiceRepository
	filemasterRepository FilemasterRepository
	membershipRepository MembershipRepository
	tableName            string
	viewQuery            string
}

func NewMembershipUserRepository(db *gorm.DB) MembershipUserRepository {
	return &membershipUserConnection{
		connection:           db,
		tableName:            models.MembershipUser.TableName(models.MembershipUser{}),
		viewQuery:            entity_view_models.EntityMembershipUserView.ViewModel(entity_view_models.EntityMembershipUserView{}),
		serviceRepository:    NewServiceRepository(db),
		filemasterRepository: NewFilemasterRepository(db),
		membershipRepository: NewMembershipRepository(db),
	}
}

func (db *membershipUserConnection) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	var records []entity_view_models.EntityMembershipUserView
	var res commons.DataTableResponse

	var conditions = ""
	var orderpart = ""
	if request.Draw == 1 && request.DataTableDefaultOrder.Column != "" {
		var column = request.DataTableDefaultOrder.Column
		orderpart = column + " " + request.DataTableDefaultOrder.Dir
	} else {
		var column = request.DataTableColumn[request.DataTableOrder[0].Column].Name
		orderpart = column + " " + request.DataTableOrder[0].Dir
	}
	start := fmt.Sprintf("%v", request.Start)
	length := fmt.Sprintf("%v", (request.Start + request.Length))

	if len(request.Filter) > 0 {
		for _, s := range request.Filter {
			conditions += " AND (" + s.Column + " = '" + s.Value + "') "
		}
	}

	if request.Search.Value != "" {
		conditions += " AND ("
		var totalFilter int = 0
		for _, s := range request.DataTableColumn {
			if s.Searchable {
				if totalFilter > 0 {
					conditions += " OR "
				}
				conditions += fmt.Sprintf("LOWER(CAST (%v AS varchar))", s.Name) + " LIKE '%" + request.Search.Value + "%' "
				totalFilter++
			}
		}
		conditions += ")"
	}

	var sql strings.Builder
	var sqlCount strings.Builder
	sql.WriteString(fmt.Sprintf("SELECT * FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s) peta_rn, ", orderpart))
	sql.WriteString(strings.Replace(db.viewQuery, "SELECT", "", -1))
	sql.WriteString(" WHERE 1 = 1 ")
	sql.WriteString(conditions)
	sql.WriteString(") peta_paged ")
	sql.WriteString(fmt.Sprintf("WHERE peta_rn > %s AND peta_rn <= %s ", start, length))
	db.connection.Raw(sql.String()).Scan(&records)

	sqlCount.WriteString(db.serviceRepository.ConvertViewQueryIntoViewCount(db.viewQuery))
	sqlCount.WriteString("WHERE 1=1")
	sqlCount.WriteString(conditions)
	db.connection.Raw(sqlCount.String()).Scan(&res.RecordsFiltered)

	res.Draw = request.Draw
	if len(records) > 0 {
		res.RecordsTotal = res.RecordsFiltered
		res.DataRow = records
	} else {
		res.RecordsTotal = 0
		res.RecordsFiltered = 0
		res.DataRow = []entity_view_models.EntityMembershipUserView{}
	}
	return res
}

func (db *membershipUserConnection) GetPagination(request commons.PaginationRequest) interface{} {
	var response commons.PaginationResponse
	var records []entity_view_models.EntityMembershipUserView

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
	for k, v := range request.Filter {
		if v != "" {
			if total_filter > 0 {
				filters += "AND "
			}
			filters += fmt.Sprintf("%v = '%v' ", k, v)
			total_filter++
		}
	}
	// #endregion

	offset := (page - 1) * pageSize
	db.connection.Where(filters).Order(orders).Offset(offset).Limit(pageSize).Find(&records)

	var count int64
	db.connection.Model(&entity_view_models.EntityMembershipUserView{}).Where(filters).Count(&count)

	response.Data = records
	response.Total = int(count)
	return response
}

func (db *membershipUserConnection) GetAll(filter map[string]interface{}) []models.MembershipUser {
	var records []models.MembershipUser
	if len(filter) == 0 {
		db.connection.Find(&records)
	} else if len(filter) != 0 {
		db.connection.Where(filter).Find(&records)
	}
	return records
}

func (db *membershipUserConnection) Insert(membershipUser models.MembershipUser, payment models.MembershipPayment) helper.Response {
	tx := db.connection.Begin()

	//-- Payment Record
	payment.Id = uuid.New().String()
	payment.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	if payment.PaymentStatus == 200 {
		payment.PaymentDate = sql.NullTime{Time: time.Now(), Valid: true}
	}
	if err := tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
	}

	// -- Get Membership Record
	var membershipRecord models.Membership
	if err := tx.First(&membershipRecord, "id = ?", membershipUser.MembershipId).Error; err != nil || membershipRecord.Id == "" {
		tx.Rollback()
		return helper.ServerResponse(false, "Membership Record not found", fmt.Sprintf("%v,", err), helper.EmptyObj{})
	}

	// Calculate Expired Date
	if payment.PaymentStatus == 200 {
		membershipUser.ExpiredDate = sql.NullTime{
			Time:  payment.PaymentDate.Time.AddDate(0, 1, 0),
			Valid: true,
		}
	}

	membershipUser.Id = uuid.New().String()
	membershipUser.MembershipPaymentId = payment.Id
	membershipUser.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	if err := tx.Create(&membershipUser).Error; err != nil {
		tx.Rollback()
		return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
	}

	tx.Commit()
	return helper.ServerResponse(true, "Ok", "", membershipUser.Id)
}

func (db *membershipUserConnection) Update(record models.MembershipUser) helper.Response {
	tx := db.connection.Begin()
	var oldRecord models.MembershipUser
	db.connection.First(&oldRecord, "id = ?", record.Id)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}

	//-- Delete Fundamental Analysis Tag
	var tags models.FundamentalAnalysisTag
	if err := tx.Where("fundamental_analysis_id = ?", record.Id).Delete(&tags).Error; err != nil {
		tx.Rollback()
		return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
	}

	record.IsActive = oldRecord.IsActive
	record.CreatedAt = oldRecord.CreatedAt
	record.CreatedBy = oldRecord.CreatedBy
	record.EntityId = oldRecord.EntityId
	record.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	res := tx.Save(&record)
	if res.RowsAffected == 0 {
		return helper.ServerResponse(false, fmt.Sprintf("%v,", res.Error), fmt.Sprintf("%v,", res.Error), helper.EmptyObj{})
	}

	tx.Commit()
	db.connection.Preload(clause.Associations).Find(&record)
	return helper.ServerResponse(true, "Ok", "", record)
}

func (db *membershipUserConnection) GetById(recordId string) helper.Response {
	var record models.MembershipUser
	db.connection.Preload("Emiten").First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	res := helper.ServerResponse(true, "Ok", "", record)
	return res
}

func (db *membershipUserConnection) DeleteById(recordId string) helper.Response {
	var record models.MembershipUser
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
