package gowl

type Data map[string]interface{}

func (d *Data) Set(key string, value interface{}) {
	if *d == nil {
		*d = make(Data)
	}
	(*d)[key] = value
}

func (d Data) Get(key string) interface{} {
	return d[key]
}

func (d Data) Has(key string) bool {
	_, ok := d[key]
	return ok
}

func (d Data) Lookup(key string) (value interface{}, ok bool) {
	value, ok = d[key]
	return
}

func (d Data) Delete(key string) {
	delete(d, key)
}

func (d Data) GetBool(key string) bool {
	v, _ := d.Get(key).(bool)
	return v
}

func (d Data) GetInt(key string) int {
	v, _ := d.Get(key).(int)
	return v
}

func (d Data) GetInt8(key string) int8 {
	v, _ := d.Get(key).(int8)
	return v
}

func (d Data) GetInt16(key string) int16 {
	v, _ := d.Get(key).(int16)
	return v
}

func (d Data) GetInt32(key string) int32 {
	v, _ := d.Get(key).(int32)
	return v
}

func (d Data) GetInt64(key string) int64 {
	v, _ := d.Get(key).(int64)
	return v
}

func (d Data) GetUint(key string) uint {
	v, _ := d.Get(key).(uint)
	return v
}

func (d Data) GetUint8(key string) uint8 {
	v, _ := d.Get(key).(uint8)
	return v
}

func (d Data) GetUint16(key string) uint16 {
	v, _ := d.Get(key).(uint16)
	return v
}

func (d Data) GetUint32(key string) uint32 {
	v, _ := d.Get(key).(uint32)
	return v
}

func (d Data) GetUint64(key string) uint64 {
	v, _ := d.Get(key).(uint64)
	return v
}

func (d Data) GetFloat32(key string) float32 {
	v, _ := d.Get(key).(float32)
	return v
}

func (d Data) GetFloat64(key string) float64 {
	v, _ := d.Get(key).(float64)
	return v
}

func (d Data) GetString(key string) string {
	v, _ := d.Get(key).(string)
	return v
}

func (d Data) GetByteSlice(key string) []byte {
	v, _ := d.Get(key).([]byte)
	return v
}

func (d Data) GetStringSlice(key string) []string {
	v, _ := d.Get(key).([]string)
	return v
}

func (d Data) GetIntSlice(key string) []int {
	v, _ := d.Get(key).([]int)
	return v
}

func (d Data) GetStringMap(key string) StringMap {
	v, _ := d.Get(key).(StringMap)
	return v
}

func (d Data) GetData(key string) Data {
	v, _ := d.Get(key).(Data)
	return v
}
