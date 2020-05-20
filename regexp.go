package goja

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"regexp"
	"unicode/utf16"
	"unicode/utf8"
)

type regexpPattern interface {
	FindSubmatchIndex(valueString, int) []int
	FindAllSubmatchIndex(valueString, int) [][]int
	FindAllSubmatchIndexUTF8(string, int) [][]int
	FindAllSubmatchIndexASCII(string, int) [][]int
	MatchString(valueString) bool
}

type regexp2Wrapper regexp2.Regexp
type regexpWrapper regexp.Regexp

type regexpObject struct {
	baseObject
	pattern regexpPattern
	source  valueString

	global, multiline, ignoreCase bool
}
//在输入字符串中搜索正则表达式匹配项
func (r *regexp2Wrapper) FindSubmatchIndex(s valueString, start int) (result []int) {
	wrapped := (*regexp2.Regexp)(r)
	var match *regexp2.Match
	var err error
	switch s := s.(type) {
	case asciiString:
		match, err = wrapped.FindStringMatch(string(s)[start:])
	case unicodeString:
		match, err = wrapped.FindRunesMatch(utf16.Decode(s[start:]))
	default:
		panic(fmt.Errorf("Unknown string type: %T", s))
	}
	if err != nil {
		return
	}

	if match == nil {
		return
	}
	groups := match.Groups()

	result = make([]int, 0, len(groups)<<1)
	for _, group := range groups {
		if len(group.Captures) > 0 {
			result = append(result, group.Index, group.Index+group.Length)
		} else {
			result = append(result, -1, 0)
		}
	}
	return
}
//在输入字符串中搜索正则表达式匹配项
func (r *regexp2Wrapper) FindAllSubmatchIndexUTF8(s string, n int) [][]int {
	wrapped := (*regexp2.Regexp)(r)
	if n < 0 {
		n = len(s) + 1
	}
	results := make([][]int, 0, n)

	idxMap := make([]int, 0, len(s))
	runes := make([]rune, 0, len(s))
	for pos, rr := range s {
		runes = append(runes, rr)
		idxMap = append(idxMap, pos)
	}
	idxMap = append(idxMap, len(s))

	match, err := wrapped.FindRunesMatch(runes)
	if err != nil {
		return nil
	}
	i := 0
	for match != nil && i < n {
		groups := match.Groups()

		result := make([]int, 0, len(groups)<<1)

		for _, group := range groups {
			if len(group.Captures) > 0 {
				result = append(result, idxMap[group.Index], idxMap[group.Index+group.Length])
			} else {
				result = append(result, -1, 0)
			}
		}

		results = append(results, result)
		match, err = wrapped.FindNextMatch(match)
		if err != nil {
			return nil
		}
		i++
	}
	return results
}
//在输入字符串中搜索正则表达式匹配项
func (r *regexp2Wrapper) FindAllSubmatchIndexASCII(s string, n int) [][]int {
	wrapped := (*regexp2.Regexp)(r)
	if n < 0 {
		n = len(s) + 1
	}
	results := make([][]int, 0, n)

	match, err := wrapped.FindStringMatch(s)
	if err != nil {
		return nil
	}
	i := 0
	for match != nil && i < n {
		groups := match.Groups()

		result := make([]int, 0, len(groups)<<1)

		for _, group := range groups {
			if len(group.Captures) > 0 {
				result = append(result, group.Index, group.Index+group.Length)
			} else {
				result = append(result, -1, 0)
			}
		}

		results = append(results, result)
		match, err = wrapped.FindNextMatch(match)
		if err != nil {
			return nil
		}
		i++
	}
	return results
}
//在输入字符串中搜索正则表达式匹配项
func (r *regexp2Wrapper) findAllSubmatchIndexUTF16(s unicodeString, n int) [][]int {
	wrapped := (*regexp2.Regexp)(r)
	if n < 0 {
		n = len(s) + 1
	}
	results := make([][]int, 0, n)

	rd := runeReaderReplace{s.reader(0)}
	posMap := make([]int, s.length()+1)
	curPos := 0
	curRuneIdx := 0
	runes := make([]rune, 0, s.length())
	for {
		rn, size, err := rd.ReadRune()
		if err != nil {
			break
		}
		runes = append(runes, rn)
		posMap[curRuneIdx] = curPos
		curRuneIdx++
		curPos += size
	}
	posMap[curRuneIdx] = curPos

	match, err := wrapped.FindRunesMatch(runes)
	if err != nil {
		return nil
	}
	for match != nil {
		groups := match.Groups()

		result := make([]int, 0, len(groups)<<1)

		for _, group := range groups {
			if len(group.Captures) > 0 {
				start := posMap[group.Index]
				end := posMap[group.Index+group.Length]
				result = append(result, start, end)
			} else {
				result = append(result, -1, 0)
			}
		}

		results = append(results, result)
		match, err = wrapped.FindNextMatch(match)
		if err != nil {
			return nil
		}
	}
	return results
}
//在输入字符串中搜索正则表达式匹配项
func (r *regexp2Wrapper) FindAllSubmatchIndex(s valueString, n int) [][]int {
	switch s := s.(type) {
	case asciiString:
		return r.FindAllSubmatchIndexASCII(string(s), n)
	case unicodeString:
		return r.findAllSubmatchIndexUTF16(s, n)
	default:
		panic("Unsupported string type")
	}
}
//如果字符串与正则表达式匹配，则MatchString返回true，如果发生超时则将设置错误
func (r *regexp2Wrapper) MatchString(s valueString) bool {
	wrapped := (*regexp2.Regexp)(r)

	switch s := s.(type) {
	case asciiString:
		matched, _ := wrapped.MatchString(string(s))
		return matched
	case unicodeString:
		matched, _ := wrapped.MatchRunes(utf16.Decode(s))
		return matched
	default:
		panic(fmt.Errorf("Unknown string type: %T", s))
	}
}
// FindReaderSubmatchIndex返回一个切片，其中包含索引对，该对对标识由RuneReader读取的文本的正则表达式的最左端匹配，
//以及其子表达式的匹配项（如果有），返回值nil表示不匹配。
func (r *regexpWrapper) FindSubmatchIndex(s valueString, start int) (result []int) {
	wrapped := (*regexp.Regexp)(r)
	return wrapped.FindReaderSubmatchIndex(runeReaderReplace{s.reader(start)})
}
// MatchReader报告RuneReader返回的文本是否包含正则表达式的任何匹配项。
func (r *regexpWrapper) MatchString(s valueString) bool {
	wrapped := (*regexp.Regexp)(r)
	return wrapped.MatchReader(runeReaderReplace{s.reader(0)})
}
// FindAllStringSubmatchIndex是FindStringSubmatchIndex的“全部”版本；
//它返回表达式的所有连续匹配的一部分，如包注释中的“全部”描述所定义。 返回值nil表示不匹配。
func (r *regexpWrapper) FindAllSubmatchIndex(s valueString, n int) [][]int {
	wrapped := (*regexp.Regexp)(r)
	switch s := s.(type) {
	case asciiString:
		return wrapped.FindAllStringSubmatchIndex(string(s), n)
	case unicodeString:
		return r.findAllSubmatchIndexUTF16(s, n)
	default:
		panic("Unsupported string type")
	}
}

func (r *regexpWrapper) FindAllSubmatchIndexUTF8(s string, n int) [][]int {
	wrapped := (*regexp.Regexp)(r)
	return wrapped.FindAllStringSubmatchIndex(s, n)
}

func (r *regexpWrapper) FindAllSubmatchIndexASCII(s string, n int) [][]int {
	return r.FindAllSubmatchIndexUTF8(s, n)
}

func (r *regexpWrapper) findAllSubmatchIndexUTF16(s unicodeString, n int) [][]int {
	wrapped := (*regexp.Regexp)(r)
	utf8Bytes := make([]byte, 0, len(s)*2)
	posMap := make(map[int]int)
	curPos := 0
	rd := runeReaderReplace{s.reader(0)}
	for {
		rn, size, err := rd.ReadRune()
		if err != nil {
			break
		}
		l := len(utf8Bytes)
		utf8Bytes = append(utf8Bytes, 0, 0, 0, 0)
		n := utf8.EncodeRune(utf8Bytes[l:], rn)
		utf8Bytes = utf8Bytes[:l+n]
		posMap[l] = curPos
		curPos += size
	}
	posMap[len(utf8Bytes)] = curPos

	rr := wrapped.FindAllSubmatchIndex(utf8Bytes, n)
	for _, res := range rr {
		for j, pos := range res {
			mapped, exists := posMap[pos]
			if !exists {
				panic("Unicode match is not on rune boundary")
			}
			res[j] = mapped
		}
	}
	return rr
}
// 将匹配结果转换为数组
func (r *regexpObject) execResultToArray(target valueString, result []int) Value {
	captureCount := len(result) >> 1
	valueArray := make([]Value, captureCount)
	matchIndex := result[0]
	lowerBound := matchIndex
	for index := 0; index < captureCount; index++ {
		offset := index << 1
		if result[offset] >= lowerBound {
			valueArray[index] = target.substring(int64(result[offset]), int64(result[offset+1]))
			lowerBound = result[offset]
		} else {
			valueArray[index] = _undefined
		}
	}
	match := r.val.runtime.newArrayValues(valueArray)
	match.self.putStr("input", target, false)
	match.self.putStr("index", intToValue(int64(matchIndex)), false)
	return match
}
// 执行正则表达式匹配
func (r *regexpObject) execRegexp(target valueString) (match bool, result []int) {
	lastIndex := int64(0)
	if p := r.getStr("lastIndex"); p != nil {
		lastIndex = p.ToInteger()
		if lastIndex < 0 {
			lastIndex = 0
		}
	}
	index := lastIndex
	if !r.global {
		index = 0
	}
	if index >= 0 && index <= target.length() {
		result = r.pattern.FindSubmatchIndex(target, int(index))
	}
	if result == nil {
		r.putStr("lastIndex", intToValue(0), true)
		return
	}
	match = true
	startIndex := index
	endIndex := int(lastIndex) + result[1]
	// We do this shift here because the .FindStringSubmatchIndex above
	// was done on a local subordinate slice of the string, not the whole string
	for index, _ := range result {
		result[index] += int(startIndex)
	}
	if r.global {
		r.putStr("lastIndex", intToValue(int64(endIndex)), true)
	}
	return
}
// 执行正则表达式匹配，并将结果转换为数组
func (r *regexpObject) exec(target valueString) Value {
	match, result := r.execRegexp(target)
	if match {
		return r.execResultToArray(target, result)
	}
	return _null
}

func (r *regexpObject) test(target valueString) bool {
	match, _ := r.execRegexp(target)
	return match
}
// 克隆一个
func (r *regexpObject) clone() *Object {
	r1 := r.val.runtime.newRegexpObject(r.prototype)
	r1.source = r.source
	r1.pattern = r.pattern
	r1.global = r.global
	r1.ignoreCase = r.ignoreCase
	r1.multiline = r.multiline
	return r1.val
}
// 初始化
func (r *regexpObject) init() {
	r.baseObject.init()
	r._putProp("lastIndex", intToValue(0), true, false, false)
}
