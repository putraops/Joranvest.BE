package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/entity_view_models"
	"joranvest/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleMemberService interface {
	GetPagination(request commons.Pagination2ndRequest) interface{}
	GetAll(filter map[string]interface{}) []entity_view_models.EntityRoleMemberView
	Save(record models.RoleMember, ctx *gin.Context) helper.Result
	GetViewById(recordId string) helper.Result
	GetById(recordId string) helper.Result
	GetByRoleId(roleId string) []entity_view_models.EntityRoleMemberView
	DeleteById(recordId string) helper.Result
	GetUsersInRole(roleId string) []entity_view_models.EntityRoleMemberView
	GetUsersNotInRole(roleId string, search string) []entity_view_models.EntityApplicationUserView

	OpenTransaction(trxHandle *gorm.DB) roleMemberService
}

type roleMemberService struct {
	DB         *gorm.DB
	jwtService JWTService
	// webinarRepository    repository.WebinarRepository
	roleMemberRepository repository.RoleMemberRepository
}

func NewRoleMemberService(db *gorm.DB, jwtService JWTService) RoleMemberService {
	return roleMemberService{
		DB:         db,
		jwtService: jwtService,
		// webinarRepository:    repository.NewWebinarRepository(db),
		roleMemberRepository: repository.NewRoleMemberRepository(db),
	}
}

func (r roleMemberService) GetPagination(request commons.Pagination2ndRequest) interface{} {
	return r.roleMemberRepository.GetPagination(request)
}

func (r roleMemberService) GetAll(filter map[string]interface{}) []entity_view_models.EntityRoleMemberView {
	return r.roleMemberRepository.GetViewAll(filter)
}

func (r roleMemberService) GetByRoleId(roleId string) []entity_view_models.EntityRoleMemberView {
	filter := make(map[string]interface{})
	filter["role_id"] = roleId

	return r.roleMemberRepository.GetViewAll(filter)
}

func (r roleMemberService) Save(record models.RoleMember, ctx *gin.Context) helper.Result {
	authHeader := ctx.GetHeader("Authorization")
	userIdentity := r.jwtService.GetUserByToken(authHeader)
	record.EntityId = &userIdentity.EntityId

	if record.Id == nil {
		record.CreatedBy = &userIdentity.UserId
		record.OwnerId = &userIdentity.UserId
		return r.roleMemberRepository.Insert(record)
	} else {
		record.UpdatedBy = &userIdentity.UserId
		return r.roleMemberRepository.Update(record)
	}
}

func (r roleMemberService) GetById(recordId string) helper.Result {
	return r.roleMemberRepository.GetById(recordId)
}

func (r roleMemberService) GetViewById(recordId string) helper.Result {
	return r.roleMemberRepository.GetViewById(recordId)
}

func (r roleMemberService) DeleteById(recordId string) helper.Result {
	return r.roleMemberRepository.DeleteById(recordId)
}

func (r roleMemberService) GetUsersInRole(roleId string) []entity_view_models.EntityRoleMemberView {
	return r.roleMemberRepository.GetUsersInRole(roleId)
}

func (r roleMemberService) GetUsersNotInRole(roleId string, search string) []entity_view_models.EntityApplicationUserView {
	return r.roleMemberRepository.GetUsersNotInRole(roleId, search)
}

func (r roleMemberService) OpenTransaction(trxHandle *gorm.DB) roleMemberService {
	r.roleMemberRepository = r.roleMemberRepository.OpenTransaction(trxHandle)
	return r
}
