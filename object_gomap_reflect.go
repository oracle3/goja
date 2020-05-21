package goja

import "reflect"

type objectGoMapReflect struct {
	objectGoReflect

	keyType, valueType reflect.Type
}
// 初始化
func (o *objectGoMapReflect) init() {
	o.objectGoReflect.init()
	o.keyType = o.value.Type().Key()
	o.valueType = o.value.Type().Elem()
}
// go的key值转js的key值
func (o *objectGoMapReflect) toKey(n Value) reflect.Value {
	key, err := o.val.runtime.toReflectValue(n, o.keyType)
	if err != nil {
		o.val.runtime.typeErrorResult(true, "map key conversion error: %v", err)
		panic("unreachable")
	}
	return key
}
// go的key值转js的key值
func (o *objectGoMapReflect) strToKey(name string) reflect.Value {
	if o.keyType.Kind() == reflect.String {
		return reflect.ValueOf(name).Convert(o.keyType)
	}
	return o.toKey(newStringValue(name))
}
// 由n获取对应的value
func (o *objectGoMapReflect) _get(n Value) Value {
	if v := o.value.MapIndex(o.toKey(n)); v.IsValid() {
		return o.val.runtime.ToValue(v.Interface())
	}

	return nil
}
// 由name获取对应的value
func (o *objectGoMapReflect) _getStr(name string) Value {
	if v := o.value.MapIndex(o.strToKey(name)); v.IsValid() {
		return o.val.runtime.ToValue(v.Interface())
	}

	return nil
}
// 由n获取对应的value
func (o *objectGoMapReflect) get(n Value) Value {
	if v := o._get(n); v != nil {
		return v
	}
	return o.objectGoReflect.get(n)
}
// 由name获取对应的value
func (o *objectGoMapReflect) getStr(name string) Value {
	if v := o._getStr(name); v != nil {
		return v
	}
	return o.objectGoReflect.getStr(name)
}
// 由n获取对应的value
func (o *objectGoMapReflect) getProp(n Value) Value {
	return o.get(n)
}
// 由name获取对应的value
func (o *objectGoMapReflect) getPropStr(name string) Value {
	return o.getStr(name)
}
// 由name获取对应的value
func (o *objectGoMapReflect) getOwnProp(name string) Value {
	if v := o._getStr(name); v != nil {
		return &valueProperty{
			value:      v,
			writable:   true,
			enumerable: true,
		}
	}
	return o.objectGoReflect.getOwnProp(name)
}
// go的value转js的value
func (o *objectGoMapReflect) toValue(val Value, throw bool) (reflect.Value, bool) {
	v, err := o.val.runtime.toReflectValue(val, o.valueType)
	if err != nil {
		o.val.runtime.typeErrorResult(throw, "map value conversion error: %v", err)
		return reflect.Value{}, false
	}

	return v, true
}
// 保存js的kv
func (o *objectGoMapReflect) put(key, val Value, throw bool) {
	k := o.toKey(key)
	v, ok := o.toValue(val, throw)
	if !ok {
		return
	}
	o.value.SetMapIndex(k, v)
}
// 保存js的kv
func (o *objectGoMapReflect) putStr(name string, val Value, throw bool) {
	k := o.strToKey(name)
	v, ok := o.toValue(val, throw)
	if !ok {
		return
	}
	o.value.SetMapIndex(k, v)
}
// 保存js的kv
func (o *objectGoMapReflect) _putProp(name string, value Value, writable, enumerable, configurable bool) Value {
	o.putStr(name, value, true)
	return value
}
// 保存js的kv
func (o *objectGoMapReflect) defineOwnProperty(n Value, descr propertyDescr, throw bool) bool {
	name := n.String()
	if !o.val.runtime.checkHostObjectPropertyDescr(name, descr, throw) {
		return false
	}

	o.put(n, descr.Value, throw)
	return true
}
// 判断是否存在name的值
func (o *objectGoMapReflect) hasOwnPropertyStr(name string) bool {
	return o.value.MapIndex(o.strToKey(name)).IsValid()
}
// 判断是否存在n的值
func (o *objectGoMapReflect) hasOwnProperty(n Value) bool {
	return o.value.MapIndex(o.toKey(n)).IsValid()
}
// 判断是否存在n的值
func (o *objectGoMapReflect) hasProperty(n Value) bool {
	if o.hasOwnProperty(n) {
		return true
	}
	return o.objectGoReflect.hasProperty(n)
}
// 判断是否存在name的值
func (o *objectGoMapReflect) hasPropertyStr(name string) bool {
	if o.hasOwnPropertyStr(name) {
		return true
	}
	return o.objectGoReflect.hasPropertyStr(name)
}
// 删除n的值
func (o *objectGoMapReflect) delete(n Value, throw bool) bool {
	o.value.SetMapIndex(o.toKey(n), reflect.Value{})
	return true
}
// 删除name的值
func (o *objectGoMapReflect) deleteStr(name string, throw bool) bool {
	o.value.SetMapIndex(o.strToKey(name), reflect.Value{})
	return true
}

type gomapReflectPropIter struct {
	o         *objectGoMapReflect
	keys      []reflect.Value
	idx       int
	recursive bool
}
// 遍历下一个
func (i *gomapReflectPropIter) next() (propIterItem, iterNextFunc) {
	for i.idx < len(i.keys) {
		key := i.keys[i.idx]
		v := i.o.value.MapIndex(key)
		i.idx++
		if v.IsValid() {
			return propIterItem{name: key.String(), enumerable: _ENUM_TRUE}, i.next
		}
	}

	if i.recursive {
		return i.o.objectGoReflect._enumerate(true)()
	}

	return propIterItem{}, nil
}

func (o *objectGoMapReflect) _enumerate(recusrive bool) iterNextFunc {
	r := &gomapReflectPropIter{
		o:         o,
		keys:      o.value.MapKeys(),
		recursive: recusrive,
	}
	return r.next
}
// 构造枚举迭代
func (o *objectGoMapReflect) enumerate(all, recursive bool) iterNextFunc {
	return (&propFilterIter{
		wrapped: o._enumerate(recursive),
		all:     all,
		seen:    make(map[string]bool),
	}).next
}
// 判断是否相等
func (o *objectGoMapReflect) equal(other objectImpl) bool {
	if other, ok := other.(*objectGoMapReflect); ok {
		return o.value.Interface() == other.value.Interface()
	}
	return false
}
