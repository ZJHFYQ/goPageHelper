/**
 * @Desc
 * @author zjhfyq
 * @data 2018/3/22 10:24.
 */
package models

//页面信息
type PageInfo struct {
	PageNum     int  //当前页
	PageSize    int  //每页的数量
	Size        int  //当前页的数量
	Total       int  //数据库中总的记录数目
	Pages       int  //总的页数
	PrePage     int  //前一页
	NextPage    int  //下一页
	IsFirstPage bool //是否为第一页
	IsLastPage  bool //是否为最后一页
	ListData    interface{} //查询出来的数据
}

//页面查询条件
type PageQueryCondition struct {
	PageNum  int //要查询的页码
	PageSize int //要查询的每页的个数
}

type DbNum struct {
	Limit  int //条数
	Offset int //偏移数
}
