package repository

import (
	"database/sql"
	"fmt"
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/entity_view_models"
	"joranvest/models/view_models"
	"strings"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WebinarSpeakerRepository interface {
	GetById(recordId string) helper.Response
	GetAll(filter map[string]interface{}) []entity_view_models.EntityWebinarSpeakerView
	GetSpeakerReviewById(recordId string) helper.Response
	GetSpeakersRatingByWebinarId(webinarId string) helper.Response
	Insert(records []models.WebinarSpeaker, speakerType int) helper.Response
}

type webinarSpeakerConnection struct {
	connection        *gorm.DB
	serviceRepository ServiceRepository
	tableName         string
	viewQuery         string
}

func NewWebinarSpeakerRepository(db *gorm.DB) WebinarSpeakerRepository {
	return &webinarSpeakerConnection{
		connection:        db,
		tableName:         models.WebinarSpeaker.TableName(models.WebinarSpeaker{}),
		viewQuery:         entity_view_models.EntityWebinarSpeakerView.ViewModel(entity_view_models.EntityWebinarSpeakerView{}),
		serviceRepository: NewServiceRepository(db),
	}
}

func (db *webinarSpeakerConnection) GetSpeakersRatingByWebinarId(webinarId string) helper.Response {
	commons.Logger()
	var records []view_models.WebinarSpeakerRatingViewModel

	var sql strings.Builder
	sql.WriteString("SELECT")
	sql.WriteString("	r.rating_master_id AS id,")
	sql.WriteString("	r.created_at,")
	sql.WriteString("	r.speaker_id,")
	sql.WriteString("	r.created_by AS user_id,")
	sql.WriteString("	r.object_rated_id,")
	sql.WriteString("	r.reference_id,")
	sql.WriteString("	r.rating,")
	sql.WriteString("	r.comment,")
	sql.WriteString("	o.name AS organization_name,")
	sql.WriteString("	concat(u3.first_name, ' ', u3.last_name) AS speaker_fullname ")
	sql.WriteString("FROM (")
	sql.WriteString("	SELECT NULL AS rating_master_id, null AS object_rated_id, null AS reference_id, 0 AS rating, '' AS comment, r.speaker_id, r.created_at, r.created_by")
	sql.WriteString("	FROM webinar_speaker r")
	sql.WriteString("	WHERE 1=1")
	sql.WriteString(fmt.Sprintf(" AND r.webinar_id = '%v'", webinarId))
	sql.WriteString("	AND r.speaker_id NOT IN (")
	sql.WriteString("		SELECT s.object_rated_id FROM rating_master s WHERE s.object_rated_id != s.reference_id AND s.reference_id = r.webinar_id")
	sql.WriteString("	)")
	sql.WriteString("	UNION")
	sql.WriteString("	SELECT s.id AS rating_master_id, s.object_rated_id, s.reference_id, s.rating, s.comment, r.speaker_id, r.created_at, r.created_by")
	sql.WriteString("	FROM webinar_speaker r")
	sql.WriteString("	INNER JOIN rating_master s ON s.reference_id = r.webinar_id AND s.object_rated_id = r.speaker_id AND s.object_rated_id != s.reference_id ")
	sql.WriteString("WHERE 1=1")
	sql.WriteString(fmt.Sprintf(" AND r.webinar_id = '%v'", webinarId))
	sql.WriteString(") AS r ")
	sql.WriteString("LEFT JOIN application_user u3 ON u3.id = r.speaker_id ")
	sql.WriteString("LEFT JOIN organization o ON o.id = r.speaker_id ")
	sql.WriteString("ORDER BY r.created_at ASC ")

	fmt.Println(sql.String())

	result := db.connection.Raw(sql.String()).Find(&records)
	if result.Error != nil {
		log.Error(fmt.Sprintf("%v,", result.Error))
		log.Error(db.serviceRepository.getCurrentFuncName())
		return helper.ServerResponse(false, fmt.Sprintf("%v,", result.Error), fmt.Sprintf("%v,", result.Error), helper.EmptyObj{})
	}

	return helper.ServerResponse(true, "Ok", "", records)
}

func (db *webinarSpeakerConnection) Insert(records []models.WebinarSpeaker, speakerType int) helper.Response {
	tx := db.connection.Begin()

	if len(records) > 0 {
		if err := tx.Where("webinar_id = ?", records[0].WebinarId).Delete(models.WebinarSpeaker{}).Error; err != nil {
			tx.Rollback()
			return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
		}

		for i := 0; i < len(records); i++ {
			records[i].Id = uuid.New().String()
			records[i].CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
		}

		if err := tx.Create(&records).Error; err != nil {
			tx.Rollback()
			return helper.ServerResponse(false, fmt.Sprintf("%v,", err), fmt.Sprintf("%v,", err), helper.EmptyObj{})
		}
	} else {
		return helper.ServerResponse(false, "Speaker cannot be empty", "Speaker cannot be empty", helper.EmptyObj{})
	}

	var webinarRecord models.Webinar
	db.connection.First(&webinarRecord, "id = ?", records[0].WebinarId)
	if webinarRecord.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	webinarRecord.SpeakerType = speakerType
	res := tx.Save(&webinarRecord)
	if res.RowsAffected == 0 {
		return helper.ServerResponse(false, fmt.Sprintf("%v,", res.Error), fmt.Sprintf("%v,", res.Error), helper.EmptyObj{})
	}

	tx.Commit()
	return helper.ServerResponse(true, "Ok", "", records)
}

func (db *webinarSpeakerConnection) GetById(recordId string) helper.Response {
	var record models.FundamentalAnalysis
	db.connection.Preload("Emiten").First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	}
	res := helper.ServerResponse(true, "Ok", "", record)
	return res
}

func (db *webinarSpeakerConnection) GetSpeakerReviewById(recordId string) helper.Response {
	var record view_models.WebinarSpeakerReviewViewModel
	db.connection.First(&record, "id = ?", recordId)
	if record.Id == "" {
		res := helper.ServerResponse(false, "Record not found", "Error", helper.EmptyObj{})
		return res
	} else {
		res := helper.ServerResponse(true, "Ok", "", record)
		return res
	}
}

func (db *webinarSpeakerConnection) GetAll(filter map[string]interface{}) []entity_view_models.EntityWebinarSpeakerView {
	var records []entity_view_models.EntityWebinarSpeakerView
	if len(filter) == 0 {
		db.connection.Find(&records)
	} else if len(filter) != 0 {
		db.connection.Where(filter).Find(&records)
	}
	return records
}
