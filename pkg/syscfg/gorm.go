package syscfg

import (
	"fmt"
	"gmicro/pkg/json"
	"os"
	path2 "path"
	"path/filepath"
)

const defaultApolloMysqlPrefix = "mysql"
const defaultDatabase = "biz"

type MysqlConf struct {
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func NewGormMysqlConf(path string) *MysqlConf {
	if path == "" {
		path = defaultMysqlConfPath
	}
	data, err := os.ReadFile(filepath.ToSlash(path2.Join(path, "mysql.json")))
	if err != nil {
		panic(err)
	}

	var v MysqlConf
	err = json.Unmarshal(data, &v)
	if err != nil {
		panic(err)
	}
	return &v
}

func (m *MysqlConf) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", m.Username, m.Password, m.Addr, m.Port, m.Database)
}
