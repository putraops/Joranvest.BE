package view_models

type EducationPlaylistByUserViewModel struct {
	Id                      string `json:"id"`
	IsActive                bool   `json:"is_active"`
	EducationId             string `json:"education_id"`
	EducationTitle          string `json:"education_title"`
	Title                   string `json:"title"`
	FileUrl                 string `json:"file_url"`
	Description             string `json:"description"`
	OrderIndex              int32  `json:"order_index"`
	EducationPlaylistUserId string `json:"education_playlist_user_id"`
	IsWatched               bool   `json:"is_watched"`
}
