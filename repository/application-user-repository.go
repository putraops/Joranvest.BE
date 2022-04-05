package repository

import (
	"database/sql"
	"fmt"
	"joranvest/commons"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/entity_view_models"
	"joranvest/models/request_models"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ApplicationUserRepository interface {
	GetPagination(request commons.Pagination2ndRequest) interface{}

	UserLookup(request helper.ReactSelectRequest) []models.ApplicationUser
	GetViewUserByEmail(username string, email string) interface{}
	GetViewUserByUsernameOrEmail(username string, email string) interface{}
	Insert(t models.ApplicationUser) (models.ApplicationUser, error)
	Update(record models.ApplicationUser) models.ApplicationUser
	UpdateProfile(dtoRecord dto.ApplicationUserDescriptionDto) helper.Response
	UpdateProfilePicture(request request_models.FileRequestDto) helper.Response
	VerifyCredential(username string, email string, password string) interface{}
	IsDuplicateEmail(email string) (tx *gorm.DB)
	GetByEmail(email string) models.ApplicationUser
	RecoverPassword(recordId string, newPassword string) helper.Response
	EmailVerificationById(userId string) helper.Response
	UserProfile(applicationUserId string) models.ApplicationUser
	GetById(applicationUserId string) helper.Response
	GetViewById(applicationUserId string) helper.Response
	GetAll() []models.ApplicationUser
	DeleteById(recordId string) helper.Response

	//-- Next ini hilang
	Lookup(req map[string]interface{}) []models.ApplicationUser
}

type applicationUserConnection struct {
	connection        *gorm.DB
	serviceRepository ServiceRepository
	tableName         string
	viewQuery         string
}

//NewApplicationUserRepository is creates a new instance of ApplicationUserRepository
func NewApplicationUserRepository(db *gorm.DB) ApplicationUserRepository {
	return &applicationUserConnection{
		connection:        db,
		serviceRepository: NewServiceRepository(db),
		tableName:         models.ApplicationUser.TableName(models.ApplicationUser{}),
		viewQuery:         entity_view_models.EntityApplicationUserView.ViewModel(entity_view_models.EntityApplicationUserView{}),
	}
}

func (db *applicationUserConnection) GetPagination(request commons.Pagination2ndRequest) interface{} {
	var response commons.PaginationResponse
	var records []entity_view_models.EntityApplicationUserView

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

	var count int64
	db.connection.Model(&entity_view_models.EntityApplicationUserView{}).Where(filters).Count(&count)

	response.Data = records
	response.Total = int(count)
	return response
}

//-- Next ini hilang
func (db *applicationUserConnection) Lookup(req map[string]interface{}) []models.ApplicationUser {
	records := []models.ApplicationUser{}
	db.connection.Order("first_name, last_name asc")

	var sqlQuery strings.Builder
	sqlQuery.WriteString("SELECT * FROM " + db.tableName)

	v, found := req["condition"]
	if found {
		sqlQuery.WriteString(" WHERE 1 = 1 AND is_admin = false ")
		requests := v.(helper.DataFilter).Request
		for _, v := range requests {
			totalFilter := 0
			if v.Operator == "like" {
				for _, valueDetail := range v.Field {
					if totalFilter == 0 {
						sqlQuery.WriteString(" AND (LOWER(" + valueDetail + ") ILIKE " + fmt.Sprint("'%", v.Value, "%'"))
					} else {
						sqlQuery.WriteString(" OR LOWER(" + valueDetail + ") ILIKE " + fmt.Sprint("'%", v.Value, "%'"))
					}
					totalFilter++
				}
			}

			if totalFilter > 0 {
				sqlQuery.WriteString(")")
			}
		}
	}

	fmt.Println("Query: ", sqlQuery.String())
	fmt.Println("Query: ", sqlQuery.String())
	db.connection.Raw(sqlQuery.String()).Scan(&records)
	return records
}

func (db *applicationUserConnection) UserLookup(request helper.ReactSelectRequest) []models.ApplicationUser {
	records := []models.ApplicationUser{}
	var orders = "first_name, last_name ASC"
	var filters = " 1 = 1 AND is_admin = false "

	if request.Q != "" {
		totalFilter := 0
		for _, field := range request.Field {
			if totalFilter == 0 {
				filters += " AND ("
			} else {
				filters += " OR "
			}

			filters += fmt.Sprintf("%v ILIKE %v", field, fmt.Sprint("'%", request.Q, "%'"))
			totalFilter++
		}
		if totalFilter > 0 {
			filters += ")"
		}
	}

	offset := (request.Page - 1) * request.Size
	db.connection.Debug().Where(filters).Order(orders).Offset(offset).Limit(request.Size).Find(&records)
	return records
}

func (db *applicationUserConnection) GetViewUserByUsernameOrEmail(username string, email string) interface{} {
	var record entity_view_models.EntityApplicationUserView
	res := db.connection.Where("LOWER(username) = ? AND (LOWER(username) <> '' OR LOWER(username) IS NULL) ", strings.ToLower(username)).Or("LOWER(email) = ?", strings.ToLower(email)).Take(&record)
	if res.Error == nil {
		return record
	}
	return nil
}
func (db *applicationUserConnection) GetViewUserByEmail(username string, email string) interface{} {
	var record entity_view_models.EntityApplicationUserView
	res := db.connection.Where("LOWER(email) = ?", strings.ToLower(email)).Take(&record)
	if res.Error == nil {
		return record
	}
	return nil
}

func (db *applicationUserConnection) Insert(record models.ApplicationUser) (models.ApplicationUser, error) {
	record.Id = uuid.New().String()
	record.IsActive = true
	record.CreatedBy = record.Id
	record.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	record.Password = helper.HashAndSalt([]byte(record.Password))
	res := db.connection.Create(&record)

	if res.RowsAffected == 0 {
		return record, fmt.Errorf("%v", res.Error)
	} else {
		return record, nil
	}
}

func (db *applicationUserConnection) Update(record models.ApplicationUser) models.ApplicationUser {
	if record.Password != "" {
		record.Password = helper.HashAndSalt([]byte(record.Password))
	} else {
		var tempUser models.ApplicationUser
		res := db.GetById(record.Id)
		tempUser = (res.Data).(models.ApplicationUser)
		record.Password = tempUser.Password
	}

	db.connection.Save(&record)
	return record
}

func (db *applicationUserConnection) UpdateProfile(dtoRecord dto.ApplicationUserDescriptionDto) helper.Response {
	commons.Logger()
	tx := db.connection.Begin()

	var record models.ApplicationUser
	record = db.GetById(dtoRecord.Id).Data.(models.ApplicationUser)

	if dtoRecord.Title != "" {
		record.Title = dtoRecord.Title
	}

	if dtoRecord.Description != "" {
		record.Description = dtoRecord.Description
	}

	record.IsActive = true
	record.UpdatedBy = dtoRecord.UpdatedBy
	record.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}

	res := tx.Save(&record)
	if res.RowsAffected == 0 {
		return helper.ServerResponse(false, fmt.Sprintf("%v,", res.Error), fmt.Sprintf("%v,", res.Error), helper.EmptyObj{})
	}

	tx.Commit()
	return helper.ServerResponse(true, "Ok", "", record)
}

func (db *applicationUserConnection) UpdateProfilePicture(request request_models.FileRequestDto) helper.Response {
	tx := db.connection.Begin()
	var record models.ApplicationUser
	db.connection.First(&record, "id = ?", request.Id)
	if request.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}

	record.Filepath = request.Filepath
	record.FilepathThumbnail = request.FilepathThumbnail
	record.Filename = request.Filename
	record.Extension = request.Extension
	record.Size = request.Size
	record.UpdatedBy = request.UpdatedBy
	record.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}

	res := tx.Save(&record)
	if res.RowsAffected == 0 {
		return helper.ServerResponse(false, fmt.Sprintf("%v,", res.Error), fmt.Sprintf("%v,", res.Error), helper.EmptyObj{})
	}

	tx.Commit()
	db.connection.Preload(clause.Associations).Find(&record)
	return helper.ServerResponse(true, "Ok", "", record)
}

func (db *applicationUserConnection) RecoverPassword(recordId string, newPassword string) helper.Response {
	tx := db.connection.Begin()
	var user models.ApplicationUser

	db.connection.First(&user, "id = ?", recordId)
	user.Password = helper.HashAndSalt([]byte(newPassword))

	res := tx.Save(&user)
	if res.RowsAffected == 0 {
		return helper.ServerResponse(false, fmt.Sprintf("%v,", res.Error), fmt.Sprintf("%v,", res.Error), helper.EmptyObj{})
	}

	tx.Commit()
	return helper.ServerResponse(true, "Ok", "", user)
}

func (db *applicationUserConnection) VerifyCredential(username string, email string, password string) interface{} {
	var record models.ApplicationUser
	res := db.connection.Where("username = ?", username).Or("email = ?", email).Take(&record)
	if res.Error == nil {
		return record
	}
	return nil
}

func (db *applicationUserConnection) IsDuplicateEmail(email string) (tx *gorm.DB) {
	var record models.ApplicationUser
	return db.connection.Where("email = ?", email).Take(&record)
}

func (db *applicationUserConnection) GetByEmail(email string) models.ApplicationUser {
	var record models.ApplicationUser
	db.connection.Where("email = ?", email).Take(&record)
	return record
}

func (db *applicationUserConnection) UserProfile(recordId string) models.ApplicationUser {
	var record models.ApplicationUser
	db.connection.Find(&record, recordId)
	return record
}

func (db *applicationUserConnection) GetAll() []models.ApplicationUser {
	var users []models.ApplicationUser
	db.connection.Where("entity_id <> ?", "").Find(&users)
	return users
}

func (db *applicationUserConnection) DeleteById(applicationUserId string) helper.Response {
	var record models.ApplicationUser
	db.connection.First(&record, "id = ?", applicationUserId)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	} else {
		db.connection.Where("id = ?", applicationUserId).Delete(&record)
		res := helper.ServerResponse(true, "Ok", "", helper.EmptyObj{})
		return res
	}
}

func (db *applicationUserConnection) GetById(applicationUserId string) helper.Response {
	var record models.ApplicationUser
	db.connection.First(&record, "id = ?", applicationUserId)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	} else {
		res := helper.ServerResponse(true, "Ok", "", record)
		return res
	}
}

func (db *applicationUserConnection) GetViewById(applicationUserId string) helper.Response {
	var record entity_view_models.EntityApplicationUserView
	db.connection.First(&record, "id = ?", applicationUserId)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	} else {
		res := helper.ServerResponse(true, "Ok", "", record)
		return res
	}
}

func (db *applicationUserConnection) EmailVerificationById(userId string) helper.Response {
	var result = db.GetById(userId)
	if !result.Status {
		return result
	}
	var record = result.Data.(models.ApplicationUser)
	record.IsEmailVerified = true

	db.connection.Model(&record).Updates(models.ApplicationUser{IsEmailVerified: true})
	return helper.ServerResponse(true, "Ok", "", record)
}
