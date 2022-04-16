package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductService interface {
	GetPagination(request commons.Pagination2ndRequest) interface{}
	// Lookup(request helper.ReactSelectRequest) helper.Response
	GetAll(filter map[string]interface{}) []models.Product
	Save(record models.Product, context *gin.Context) helper.Result
	GetViewById(recordId string) helper.Result
	GetById(recordId string) helper.Result
	GetProductByRecordId(recordId string) helper.Result
	GetByProductType(product_type string) helper.Result
	DeleteById(recordId string) helper.Result

	OpenTransaction(*gorm.DB) productService
}

type productService struct {
	DB                *gorm.DB
	jwtService        JWTService
	productRepository repository.ProductRepository
}

func NewProductService(db *gorm.DB, jwtService JWTService) ProductService {
	return productService{
		DB:                db,
		jwtService:        jwtService,
		productRepository: repository.NewProductRepository(db),
	}
}

func (r productService) GetPagination(request commons.Pagination2ndRequest) interface{} {
	return r.productRepository.GetPagination(request)
}

func (r productService) GetAll(filter map[string]interface{}) []models.Product {
	return r.productRepository.GetAll(filter)
}

func (r productService) Save(record models.Product, context *gin.Context) helper.Result {
	authHeader := context.GetHeader("Authorization")
	userIdentity := r.jwtService.GetUserByToken(authHeader)

	if record.Id == nil {
		record.CreatedBy = &userIdentity.UserId
		return r.productRepository.Insert(record)
	} else {
		record.UpdatedBy = &userIdentity.UserId
		return r.productRepository.Update(record)
	}
}

func (r productService) GetById(recordId string) helper.Result {
	return r.productRepository.GetById(recordId)
}

func (r productService) GetByProductType(product_type string) helper.Result {
	return r.productRepository.GetByProductType(product_type)
}

func (r productService) GetViewById(recordId string) helper.Result {
	return r.productRepository.GetViewById(recordId)
}

//-- To get product (Product, Webinar or Membership)
func (r productService) GetProductByRecordId(recordId string) helper.Result {
	return r.productRepository.GetProductByRecordId(recordId)
}

func (r productService) DeleteById(recordId string) helper.Result {
	return r.productRepository.DeleteById(recordId)
}

func (r productService) OpenTransaction(trxHandle *gorm.DB) productService {
	r.productRepository = r.productRepository.OpenTransaction(trxHandle)
	return r
}
