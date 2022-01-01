package request_models

type PatchingEmiten struct {
	EmitenName     string `json:"emiten_name" form:"emiten_name" binding:"required"`
	EmitenCode     string `json:"emiten_code" form:"emiten_code" binding:"required"`
	EmitenCategory string `json:"emiten_category" form:"emiten_category" binding:"required"`
	EmitenSector   string `json:"emiten_sector" form:"emiten_sector" binding:"required"`
}
