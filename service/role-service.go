package service

import (
	"joranvest/commons"
	"joranvest/dto"
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

	//-- Configuration
	SetDashboardAccess(dto dto.RoleAccessDto, ctx *gin.Context) helper.Result
	SetFullAccess(dto dto.RoleAccessDto, ctx *gin.Context) helper.Result
	SetPaymentNotification(dto dto.PaymentNotificationDto, ctx *gin.Context) helper.Result
	GetNotificationConfigurationByRoleId(roleId string) helper.Result

	OpenTransaction(trxHandle *gorm.DB) roleService
}

type roleService struct {
	DB                         *gorm.DB
	jwtService                 JWTService
	roleRepository             repository.RoleRepository
	roleNotificationRepository repository.RoleNotificationRepository
	serviceRepository          repository.ServiceRepository
}

func NewRoleService(db *gorm.DB, jwtService JWTService) RoleService {
	return roleService{
		DB:                         db,
		jwtService:                 jwtService,
		roleRepository:             repository.NewRoleRepository(db),
		roleNotificationRepository: repository.NewRoleNotificationRepository(db),
		serviceRepository:          repository.NewServiceRepository(db),
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

func (r roleService) SetFullAccess(dto dto.RoleAccessDto, ctx *gin.Context) helper.Result {
	authHeader := ctx.GetHeader("Authorization")
	userIdentity := r.jwtService.GetUserByToken(authHeader)

	var record models.Role
	result := r.roleRepository.GetById(dto.RoleId)
	if !result.Status {
		return result
	}

	record = result.Data.(models.Role)
	record.UpdatedBy = &userIdentity.UserId
	record.HasFullAccess = &dto.IsChecked
	return r.roleRepository.Update(record)
}

func (r roleService) SetDashboardAccess(dto dto.RoleAccessDto, ctx *gin.Context) helper.Result {
	authHeader := ctx.GetHeader("Authorization")
	userIdentity := r.jwtService.GetUserByToken(authHeader)

	var record models.Role
	result := r.roleRepository.GetById(dto.RoleId)
	if !result.Status {
		return result
	}

	record = result.Data.(models.Role)
	record.UpdatedBy = &userIdentity.UserId
	record.HasDashboardAccess = &dto.IsChecked
	return r.roleRepository.Update(record)
}

func (r roleService) SetPaymentNotification(dto dto.PaymentNotificationDto, ctx *gin.Context) helper.Result {
	notificationType := "Email"
	authHeader := ctx.GetHeader("Authorization")
	userIdentity := r.jwtService.GetUserByToken(authHeader)

	var record models.RoleNotification
	result := r.roleNotificationRepository.GetRoleById(dto.RoleId)
	if result.Status {
		record = result.Data.(models.RoleNotification)
	}

	if record.Id == nil {
		record.CreatedBy = &userIdentity.UserId
	} else {
		record.UpdatedBy = &userIdentity.UserId
	}

	record.RoleId = &dto.RoleId
	record.PaymentNotificationType = &notificationType
	if !dto.IsChecked {
		record.PaymentNotificationType = nil
	}
	record.HasPaymentNotification = &dto.IsChecked

	return r.roleNotificationRepository.SetNotification(record)
}

func (r roleService) GetNotificationConfigurationByRoleId(roleId string) helper.Result {
	return r.roleNotificationRepository.GetRoleById(roleId)
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
