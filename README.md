#  go上的pageHelper
 &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;做java开发的时候，很喜欢mybatis的分页插件pageHelper，来go之后，发现go上并没有类似的分页的插件,于是
动手仿照pageHelper的设计思想，粗超的完成了这个go上的goPageHelper。
## 2018-3-22
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;目前只是简单的实现了单表查询，连表查询的功能。在使用pageInfo的时候，实现了分页的合理化。
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;目前存在的问题：在使用pageInfo的时候，会执行select count(*) from ***的语句，对于连表查询而言，如果数据过大，可能造成性能上的影响，单表查询则不存在。当然了，你也可以不使用pageInfo。
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;对于上面这个问题，目前想法是加上一个开关来实现是否启动totalCount。