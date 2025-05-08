/**
 * @Author: zjj
 * @Date: 2025/8/7
 * @Desc:
**/

package option

import (
	"fmt"
	"github.com/emicklei/proto"
	utils "gmicro/common"
	"gmicro/pkg/log"
	parse2 "gmicro/proto/tool/parse"
	"gmicro/proto/tool/vo"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/rawsql"
	"path"
	"strings"
)

/*

// 文件表
message ModelFile {
    uint64 id = 1;
    uint32 created_at = 2;
    uint32 updated_at = 3;
    uint32 deleted_at = 4;
    uint64 creator_id = 5;
    int64 size = 6;
    // @desc: 原文件名
    string name = 7;
    // @desc: 文件重命名
    string rename = 8;
    // @desc: 文件路径
    // @gorm:"index:idx_path"
    string path = 9;
    // @desc: 存储桶
    // @gorm:"index:idx_bucket_domain;type:VARCHAR(511)"
    string bucket = 10;
    // @desc: 域名
    // @gorm:"index:idx_bucket_domain;type:VARCHAR(511)"
    string domain = 11;
}

*/

func WithGenServerGormGen() Option {
	return func(pbCtx *parse2.PbContext, req *vo.CodeFuncParams) {
		serverPath := path.Join(req.OutputDir, pbCtx.ServiceName+"server")
		if !parse2.FileExists(serverPath) {
			utils.CreateDir(serverPath)
		}
		implPath := path.Join(serverPath, "impl")
		scriptPath := path.Join(implPath, "script")
		if !parse2.FileExists(scriptPath) {
			utils.CreateDir(scriptPath)
		}
		for _, message := range pbCtx.MsgList {
			if !strings.HasPrefix(message.Name, "Model") {
				continue
			}
			sql := ParseProtoToSQL(message)
			err := utils.CreateAndWriteFile(path.Join(scriptPath, toSnakeCase(message.Name)+".sql"), sql)
			if err != nil {
				log.Errorf("err:%v", err)
			}
		}
		err := executeAndGenerate(scriptPath, implPath)
		if err != nil {
			log.Errorf("err:%v", err)
		}
	}
}

type Querier interface {
	// SELECT * FROM @@table WHERE id=@id and (deleted_at=0 OR deleted_at IS NULL)
	GetById(id int) (gen.T, error)
}

func executeAndGenerate(scriptPath, implPath string) error {
	gormDB, err := gorm.Open(rawsql.New(rawsql.Config{
		FilePath: []string{scriptPath}, // 建表sql目录
	}))
	if err != nil {
		log.Fatalf("err:%v", err)
	}
	fieldOpts := []gen.ModelOpt{
		gen.FieldGORMTag("updated_at", func(tag field.GormTag) field.GormTag {
			tag.Set("autoUpdateTime", "")
			return tag
		}),
		gen.FieldGORMTag("created_at", func(tag field.GormTag) field.GormTag {
			tag.Set("autoCreateTime", "")
			return tag
		}),
		gen.FieldType("deleted_at", "soft_delete.DeletedAt"),
	}
	g := gen.NewGenerator(gen.Config{
		OutPath:           path.Join(implPath, "query"),
		Mode:              gen.WithDefaultQuery | gen.WithoutContext | gen.WithQueryInterface,
		FieldCoverable:    true,
		FieldWithTypeTag:  true,
		FieldWithIndexTag: true,
		FieldSignable:     true,
	})
	g.UseDB(gormDB)
	models := g.GenerateAllTable(fieldOpts...)
	g.ApplyBasic(models...)
	g.ApplyInterface(func(Querier) {}, models...)
	g.Execute()

	return nil
}

type Field struct {
	Name    string
	Type    string
	Comment string
	Default string
}

type IndexMap map[string][]string // index name -> []fields

const (
	SQLLine        = "\t`%s` %s %s\t"
	CommentLine    = "\tCOMMENT '%s'\t"
	PRIMARYKeyLine = "\tPRIMARY KEY (`%s`),\n"
	IndexLine      = "\tINDEX `%s` (%s),\n"
	MySqlEngine    = "INNODB"
	CreateSQLHead  = "CREATE TABLE `%s` (" + "\n"
	CreateSQLTail  = "\n)\nENGINE=" + MySqlEngine + " DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT '%s';"
	SpaceN         = ",\n"
	PKeyDesc       = "\tNOT NULL AUTO_INCREMENT"
)

func ParseProtoToSQL(m *proto.Message) string {
	// 提取 @desc 注释
	var extractDesc = func(comment *proto.Comment) string {
		if comment == nil {
			return ""
		}
		if strings.Contains(comment.Message(), "@desc:") {
			return strings.TrimSpace(strings.Replace(comment.Message(), "@desc:", "", 1))
		}
		return ""
	}

	// 从 @gorm 标签提取索引名
	var extractIndex = func(comment *proto.Comment) string {
		if comment == nil {
			return ""
		}
		for _, line := range comment.Lines {
			line = strings.TrimSpace(line)
			if strings.Contains(line, "@gorm:") {
				start := strings.Index(line, `@gorm:"`)
				if start >= 0 {
					tagContent := line[start+7:]
					end := strings.Index(tagContent, `"`)
					if end > 0 {
						tagContent = tagContent[:end]
						parts := strings.Split(tagContent, ";")
						for _, p := range parts {
							if strings.HasPrefix(strings.ToLower(p), "index:") {
								return strings.TrimPrefix(p, "index:")
							}
						}
					}
				}
			}
		}
		return ""
	}

	// 从 @gorm 标签提取自定义 SQL 类型
	var extractSQLType = func(comment *proto.Comment) string {
		if comment == nil {
			return ""
		}
		for _, line := range comment.Lines {
			line = strings.TrimSpace(line)
			if strings.Contains(line, "@gorm:") {
				start := strings.Index(line, `@gorm:"`)
				if start >= 0 {
					tagContent := line[start+7:]
					end := strings.Index(tagContent, `"`)
					if end > 0 {
						tagContent = tagContent[:end]
						parts := strings.Split(tagContent, ";")
						for _, p := range parts {
							if strings.HasPrefix(strings.ToLower(p), "type:") {
								return strings.TrimPrefix(p, "type:")
							}
						}
					}
				}
			}
		}
		return ""
	}

	// SQL 类型映射（优先 @gorm type）
	var mapProtoToSQL = func(protoType string, comment *proto.Comment) string {
		if sqlType := extractSQLType(comment); sqlType != "" {
			return sqlType
		}
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

	// 默认值
	var defaultValue = func(protoType string) string {
		var d = "\tDEFAULT\t"
		var v = ""
		switch protoType {
		case "uint64", "uint32", "int64", "int32":
			v = "0"
		case "string":
			v = "''"
		case "bool":
			v = "false"
		default:
			v = "NULL"
		}
		return d + v
	}

	// backtick 拼接列
	var backtickJoin = func(cols []string) string {
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

	var fields []Field
	tableName := toSnakeCase(m.Name)
	indexes := make(IndexMap)
	var primaryKey string

	for _, element := range m.Elements {
		normalField, ok := element.(*proto.NormalField)
		if !ok {
			continue
		}

		comment := extractDesc(normalField.Comment)
		sqlType := mapProtoToSQL(normalField.Type, normalField.Comment)
		defaultVal := defaultValue(normalField.Type)

		// 主键判断
		if strings.EqualFold(normalField.Name, "id") {
			primaryKey = normalField.Name
			defaultVal = PKeyDesc
		}

		// 索引
		if idxName := extractIndex(normalField.Comment); idxName != "" {
			indexes[idxName] = append(indexes[idxName], normalField.Name)
		}

		fields = append(fields, Field{
			Name:    normalField.Name,
			Type:    sqlType,
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

	if primaryKey != "" {
		sql += fmt.Sprintf(PRIMARYKeyLine, primaryKey)
	}

	for idxName, cols := range indexes {
		if len(cols) > 0 {
			sql += fmt.Sprintf(IndexLine, idxName, backtickJoin(cols))
		}
	}

	sql = strings.TrimSuffix(sql, SpaceN) + "\n"
	sql += fmt.Sprintf(CreateSQLTail, m.Name)

	log.Infof("📄 Generated SQL: %s", sql)
	return sql
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
