package service

import (
	"joranvest/commons"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
)

type PaymentService interface {
	GetPagination(request commons.Pagination2ndRequest) interface{}
	GetAll(filter map[string]interface{}) []models.Payment
	GetUniqueNumber() int
	MembershipPayment(record models.Payment) helper.Response
	WebinarPayment(record models.Payment) helper.Response
	Update(record models.Payment) helper.Response
	UpdatePaymentStatus(req dto.UpdatePaymentStatusDto) helper.Response
	GetById(recordId string) helper.Response
	DeleteById(recordId string) helper.Response
}

type paymentService struct {
	paymentRepository repository.PaymentRepository
	helper.AppSession
}

func NewPaymentService(repo repository.PaymentRepository) PaymentService {
	return &paymentService{
		paymentRepository: repo,
	}
}

func (service *paymentService) GetPagination(request commons.Pagination2ndRequest) interface{} {
	return service.paymentRepository.GetPagination(request)
}

func (service *paymentService) GetAll(filter map[string]interface{}) []models.Payment {
	return service.paymentRepository.GetAll(filter)
}

func (service *paymentService) GetUniqueNumber() int {
	return service.paymentRepository.GetUniqueNumber()
}

func (service *paymentService) MembershipPayment(record models.Payment) helper.Response {
	return service.paymentRepository.MembershipPayment(record)
}
func (service *paymentService) WebinarPayment(record models.Payment) helper.Response {
	return service.paymentRepository.WebinarPayment(record)
}

func (service *paymentService) Update(record models.Payment) helper.Response {
	return service.paymentRepository.Update(record)
}

func (service *paymentService) UpdatePaymentStatus(req dto.UpdatePaymentStatusDto) helper.Response {
	return service.paymentRepository.UpdatePaymentStatus(req)
}

func (service *paymentService) GetById(recordId string) helper.Response {
	return service.paymentRepository.GetById(recordId)
}

func (service *paymentService) DeleteById(recordId string) helper.Response {
	return service.paymentRepository.DeleteById(recordId)
}
