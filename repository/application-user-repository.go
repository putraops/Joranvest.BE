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
	GetDatatables(request commons.DataTableRequest) commons.DataTableResponse
	Lookup(req map[string]interface{}) []models.ApplicationUser
	GetViewUserByEmail(username string, email string) interface{}
	GetViewUserByUsernameOrEmail(username string, email string) interface{}
	Insert(t models.ApplicationUser) (models.ApplicationUser, error)
	Update(record models.ApplicationUser) models.ApplicationUser
	UpdateProfile(dtoRecord dto.ApplicationUserDescriptionDto) helper.Response
	UpdateProfilePicture(request request_models.FileRequestDto) helper.Response
	VerifyCredential(username string, email string, password string) interface{}
	IsDuplicateEmail(email string) (tx *gorm.DB)
	GetByEmail(email string) models.ApplicationUser
	RecoverPassword(recordId string, oldPassword string) helper.Response
	EmailVerificationById(userId string) helper.Response
	UserProfile(applicationUserId string) models.ApplicationUser
	GetById(applicationUserId string) helper.Response
	GetViewById(applicationUserId string) helper.Response
	GetAll() []models.ApplicationUser
	DeleteById(recordId string) helper.Response
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

func (db *applicationUserConnection) GetDatatables(request commons.DataTableRequest) commons.DataTableResponse {
	var records []entity_view_models.EntityApplicationUserView
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
		res.DataRow = []entity_view_models.EntityMembershipView{}
	}
	return res
}

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
						sqlQuery.WriteString(" AND (LOWER(" + valueDetail + ") LIKE " + fmt.Sprint("'%", v.Value, "%'"))
					} else {
						sqlQuery.WriteString(" OR LOWER(" + valueDetail + ") LIKE " + fmt.Sprint("'%", v.Value, "%'"))
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
	db.connection.Raw(sqlQuery.String()).Scan(&records)
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
	record.IsEmailVerified = true
	record.CreatedBy = record.Id
	record.UpdatedAt = sql.NullTime{Time: time.Now().Local().UTC(), Valid: true}
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
	record.UpdatedAt = sql.NullTime{Time: time.Now().Local().UTC(), Valid: true}

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
	record.UpdatedAt = sql.NullTime{Time: time.Now().Local().UTC(), Valid: true}

	res := tx.Save(&record)
	if res.RowsAffected == 0 {
		return helper.ServerResponse(false, fmt.Sprintf("%v,", res.Error), fmt.Sprintf("%v,", res.Error), helper.EmptyObj{})
	}

	tx.Commit()
	db.connection.Preload(clause.Associations).Find(&record)
	return helper.ServerResponse(true, "Ok", "", record)
}

func (db *applicationUserConnection) RecoverPassword(recordId string, oldPassword string) helper.Response {
	tx := db.connection.Begin()
	var user models.ApplicationUser
	db.connection.Find(&user, recordId)
	user.Password = helper.HashAndSalt([]byte(oldPassword))

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
