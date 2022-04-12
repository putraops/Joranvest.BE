package repository

import (
	"encoding/json"
	"fmt"
	"joranvest/helper"
	"joranvest/models"
	"runtime"
	"strings"

	"gorm.io/gorm"
)

type ServiceRepository interface {
	GetAll(filter map[string]interface{}) interface{}
	ConvertViewQueryIntoViewCount(param string) string
	ConvertViewQueryIntoViewCountByPublic(param string, tableName string) string
	getCurrentFuncName() string
	MapFields(oldRecord interface{}, newRecord interface{}) helper.Result
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

func (db *serviceConnection) MapFields(oldRecord interface{}, newRecord interface{}) helper.Result {
	var oldRecordMap map[string]interface{}
	oldRecordJson, _ := json.Marshal(oldRecord)
	json.Unmarshal(oldRecordJson, &oldRecordMap)

	var newRecordMap map[string]interface{}
	newRecordJson, _ := json.Marshal(newRecord)
	json.Unmarshal(newRecordJson, &newRecordMap)

	for key, newValue := range newRecordMap {
		switch c := newValue.(type) {
		case bool:
			if newRecordMap[key] != oldRecordMap[key] {
				oldRecordMap[key] = newRecordMap[key]
			}
		case string, int, int8, int16, int32, int64, float32, float64:
			if newRecordMap[key] != "" && newRecordMap[key] != oldRecordMap[key] {
				oldRecordMap[key] = newRecordMap[key]
			}
		case map[string]interface{}:
			//-- Check map is a sql.NullTime
			if _, isTime := c["Time"]; isTime {
				if _, isSqlNullTime := c["Valid"]; isSqlNullTime {
					if c["Valid"].(bool) && newRecordMap[key].(map[string]interface{})["Time"] != oldRecordMap[key].(map[string]interface{})["Time"] {
						oldRecordMap[key] = newRecordMap[key]
					}
				}
			}
		default:
			helper.StandartResult(false, fmt.Sprintf("There is no map for key %v", c), helper.EmptyObj{})
		}
	}

	resultObj, err := json.Marshal(oldRecordMap)
	if err != nil {
		return helper.StandartResult(false, err.Error(), helper.EmptyObj{})
	}
	return helper.StandartResult(true, "Ok", resultObj)
}

func (db *serviceConnection) getCurrentFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
