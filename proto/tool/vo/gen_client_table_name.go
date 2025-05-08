/**
 * @Author: zjj
 * @Date: 2025/5/6
 * @Desc:
**/

package vo

type GenClientTableName struct {
	ProtoBase
	Items []*GenClientTableNameItem
}

type GenClientTableNameItem struct {
	MessageName string
	TableName   string
}
