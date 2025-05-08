package syscfg

import (
	"github.com/spf13/viper"
	"gmicro/pkg/json"
	"gmicro/pkg/log"
)

type SysCfg struct {
	V *viper.Viper

	ServerConf *ServerConf
}

var sc *SysCfg

func New(viper *viper.Viper, option ...Option) (*SysCfg, error) {
	sc = &SysCfg{
		V: viper,
	}
	option = append(option, OptionWithServer())
	// 初始化组件
	for _, o := range option {
		o(sc)
	}
	return sc, nil
}

func JsonConvertStruct(re interface{}, out interface{}) error {
	marshal, err := json.Marshal(re)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}

	err = json.Unmarshal(marshal, out)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	return nil
}
