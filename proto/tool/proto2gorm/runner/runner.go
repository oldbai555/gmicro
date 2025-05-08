package runner

import (
	"gorm.io/gen/field"
	"log"
	_ "os/exec"

	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/rawsql"
)

func ExecuteAndGenerate() error {
	gormDB, err := gorm.Open(rawsql.New(rawsql.Config{
		FilePath: []string{
			"./mysql", // 建表sql目录
		},
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
		gen.FieldType("deleted_at", "gorm.DeletedAt"),
	}
	g := gen.NewGenerator(gen.Config{
		OutPath:           "./dao",
		Mode:              gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldCoverable:    true,
		FieldWithTypeTag:  true,
		FieldWithIndexTag: true,
	})
	g.UseDB(gormDB)
	models := g.GenerateAllTable(fieldOpts...)
	g.ApplyBasic(models...)
	g.Execute()

	return nil
}
