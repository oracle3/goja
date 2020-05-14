// Package file encapsulates the file abstractions used by the ast & parser.
//
package file

import (
	"fmt"
	"strings"
)

// Idx is a compact encoding of a source position within a file set.
// It can be converted into a Position for a more convenient, but much
// larger, representation.
//Idx是文件集中源位置的压缩编码。
//它可以转换为一个更方便，但更大的表示位置。
type Idx int

// Position describes an arbitrary source position
// including the filename, line, and column location.
//位置描述任意源位置
//包括文件名、行和列位置。
type Position struct {
	Filename string // The filename where the error occurred, if any 发生错误的文件名（如果有）
	Offset   int    // The src offset 文件偏移量
	Line     int    // The line number, starting at 1 行号，从1开始
	Column   int    // The column number, starting at 1 (The character count) 列号，从1开始（字符计数）

}

// A Position is valid if the line number is > 0.
//如果行号大于0，则位置有效。
func (self *Position) isValid() bool {
	return self.Line > 0
}

// String returns a string in one of several forms:
//
//	file:line:column    A valid position with filename
//	line:column         A valid position without filename
//	file                An invalid position with filename
//	-                   An invalid position without filename
//
//String以以下几种形式之一返回字符串：
//
//文件：行：列 有文件名的有效位置
//行：列 没有文件名的有效位置
//文件 有文件名的无效位置
//- 没有文件名的无效位置
//

func (self *Position) String() string {
	str := self.Filename
	if self.isValid() {
		if str != "" {
			str += ":"
		}
		str += fmt.Sprintf("%d:%d", self.Line, self.Column)
	}
	if str == "" {
		str = "-"
	}
	return str
}

// FileSet

// A FileSet represents a set of source files.
//文件集表示一组源文件。
type FileSet struct {
	files []*File
	last  *File
}

// AddFile adds a new file with the given filename and src.
//
// This an internal method, but exported for cross-package use.
//AddFile添加具有给定文件名和src的新文件。
//
//这是一个内部方法，但导出用于跨包使用。

func (self *FileSet) AddFile(filename, src string) int {
	base := self.nextBase()
	file := &File{
		name: filename,
		src:  src,
		base: base,
	}
	self.files = append(self.files, file)
	self.last = file
	return base
}
// 返回文件集的文件索引
func (self *FileSet) nextBase() int {
	if self.last == nil {
		return 1
	}
	return self.last.base + len(self.last.src) + 1
}
// 返回指定位置的文件
func (self *FileSet) File(idx Idx) *File {
	for _, file := range self.files {
		if idx <= Idx(file.base+len(file.src)) {
			return file
		}
	}
	return nil
}

// Position converts an Idx in the FileSet into a Position.
// Position将文件集中的Idx转换为Position。
func (self *FileSet) Position(idx Idx) *Position {
	position := &Position{}
	for _, file := range self.files {
		if idx <= Idx(file.base+len(file.src)) {
			offset := int(idx) - file.base
			src := file.src[:offset]
			position.Filename = file.name
			position.Offset = offset
			position.Line = 1 + strings.Count(src, "\n")
			if index := strings.LastIndex(src, "\n"); index >= 0 {
				position.Column = offset - index
			} else {
				position.Column = 1 + len(src)
			}
		}
	}
	return position
}

type File struct {
	name string
	src  string
	base int // This will always be 1 or greater
}

func NewFile(filename, src string, base int) *File {
	return &File{
		name: filename,
		src:  src,
		base: base,
	}
}

func (fl *File) Name() string {
	return fl.name
}

func (fl *File) Source() string {
	return fl.src
}

func (fl *File) Base() int {
	return fl.base
}
