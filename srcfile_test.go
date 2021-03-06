package goja

import "testing"
//测试SrcFile
func TestPosition(t *testing.T) {
	const SRC = `line1
line2
line3`
	f := NewSrcFile("", SRC, nil)

	tests := []struct {
		offset int
		line   int
		col    int
	}{
		{0, 1, 1},
		{2, 1, 3},
		{2, 1, 3},
		{6, 2, 1},
		{7, 2, 2},
		{12, 3, 1},
		{12, 3, 1},
		{13, 3, 2},
		{13, 3, 2},
		{16, 3, 5},
		{17, 3, 6},
	}

	for i, test := range tests {
		if p := f.Position(test.offset); p.Line != test.line || p.Col != test.col {
			t.Fatalf("%d. Line: %d, col: %d", i, p.Line, p.Col)
		}
	}
}
//测试SrcFile并发性
func TestSrcFileConcurrency(t *testing.T) {
	const SRC = `line1
line2
line3`
	f := NewSrcFile("", SRC, nil)
	go func() {
		f.Position(12)
	}()
	f.Position(2)
}
