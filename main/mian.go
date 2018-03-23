/**
 * @Desc
 * @author zjhfyq 
 * @data 2018/3/22 11:48.
 */
package main

import (
	_"github.com/go-sql-driver/mysql"
	"goPageHelper/db"
	"goPageHelper/models"
	"fmt"
	"encoding/json"
)



func main() {

	////打开数据库连接
	dbp := db.GetDbPointer("mysql","zct:123456@tcp(10.0.0.101:3306)/go_generater?charset=utf8")

	defer  dbp.Close()

	//单表查询
	//result :=db.Query(dbp ,models.Users{},"users",&models.PageQueryCondition{PageNum:308,PageSize:20})
	result := db.QueryBySql(dbp,"select  Username,Password,Age from users ","users",&models.PageQueryCondition{PageNum:308,PageSize:20},models.Users{})
	jsons , err :=json.MarshalIndent(result,"	"," ")
	if err == nil {
		fmt.Printf(string(jsons))
	}
	//
	//var users []models.Users
	//for _ , inter :=range result  {
	//	if value, ok := inter.(models.Users);ok {
	//		users = append(users, value)
	//	}
	//}
	//for index , user := range users{
	//	fmt.Printf("index : %d , username : %s ,password : %s ,age : %d \n",index,user.Username,user.Password,user.Age)
	//}
}


