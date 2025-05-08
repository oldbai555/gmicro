/**
 * @Author: zjj
 * @Date: 2025/5/6
 * @Desc:
**/

package vo

type GenClientRpc struct {
	ProtoBase
	Items []*GenClientRpcItem
}

type GenClientRpcItem struct {
	RpcName   string
	ProtoName string
}
