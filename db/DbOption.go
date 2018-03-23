/**
 * @Desc
 * @author zjhfyq
 * @data 2018/3/22 13:13.
 */
package db

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"goPageHelper/models"
	"goPageHelper/process"
	"log"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

//适合连表查询
//sqlStr sql查询语句
//sqlAfterFrom sql语句查询的表名和条件
//targetModel 查询结果封装到的目标对象
func GetDbPointer(driverName string, dataSourceName string) (db *sql.DB) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(errors.New(err.Error()))
		log.Println(debug.Stack())
	}
	return
}

//适合单表查询
//model 用于生成sql ，以及查询结果封装
func QueryByModel(db *sql.DB, model interface{}, tableName string, pageCondition *models.PageQueryCondition) (pageInfo *models.PageInfo) {
	pageInfo = GetPageInfo(db, tableName, pageCondition)
	result := SelectByModel(db, model, tableName, pageCondition)
	pageInfo.ListData = result
	pageInfo.Size = len(result)
	return
}


//适合连表查询
//sqlStr sql查询语句
//sqlAfterFrom sql语句查询的表名和条件
//targetModel 查询结果封装到的目标对象
func QueryBySql(db *sql.DB, sqlStr string,sqlAfterFrom string,pageCondition *models.PageQueryCondition,targetModel interface{}) (pageInfo *models.PageInfo) {
	pageInfo = GetPageInfo(db, sqlAfterFrom, pageCondition)
	result := SelectBySql(db,sqlStr,pageCondition,targetModel)
	pageInfo.ListData = result
	pageInfo.Size = len(result)
	return
}

//适合单表查询
//model 用于生成sql ，以及查询结果封装
func SelectByModel(db *sql.DB, model interface{}, tableName string, pageCondition *models.PageQueryCondition) (result []interface{}) {
	rtype := reflect.TypeOf(model)
	if rtype.Kind() == reflect.Ptr {
		rtype = rtype.Elem()
	}
	if rtype.Kind() == reflect.Slice {
		rtype = rtype.Elem()
	}
	sqlStr := "select  "
	fieldNum := rtype.NumField()
	for i := 0; i < fieldNum; i++ {
		if i != fieldNum-1 {
			sqlStr += rtype.Field(i).Name + ","
		} else {
			sqlStr += rtype.Field(i).Name
		}
	}
	sqlStr += " from " + tableName
	if pageCondition.PageNum != 0 {
		dbnum := process.GetDbNum(*pageCondition)
		sqlStr += " limit " + strconv.Itoa(dbnum.Limit) + " offset " + strconv.Itoa(dbnum.Offset)
	}
	log.Println(sqlStr)
	rows, err := db.Query(sqlStr)
	defer  rows.Close()
	if logError(err) {
		columns, err := rows.Columns()
		clen := len(columns)
		if logError(err) {
			for rows.Next() {
				midv := SetValue(clen,rtype,rows,columns)
				result = append(result, midv)
			}
		}
	}
	return
}



func SelectBySql(db *sql.DB, sqlStr string ,pageCondition *models.PageQueryCondition,targetModel interface{})(result []interface{}){
	rtype := reflect.TypeOf(targetModel)
	if rtype.Kind() == reflect.Ptr {
		rtype = rtype.Elem()
	}
	if rtype.Kind() == reflect.Slice {
		rtype = rtype.Elem()
	}

	if pageCondition.PageNum != 0 {
		dbnum := process.GetDbNum(*pageCondition)
		sqlStr += " limit " + strconv.Itoa(dbnum.Limit) + " offset " + strconv.Itoa(dbnum.Offset)
	}
	log.Println(sqlStr)
	rows, err := db.Query(sqlStr)
	defer  rows.Close()
	if logError(err) {
		columns, err := rows.Columns()
		clen := len(columns)
		if logError(err) {
			for rows.Next() {
				midv := SetValue(clen,rtype,rows,columns)
				result = append(result, midv)
			}
		}
	}
	return
}

//封装查询结果
func SetValue(clen int,rtype reflect.Type,row *sql.Rows,columns []string) (result interface{}) {
	values := make([]interface{}, clen)
	scanArgs := make([]interface{}, clen)
	midModel := reflect.New(rtype)
	midv := midModel.Elem()

	for i := 0; i < clen; i++ {
		scanArgs[i] = &values[i]
	}
	row.Scan(scanArgs...)
	for i := 0; i < clen; i++ {
		var valueStr = ""
		//实际上数据库返回回来的数据  都是[]uint8类型的
		if value, ok := values[i].([]uint8); ok {
			valueStr = string(value)
		}
		if strings.Contains(midv.FieldByName(columns[i]).Type().String(), "int") {
			valueInt, err := strconv.Atoi(valueStr)
			if logError(err) {
				midv.FieldByName(columns[i]).SetInt(int64(valueInt))
			}
		} else if strings.Contains(midv.FieldByName(columns[i]).Type().String(), "string") {
			midv.FieldByName(columns[i]).SetString(valueStr)
		} else if strings.Contains(midv.FieldByName(columns[i]).Type().String(), "bool") {
			valueBool, err := strconv.ParseBool(valueStr)
			if logError(err) {
				midv.FieldByName(columns[i]).SetBool(valueBool)
			}
		} else if strings.Contains(midv.FieldByName(columns[i]).Type().String(), "float") {
			valueFloat64, err := strconv.ParseFloat(valueStr, 64)
			if logError(err) {
				midv.FieldByName(columns[i]).SetFloat(valueFloat64)
			}
		} else if strings.Contains(midv.FieldByName(columns[i]).Type().String(), "time") {
			valueTime, err := time.Parse("2006-1-2 15:4:5", valueStr)
			if logError(err) {
				midv.FieldByName(columns[i]).Set(reflect.ValueOf(valueTime))
			}
		}
	}
	result = midv.Interface()
	return
}























//传指针是为了合理化页码
func GetPageInfo(db *sql.DB, tableName string, pageCondition *models.PageQueryCondition) (pageInfo *models.PageInfo) {

	pageInfo = &models.PageInfo{}
	sqlCount := "select count(*) from " + tableName

	count := 0
	rows, err := db.Query(sqlCount)
	defer  rows.Close()
	if logError(err) {
		for rows.Next() {
			rows.Scan(&count)
		}
	}
	//设置默页数
	if pageCondition.PageNum <= 0 {
		log.Printf("pageNum: %s is invalid ",pageCondition.PageNum)
		pageCondition.PageNum = 1
	}
	//设置默认的页面大小
	if pageCondition.PageSize <= 0 {
		log.Printf("pageSize: %s is invalid ",pageCondition.PageSize)
		pageCondition.PageSize = 15
	}

	pageInfo.Total = count

	pages := 0
	if count%pageCondition.PageSize == 0 {
		pages = count / pageCondition.PageSize
	} else {
		pages = count/pageCondition.PageSize + 1
	}
	pageInfo.Pages = pages

	//合理化查询的页码
	if pageCondition.PageNum > pages {
		pageCondition.PageNum = pages
	}
	if pageCondition.PageNum <= 0 {
		pageCondition.PageNum = 1
	}

	pageInfo.PageNum = pageCondition.PageNum
	pageInfo.PageSize = pageCondition.PageSize

	if pages > pageInfo.PageNum {
		pageInfo.IsLastPage = false
		pageInfo.NextPage = pageInfo.PageNum + 1
	} else {
		pageInfo.IsLastPage = true
		pageInfo.NextPage = pageInfo.PageNum
	}

	if pageInfo.PageNum > 1 {
		pageInfo.IsFirstPage = false
		pageInfo.PrePage = pageInfo.PageNum - 1
	} else {
		pageInfo.IsFirstPage = true
		pageInfo.PrePage = 1
	}

	dbnum := process.GetDbNum(*pageCondition)
	//如果偏移大于总数
	if dbnum.Offset > pageInfo.Total {
		pageCondition.PageNum = 1
	}

	return
}

func logError(err error) bool {
	if err != nil {
		log.Println("ERROR : ", err)
		log.Println(debug.Stack())
		return false
	} else {
		return true
	}
}
