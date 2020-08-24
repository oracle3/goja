package goja

import (
	"fmt"
	"github.com/oracle3/goja/parser"
	"testing"
)

func TestVM1(t *testing.T) {
	r := &Runtime{}
	r.init()

	vm := r.vm

	vm.prg = &Program{
		values: []Value{valueInt(2), valueInt(3), asciiString("test")},
		code: []instruction{
			bindName("v"),
			newObject,
			setGlobal("v"),
			loadVal(2),
			loadVal(1),
			loadVal(0),
			add,
			setElem,
			pop,
			getVar1("v"),
			halt,
		},
	}

	vm.run()

	rv := vm.pop()

	if obj, ok := rv.(*Object); ok {
		if v := obj.self.getStr("test").ToInteger(); v != 5 {
			t.Fatalf("Unexpected property value: %v", v)
		}
	} else {
		t.Fatalf("Unexpected result: %v", rv)
	}

}

func TestEvalVar(t *testing.T) {
	const SCRIPT = `
	function test() {
		var a;
		return eval("var a = 'yes'; var z = 'no'; a;") === "yes" && a === "yes";
	}
	test();
	`

	testScript1(SCRIPT, valueTrue, t)
}

var jumptable = []func(*vm, *instr){
	f_jump,
	f_halt,
}

func f_jump(vm *vm, i *instr) {
	vm.pc += i.prim
}

func f_halt(vm *vm, i *instr) {
	vm.halt = true
}

func f_loadVal(vm *vm, i *instr) {
	vm.push(vm.prg.values[i.prim])
	vm.pc++
}

func f_add(vm *vm) {
	right := vm.stack[vm.sp-1]
	left := vm.stack[vm.sp-2]

	if o, ok := left.(*Object); ok {
		left = o.self.toPrimitive()
	}

	if o, ok := right.(*Object); ok {
		right = o.self.toPrimitive()
	}

	var ret Value

	leftString, isLeftString := left.assertString()
	rightString, isRightString := right.assertString()

	if isLeftString || isRightString {
		if !isLeftString {
			leftString = left.ToString()
		}
		if !isRightString {
			rightString = right.ToString()
		}
		ret = leftString.concat(rightString)
	} else {
		if leftInt, ok := left.assertInt(); ok {
			if rightInt, ok := right.assertInt(); ok {
				ret = intToValue(int64(leftInt) + int64(rightInt))
			} else {
				ret = floatToValue(float64(leftInt) + right.ToFloat())
			}
		} else {
			ret = floatToValue(left.ToFloat() + right.ToFloat())
		}
	}

	vm.stack[vm.sp-2] = ret
	vm.sp--
	vm.pc++
}

type instr struct {
	code int
	prim int
	arg  interface{}
}

type jumparg struct {
	offset int
	other  string
}

func BenchmarkVmNOP2(b *testing.B) {
	prg := []func(*vm){
		//loadVal(0).exec,
		//loadVal(1).exec,
		//add.exec,
		jump(1).exec,
		halt.exec,
	}

	r := &Runtime{}
	r.init()

	vm := r.vm
	vm.prg = &Program{
		values: []Value{intToValue(2), intToValue(3)},
	}

	for i := 0; i < b.N; i++ {
		vm.halt = false
		vm.pc = 0
		for !vm.halt {
			prg[vm.pc](vm)
		}
		//vm.sp--
		/*r := vm.pop()
		if r.ToInteger() != 5 {
			b.Fatalf("Unexpected result: %+v", r)
		}
		if vm.sp != 0 {
			b.Fatalf("Unexpected sp: %d", vm.sp)
		}*/
	}
}

func BenchmarkVmNOP1(b *testing.B) {
	prg := []instr{
		{code: 2, prim: 0},
		{code: 2, prim: 1},
		{code: 3},
		{code: 1},
	}

	r := &Runtime{}
	r.init()

	vm := r.vm
	vm.prg = &Program{
		values: []Value{intToValue(2), intToValue(3)},
	}
	for i := 0; i < b.N; i++ {
		vm.halt = false
		vm.pc = 0
	L:
		for {
			instr := &prg[vm.pc]
			//jumptable[instr.code](vm, instr)
			switch instr.code {
			case 10:
				vm.pc += 1
			case 11:
				vm.pc += 2
			case 12:
				vm.pc += 3
			case 13:
				vm.pc += 4
			case 14:
				vm.pc += 5
			case 15:
				vm.pc += 6
			case 16:
				vm.pc += 7
			case 17:
				vm.pc += 8
			case 18:
				vm.pc += 9
			case 19:
				vm.pc += 10
			case 20:
				vm.pc += 11
			case 21:
				vm.pc += 12
			case 22:
				vm.pc += 13
			case 23:
				vm.pc += 14
			case 24:
				vm.pc += 15
			case 25:
				vm.pc += 16
			case 0:
				//vm.pc += instr.prim
				f_jump(vm, instr)
			case 1:
				break L
			case 2:
				f_loadVal(vm, instr)
			case 3:
				f_add(vm)
			default:
				jumptable[instr.code](vm, instr)
			}

		}
		r := vm.pop()
		if r.ToInteger() != 5 {
			b.Fatalf("Unexpected result: %+v", r)
		}
		if vm.sp != 0 {
			b.Fatalf("Unexpected sp: %d", vm.sp)
		}

		//vm.sp -= 1
	}
}

func BenchmarkVmNOP(b *testing.B) {
	r := &Runtime{}
	r.init()

	vm := r.vm
	vm.prg = &Program{
		code: []instruction{
			jump(1),
			//jump(1),
			halt,
		},
	}

	for i := 0; i < b.N; i++ {
		vm.pc = 0
		vm.run()
	}

}

func BenchmarkVm1(b *testing.B) {
	r := &Runtime{}
	r.init()

	vm := r.vm

	//ins1 := loadVal1(0)
	//ins2 := loadVal1(1)

	vm.prg = &Program{
		values: []Value{valueInt(2), valueInt(3)},
		code: []instruction{
			loadVal(0),
			loadVal(1),
			add,
			halt,
		},
	}

	for i := 0; i < b.N; i++ {
		vm.pc = 0
		vm.run()
		r := vm.pop()
		if r.ToInteger() != 5 {
			b.Fatalf("Unexpected result: %+v", r)
		}
		if vm.sp != 0 {
			b.Fatalf("Unexpected sp: %d", vm.sp)
		}
	}
}

func BenchmarkFib(b *testing.B) {
	const TEST_FIB = `
function fib(n) {
if (n < 2) return n;
return fib(n - 2) + fib(n - 1);
}

fib(35);
`
	b.StopTimer()
	prg, err := parser.ParseFile(nil, "test.js", TEST_FIB, 0)
	if err != nil {
		b.Fatal(err)
	}

	c := newCompiler()
	c.compile(prg)
	c.p.dumpCode(b.Logf)

	r := &Runtime{}
	r.init()

	vm := r.vm

	var expectedResult Value = valueInt(9227465)

	b.StartTimer()

	vm.prg = c.p
	vm.run()
	v := vm.pop()

	b.Logf("stack size: %d", len(vm.stack))
	b.Logf("stashAllocs: %d", vm.stashAllocs)

	if !v.SameAs(expectedResult) {
		b.Fatalf("Result: %+v, expected: %+v", v, expectedResult)
	}

}

func BenchmarkEmptyLoop(b *testing.B) {
	const SCRIPT = `
	function f() {
		for (var i = 0; i < 100; i++) {
		}
	}
	f()
	`
	b.StopTimer()
	vm := New()
	prg := MustCompile("test.js", SCRIPT, false)
	// prg.dumpCode(log.Printf)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		vm.RunProgram(prg)
	}
}

func BenchmarkVMAdd(b *testing.B) {
	vm := &vm{}
	vm.stack = append(vm.stack, nil, nil)
	vm.sp = len(vm.stack)

	var v1 Value = valueInt(3)
	var v2 Value = valueInt(5)

	for i := 0; i < b.N; i++ {
		vm.stack[0] = v1
		vm.stack[1] = v2
		add.exec(vm)
		vm.sp++
	}
}

func BenchmarkFuncCall(b *testing.B) {
	const SCRIPT = `
	function f(a, b, c, d) {
	}
	`

	b.StopTimer()

	vm := New()
	prg := MustCompile("test.js", SCRIPT, false)

	vm.RunProgram(prg)
	if f, ok := AssertFunction(vm.Get("f")); ok {
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			f(nil, nil, intToValue(1), intToValue(2), intToValue(3), intToValue(4), intToValue(5), intToValue(6))
		}
	} else {
		b.Fatal("f is not a function")
	}
}

func log(call FunctionCall) Value {
	str := call.Argument(0)
	fmt.Println(str.String())
	return str
}

func TestArray(t *testing.T) {
	const SCRIPT = `
var animals = ['pigs', 'goats', 'sheep'];
var len = animals.length;
//console.log(animals.length);
//var count = animals.push('cows');
//console.log(count);
//// expected output: 4
//console.log(animals);
//// expected output: Array ["pigs", "goats", "sheep", "cows"]
//
//animals.push('chickens', 'cats', 'dogs');
//console.log(animals);
	`

	vm := New()
	console := vm.NewObject()
	console.Set("log", log)
	vm.Set("console", console)

	prg := MustCompile("test.js", SCRIPT, false)

	result,err := vm.RunProgram(prg)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result.String())
	}

}