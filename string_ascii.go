package goja

import (
	"fmt"
	"io"
	"math"
	"reflect"
	"strconv"
	"strings"
)

type asciiString string

type asciiRuneReader struct {
	s   asciiString
	pos int
}
// 读取一个字符，转了int32
func (rr *asciiRuneReader) ReadRune() (r rune, size int, err error) {
	if rr.pos < len(rr.s) {
		r = rune(rr.s[rr.pos])
		size = 1
		rr.pos++
	} else {
		err = io.EOF
	}
	return
}
// 从start切片构造一个string
func (s asciiString) reader(start int) io.RuneReader {
	return &asciiRuneReader{
		s: s[start:],
	}
}

// ss must be trimmed 字符串转int
func strToInt(ss string) (int64, error) {
	if ss == "" {
		return 0, nil
	}
	if ss == "-0" {
		return 0, strconv.ErrSyntax
	}
	if len(ss) > 2 {
		switch ss[:2] {
		case "0x", "0X":
			i, _ := strconv.ParseInt(ss[2:], 16, 64)
			return i, nil
		case "0b", "0B":
			i, _ := strconv.ParseInt(ss[2:], 2, 64)
			return i, nil
		case "0o", "0O":
			i, _ := strconv.ParseInt(ss[2:], 8, 64)
			return i, nil
		}
	}
	return strconv.ParseInt(ss, 10, 64)
}
//字符串转int
func (s asciiString) _toInt() (int64, error) {
	return strToInt(strings.TrimSpace(string(s)))
}
// 检查是否存在范围错误
func isRangeErr(err error) bool {
	if err, ok := err.(*strconv.NumError); ok {
		return err.Err == strconv.ErrRange
	}
	return false
}
// 字符串转float
func (s asciiString) _toFloat() (float64, error) {
	ss := strings.TrimSpace(string(s))
	if ss == "" {
		return 0, nil
	}
	if ss == "-0" {
		var f float64
		return -f, nil
	}
	f, err := strconv.ParseFloat(ss, 64)
	if isRangeErr(err) {
		err = nil
	}
	return f, err
}
// 字符串转int64
func (s asciiString) ToInteger() int64 {
	if s == "" {
		return 0
	}
	if s == "Infinity" || s == "+Infinity" {
		return math.MaxInt64
	}
	if s == "-Infinity" {
		return math.MinInt64
	}
	i, err := s._toInt()
	if err != nil {
		f, err := s._toFloat()
		if err == nil {
			return int64(f)
		}
	}
	return i
}
// 字符串返回
func (s asciiString) ToString() valueString {
	return s
}
// 字符串返回
func (s asciiString) String() string {
	return string(s)
}
// 字符串转float
func (s asciiString) ToFloat() float64 {
	if s == "" {
		return 0
	}
	if s == "Infinity" || s == "+Infinity" {
		return math.Inf(1)
	}
	if s == "-Infinity" {
		return math.Inf(-1)
	}
	f, err := s._toFloat()
	if err != nil {
		i, err := s._toInt()
		if err == nil {
			return float64(i)
		}
		f = math.NaN()
	}
	return f
}
// 字符串不为空就是true
func (s asciiString) ToBoolean() bool {
	return s != ""
}
// 字符串转数字，可能是int或float
func (s asciiString) ToNumber() Value {
	if s == "" {
		return intToValue(0)
	}
	if s == "Infinity" || s == "+Infinity" {
		return _positiveInf
	}
	if s == "-Infinity" {
		return _negativeInf
	}

	if i, err := s._toInt(); err == nil {
		return intToValue(i)
	}

	if f, err := s._toFloat(); err == nil {
		return floatToValue(f)
	}

	return _NaN
}
// 字符串转对象
func (s asciiString) ToObject(r *Runtime) *Object {
	return r._newString(s)
}
// 比较字符串是否相等
func (s asciiString) SameAs(other Value) bool {
	if otherStr, ok := other.(asciiString); ok {
		return s == otherStr
	}
	return false
}
// 比较字符串是否相等
func (s asciiString) Equals(other Value) bool {
	if o, ok := other.(asciiString); ok {
		return s == o
	}

	if o, ok := other.assertInt(); ok {
		if o1, e := s._toInt(); e == nil {
			return o1 == o
		}
		return false
	}

	if o, ok := other.assertFloat(); ok {
		return s.ToFloat() == o
	}

	if o, ok := other.(valueBool); ok {
		if o1, e := s._toFloat(); e == nil {
			return o1 == o.ToFloat()
		}
		return false
	}

	if o, ok := other.(*Object); ok {
		return s.Equals(o.self.toPrimitive())
	}
	return false
}
// 比较字符串是否相等
func (s asciiString) StrictEquals(other Value) bool {
	if otherStr, ok := other.(asciiString); ok {
		return s == otherStr
	}
	return false
}
// 字符串不是整数
func (s asciiString) assertInt() (int64, bool) {
	return 0, false
}
// 字符串不是float
func (s asciiString) assertFloat() (float64, bool) {
	return 0, false
}
// 字符串是字符串
func (s asciiString) assertString() (valueString, bool) {
	return s, true
}
// 返回对应的对象
func (s asciiString) baseObject(r *Runtime) *Object {
	ss := r.stringSingleton
	ss.value = s
	ss.setLength()
	return ss.val
}
// 返回idx位置的字符
func (s asciiString) charAt(idx int64) rune {
	return rune(s[idx])
}
// 返回字符串长度
func (s asciiString) length() int64 {
	return int64(len(s))
}
// 拼接两个字符串
func (s asciiString) concat(other valueString) valueString {
	switch other := other.(type) {
	case asciiString:
		b := make([]byte, len(s)+len(other))
		copy(b, s)
		copy(b[len(s):], other)
		return asciiString(b)
		//return asciiString(string(s) + string(other))
	case unicodeString:
		b := make([]uint16, len(s)+len(other))
		for i := 0; i < len(s); i++ {
			b[i] = uint16(s[i])
		}
		copy(b[len(s):], other)
		return unicodeString(b)
	default:
		panic(fmt.Errorf("Unknown string type: %T", other))
	}
}
// 截取字符串
func (s asciiString) substring(start, end int64) valueString {
	return asciiString(s[start:end])
}
// 比较字符串
func (s asciiString) compareTo(other valueString) int {
	switch other := other.(type) {
	case asciiString:
		return strings.Compare(string(s), string(other))
	case unicodeString:
		return strings.Compare(string(s), other.String())
	default:
		panic(fmt.Errorf("Unknown string type: %T", other))
	}
}
//Index返回s中substr的第一个实例的索引，如果s中不存在substr，则返回-1。
func (s asciiString) index(substr valueString, start int64) int64 {
	if substr, ok := substr.(asciiString); ok {
		p := int64(strings.Index(string(s[start:]), string(substr)))
		if p >= 0 {
			return p + start
		}
	}
	return -1
}
//LastIndex返回s中substr的最后一个实例的索引，如果s中不存在substr，则返回-1。
func (s asciiString) lastIndex(substr valueString, pos int64) int64 {
	if substr, ok := substr.(asciiString); ok {
		end := pos + int64(len(substr))
		var ss string
		if end > int64(len(s)) {
			ss = string(s)
		} else {
			ss = string(s[:end])
		}
		return int64(strings.LastIndex(ss, string(substr)))
	}
	return -1
}
// 转小写
func (s asciiString) toLower() valueString {
	return asciiString(strings.ToLower(string(s)))
}
// 转大写
func (s asciiString) toUpper() valueString {
	return asciiString(strings.ToUpper(string(s)))
}
// 过滤空格
func (s asciiString) toTrimmedUTF8() string {
	return strings.TrimSpace(string(s))
}
// 导出字符串
func (s asciiString) Export() interface{} {
	return string(s)
}
// 导出字符串类型
func (s asciiString) ExportType() reflect.Type {
	return reflectTypeString
}
