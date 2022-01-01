package entity_view_models

import (
	"joranvest/models"
	"strings"
)

type EntityRatingMasterView struct {
	models.RatingMaster
	RaterFirstName   string `json:"rater_first_name"`
	RaterLastName    string `json:"rater_last_name"`
	RaterFullName    string `json:"rater_full_name"`
	RaterInitialName string `json:"rater_initial_name"`
	UserCreate       string `json:"user_create"`
	UserUpdate       string `json:"user_update"`
}

func (EntityRatingMasterView) TableName() string {
	return "vw_rating_master"
}

func (EntityRatingMasterView) ViewModel() string {
	var sql strings.Builder
	sql.WriteString("SELECT")
	sql.WriteString("  r.id,")
	sql.WriteString("  r.is_active,")
	sql.WriteString("  r.is_locked,")
	sql.WriteString("  r.is_default,")
	sql.WriteString("  r.created_at,")
	sql.WriteString("  r.created_by,")
	sql.WriteString("  r.updated_at,")
	sql.WriteString("  r.updated_by,")
	sql.WriteString("  r.approved_at,")
	sql.WriteString("  r.approved_by,")
	sql.WriteString("  r.submitted_at,")
	sql.WriteString("  r.submitted_by,")
	sql.WriteString("  r.entity_id,")
	sql.WriteString("  r.user_id,")
	sql.WriteString("  r.object_rated_id,")
	sql.WriteString("  r.reference_id,")
	sql.WriteString("  u3.first_name AS rater_first_name,")
	sql.WriteString("  u3.last_name AS rater_last_name,")
	sql.WriteString("  CONCAT(u3.first_name, ' ', u3.last_name) AS rater_full_name,")
	sql.WriteString("  CONCAT(UPPER(LEFT(u3.first_name, 1)), '', UPPER(LEFT(u3.last_name, 1))) AS rater_initial_name,")
	sql.WriteString("  r.rating,")
	sql.WriteString("  r.comment,")
	sql.WriteString("  CONCAT(u1.first_name, ' ', u1.last_name) AS user_create,")
	sql.WriteString("  CONCAT(u2.first_name, ' ', u2.last_name) AS user_update ")
	sql.WriteString(" FROM rating_master r ")
	sql.WriteString("  LEFT JOIN application_user u3 ON u3.id = r.user_id")
	sql.WriteString("  LEFT JOIN application_user u1 ON u1.id = r.created_by")
	sql.WriteString("  LEFT JOIN application_user u2 ON u2.id = r.updated_by")
	return sql.String()
}
func (EntityRatingMasterView) Migration() map[string]string {
	var view = EntityRatingMasterView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
