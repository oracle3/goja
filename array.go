package goja

import (
	"math"
	"reflect"
	"strconv"
)

type arrayIterObject struct {
	baseObject
	obj     *Object
	nextIdx int64
	kind    iterationKind
}

func (ai *arrayIterObject) next() Value {
	if ai.obj == nil {
		return ai.val.runtime.createIterResultObject(_undefined, true)
	}
	l := toLength(ai.obj.self.getStr("length", nil))
	index := ai.nextIdx
	if index >= l {
		ai.obj = nil
		return ai.val.runtime.createIterResultObject(_undefined, true)
	}
	ai.nextIdx++
	idxVal := intToValue(index)
	if ai.kind == iterationKindKey {
		return ai.val.runtime.createIterResultObject(idxVal, false)
	}
	elementValue := ai.obj.self.get(idxVal, nil)
	var result Value
	if ai.kind == iterationKindValue {
		result = elementValue
	} else {
		result = ai.val.runtime.newArrayValues([]Value{idxVal, elementValue})
	}
	return ai.val.runtime.createIterResultObject(result, false)
}

func (r *Runtime) createArrayIterator(iterObj *Object, kind iterationKind) Value {
	o := &Object{runtime: r}

	ai := &arrayIterObject{
		obj:  iterObj,
		kind: kind,
	}
	ai.class = classArrayIterator
	ai.val = o
	ai.extensible = true
	o.self = ai
	ai.prototype = r.global.ArrayIteratorPrototype
	ai.init()

	return o
}

type arrayObject struct {
	baseObject
	values         []Value
	length         int64
	objCount       int64
	propValueCount int
	lengthProp     valueProperty
}

func (a *arrayObject) init() {
	a.baseObject.init()
	a.lengthProp.writable = true

	a._put("length", &a.lengthProp)
}

func (a *arrayObject) _setLengthInt(l int64, throw bool) bool {
	if l >= 0 && l <= math.MaxUint32 {
		ret := true
		if l <= a.length {
			if a.propValueCount > 0 {
				// Slow path
				var s int64
				if a.length < int64(len(a.values)) {
					s = a.length - 1
				} else {
					s = int64(len(a.values)) - 1
				}
				for i := s; i >= l; i-- {
					if prop, ok := a.values[i].(*valueProperty); ok {
						if !prop.configurable {
							l = i + 1
							ret = false
							break
						}
						a.propValueCount--
					}
				}
			}
		}
		if l <= int64(len(a.values)) {
			if l >= 16 && l < int64(cap(a.values))>>2 {
				ar := make([]Value, l)
				copy(ar, a.values)
				a.values = ar
			} else {
				ar := a.values[l:len(a.values)]
				for i := range ar {
					ar[i] = nil
				}
				a.values = a.values[:l]
			}
		}
		a.length = l
		if !ret {
			a.val.runtime.typeErrorResult(throw, "Cannot redefine property: length")
		}
		return ret
	}
	panic(a.val.runtime.newError(a.val.runtime.global.RangeError, "Invalid array length"))
}

func (a *arrayObject) setLengthInt(l int64, throw bool) bool {
	if l == a.length {
		return true
	}
	if !a.lengthProp.writable {
		a.val.runtime.typeErrorResult(throw, "length is not writable")
		return false
	}
	return a._setLengthInt(l, throw)
}

func (a *arrayObject) setLength(v Value, throw bool) bool {
	l, ok := toIntIgnoreNegZero(v)
	if ok && l == a.length {
		return true
	}
	if !a.lengthProp.writable {
		a.val.runtime.typeErrorResult(throw, "length is not writable")
		return false
	}
	if ok {
		return a._setLengthInt(l, throw)
	}
	panic(a.val.runtime.newError(a.val.runtime.global.RangeError, "Invalid array length"))
}

func (a *arrayObject) getIdx(idx int64) Value {
	if idx >= 0 && idx < int64(len(a.values)) {
		return a.values[idx]
	}
	return nil
}

func (a *arrayObject) sortLen() int64 {
	return int64(len(a.values))
}

func (a *arrayObject) sortGet(i int64) Value {
	v := a.values[i]
	if p, ok := v.(*valueProperty); ok {
		v = p.get(a.val)
	}
	return v
}

func (a *arrayObject) swap(i, j int64) {
	a.values[i], a.values[j] = a.values[j], a.values[i]
}

func toIdx(v Value) (idx int64) {
	idx = -1
	if idxVal, ok1 := v.(valueInt); ok1 {
		idx = int64(idxVal)
	} else {
		if i, err := strconv.ParseInt(v.String(), 10, 64); err == nil {
			idx = i
		}
	}
	if idx >= 0 && idx < math.MaxUint32 {
		return
	}
	return -1
}

func strToIdx(s string) (idx int64) {
	idx = -1
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		idx = i
	}

	if idx >= 0 && idx < math.MaxUint32 {
		return
	}
	return -1
}

func (a *arrayObject) get(p Value, receiver Value) Value {
	return a.getWithOwnProp(a.getOwnProp(p), p, receiver)
}

func (a *arrayObject) getStr(name string, receiver Value) Value {
	return a.getStrWithOwnProp(a.getOwnPropStr(name), name, receiver)
}

func (a *arrayObject) getLengthProp() Value {
	a.lengthProp.value = intToValue(a.length)
	return &a.lengthProp
}

func (a *arrayObject) getOwnProp(n Value) Value {
	if s, ok := n.(*valueSymbol); ok {
		return a.getOwnPropSym(s)
	}
	if idx := toIdx(n); idx >= 0 {
		return a.getIdx(idx)
	}
	s := n.String()
	if s == "length" {
		return a.getLengthProp()
	}

	return a.baseObject.getOwnPropStr(s)
}

func (a *arrayObject) getOwnPropStr(name string) Value {
	if i := strToIdx(name); i >= 0 {
		if i >= 0 && i < int64(len(a.values)) {
			return a.values[i]
		}
	}
	if name == "length" {
		return a.getLengthProp()
	}
	return a.baseObject.getOwnPropStr(name)
}

func (a *arrayObject) setIdx(idx int64, val Value, throw bool, origNameStr string, origName Value) {
	var prop Value
	if idx < int64(len(a.values)) {
		prop = a.values[idx]
	}

	if prop == nil {
		if proto := a.prototype; proto != nil {
			// we know it's foreign because prototype loops are not allowed
			var b bool
			if origName != nil {
				b = proto.self.setForeign(origName, val, a.val, throw)
			} else {
				b = proto.self.setForeignStr(origNameStr, val, a.val, throw)
			}
			if b {
				return
			}
		}
		// new property
		if !a.extensible {
			a.val.runtime.typeErrorResult(throw, "Cannot add property %d, object is not extensible", idx)
			return
		} else {
			if idx >= a.length {
				if !a.setLengthInt(idx+1, throw) {
					return
				}
			}
			if idx >= int64(len(a.values)) {
				if !a.expand(idx) {
					a.val.self.(*sparseArrayObject).add(idx, val)
					return
				}
			}
			a.objCount++
		}
	} else {
		if prop, ok := prop.(*valueProperty); ok {
			if !prop.isWritable() {
				a.val.runtime.typeErrorResult(throw)
				return
			}
			prop.set(a.val, val)
			return
		}
	}
	a.values[idx] = val
}

func (a *arrayObject) setOwn(n Value, val Value, throw bool) {
	if s, ok := n.(*valueSymbol); ok {
		a.setOwnSym(s, val, throw)
		return
	}
	if idx := toIdx(n); idx >= 0 {
		a.setIdx(idx, val, throw, "", n)
	} else {
		name := n.String()
		if name == "length" {
			a.setLength(val, throw)
		} else {
			a.baseObject.setOwnStr(name, val, throw)
		}
	}
}

func (a *arrayObject) setOwnStr(name string, val Value, throw bool) {
	if idx := strToIdx(name); idx >= 0 {
		a.setIdx(idx, val, throw, name, nil)
	} else {
		if name == "length" {
			a.setLength(val, throw)
		} else {
			a.baseObject.setOwnStr(name, val, throw)
		}
	}
}

func (a *arrayObject) setForeign(name Value, val, receiver Value, throw bool) bool {
	return a._setForeign(name, a.getOwnProp(name), val, receiver, throw)
}

func (a *arrayObject) setForeignStr(name string, val, receiver Value, throw bool) bool {
	return a._setForeignStr(name, a.getOwnPropStr(name), val, receiver, throw)
}

type arrayPropIter struct {
	a   *arrayObject
	idx int
}

func (i *arrayPropIter) next() (propIterItem, iterNextFunc) {
	for i.idx < len(i.a.values) {
		name := strconv.Itoa(i.idx)
		prop := i.a.values[i.idx]
		i.idx++
		if prop != nil {
			return propIterItem{name: name, value: prop}, i.next
		}
	}

	return i.a.baseObject.enumerateUnfiltered()()
}

func (a *arrayObject) enumerateUnfiltered() iterNextFunc {
	return (&arrayPropIter{
		a: a,
	}).next
}

func (a *arrayObject) ownKeys(all bool, accum []Value) []Value {
	for i, prop := range a.values {
		name := strconv.Itoa(i)
		if prop != nil {
			if !all {
				if prop, ok := prop.(*valueProperty); ok && !prop.enumerable {
					continue
				}
			}
			accum = append(accum, asciiString(name))
		}
	}
	return a.baseObject.ownKeys(all, accum)
}

func (a *arrayObject) hasOwnProperty(n Value) bool {
	if s, ok := n.(*valueSymbol); ok {
		return a.hasOwnSym(s)
	}
	if idx := toIdx(n); idx >= 0 {
		return idx < int64(len(a.values)) && a.values[idx] != nil
	} else {
		return a.baseObject.hasOwnProperty(n)
	}
}

func (a *arrayObject) hasOwnPropertyStr(name string) bool {
	if idx := strToIdx(name); idx >= 0 {
		return idx < int64(len(a.values)) && a.values[idx] != nil
	} else {
		return a.baseObject.hasOwnPropertyStr(name)
	}
}

func (a *arrayObject) expand(idx int64) bool {
	targetLen := idx + 1
	if targetLen > int64(len(a.values)) {
		if targetLen < int64(cap(a.values)) {
			a.values = a.values[:targetLen]
		} else {
			if idx > 4096 && (a.objCount == 0 || idx/a.objCount > 10) {
				//log.Println("Switching standard->sparse")
				sa := &sparseArrayObject{
					baseObject:     a.baseObject,
					length:         a.length,
					propValueCount: a.propValueCount,
				}
				sa.setValues(a.values, a.objCount+1)
				sa.val.self = sa
				sa.init()
				sa.lengthProp.writable = a.lengthProp.writable
				return false
			} else {
				// Use the same algorithm as in runtime.growSlice
				newcap := int64(cap(a.values))
				doublecap := newcap + newcap
				if targetLen > doublecap {
					newcap = targetLen
				} else {
					if len(a.values) < 1024 {
						newcap = doublecap
					} else {
						for newcap < targetLen {
							newcap += newcap / 4
						}
					}
				}
				newValues := make([]Value, targetLen, newcap)
				copy(newValues, a.values)
				a.values = newValues
			}
		}
	}
	return true
}

func (r *Runtime) defineArrayLength(prop *valueProperty, descr PropertyDescriptor, setter func(Value, bool) bool, throw bool) bool {
	ret := true

	if descr.Configurable == FLAG_TRUE || descr.Enumerable == FLAG_TRUE || descr.Getter != nil || descr.Setter != nil {
		ret = false
		goto Reject
	}

	if newLen := descr.Value; newLen != nil {
		ret = setter(newLen, false)
	} else {
		ret = true
	}

	if descr.Writable != FLAG_NOT_SET {
		w := descr.Writable.Bool()
		if prop.writable {
			prop.writable = w
		} else {
			if w {
				ret = false
				goto Reject
			}
		}
	}

Reject:
	if !ret {
		r.typeErrorResult(throw, "Cannot redefine property: length")
	}

	return ret
}

func (a *arrayObject) defineOwnProperty(n Value, descr PropertyDescriptor, throw bool) bool {
	if s, ok := n.(*valueSymbol); ok {
		return a.defineOwnPropertySym(s, descr, throw)
	}
	if idx := toIdx(n); idx >= 0 {
		var existing Value
		if idx < int64(len(a.values)) {
			existing = a.values[idx]
		}
		prop, ok := a.baseObject._defineOwnProperty(n.String(), existing, descr, throw)
		if ok {
			if idx >= a.length {
				if !a.setLengthInt(idx+1, throw) {
					return false
				}
			}
			if a.expand(idx) {
				a.values[idx] = prop
				a.objCount++
				if _, ok := prop.(*valueProperty); ok {
					a.propValueCount++
				}
			} else {
				a.val.self.(*sparseArrayObject).add(idx, prop)
			}
		}
		return ok
	} else {
		name := n.String()
		if name == "length" {
			return a.val.runtime.defineArrayLength(&a.lengthProp, descr, a.setLength, throw)
		}
		return a.defineOwnPropertyStr(name, descr, throw)
	}
}

func (a *arrayObject) _deleteProp(idx int64, throw bool) bool {
	if idx < int64(len(a.values)) {
		if v := a.values[idx]; v != nil {
			if p, ok := v.(*valueProperty); ok {
				if !p.configurable {
					a.val.runtime.typeErrorResult(throw, "Cannot delete property '%d' of %s", idx, a.val.toString())
					return false
				}
				a.propValueCount--
			}
			a.values[idx] = nil
			a.objCount--
		}
	}
	return true
}

func (a *arrayObject) delete(n Value, throw bool) bool {
	if idx := toIdx(n); idx >= 0 {
		return a._deleteProp(idx, throw)
	}
	return a.baseObject.delete(n, throw)
}

func (a *arrayObject) deleteStr(name string, throw bool) bool {
	if idx := strToIdx(name); idx >= 0 {
		return a._deleteProp(idx, throw)
	}
	return a.baseObject.deleteStr(name, throw)
}

func (a *arrayObject) export() interface{} {
	arr := make([]interface{}, a.length)
	for i, v := range a.values {
		if v != nil {
			arr[i] = v.Export()
		}
	}

	return arr
}

func (a *arrayObject) exportType() reflect.Type {
	return reflectTypeArray
}

func (a *arrayObject) setValuesFromSparse(items []sparseArrayItem, newMaxIdx int64) {
	a.values = make([]Value, newMaxIdx+1)
	for _, item := range items {
		a.values[item.idx] = item.value
	}
	a.objCount = int64(len(items))
}
