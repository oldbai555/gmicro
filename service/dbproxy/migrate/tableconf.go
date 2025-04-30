package migrate

type ColumnConf struct {
	Name          string
	Type          string
	Comment       string
	AutoIncrement bool
	Size          int
	Unsigned      bool
	Default       string
}

type TableConf struct {
	Scaled      bool
	ScaleMinSeq uint32
	ScaleMaxSeq uint32
	Name        string
	Comment     string
	Fields      []*ColumnConf
	Keys        []*ColumnKeySt
}
