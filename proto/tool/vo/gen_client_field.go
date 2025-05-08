/**
 * @Author: zjj
 * @Date: 2025/5/6
 * @Desc:
**/

package vo

type GenClientField struct {
	ProtoBase
	Items []*GenClientFieldItem
}

type GenClientFieldItem struct {
	Name string
	CV   string
	UV   string
}
