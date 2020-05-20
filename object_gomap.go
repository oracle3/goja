package goja

import (
	"reflect"
	"strconv"
)

type objectGoMapSimple struct {
	baseObject
	data map[string]interface{}
}
// 初始化
func (o *objectGoMapSimple) init() {
	o.baseObject.init()
	o.prototype = o.val.runtime.global.ObjectPrototype
	o.class = classObject
	o.extensible = true
}
// 获取map中n的值
func (o *objectGoMapSimple) _get(n Value) Value {
	return o._getStr(n.String())
}
// 获取map中name的值
func (o *objectGoMapSimple) _getStr(name string) Value {
	v, exists := o.data[name]
	if !exists {
		return nil
	}
	return o.val.runtime.ToValue(v)
}
// 获取map中n的值
func (o *objectGoMapSimple) get(n Value) Value {
	return o.getStr(n.String())
}
// 获取map中n的值
func (o *objectGoMapSimple) getProp(n Value) Value {
	return o.getPropStr(n.String())
}
// 获取map中n的值
func (o *objectGoMapSimple) getPropStr(name string) Value {
	if v := o._getStr(name); v != nil {
		return v
	}
	return o.baseObject.getPropStr(name)
}
// 获取map中name的值
func (o *objectGoMapSimple) getStr(name string) Value {
	if v := o._getStr(name); v != nil {
		return v
	}
	return o.baseObject._getStr(name)
}
// 获取map中n的值
func (o *objectGoMapSimple) getOwnProp(name string) Value {
	if v := o._getStr(name); v != nil {
		return v
	}
	return o.baseObject.getOwnProp(name)
}
// 保存n和val到map中
func (o *objectGoMapSimple) put(n Value, val Value, throw bool) {
	o.putStr(n.String(), val, throw)
}
// 判断map中是否包含name
func (o *objectGoMapSimple) _hasStr(name string) bool {
	_, exists := o.data[name]
	return exists
}
// 判断map中是否包含n
func (o *objectGoMapSimple) _has(n Value) bool {
	return o._hasStr(n.String())
}
// 保存name和value到map中
func (o *objectGoMapSimple) putStr(name string, val Value, throw bool) {
	if o.extensible || o._hasStr(name) {
		o.data[name] = val.Export()
	} else {
		o.val.runtime.typeErrorResult(throw, "Host object is not extensible")
	}
}
// 判断map中是否包含n
func (o *objectGoMapSimple) hasProperty(n Value) bool {
	if o._has(n) {
		return true
	}
	return o.baseObject.hasProperty(n)
}
// 判断map中是否包含name
func (o *objectGoMapSimple) hasPropertyStr(name string) bool {
	if o._hasStr(name) {
		return true
	}
	return o.baseObject.hasOwnPropertyStr(name)
}
// 判断map中是否包含n
func (o *objectGoMapSimple) hasOwnProperty(n Value) bool {
	return o._has(n)
}
// 判断map中是否包含name
func (o *objectGoMapSimple) hasOwnPropertyStr(name string) bool {
	return o._hasStr(name)
}
// 保存name和value到map中
func (o *objectGoMapSimple) _putProp(name string, value Value, writable, enumerable, configurable bool) Value {
	o.putStr(name, value, false)
	return value
}
// 保存name和descr.Value到map中
func (o *objectGoMapSimple) defineOwnProperty(name Value, descr propertyDescr, throw bool) bool {
	if descr.Getter != nil || descr.Setter != nil {
		o.val.runtime.typeErrorResult(throw, "Host objects do not support accessor properties")
		return false
	}
	o.put(name, descr.Value, throw)
	return true
}

/*
func (o *objectGoMapSimple) toPrimitiveNumber() Value {
	return o.toPrimitiveString()
}

func (o *objectGoMapSimple) toPrimitiveString() Value {
	return stringObjectObject
}

func (o *objectGoMapSimple) toPrimitive() Value {
	return o.toPrimitiveString()
}

func (o *objectGoMapSimple) assertCallable() (call func(FunctionCall) Value, ok bool) {
	return nil, false
}
*/
// 删除map中的name
func (o *objectGoMapSimple) deleteStr(name string, throw bool) bool {
	delete(o.data, name)
	return true
}
// 删除map中的name
func (o *objectGoMapSimple) delete(name Value, throw bool) bool {
	return o.deleteStr(name.String(), throw)
}

type gomapPropIter struct {
	o         *objectGoMapSimple
	propNames []string
	recursive bool
	idx       int
}
// 下一个
func (i *gomapPropIter) next() (propIterItem, iterNextFunc) {
	for i.idx < len(i.propNames) {
		name := i.propNames[i.idx]
		i.idx++
		if _, exists := i.o.data[name]; exists {
			return propIterItem{name: name, enumerable: _ENUM_TRUE}, i.next
		}
	}

	if i.recursive {
		return i.o.prototype.self._enumerate(true)()
	}

	return propIterItem{}, nil
}
// 构造枚举迭代
func (o *objectGoMapSimple) enumerate(all, recursive bool) iterNextFunc {
	return (&propFilterIter{
		wrapped: o._enumerate(recursive),
		all:     all,
		seen:    make(map[string]bool),
	}).next
}

func (o *objectGoMapSimple) _enumerate(recursive bool) iterNextFunc {
	propNames := make([]string, len(o.data))
	i := 0
	for key, _ := range o.data {
		propNames[i] = key
		i++
	}
	return (&gomapPropIter{
		o:         o,
		propNames: propNames,
		recursive: recursive,
	}).next
}
// 导出map的数据
func (o *objectGoMapSimple) export() interface{} {
	return o.data
}
// 导出map类型
func (o *objectGoMapSimple) exportType() reflect.Type {
	return reflectTypeMap
}
// 判断map是否相等
func (o *objectGoMapSimple) equal(other objectImpl) bool {
	if other, ok := other.(*objectGoMapSimple); ok {
		return o == other
	}
	return false
}
// 返回map的数据数量
func (o *objectGoMapSimple) sortLen() int64 {
	return int64(len(o.data))
}
// 获取name为i的数据
func (o *objectGoMapSimple) sortGet(i int64) Value {
	return o.getStr(strconv.FormatInt(i, 10))
}
// i，j的数据互换
func (o *objectGoMapSimple) swap(i, j int64) {
	ii := strconv.FormatInt(i, 10)
	jj := strconv.FormatInt(j, 10)
	x := o.getStr(ii)
	y := o.getStr(jj)

	o.putStr(ii, y, false)
	o.putStr(jj, x, false)
}
