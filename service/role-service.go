package service

import (
	"joranvest/commons"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"

	"github.com/mashingan/smapping"
	"gorm.io/gorm"
)

type RoleService interface {
	GetPagination(request commons.Pagination2ndRequest) interface{}
	GetAll(filter map[string]interface{}) []models.Role
	Insert(dto dto.RoleDto) helper.Result
	Update(record models.Role) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response

	OpenTransaction(*gorm.DB) roleService
}

type roleService struct {
	DB                 *gorm.DB
	jwtService         JWTService
	roleRepository     repository.RoleRepository
	teamRepository     repository.TeamRepository
	teamRoleRepository repository.TeamRoleRepository
}

func NewRoleService(db *gorm.DB, jwtService JWTService) RoleService {
	return roleService{
		DB:                 db,
		jwtService:         jwtService,
		roleRepository:     repository.NewRoleRepository(db),
		teamRepository:     repository.NewTeamRepository(db),
		teamRoleRepository: repository.NewTeamRoleRepository(db),
	}
}

func (r roleService) GetPagination(request commons.Pagination2ndRequest) interface{} {
	return r.roleRepository.GetPagination(request)
}

func (r roleService) GetAll(filter map[string]interface{}) []models.Role {
	return r.roleRepository.GetAll(filter)
}

func (r roleService) Insert(dto dto.RoleDto) helper.Result {
	authHeader := dto.Context.GetHeader("Authorization")
	userIdentity := r.jwtService.GetUserByToken(authHeader)

	//-- Map Dto to Role
	var roleRecord models.Role
	smapping.FillStruct(&roleRecord, smapping.MapFields(&dto))
	roleRecord.CreatedBy = userIdentity.UserId
	roleRecord.EntityId = userIdentity.UserId

	return r.roleRepository.Insert(roleRecord)
}

func (r roleService) Update(record models.Role) helper.Response {
	return r.roleRepository.Update(record)
}

func (r roleService) GetById(recordId string) helper.Response {
	return r.roleRepository.GetById(recordId)
}

func (r roleService) DeleteById(recordId string) helper.Response {
	return r.roleRepository.DeleteById(recordId)
}

func (r roleService) OpenTransaction(trxHandle *gorm.DB) roleService {
	r.roleRepository = r.roleRepository.OpenTransaction(trxHandle)
	return r
}
