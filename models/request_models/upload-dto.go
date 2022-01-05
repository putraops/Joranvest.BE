package request_models

type FileRequestDto struct {
	Id                string `json:"id" form:"id"`
	Filepath          string `json:"filepath" form:"payment_id"`
	FilepathThumbnail string `json:"filepath_thumbnail" form:"payment_id"`
	Filename          string `json:"filename" form:"payment_id"`
	FileType          int    `json:"file_type" form:"payment_id"`
	Extension         string `json:"extension" form:"payment_id"`
	Size              string `json:"size"`

	EntityId  string `json:"-"`
	UpdatedBy string
}
