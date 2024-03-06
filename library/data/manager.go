package data

import "github.com/finishy1995/go-library/storage"

type Manager struct {
	poolMap map[string]*pool
	st      storage.Storage
}

// InitFunc Item 对象的初始化函数。需要一个 value 就能初始化
type InitFunc func() Item

var (
	managerInstance = NewManager(nil)
)

func NewManager(st storage.Storage) *Manager {
	return &Manager{poolMap: make(map[string]*pool), st: st}
}

func (m *Manager) SetStorage(st storage.Storage) {
	m.st = st
}

// SetupPool 注册“名称、初始化函数、需要的数据、加载数据的配置”，批量生成指定数量的 Item，并全部塞入对应的 pool
func (m *Manager) SetupPool(name string, initFunc InitFunc, value Config, config *LoaderConfig) error {
	loader := getLoader(m.st, config)
	results, err := loader.Find(value)
	if err != nil {
		return err
	}

	p, ok := m.poolMap[name]
	if !ok {
		p = &pool{cache: NewCache[Item](), loader: loader}
		m.poolMap[name] = p
	}
	return p.load(initFunc, results)
}

func (m *Manager) Output2File(path string) {
	// TODO: 导出最新配置到文件中，迁移用 （需要用时再实现）
}

func (m *Manager) Save(cacheName string, value Config) error {
	if c, ok := m.poolMap[cacheName]; ok {
		return c.loader.Save(value)
	}
	return ErrParamInvalid
}

func (m *Manager) Create(cacheName string, value Config) error {
	if c, ok := m.poolMap[cacheName]; ok {
		return c.loader.Create(value)
	}
	return ErrParamInvalid
}

func (m *Manager) GetAllItems(cacheName string) map[string]Item {
	if c, ok := m.poolMap[cacheName]; ok {
		return c.cache.Items()
	}
	return nil
}

func (m *Manager) GetOneItem(cacheName, itemID string) (Item, bool) {
	if c, ok := m.poolMap[cacheName]; ok {
		return c.cache.Get(itemID)
	}
	return nil, false
}

// SetupPool 设置池
func SetupPool(name string, initFunc InitFunc, value Config, config *LoaderConfig) error {
	return managerInstance.SetupPool(name, initFunc, value, config)
}

// GetAllItems 获取指定名称的缓存池中的所有项目
func GetAllItems(cacheName string) map[string]Item {
	return managerInstance.GetAllItems(cacheName)
}

// GetOneItem 获取缓存池中指定 ID 的单个项目
func GetOneItem(cacheName, itemID string) (Item, bool) {
	return managerInstance.GetOneItem(cacheName, itemID)
}

func GetAllTypeItems(cacheName string, typ string) map[string]Item {
	items := GetAllItems(cacheName)
	results := make(map[string]Item)
	for key, value := range items {
		if value.GetConfig().GetType() == typ {
			results[key] = value
		}
	}
	return results
}

func SetStorage(st storage.Storage) {
	managerInstance.SetStorage(st)
}

func ConfigCreate(cacheName string, value Config) error {
	return managerInstance.Create(cacheName, value)
}

func ConfigSave(cacheName string, value Config) error {
	return managerInstance.Save(cacheName, value)
}
