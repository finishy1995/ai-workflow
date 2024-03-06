package data

type pool struct {
	cache  *Cache[Item]
	loader Loader
}

func (p *pool) load(initFunc InitFunc, value []Config) error {
	for _, val := range value {
		if !val.Init() {
			return ErrCannotCreateInstance
		}
		item := initFunc()
		if item == nil {
			return ErrCannotCreateInstance
		}
		if !item.Init(val) {
			return ErrCannotCreateInstance
		}

		p.cache.Set(item.GetID(), item)
	}

	return nil
}
