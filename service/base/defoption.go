/**
 * @Author: zjj
 * @Date: 2025/9/17
 * @Desc:
**/

package base

import (
	"gmicro/pkg/gerr"
	"gmicro/pkg/log"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

func DefaultOptionsProcessor(db gen.Dao, listOption *ListOption) error {
	if listOption == nil {
		return gerr.NewInvalidArg("not found option")
	}
	tableName := db.TableName()
	err := NewOptionsProcessor(listOption).
		AddOffsetLimit(DefaultListOption_DefaultListOptionOffsetLimit, func(offset, limit int) error {
			db.Offset(offset).Limit(limit)
			return nil
		}).
		AddStringList(
			DefaultListOption_DefaultListOptionSelect,
			func(valList []string) error {
				var fields []field.Expr
				for _, val := range valList {
					fields = append(fields, field.NewField(tableName, val))
				}
				db.Select(fields...)
				return nil
			}).
		AddUint32(
			DefaultListOption_DefaultListOptionOrderBy,
			func(val uint32) error {
				if val == uint32(DefaultOrderBy_DefaultOrderByCreatedAtDesc) {
					newField := field.NewField(tableName, "created_at")
					db.Order(newField.Desc())
				} else if val == uint32(DefaultOrderBy_DefaultOrderByCreatedAtAcs) {
					newField := field.NewField(tableName, "created_at")
					db.Order(newField.Asc())
				} else if val == uint32(DefaultOrderBy_DefaultOrderByIdDesc) {
					newField := field.NewField(tableName, "id")
					db.Order(newField.Desc())
				}
				return nil
			}).
		AddStringList(
			DefaultListOption_DefaultListOptionGroupBy,
			func(valList []string) error {
				var fields []field.Expr
				for _, val := range valList {
					fields = append(fields, field.NewField(tableName, val))
				}
				db.Group(fields...)
				return nil
			}).
		AddBool(
			DefaultListOption_DefaultListOptionWithTrash,
			func(val bool) error {
				if val {
					db.Unscoped()
				}
				return nil
			}).
		AddUint64List(
			DefaultListOption_DefaultListOptionIdList,
			func(valList []uint64) error {
				newField := field.NewUint64(tableName, "id")
				if len(valList) == 1 {
					db.Where(newField.Eq(valList[0]))
				} else {
					db.Where(newField.In(valList...))
				}
				return nil
			}).
		AddTimeStampRange(
			DefaultListOption_DefaultListOptionCreatedAt,
			func(begin, end uint32) error {
				newField := field.NewUint32(tableName, "created_at")
				db.Where(newField.Between(begin, end))
				return nil
			}).
		AddUint64List(
			DefaultListOption_DefaultListOptionCreatorIdList,
			func(valList []uint64) error {
				newField := field.NewUint64(tableName, "creator_id")
				if len(valList) == 1 {
					db.Where(newField.Eq(valList[0]))
				} else {
					db.Where(newField.In(valList...))
				}
				return nil
			}).
		Process()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}
