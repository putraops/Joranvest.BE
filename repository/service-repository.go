package repository

import (
	"fmt"
	"joranvest/helper"
	"joranvest/models"
	"strings"

	"gorm.io/gorm"
)

type ServiceRepository interface {
	GetAll(filter map[string]interface{}) interface{}
	Lookup(req map[string]interface{}, r helper.Select2Request) []models.ApplicationUser
	ConvertViewQueryIntoViewCount(param string) string
	ConvertViewQueryIntoViewCountByPublic(param string, tableName string) string
}

type serviceConnection struct {
	connection *gorm.DB
}

func NewServiceRepository(db *gorm.DB) ServiceRepository {
	return &serviceConnection{
		connection: db,
	}
}

func (db *serviceConnection) ConvertViewQueryIntoViewCount(request string) string {
	var response = ""
	strParts := strings.Split(strings.ToLower(request), "from")
	response = "SELECT COUNT(*) FROM " + strParts[1] + " "
	return response
}

func (db *serviceConnection) ConvertViewQueryIntoViewCountByPublic(request string, tableName string) string {
	var response = ""
	strParts := strings.Split(strings.ToLower(request), "from public."+tableName+" r")
	response = "SELECT COUNT(*) FROM " + tableName + " AS r " + strParts[1] + " "
	return response
}

func (db *serviceConnection) GetAll(filter map[string]interface{}) interface{} {
	var sqlQuery strings.Builder
	var res interface{}
	sqlQuery.WriteString("SELECT id FROM " + fmt.Sprint(filter["table_name"]))
	fmt.Println("Query: ", sqlQuery.String())
	return db.connection.Raw(sqlQuery.String()).Scan(&res)
}

func (db *serviceConnection) Lookup(req map[string]interface{}, r helper.Select2Request) []models.ApplicationUser {
	var records []models.ApplicationUser
	db.connection.Order("first_name asc")

	var sqlQuery strings.Builder

	v, found := req["condition"]
	sqlQuery.WriteString("SELECT * FROM customer")
	if found {
		sqlQuery.WriteString(" WHERE 1 = 1")
		requests := v.(helper.DataFilter).Request
		for _, v := range requests {
			totalFilter := 0
			if v.Operator == "like" {
				for _, valueDetail := range v.Field {
					if totalFilter == 0 {
						sqlQuery.WriteString(" AND ( " + valueDetail + " LIKE " + fmt.Sprint("'%", v.Value, "%'"))
					} else {
						sqlQuery.WriteString(" OR " + valueDetail + " LIKE " + fmt.Sprint("'%", v.Value, "%'"))
					}
					totalFilter++
				}
			}

			if totalFilter > 0 {
				sqlQuery.WriteString(")")
			}
		}
	}

	fmt.Println("Query: ", sqlQuery.String())

	// for k, v := range req {
	// 	if k == "condition" {
	// 		// len(s)
	// 		res := v.(helper.DataFilter).Request
	// 		for _, valueDetail := range res {
	// 			//if valueDetail.Operator == "like" {
	// 			fmt.Println(valueDetail.Operator)
	// 			//}
	// 		}

	// 		// fmt.Println("---")
	// 		// fmt.Printf("Condition: %v", k)
	// 		// fmt.Println("---")
	// 		// fmt.Printf("Length: %v: ", len(res))
	// 		// fmt.Println("---")
	// 	}
	// 	//fmt.Printf("key[%s] value[%s]\n", k, v)
	// }

	db.connection.Raw(sqlQuery.String()).Scan(&records)

	// db.connection.Find(&records)

	// for k, v := range req {
	// 	if k == "condition" {
	// 		// len(s)
	// 		res := v.(helper.DataFilter).Request

	// 		fmt.Println("---")
	// 		fmt.Printf("Condition: %v", k)
	// 		fmt.Println("---")
	// 		fmt.Printf("Length: %v: ", len(res))
	// 		fmt.Println("---")
	// 	}
	// 	fmt.Printf("key[%s] value[%s]\n", k, v)
	// }
	// fmt.Println("---")

	//var p map[string]interface{}
	//v, found := req["condition"]

	// fmt.Println("Get")
	// if v, found := req["condition"]; found {
	// 	fmt.Println(v)
	// }

	//fmt.Println("Get Query: ")
	//db.connection.Debug().Where("name = ?", "jinzhu").Scan(&records)
	//fmt.Println("End")
	//db.connection.Raw("SELECT id, first_name, last_name, phone FROM customers").Scan(&records)
	//db.connection.Where(req["condition"])
	// db.connection

	//fmt.Println("End Get")

	return records
}
