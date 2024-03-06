package data

import (
	"encoding/json"
	"fmt"
	"github.com/finishy1995/go-library/storage"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"reflect"
)

type LoaderType uint8

const (
	// Memory 加载器，什么都不做
	Memory LoaderType = iota
	// StorageWithFileInit 根据 storage 中的内容进行加载，同时会阅读指定文件中的内容。如果文件中有数据库中不存在的条目，则强制更新入数据库
	StorageWithFileInit
)

var (
	defaultLoaderInstance = &memoryLoader{}
)

type LoaderConfig struct {
	Type      LoaderType
	FilePath  string
	Expr      string
	TableName string
}

type Loader interface {
	First(value Config) error
	Save(value Config) error
	Create(value Config) error
	Find(value Config) ([]Config, error)
}

func getLoader(st storage.Storage, config *LoaderConfig) Loader {
	if config.Type == StorageWithFileInit {
		return newStorageWithFileInitLoader(st, config.FilePath, config.Expr, config.TableName)
	}
	return defaultLoaderInstance
}

type memoryLoader struct {
}

func (ml *memoryLoader) First(value Config) error {
	return nil
}

func (ml *memoryLoader) Save(_ Config) error {
	return nil
}

func (ml *memoryLoader) Create(_ Config) error {
	return nil
}

func (ml *memoryLoader) Find(_ Config) ([]Config, error) {
	return nil, nil
}

type storageWithFileInitLoader struct {
	st        storage.Storage
	filePath  string
	expr      string
	tableName string
}

func newStorageWithFileInitLoader(st storage.Storage, filePath string, expr string, tableName string) *storageWithFileInitLoader {
	return &storageWithFileInitLoader{st: st, filePath: filePath, expr: expr, tableName: tableName}
}

func (sl *storageWithFileInitLoader) getStorage() storage.Storage {
	if sl.st != nil {
		return sl.st
	}
	panic(ErrMustInit)
}

func (sl *storageWithFileInitLoader) First(value Config) error {
	return sl.getStorage().First(value, sl.tableName, value.GetID())
}

func (sl *storageWithFileInitLoader) Save(value Config) error {
	return sl.getStorage().Save(value, sl.tableName)
}

func (sl *storageWithFileInitLoader) Create(value Config) error {
	// 检查value是否为结构体指针
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return ErrParamInvalid
	}

	// 解引用获得结构体
	structVal := val.Elem().Interface()

	return retry3Times(func() error {
		return sl.getStorage().Create(structVal, sl.tableName)
	})
}

func (sl *storageWithFileInitLoader) Find(value Config) ([]Config, error) {
	// 检查value是否为结构体指针
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return nil, ErrParamInvalid
	}

	// 创建一个空的结构体切片
	sliceType := reflect.SliceOf(reflect.TypeOf(value).Elem())
	sliceValue := reflect.MakeSlice(sliceType, 0, 0)
	sliceForFileValue := reflect.MakeSlice(sliceType, 0, 0)

	// 创建一个指向这个切片的指针
	slicePtr := reflect.New(sliceType).Elem()
	slicePtr.Set(sliceValue)
	sliceForFilePtr := reflect.New(sliceType).Elem()
	sliceForFilePtr.Set(sliceForFileValue)

	err := retry3Times(func() error {
		return sl.getStorage().Find(slicePtr.Addr().Interface(), sl.tableName, -1, sl.expr)
	})
	if err != nil {
		return nil, err
	}

	existingIDs := make(map[string]bool)
	results := make([]Config, 0, slicePtr.Len())
	for i := 0; i < slicePtr.Len(); i++ {
		item := slicePtr.Index(i).Addr().Interface().(Config)
		existingIDs[item.GetID()] = true
		item.SetLoader(sl)
		results = append(results, item)
	}

	// 加载文件内容
	err = loadFile(sl.filePath, sliceForFilePtr.Addr().Interface())
	if err != nil {
		logx.Alert(fmt.Sprintf("File [%s] cannot load, error: %s", sl.filePath, err.Error()))
		err = nil
	} else {
		for i := 0; i < sliceForFilePtr.Len(); i++ {
			item := sliceForFilePtr.Index(i).Addr().Interface().(Config)
			if value.GetType() != "" && item.GetType() != value.GetType() {
				continue
			}

			if _, exists := existingIDs[item.GetID()]; !exists {
				item.SetLoader(sl)
				results = append(results, item)
				// 尝试写数据库
				err = retry3Times(func() error {
					return sl.Create(item)
				})
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return results, err
}

func loadFile(filePath string, sliceForFilePtr interface{}) error {
	// 读取文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// 解析 JSON 到 sliceForFilePtr 指向的切片中
	err = json.Unmarshal(data, sliceForFilePtr)
	if err != nil {
		return err
	}

	return nil
}
