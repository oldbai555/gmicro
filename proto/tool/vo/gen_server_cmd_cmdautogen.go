/**
 * @Author: zjj
 * @Date: 2025/5/7
 * @Desc:
**/

package vo

type GenServerCmdCmdAutoGen struct {
	ProtoBase
	Items []*GenServerCmdCmdAutoGenItem
}

type GenServerCmdCmdAutoGenItem struct {
	ProtoName string
	RpcName   string
}
