package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleService interface {
	GetPagination(request commons.Pagination2ndRequest) interface{}
	GetAll(filter map[string]interface{}) []models.Role
	Save(record models.Role, ctx *gin.Context) helper.Result
	GetViewById(recordId string) helper.Result
	GetById(recordId string) helper.Result
	DeleteById(recordId string) helper.Result

	OpenTransaction(trxHandle *gorm.DB) roleService
}

type roleService struct {
	DB             *gorm.DB
	jwtService     JWTService
	roleRepository repository.RoleRepository
}

func NewRoleService(db *gorm.DB, jwtService JWTService) RoleService {
	return roleService{
		DB:             db,
		jwtService:     jwtService,
		roleRepository: repository.NewRoleRepository(db),
	}
}

func (r roleService) GetPagination(request commons.Pagination2ndRequest) interface{} {
	return r.roleRepository.GetPagination(request)
}

func (r roleService) GetAll(filter map[string]interface{}) []models.Role {
	return r.roleRepository.GetAll(filter)
}

func (r roleService) Save(record models.Role, ctx *gin.Context) helper.Result {
	authHeader := ctx.GetHeader("Authorization")
	userIdentity := r.jwtService.GetUserByToken(authHeader)
	record.EntityId = &userIdentity.EntityId

	if record.Id == nil {
		record.CreatedBy = &userIdentity.UserId
		record.OwnerId = &userIdentity.UserId
		return r.roleRepository.Insert(record)
	} else {
		record.UpdatedBy = &userIdentity.UserId
		return r.roleRepository.Update(record)
	}
}

func (r roleService) GetById(recordId string) helper.Result {
	return r.roleRepository.GetById(recordId)
}

func (r roleService) GetViewById(recordId string) helper.Result {
	return r.roleRepository.GetViewById(recordId)
}

func (r roleService) DeleteById(recordId string) helper.Result {
	return r.roleRepository.DeleteById(recordId)
}

func (r roleService) OpenTransaction(trxHandle *gorm.DB) roleService {
	r.roleRepository = r.roleRepository.OpenTransaction(trxHandle)
	return r
}
