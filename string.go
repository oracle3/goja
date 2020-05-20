package goja

import (
	"io"
	"strconv"
	"unicode/utf16"
	"unicode/utf8"
)

const (
	__proto__ = "__proto__"
)

var (
	stringTrue         valueString = asciiString("true")
	stringFalse        valueString = asciiString("false")
	stringNull         valueString = asciiString("null")
	stringUndefined    valueString = asciiString("undefined")
	stringObjectC      valueString = asciiString("object")
	stringFunction     valueString = asciiString("function")
	stringBoolean      valueString = asciiString("boolean")
	stringString       valueString = asciiString("string")
	stringNumber       valueString = asciiString("number")
	stringNaN          valueString = asciiString("NaN")
	stringInfinity                 = asciiString("Infinity")
	stringPlusInfinity             = asciiString("+Infinity")
	stringNegInfinity              = asciiString("-Infinity")
	stringEmpty        valueString = asciiString("")
	string__proto__    valueString = asciiString(__proto__)

	stringError          valueString = asciiString("Error")
	stringTypeError      valueString = asciiString("TypeError")
	stringReferenceError valueString = asciiString("ReferenceError")
	stringSyntaxError    valueString = asciiString("SyntaxError")
	stringRangeError     valueString = asciiString("RangeError")
	stringEvalError      valueString = asciiString("EvalError")
	stringURIError       valueString = asciiString("URIError")
	stringGoError        valueString = asciiString("GoError")

	stringObjectNull      valueString = asciiString("[object Null]")
	stringObjectObject    valueString = asciiString("[object Object]")
	stringObjectUndefined valueString = asciiString("[object Undefined]")
	stringGlobalObject    valueString = asciiString("Global Object")
	stringInvalidDate     valueString = asciiString("Invalid Date")
)

type valueString interface {
	Value
	charAt(int64) rune
	length() int64
	concat(valueString) valueString
	substring(start, end int64) valueString
	compareTo(valueString) int
	reader(start int) io.RuneReader
	index(valueString, int64) int64
	lastIndex(valueString, int64) int64
	toLower() valueString
	toUpper() valueString
	toTrimmedUTF8() string
}

type stringObject struct {
	baseObject
	value      valueString
	length     int64
	lengthProp valueProperty
}
// 构造一个Unicode字符串
func newUnicodeString(s string) valueString {
	return unicodeString(utf16.Encode([]rune(s)))
}
// 构造一个字符串，可能是Unicode或asc格式
func newStringValue(s string) valueString {
	for _, chr := range s {
		if chr >= utf8.RuneSelf {
			return newUnicodeString(s)
		}
	}
	return asciiString(s)
}
// 字符串对象初始化
func (s *stringObject) init() {
	s.baseObject.init()
	s.setLength()
}
// 设置长度属性
func (s *stringObject) setLength() {
	if s.value != nil {
		s.length = s.value.length()
	}
	s.lengthProp.value = intToValue(s.length)
	s._put("length", &s.lengthProp)
}
// 获取指定位置的字符
func (s *stringObject) get(n Value) Value {
	if idx := toIdx(n); idx >= 0 && idx < s.length {
		return s.getIdx(idx)
	}
	return s.baseObject.get(n)
}
// 获取指定位置的字符
func (s *stringObject) getStr(name string) Value {
	if i := strToIdx(name); i >= 0 && i < s.length {
		return s.getIdx(i)
	}
	return s.baseObject.getStr(name)
}
// 获取指定位置的字符或者属性
func (s *stringObject) getPropStr(name string) Value {
	if i := strToIdx(name); i >= 0 && i < s.length {
		return s.getIdx(i)
	}
	return s.baseObject.getPropStr(name)
}
// 获取指定位置的字符或者属性
func (s *stringObject) getProp(n Value) Value {
	if i := toIdx(n); i >= 0 && i < s.length {
		return s.getIdx(i)
	}
	return s.baseObject.getProp(n)
}
// 获取指定位置的字符或者属性
func (s *stringObject) getOwnProp(name string) Value {
	if i := strToIdx(name); i >= 0 && i < s.length {
		val := s.getIdx(i)
		return &valueProperty{
			value:      val,
			enumerable: true,
		}
	}

	return s.baseObject.getOwnProp(name)
}
// 获取指定位置的字符
func (s *stringObject) getIdx(idx int64) Value {
	return s.value.substring(idx, idx+1)
}
// 对字符串put会异常
func (s *stringObject) put(n Value, val Value, throw bool) {
	if i := toIdx(n); i >= 0 && i < s.length {
		s.val.runtime.typeErrorResult(throw, "Cannot assign to read only property '%d' of a String", i)
		return
	}

	s.baseObject.put(n, val, throw)
}
// 对字符串put会异常
func (s *stringObject) putStr(name string, val Value, throw bool) {
	if i := strToIdx(name); i >= 0 && i < s.length {
		s.val.runtime.typeErrorResult(throw, "Cannot assign to read only property '%d' of a String", i)
		return
	}

	s.baseObject.putStr(name, val, throw)
}
// 对字符串定义属性会异常
func (s *stringObject) defineOwnProperty(n Value, descr propertyDescr, throw bool) bool {
	if i := toIdx(n); i >= 0 && i < s.length {
		s.val.runtime.typeErrorResult(throw, "Cannot redefine property: %d", i)
		return false
	}

	return s.baseObject.defineOwnProperty(n, descr, throw)
}

type stringPropIter struct {
	str         valueString // separate, because obj can be the singleton
	obj         *stringObject
	idx, length int64
	recursive   bool
}
// 获取下一个
func (i *stringPropIter) next() (propIterItem, iterNextFunc) {
	if i.idx < i.length {
		name := strconv.FormatInt(i.idx, 10)
		i.idx++
		return propIterItem{name: name, enumerable: _ENUM_TRUE}, i.next
	}

	return i.obj.baseObject._enumerate(i.recursive)()
}

func (s *stringObject) _enumerate(recursive bool) iterNextFunc {
	return (&stringPropIter{
		str:       s.value,
		obj:       s,
		length:    s.length,
		recursive: recursive,
	}).next
}
// 构造枚举循环
func (s *stringObject) enumerate(all, recursive bool) iterNextFunc {
	return (&propFilterIter{
		wrapped: s._enumerate(recursive),
		all:     all,
		seen:    make(map[string]bool),
	}).next
}
// 不允许删除字符串
func (s *stringObject) deleteStr(name string, throw bool) bool {
	if i := strToIdx(name); i >= 0 && i < s.length {
		s.val.runtime.typeErrorResult(throw, "Cannot delete property '%d' of a String", i)
		return false
	}

	return s.baseObject.deleteStr(name, throw)
}
// 不允许删除字符串
func (s *stringObject) delete(n Value, throw bool) bool {
	if i := toIdx(n); i >= 0 && i < s.length {
		s.val.runtime.typeErrorResult(throw, "Cannot delete property '%d' of a String", i)
		return false
	}

	return s.baseObject.delete(n, throw)
}
// 获取指定位置的属性都是true
func (s *stringObject) hasOwnProperty(n Value) bool {
	if i := toIdx(n); i >= 0 && i < s.length {
		return true
	}
	return s.baseObject.hasOwnProperty(n)
}
// 获取指定位置的属性都是true
func (s *stringObject) hasOwnPropertyStr(name string) bool {
	if i := strToIdx(name); i >= 0 && i < s.length {
		return true
	}
	return s.baseObject.hasOwnPropertyStr(name)
}
