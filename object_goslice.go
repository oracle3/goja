package goja

import (
	"reflect"
	"strconv"
)

type objectGoSlice struct {
	baseObject
	data            *[]interface{}
	lengthProp      valueProperty
	sliceExtensible bool
}
// 初始化
func (o *objectGoSlice) init() {
	o.baseObject.init()
	o.class = classArray
	o.prototype = o.val.runtime.global.ArrayPrototype
	o.lengthProp.writable = o.sliceExtensible
	o._setLen()
	o.baseObject._put("length", &o.lengthProp)
}
// 设置数据长度
func (o *objectGoSlice) _setLen() {
	o.lengthProp.value = intToValue(int64(len(*o.data)))
}
// 获取idx位置的值
func (o *objectGoSlice) getIdx(idx int64) Value {
	if idx < int64(len(*o.data)) {
		return o.val.runtime.ToValue((*o.data)[idx])
	}
	return nil
}
// 获取n位置的值
func (o *objectGoSlice) _get(n Value) Value {
	if idx := toIdx(n); idx >= 0 {
		return o.getIdx(idx)
	}
	return nil
}
// 获取name位置的值
func (o *objectGoSlice) _getStr(name string) Value {
	if idx := strToIdx(name); idx >= 0 {
		return o.getIdx(idx)
	}
	return nil
}
// 获取n位置的值
func (o *objectGoSlice) get(n Value) Value {
	if v := o._get(n); v != nil {
		return v
	}
	return o.baseObject._getStr(n.String())
}
// 获取name位置的值
func (o *objectGoSlice) getStr(name string) Value {
	if v := o._getStr(name); v != nil {
		return v
	}
	return o.baseObject._getStr(name)
}
// 获取n位置的值
func (o *objectGoSlice) getProp(n Value) Value {
	if v := o._get(n); v != nil {
		return v
	}
	return o.baseObject.getPropStr(n.String())
}
// 获取name位置的值
func (o *objectGoSlice) getPropStr(name string) Value {
	if v := o._getStr(name); v != nil {
		return v
	}
	return o.baseObject.getPropStr(name)
}
// 获取name位置的值
func (o *objectGoSlice) getOwnProp(name string) Value {
	if v := o._getStr(name); v != nil {
		return &valueProperty{
			value:      v,
			writable:   true,
			enumerable: true,
		}
	}
	return o.baseObject.getOwnProp(name)
}
// 空间扩展
func (o *objectGoSlice) grow(size int64) {
	newcap := int64(cap(*o.data))
	if newcap < size {
		// Use the same algorithm as in runtime.growSlice 使用与runtime.growSlice中相同的算法
		doublecap := newcap + newcap
		if size > doublecap {
			newcap = size
		} else {
			if len(*o.data) < 1024 {
				newcap = doublecap
			} else {
				for newcap < size {
					newcap += newcap / 4
				}
			}
		}

		n := make([]interface{}, size, newcap)
		copy(n, *o.data)
		*o.data = n
	} else {
		*o.data = (*o.data)[:size]
	}
	o._setLen()
}
// 在idx位置保存v值
func (o *objectGoSlice) putIdx(idx int64, v Value, throw bool) {
	if idx >= int64(len(*o.data)) {
		if !o.sliceExtensible {
			o.val.runtime.typeErrorResult(throw, "Cannot extend Go slice")
			return
		}
		o.grow(idx + 1)
	}
	(*o.data)[idx] = v.Export()
}
// 在n位置保存val值
func (o *objectGoSlice) put(n Value, val Value, throw bool) {
	if idx := toIdx(n); idx >= 0 {
		o.putIdx(idx, val, throw)
		return
	}
	// TODO: length
	o.baseObject.put(n, val, throw)
}
// 在name位置保存val值
func (o *objectGoSlice) putStr(name string, val Value, throw bool) {
	if idx := strToIdx(name); idx >= 0 {
		o.putIdx(idx, val, throw)
		return
	}
	// TODO: length
	o.baseObject.putStr(name, val, throw)
}
// 判断n值是否存在
func (o *objectGoSlice) _has(n Value) bool {
	if idx := toIdx(n); idx >= 0 {
		return idx < int64(len(*o.data))
	}
	return false
}
// 判断name值是否存在
func (o *objectGoSlice) _hasStr(name string) bool {
	if idx := strToIdx(name); idx >= 0 {
		return idx < int64(len(*o.data))
	}
	return false
}
// 判断n值是否存在
func (o *objectGoSlice) hasProperty(n Value) bool {
	if o._has(n) {
		return true
	}
	return o.baseObject.hasProperty(n)
}
// 判断name值是否存在
func (o *objectGoSlice) hasPropertyStr(name string) bool {
	if o._hasStr(name) {
		return true
	}
	return o.baseObject.hasPropertyStr(name)
}
// 判断n值是否存在
func (o *objectGoSlice) hasOwnProperty(n Value) bool {
	if o._has(n) {
		return true
	}
	return o.baseObject.hasOwnProperty(n)
}
// 判断name值是否存在
func (o *objectGoSlice) hasOwnPropertyStr(name string) bool {
	if o._hasStr(name) {
		return true
	}
	return o.baseObject.hasOwnPropertyStr(name)
}
// 在name位置保存val值
func (o *objectGoSlice) _putProp(name string, value Value, writable, enumerable, configurable bool) Value {
	o.putStr(name, value, false)
	return value
}
// 在n位置保存descr.Value值
func (o *objectGoSlice) defineOwnProperty(n Value, descr propertyDescr, throw bool) bool {
	if idx := toIdx(n); idx >= 0 {
		if !o.val.runtime.checkHostObjectPropertyDescr(n.String(), descr, throw) {
			return false
		}
		val := descr.Value
		if val == nil {
			val = _undefined
		}
		o.putIdx(idx, val, throw)
		return true
	}
	return o.baseObject.defineOwnProperty(n, descr, throw)
}
//将对象的所有元素连接成一个字符串并返回这个字符串。
func (o *objectGoSlice) toPrimitiveNumber() Value {
	return o.toPrimitiveString()
}
//将对象的所有元素连接成一个字符串并返回这个字符串。
func (o *objectGoSlice) toPrimitiveString() Value {
	return o.val.runtime.arrayproto_join(FunctionCall{
		This: o.val,
	})
}
//将对象的所有元素连接成一个字符串并返回这个字符串。
func (o *objectGoSlice) toPrimitive() Value {
	return o.toPrimitiveString()
}
// 删除name位置的值
func (o *objectGoSlice) deleteStr(name string, throw bool) bool {
	if idx := strToIdx(name); idx >= 0 && idx < int64(len(*o.data)) {
		(*o.data)[idx] = nil
		return true
	}
	return o.baseObject.deleteStr(name, throw)
}
// 删除name位置的值
func (o *objectGoSlice) delete(name Value, throw bool) bool {
	if idx := toIdx(name); idx >= 0 && idx < int64(len(*o.data)) {
		(*o.data)[idx] = nil
		return true
	}
	return o.baseObject.delete(name, throw)
}

type goslicePropIter struct {
	o          *objectGoSlice
	recursive  bool
	idx, limit int
}
// 遍历下一个
func (i *goslicePropIter) next() (propIterItem, iterNextFunc) {
	if i.idx < i.limit && i.idx < len(*i.o.data) {
		name := strconv.Itoa(i.idx)
		i.idx++
		return propIterItem{name: name, enumerable: _ENUM_TRUE}, i.next
	}

	if i.recursive {
		return i.o.prototype.self._enumerate(i.recursive)()
	}

	return propIterItem{}, nil
}
// 构造枚举迭代
func (o *objectGoSlice) enumerate(all, recursive bool) iterNextFunc {
	return (&propFilterIter{
		wrapped: o._enumerate(recursive),
		all:     all,
		seen:    make(map[string]bool),
	}).next

}

func (o *objectGoSlice) _enumerate(recursive bool) iterNextFunc {
	return (&goslicePropIter{
		o:         o,
		recursive: recursive,
		limit:     len(*o.data),
	}).next
}
// 导出当前数据
func (o *objectGoSlice) export() interface{} {
	return *o.data
}
// 当前数据类型是map
func (o *objectGoSlice) exportType() reflect.Type {
	return reflectTypeArray
}
// 判断两个是否相等
func (o *objectGoSlice) equal(other objectImpl) bool {
	if other, ok := other.(*objectGoSlice); ok {
		return o.data == other.data
	}
	return false
}

func (o *objectGoSlice) sortLen() int64 {
	return int64(len(*o.data))
}

func (o *objectGoSlice) sortGet(i int64) Value {
	return o.get(intToValue(i))
}

func (o *objectGoSlice) swap(i, j int64) {
	ii := intToValue(i)
	jj := intToValue(j)
	x := o.get(ii)
	y := o.get(jj)

	o.put(ii, y, false)
	o.put(jj, x, false)
}
