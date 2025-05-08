/**
 * @Author: zjj
 * @Date: 2025/7/16
 * @Desc:
**/

package vo

type GenServerAutoDbAutoGen struct {
	Name        string
	ServiceName string
	Field       []*GenServerAutoDbAutoGenField
}

type GenServerAutoDbAutoGenField struct {
	FieldName  string
	Type       string
	OrmType    string
	Comment    string
	IsArray    bool
	IsUnsigned bool
}
