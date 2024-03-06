package data

import (
	"encoding/json"
	"github.com/finishy1995/effibot-core/library/id"
	"github.com/finishy1995/go-library/storage"
	"github.com/zeromicro/go-zero/core/logx"
	"reflect"
	"sync"
)

type VersionType uint32

const (
	NilVersion   VersionType = 0
	StartVersion VersionType = 1
)

type BaseConfig[T ConfigVersion] struct {
	storage.Model
	mutex         sync.Mutex
	ID            string                 `dynamo:",hash" json:"id"`
	Name          string                 `json:"name"`
	Type          string                 `json:"type"`
	ID2VersionMap map[string]VersionType `json:"id_2_version_map"`
	Versions      []T                    `json:"versions"`
	versionIndex  map[VersionType]int
	Latest        VersionType `json:"latest"`
	envArray      []*Env
	loader        Loader
}

func (bc *BaseConfig[T]) newT(needInit bool) (T, bool) {
	// 使用反射来创建 T 的实例
	var defaultValue T
	tType := reflect.TypeOf(defaultValue)
	if tType.Kind() == reflect.Ptr && tType.Elem().Kind() == reflect.Struct {
		// 如果 T 是结构体指针，创建一个新的结构体实例
		defaultValue = reflect.New(tType.Elem()).Interface().(T)
	} else {
		// 如果 T 不是结构体指针，则返回错误或执行其他操作
		// 因为我们强制要求 T 必须是结构体指针
		return defaultValue, false
	}

	ok := true
	if needInit {
		ok = defaultValue.Init()
	}
	return defaultValue, ok
}

func (bc *BaseConfig[T]) Init() bool {
	// 如果 VersionMap 为 nil，则初始化为一个空映射，并添加一个键值对
	if bc.Versions == nil {
		bc.Versions = make([]T, 0, 1)
	}
	if len(bc.Versions) == 0 {
		defaultValue, ok := bc.newT(true)
		if !ok {
			return false
		}
		bc.Latest = defaultValue.GetVersion()
		bc.Versions = append(bc.Versions, defaultValue)
		bc.versionIndex = map[VersionType]int{bc.Latest: 0}
	} else {
		flag := false // 是否需要初始化 version
		if bc.Latest == 0 {
			flag = true
		}
		bc.versionIndex = make(map[VersionType]int, len(bc.Versions))
		for index, version := range bc.Versions {
			if !version.Init() {
				return false
			}
			bc.versionIndex[version.GetVersion()] = index
			if flag && version.GetVersion() > bc.Latest {
				// 将 Latest 赋值为最新的版本
				bc.Latest = version.GetVersion()
			}
		}
	}

	if bc.ID == "" {
		bc.ID = id.GenerateID().String()
	}
	if bc.Name == "" {
		bc.Name = bc.ID
	}
	if bc.ID2VersionMap == nil {
		bc.ID2VersionMap = make(map[string]VersionType)
	}
	if bc.envArray == nil {
		bc.envArray = make([]*Env, 0)
	}

	return true
}

func (bc *BaseConfig[T]) GetName() string {
	return bc.Name
}

func (bc *BaseConfig[T]) GetID() string {
	return bc.ID
}

func (bc *BaseConfig[T]) GetType() string {
	return bc.Type
}

func (bc *BaseConfig[T]) SetType(typ string) {
	bc.Type = typ
}

func (bc *BaseConfig[T]) GetEnv() []*Env {
	return bc.envArray
}

func (bc *BaseConfig[T]) SetEnv(env []*Env) {
	bc.envArray = env
}

func (bc *BaseConfig[T]) GetVersionByUserID(id string) ConfigVersion {
	if versionID, ok := bc.ID2VersionMap[id]; ok {
		if version, ok := bc.versionIndex[versionID]; ok {
			return bc.Versions[version]
		}
	}
	if version, ok := bc.versionIndex[bc.Latest]; ok {
		return bc.Versions[version]
	}
	return bc.Versions[len(bc.Versions)-1]
}

func (bc *BaseConfig[T]) NewVersion(change []byte) (VersionType, error) {
	t, ok := bc.newT(false)
	if !ok {
		return NilVersion, ErrCannotCreateInstance
	}
	err := json.Unmarshal(change, t)
	if err != nil {
		return NilVersion, err
	}

	max := StartVersion
	for _, value := range bc.Versions {
		if value.GetVersion() > max {
			max = value.GetVersion()
		}
	}
	t.SetVersion(max + 1)
	if !t.Init() {
		return NilVersion, ErrCannotCreateInstance
	}

	if bc.loader != nil {
		type TmpStruct struct {
			BaseConfig[T]
		}
		tmpBC := &TmpStruct{
			BaseConfig[T]{
				ID: bc.ID,
			},
		}
		err = bc.loader.First(tmpBC)
		if err != nil {
			return NilVersion, err
		}
		tmpBC.Versions = bc.Versions
		tmpBC.versionIndex = bc.versionIndex
		tmpBC.Latest = bc.Latest
		tmpBC.addVersion(t)

		err = bc.loader.Save(tmpBC)
		if err != nil {
			return NilVersion, err
		}
	}
	bc.addVersion(t)

	return t.GetVersion(), nil
}

func (bc *BaseConfig[T]) addVersion(t T) {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()
	bc.Versions = append(bc.Versions, t)
	bc.versionIndex[t.GetVersion()] = len(bc.Versions) - 1
}

func (bc *BaseConfig[T]) verifyVersion(version VersionType) bool {
	_, ok := bc.versionIndex[version]
	return ok
}

func (bc *BaseConfig[T]) deleteVersion(version VersionType) {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()
	if index, ok := bc.versionIndex[version]; ok {
		bc.Versions = append(bc.Versions[:index], bc.Versions[index+1:]...)
		delete(bc.versionIndex, version)
	}
}

func (bc *BaseConfig[T]) DeleteVersion(version VersionType) bool {
	if bc.loader != nil {
		type TmpStruct struct {
			BaseConfig[T]
		}
		tmpBC := &TmpStruct{
			BaseConfig[T]{
				ID: bc.ID,
			},
		}
		err := bc.loader.First(tmpBC)
		if err != nil {
			logx.Errorf("cannot first database for DeleteVersion, error: %s", err.Error())
			return false
		}
		tmpBC.Versions = bc.Versions
		tmpBC.versionIndex = bc.versionIndex
		tmpBC.Latest = bc.Latest
		tmpBC.deleteVersion(version)

		err = bc.loader.Save(tmpBC)
		if err != nil {
			logx.Errorf("cannot save to database for DeleteVersion, error: %s", err.Error())
			return false
		}
	}

	bc.deleteVersion(version)
	return true
}

func (bc *BaseConfig[T]) ActiveVersion(version VersionType) bool {
	if !bc.verifyVersion(version) {
		return false
	}

	if bc.loader != nil {
		type TmpStruct struct {
			BaseConfig[T]
		}
		tmpBC := &TmpStruct{
			BaseConfig[T]{
				ID: bc.ID,
			},
		}
		err := bc.loader.First(tmpBC)
		if err != nil {
			logx.Errorf("cannot first database for ActiveVersion, error: %s", err.Error())
			return false
		}
		tmpBC.Versions = bc.Versions
		tmpBC.versionIndex = bc.versionIndex
		tmpBC.Latest = version

		err = bc.loader.Save(tmpBC)
		if err != nil {
			logx.Errorf("cannot save to database for ActiveVersion, error: %s", err.Error())
			return false
		}
	}

	bc.Latest = version
	return true
}

func (bc *BaseConfig[T]) SetLoader(loader Loader) {
	bc.loader = loader
}

type Config interface {
	Init() bool
	SetLoader(loader Loader)
	GetName() string
	GetID() string
	GetType() string
	SetType(string)
	GetEnv() []*Env
	SetEnv(env []*Env)
	GetVersionByUserID(id string) ConfigVersion
	NewVersion(change []byte) (VersionType, error)
	DeleteVersion(version VersionType) bool
	ActiveVersion(version VersionType) bool
}

type ConfigVersion interface {
	Init() bool
	GetVersion() VersionType // 获取版本号
	SetVersion(version VersionType)
	GetDescription() string // 获取备注、描述
}

type BaseConfigVersion struct {
	Version     VersionType `json:"version"`
	Description string      `json:"description"`
}

func (bcv *BaseConfigVersion) GetVersion() VersionType {
	return bcv.Version
}
func (bcv *BaseConfigVersion) SetVersion(version VersionType) {
	bcv.Version = version
}

func (bcv *BaseConfigVersion) GetDescription() string {
	return bcv.Description
}

func (bcv *BaseConfigVersion) Init() bool {
	if bcv.Version == NilVersion {
		bcv.Version = StartVersion
	}

	return true
}
