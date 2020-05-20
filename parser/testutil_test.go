package parser

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
)

// Quick and dirty replacement for terst
// 调用函数f()，并捕获异常，打印信息
func tt(t *testing.T, f func()) {
	defer func() {
		if x := recover(); x != nil {
			_, file, line, _ := runtime.Caller(4)
			t.Errorf("Error at %s:%d: %v", filepath.Base(file), line, x)
		}
	}()

	f()
}
// 判断a，b转换为字符串后是否相等
func is(a, b interface{}) {
	as := fmt.Sprintf("%v", a)
	bs := fmt.Sprintf("%v", b)
	if as != bs {
		panic(fmt.Errorf("%+v(%T) != %+v(%T)", a, a, b, b))
	}
}
