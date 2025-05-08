/**
 * @Author: zjj
 * @Date: 2025/5/6
 * @Desc:
**/

package vo

type GenClientCMD struct {
	ProtoBase
	Items []*GenClientCMDItem
}

type GenClientCMDItem struct {
	ProtoName string
	RpcName   string
}
