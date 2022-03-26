package controllers

import (
	"net/http"

	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
	"joranvest/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductController interface {
	GetPagination(context *gin.Context)
	// Lookup(context *gin.Context)
	Save(context *gin.Context)
	GetById(context *gin.Context)
	GetProductByRecordId(context *gin.Context)
	GetByProductType(context *gin.Context)
	GetViewById(context *gin.Context)
	DeleteById(context *gin.Context)
}

type productController struct {
	productService    service.ProductService
	productRepository repository.ProductRepository
	jwtService        service.JWTService
	db                *gorm.DB
}

func NewProductController(db *gorm.DB, jwtService service.JWTService) ProductController {
	return &productController{
		db:                db,
		jwtService:        jwtService,
		productService:    service.NewProductService(db, jwtService),
		productRepository: repository.NewProductRepository(db),
	}
}

// @Tags         Product
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        body body commons.Pagination2ndRequest true "body"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /product/getPagination [post]
func (c productController) GetPagination(context *gin.Context) {
	var req commons.Pagination2ndRequest
	errDTO := context.Bind(&req)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.productService.GetPagination(req)
	context.JSON(http.StatusOK, result)
}

// // @Tags         Product
// // @Security 	 ApiKeyAuth
// // @Accept       json
// // @Produce      json
// // @Param        body body helper.ReactSelectRequest true "body"
// // @Param        q query string false "id"
// // @Success      200 {object} object
// // @Failure 	 400,404 {object} object
// // @Router       /product/lookup [post]
// func (c productController) Lookup(context *gin.Context) {
// 	var request helper.ReactSelectRequest

// 	errDTO := context.Bind(&request)
// 	if errDTO != nil {
// 		res := helper.StandartResult(false, errDTO.Error(), helper.EmptyObj{})
// 		context.JSON(http.StatusBadRequest, res)
// 		return
// 	}

// 	var result = c.productService.Lookup(request)
// 	response := helper.StandartResult(true, "Ok", result.Data)
// 	context.JSON(http.StatusOK, response)
// }

// @Tags         Product
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        body body models.Product true "record"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /product/save [post]
func (r productController) Save(c *gin.Context) {
	var result helper.Result
	var record models.Product
	//dto.Context = c

	errDto := c.Bind(&record)
	if errDto != nil {
		res := helper.StandartResult(false, errDto.Error(), helper.EmptyObj{})
		c.JSON(http.StatusBadRequest, res)
		return
	}

	result = r.productService.Save(record, c)
	c.JSON(http.StatusOK, helper.StandartResult(result.Status, result.Message, result.Data))
	return
}

// @Tags         Product
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "id"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /product/getById/{id} [get]
func (c productController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
		return
	}
	result := c.productService.GetById(id)
	context.JSON(http.StatusOK, result)
}

// @Tags         Product
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        record_id path string true "record_id"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /product/getProductByRecordId/{record_id} [get]
func (c productController) GetProductByRecordId(context *gin.Context) {
	record_id := context.Param("record_id")
	if record_id == "" {
		response := helper.BuildErrorResponse("Failed to get record_id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
		return
	}
	result := c.productService.GetProductByRecordId(record_id)
	context.JSON(http.StatusOK, result)
}

// @Tags         Product
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        product_type path string true "product_type"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /product/getByProductType/{product_type} [get]
func (c productController) GetByProductType(context *gin.Context) {
	product_type := context.Param("product_type")
	if product_type == "" {
		response := helper.BuildErrorResponse("Failed to get product_type", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
		return
	}
	result := c.productService.GetByProductType(product_type)
	context.JSON(http.StatusOK, result)
}

// @Tags         Product
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "id"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /product/getViewById/{id} [get]
func (c productController) GetViewById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
		return
	}
	result := c.productService.GetViewById(id)
	context.JSON(http.StatusOK, result)
}

// @Tags         Product
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "id"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /product/deleteById/{id} [delete]
func (c productController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
		return
	}
	result := c.productService.DeleteById(id)
	context.JSON(http.StatusOK, result)
}
