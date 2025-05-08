package parser

import (
	"fmt"
	"github.com/emicklei/proto"
	"log"
	"os"
	"path"
	"proto2gorm/common"
	"strings"
)

type Field struct {
	Name    string
	Type    string
	Comment string
	Default string
}

type IndexMap map[string][]string // index name -> []fields

const (
	SQLLine        = "\t`%s` %s DEFAULT %s"
	CommentLine    = "\tCOMMENT '%s'"
	PRIMARYKeyLine = "\tPRIMARY KEY (`%s`),\n"
	IndexLine      = "\tINDEX `%s` (%s),\n"
	MySqlEngine    = "INNODB"
	CreateSQLHead  = "CREATE TABLE `%s` (" + "\n"
	CreateSQLTail  = "\n)\nENGINE=" + MySqlEngine + " DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT '%s';"
	SpaceN         = ",\n"
)

func ParseProtoToSQL(protoPath string) error {
	reader, err := os.Open(protoPath)
	if err != nil {
		log.Printf("can not open proto file %s, error: %v", protoPath, err)
		return err
	}
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		log.Printf("failed to parse proto: %v", err)
		return err
	}

	proto.Walk(definition,
		proto.WithMessage(func(m *proto.Message) {
			if !strings.HasPrefix(m.Name, "Model") {
				return
			}

			var fields []Field
			tableName := toSnakeCase(m.Name)
			indexes := make(IndexMap)
			var primaryKey string
			for _, element := range m.Elements {
				field, ok := element.(*proto.NormalField)
				if !ok {
					continue
				}

				comment := extractDesc(field.Comment)
				defaultVal := defaultValue(field.Type)

				// ✅ 判断是否是主键
				if strings.EqualFold(field.Name, "id") {
					primaryKey = field.Name
				}

				// 提取索引信息
				if idxName := extractIndex(field.Comment); idxName != "" {
					indexes[idxName] = append(indexes[idxName], field.Name)
				}

				fields = append(fields, Field{
					Name:    field.Name,
					Type:    mapProtoToSQL(field.Type),
					Comment: comment,
					Default: defaultVal,
				})
			}

			sql := fmt.Sprintf(CreateSQLHead, tableName)
			for _, f := range fields {
				sql += fmt.Sprintf(SQLLine, f.Name, f.Type, f.Default)
				if f.Comment != "" {
					sql += fmt.Sprintf(CommentLine, f.Comment)
				}
				sql += SpaceN
			}

			// 添加主键字段
			if primaryKey != "" {
				sql += fmt.Sprintf(PRIMARYKeyLine, primaryKey)
			}

			// 添加索引语句
			idxCount := 0
			for idxName, cols := range indexes {
				if len(cols) > 0 {
					sql += fmt.Sprintf(IndexLine, idxName, backtickJoin(cols))
					idxCount++
				}
			}

			// 移除最后多余的逗号
			sql = strings.TrimSuffix(sql, SpaceN) + "\n"

			// 表尾信息
			sql += fmt.Sprintf(CreateSQLTail, m.Name)

			fmt.Println("📄 Generated SQL:", sql)

			filePath := path.Join("mysql", tableName+".sql")
			if err := common.CreateAndWriteFile(filePath, sql); err != nil {
				log.Fatalf("write file error: %v", err)
			}
		}),
	)

	return nil
}

// 提取 @desc 注释
func extractDesc(comment *proto.Comment) string {
	if comment == nil {
		return ""
	}
	if strings.Contains(comment.Message(), "@desc:") {
		return strings.Replace(strings.TrimSpace(comment.Message()), "@desc:", "", -1)
	}
	return ""
}

// 提取 @index:"xxx" 索引名
func extractIndex(comment *proto.Comment) string {
	if comment == nil {
		return ""
	}
	for _, msg := range comment.Lines {
		msg = strings.TrimSpace(msg)
		if strings.Contains(msg, "@index:") {
			start := strings.Index(msg, `@index:"`)
			if start >= 0 {
				sub := msg[start+8:]
				end := strings.Index(sub, `"`)
				if end > 0 {
					return sub[:end]
				}
			}
		}
	}

	return ""
}

// SQL 类型映射
func mapProtoToSQL(protoType string) string {
	switch protoType {
	case "uint64":
		return "BIGINT UNSIGNED"
	case "int64":
		return "BIGINT"
	case "uint32":
		return "INT UNSIGNED"
	case "int32":
		return "INT"
	case "string":
		return "VARCHAR(255)"
	case "bool":
		return "BOOLEAN"
	default:
		return "TEXT"
	}
}

// 默认值设置
func defaultValue(protoType string) string {
	switch protoType {
	case "uint64", "uint32", "int64", "int32":
		return "0"
	case "string":
		return "''"
	case "bool":
		return "false"
	default:
		return "NULL"
	}
}

// 转为 backtick 包裹字段名
func backtickJoin(cols []string) string {
	var b strings.Builder
	for i, c := range cols {
		b.WriteString("`")
		b.WriteString(c)
		b.WriteString("`")
		if i < len(cols)-1 {
			b.WriteString(", ")
		}
	}
	return b.String()
}

// 驼峰转蛇形
func toSnakeCase(name string) string {
	var result []rune
	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
