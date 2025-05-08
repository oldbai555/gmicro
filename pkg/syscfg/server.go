package syscfg

import (
	"fmt"
	"github.com/spf13/viper"
	"gmicro/pkg/gerr"
	"gmicro/pkg/log"
)

const defaultApolloServerPrefix = "server"

type ServerConf struct {
	Name string `json:"name"`
	Port uint32 `json:"port"`
	Ip   string `json:"ip"`
}

func NewServerConf(viper *viper.Viper) *ServerConf {
	var v ServerConf
	val := viper.Get(defaultApolloServerPrefix)
	err := JsonConvertStruct(val, &v)
	if err != nil {
		log.Errorf("err is %v", err)
		panic(err)
	}
	return &v
}

func GetServerConf() (*ServerConf, error) {
	if Global == nil {
		return nil, gerr.NewInvalidArg("not found Global")
	}
	if Global.ServerConf == nil {
		return nil, gerr.NewInvalidArg("not found server conf")
	}
	return Global.ServerConf, nil
}

func GetSrvAddr() (string, error) {
	conf, err := GetServerConf()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", conf.Ip, conf.Port), nil
}
