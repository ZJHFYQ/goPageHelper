/**
 * @Desc
 * @author zjhfyq 
 * @data 2018/3/22 11:36.
 */
package process

import "goPageHelper/models"


//计算offset 和 limit
func GetDbNum(condition models.PageQueryCondition)(dbnum models.DbNum){
	if condition.PageNum <=1 {
		condition.PageNum = 1
	}
	if condition.PageSize <= 1 {
		condition.PageSize =10
	}
	dbnum.Limit = condition.PageSize
	dbnum.Offset = (condition.PageNum - 1) * condition.PageSize
	return
}


