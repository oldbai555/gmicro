/**
 * @Author: zjj
 * @Date: 2025/7/16
 * @Desc:
**/

package option

import (
	"github.com/emicklei/proto"
	utils "gmicro/common"
	"gmicro/pkg/log"
	parse2 "gmicro/proto/tool/parse"
	"gmicro/proto/tool/temp"
	"gmicro/proto/tool/temp/gtpl"
	"gmicro/proto/tool/vo"
	"path"
	"strings"
)

func WithGenServerAutoDbAutoGen() Option {
	return func(pbCtx *parse2.PbContext, req *vo.CodeFuncParams) {
		var isBuiltinType = func(typ string) bool {
			switch typ {
			case "string", "uint32", "int32", "uint64", "int64", "bool", "bytes", "float", "double":
				return true
			}
			return false
		}

		var mysqlFieldType = func(typ string) string {
			// 根据类型返回对应的 MySQL 字段类型
			switch typ {
			case "string":
				return "varchar(255)"
			case "uint32", "int32":
				return "int"
			case "uint64", "int64":
				return "bigint"
			case "bool":
				return "tinyint(1)"
			case "bytes":
				return "blob"
			case "float":
				return "float"
			case "double":
				return "double"
			default:
				return "json" // 默认使用 TEXT 类型
			}
		}

		var isUnsigned = func(typ string) bool {
			// 这里可以根据需要实现是否需要无符号的逻辑
			switch typ {
			case "uint32", "uint64":
				return true
			}
			return false
		}

		var getComment = func(field *proto.NormalField) string {
			// 获取消息的注释
			if field.Comment != nil {
				for _, line := range field.Comment.Lines {
					if strings.HasPrefix("@desc:", line) {
						return strings.TrimSpace(strings.TrimPrefix(line, "@desc:"))
					}
				}
				return ""
			}
			return ""
		}

		serverPath := path.Join(req.OutputDir, pbCtx.ServiceName+"server")
		if !parse2.FileExists(serverPath) {
			utils.CreateDir(serverPath)
		}
		autoDbPath := path.Join(serverPath, "autodb")
		if !parse2.FileExists(autoDbPath) {
			utils.CreateDir(autoDbPath)
		}

		for _, message := range pbCtx.MsgList {
			if !parse2.IsModel(message) {
				continue
			}
			var genVo = vo.GenServerAutoDbAutoGen{
				Name:        message.Name,
				ServiceName: pbCtx.ServiceName,
				Field:       nil,
			}
			fields := parse2.GetMsgFields(message)
			for _, field := range fields.NormalFields {
				var typ = "object"
				if isBuiltinType(field.Type) {
					typ = field.Type
				}
				genVo.Field = append(genVo.Field, &vo.GenServerAutoDbAutoGenField{
					FieldName:  field.Name,
					Type:       typ,
					OrmType:    mysqlFieldType(field.Type),
					Comment:    getComment(field),
					IsArray:    field.Repeated,
					IsUnsigned: isUnsigned(field.Type),
				})
			}
			mainPath := path.Join(autoDbPath, strings.ToLower(message.Name)+"_autogen.go")
			template, err := temp.GenCodeByTemplate(gtpl.GenServerAutoDbAutogen, &genVo)
			if err != nil {
				log.Errorf("gen server auto db autogen error: %v", err)
				return
			}
			err = utils.CreateAndWriteFile(mainPath, template)
			if err != nil {
				log.Errorf("create and write file error: %v", err)
				return
			}
		}

	}
}
