package syscfg

type Option func(tool *SysCfg)

func OptionWithServer() Option {
	return func(tool *SysCfg) {
		tool.ServerConf = NewServerConf(tool.V)
	}
}
