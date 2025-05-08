package syscfg

var Global *SysCfg

func InitGlobal(path string, opt ...Option) {
	Global = LoadSysCfgByYaml(path, opt...)
}
