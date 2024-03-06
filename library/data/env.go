package data

type EnvType uint8

const (
	EnvTypeString EnvType = iota
	EnvTypeInt
	EnvTypeFloat
	EnvTypeStringSlice
	EnvTypeStringStringMap
)

type Env struct {
	Type        EnvType
	Name        string
	Description string
}
