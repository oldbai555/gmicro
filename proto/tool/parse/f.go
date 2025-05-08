/**
 * @Author: zjj
 * @Date: 2025/5/9
 * @Desc:
**/

package parse

import (
	"fmt"
	utils "gmicro/common"
	"strings"
)

const (
	ProtoFileNameSuffix = ".proto"
	MessagePrefixModel  = "Model"
)

var ProtoFileNamePrefix = "lb"

func SetProtoFileNamePrefix(s string) {
	ProtoFileNamePrefix = s
}

// TrimProtoFileNameSuffix 去除proto文件的后缀 .proto
func TrimProtoFileNameSuffix(protoFileName string) string {
	return strings.TrimSuffix(protoFileName, ProtoFileNameSuffix)
}

// TrimProtoFileNamePrefix 去除proto文件 前缀
func TrimProtoFileNamePrefix(protoFileName string) string {
	if protoFileName == ProtoFileNamePrefix {
		return protoFileName
	}
	return strings.TrimPrefix(protoFileName, ProtoFileNamePrefix)
}

// TrimPrefixMessageNameWithModel 去除 Message 前缀 Model
func TrimPrefixMessageNameWithModel(msgName string) string {
	return strings.TrimPrefix(msgName, MessagePrefixModel)
}

// GenTableName 生成表名
func GenTableName(protoFileName, msgName string) string {
	return fmt.Sprintf("%s_%s", TrimProtoFileNameSuffix(protoFileName), utils.Camel2UnderScore(TrimPrefixMessageNameWithModel(msgName)))
}

// GetPbNameWithSuffix 获取Proto全称 包含后缀
func GetPbNameWithSuffix(pbName string) string {
	if !strings.HasSuffix(ProtoFileNameSuffix, pbName) {
		pbName = pbName + ProtoFileNameSuffix
	}
	return pbName
}
