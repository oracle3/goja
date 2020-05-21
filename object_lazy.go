package goja

import "reflect"

type lazyObject struct {
	val    *Object
	create func(*Object) objectImpl
}
// 返回类名
func (o *lazyObject) className() string {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.className()
}
// 获取n对应的值
func (o *lazyObject) get(n Value) Value {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.get(n)
}
// 获取n对应的值
func (o *lazyObject) getProp(n Value) Value {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.getProp(n)
}
// 获取name对应的值
func (o *lazyObject) getPropStr(name string) Value {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.getPropStr(name)
}
// 获取name对应的值
func (o *lazyObject) getStr(name string) Value {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.getStr(name)
}
// 获取name对应的值
func (o *lazyObject) getOwnProp(name string) Value {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.getOwnProp(name)
}
// 保存n位置值val
func (o *lazyObject) put(n Value, val Value, throw bool) {
	obj := o.create(o.val)
	o.val.self = obj
	obj.put(n, val, throw)
}
// 保存name位置值val
func (o *lazyObject) putStr(name string, val Value, throw bool) {
	obj := o.create(o.val)
	o.val.self = obj
	obj.putStr(name, val, throw)
}
// 判断n位置是否有值
func (o *lazyObject) hasProperty(n Value) bool {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.hasProperty(n)
}
// 判断name位置是否有值
func (o *lazyObject) hasPropertyStr(name string) bool {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.hasPropertyStr(name)
}
// 判断n位置是否有值
func (o *lazyObject) hasOwnProperty(n Value) bool {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.hasOwnProperty(n)
}
// 判断name位置是否有值
func (o *lazyObject) hasOwnPropertyStr(name string) bool {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.hasOwnPropertyStr(name)
}
// 保存name位置值为value
func (o *lazyObject) _putProp(name string, value Value, writable, enumerable, configurable bool) Value {
	obj := o.create(o.val)
	o.val.self = obj
	return obj._putProp(name, value, writable, enumerable, configurable)
}
// 保存name位置值为descr.value
func (o *lazyObject) defineOwnProperty(name Value, descr propertyDescr, throw bool) bool {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.defineOwnProperty(name, descr, throw)
}
// 把对象数据转为number
func (o *lazyObject) toPrimitiveNumber() Value {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.toPrimitiveNumber()
}
// 把对象数据转为字符串
func (o *lazyObject) toPrimitiveString() Value {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.toPrimitiveString()
}
// 把对象数据转为字符串
func (o *lazyObject) toPrimitive() Value {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.toPrimitive()
}

func (o *lazyObject) assertCallable() (call func(FunctionCall) Value, ok bool) {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.assertCallable()
}
// 删除name位置的值
func (o *lazyObject) deleteStr(name string, throw bool) bool {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.deleteStr(name, throw)
}
// 删除name位置的值
func (o *lazyObject) delete(name Value, throw bool) bool {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.delete(name, throw)
}

func (o *lazyObject) proto() *Object {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.proto()
}

func (o *lazyObject) hasInstance(v Value) bool {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.hasInstance(v)
}

func (o *lazyObject) isExtensible() bool {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.isExtensible()
}

func (o *lazyObject) preventExtensions() {
	obj := o.create(o.val)
	o.val.self = obj
	obj.preventExtensions()
}
// 构造枚举迭代
func (o *lazyObject) enumerate(all, recusrive bool) iterNextFunc {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.enumerate(all, recusrive)
}

func (o *lazyObject) _enumerate(recursive bool) iterNextFunc {
	obj := o.create(o.val)
	o.val.self = obj
	return obj._enumerate(recursive)
}
// 导出数据
func (o *lazyObject) export() interface{} {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.export()
}
// 导出数据类型
func (o *lazyObject) exportType() reflect.Type {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.exportType()
}
// 判断是否相等
func (o *lazyObject) equal(other objectImpl) bool {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.equal(other)
}
// 获取数据长度
func (o *lazyObject) sortLen() int64 {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.sortLen()
}
// 获取数据长度
func (o *lazyObject) sortGet(i int64) Value {
	obj := o.create(o.val)
	o.val.self = obj
	return obj.sortGet(i)
}
// i,j位置数据互换
func (o *lazyObject) swap(i, j int64) {
	obj := o.create(o.val)
	o.val.self = obj
	obj.swap(i, j)
}
