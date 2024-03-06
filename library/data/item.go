package data

type Item interface {
	// Init 初始化
	Init(config Config) bool
	GetConfig() Config
	GetName() string
	GetID() string
}

type BaseItem struct {
	config Config
}

func (bi *BaseItem) Init(config Config) bool {
	bi.config = config
	return true
}

func (bi *BaseItem) GetConfig() Config {
	return bi.config
}

func (bi *BaseItem) GetName() string {
	return bi.config.GetName()
}

func (bi *BaseItem) GetID() string {
	return bi.config.GetID()
}
