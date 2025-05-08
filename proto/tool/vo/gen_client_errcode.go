/**
 * @Author: zjj
 * @Date: 2025/5/6
 * @Desc:
**/

package vo

type GenClientErrCode struct {
	ProtoBase
	Items []*GenClientErrCodeItem
}

type GenClientErrCodeItem struct {
	Name    string
	Value   int32
	Comment string
}
