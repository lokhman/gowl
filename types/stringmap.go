package types

type StringMap map[string]string

func (sm StringMap) Set(key string, value string) {
	sm[key] = value
}

func (sm StringMap) Get(key string) string {
	return sm[key]
}

func (sm StringMap) Has(key string) bool {
	_, ok := sm[key]
	return ok
}

func (sm StringMap) Lookup(key string) (value string, ok bool) {
	value, ok = sm[key]
	return
}

func (sm StringMap) Delete(key string) {
	delete(sm, key)
}

func (sm StringMap) Copy() StringMap {
	smc := make(StringMap, len(sm))
	for key, value := range smc {
		smc[key] = value
	}
	return smc
}
