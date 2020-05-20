package goja

import (
	"fmt"
)
// 构造一个匿名函数
func (r *Runtime) builtin_Function(args []Value, proto *Object) *Object {
	src := "(function anonymous("
	if len(args) > 1 {
		for _, arg := range args[:len(args)-1] {
			src += arg.String() + ","
		}
		src = src[:len(src)-1]
	}
	body := ""
	if len(args) > 0 {
		body = args[len(args)-1].String()
	}
	src += "){" + body + "})"

	return r.toObject(r.eval(src, false, false, _undefined))
}
//Function.prototype.toString()
//toString() 方法返回一个表示当前函数源代码的字符串。
func (r *Runtime) functionproto_toString(call FunctionCall) Value {
	obj := r.toObject(call.This)
repeat:
	switch f := obj.self.(type) {
	case *funcObject:
		return newStringValue(f.src)
	case *nativeFuncObject:
		return newStringValue(fmt.Sprintf("function %s() { [native code] }", f.nameProp.get(call.This).ToString()))
	case *boundFuncObject:
		return newStringValue(fmt.Sprintf("function %s() { [native code] }", f.nameProp.get(call.This).ToString()))
	case *lazyObject:
		obj.self = f.create(obj)
		goto repeat
	}

	r.typeErrorResult(true, "Object is not a function")
	return nil
}

func (r *Runtime) toValueArray(a Value) []Value {
	obj := r.toObject(a)
	l := toUInt32(obj.self.getStr("length"))
	ret := make([]Value, l)
	for i := uint32(0); i < l; i++ {
		ret[i] = obj.self.get(valueInt(i))
	}
	return ret
}
//Function.prototype.apply()
//apply() 方法调用一个具有给定this值的函数，以及作为一个数组（或类似数组对象）提供的参数。
func (r *Runtime) functionproto_apply(call FunctionCall) Value {
	f := r.toCallable(call.This)
	var args []Value
	if len(call.Arguments) >= 2 {
		args = r.toValueArray(call.Arguments[1])
	}
	return f(FunctionCall{
		This:      call.Argument(0),
		Arguments: args,
	})
}
//Function.prototype.call()
//call() 方法使用一个指定的 this 值和单独给出的一个或多个参数来调用一个函数。
func (r *Runtime) functionproto_call(call FunctionCall) Value {
	f := r.toCallable(call.This)
	var args []Value
	if len(call.Arguments) > 0 {
		args = call.Arguments[1:]
	}
	return f(FunctionCall{
		This:      call.Argument(0),
		Arguments: args,
	})
}

func (r *Runtime) boundCallable(target func(FunctionCall) Value, boundArgs []Value) func(FunctionCall) Value {
	var this Value
	var args []Value
	if len(boundArgs) > 0 {
		this = boundArgs[0]
		args = make([]Value, len(boundArgs)-1)
		copy(args, boundArgs[1:])
	} else {
		this = _undefined
	}
	return func(call FunctionCall) Value {
		a := append(args, call.Arguments...)
		return target(FunctionCall{
			This:      this,
			Arguments: a,
		})
	}
}
// 构造bound的数据
func (r *Runtime) boundConstruct(target func([]Value) *Object, boundArgs []Value) func([]Value) *Object {
	if target == nil {
		return nil
	}
	var args []Value
	if len(boundArgs) > 1 {
		args = make([]Value, len(boundArgs)-1)
		copy(args, boundArgs[1:])
	}
	return func(fargs []Value) *Object {
		a := append(args, fargs...)
		copy(a, args)
		return target(a)
	}
}
//Function.prototype.bind()
//bind() 方法创建一个新的函数，在 bind() 被调用时，这个新函数的 this 被指定为 bind() 的第一个参数，
//而其余参数将作为新函数的参数，供调用时使用。
func (r *Runtime) functionproto_bind(call FunctionCall) Value {
	obj := r.toObject(call.This)
	f := obj.self
	var fcall func(FunctionCall) Value
	var construct func([]Value) *Object
repeat:
	switch ff := f.(type) {
	case *funcObject:
		fcall = ff.Call
		construct = ff.construct
	case *nativeFuncObject:
		fcall = ff.f
		construct = ff.construct
	case *boundFuncObject:
		f = &ff.nativeFuncObject
		goto repeat
	case *lazyObject:
		f = ff.create(obj)
		goto repeat
	default:
		r.typeErrorResult(true, "Value is not callable: %s", obj.ToString())
	}

	l := int(toUInt32(obj.self.getStr("length")))
	l -= len(call.Arguments) - 1
	if l < 0 {
		l = 0
	}

	v := &Object{runtime: r}

	ff := r.newNativeFuncObj(v, r.boundCallable(fcall, call.Arguments), r.boundConstruct(construct, call.Arguments), "", nil, l)
	v.self = &boundFuncObject{
		nativeFuncObject: *ff,
	}

	//ret := r.newNativeFunc(r.boundCallable(f, call.Arguments), nil, "", nil, l)
	//o := ret.self
	//o.putStr("caller", r.global.throwerProperty, false)
	//o.putStr("arguments", r.global.throwerProperty, false)
	return v
}
// Function类的实现
func (r *Runtime) initFunction() {
	o := r.global.FunctionPrototype.self
	o.(*nativeFuncObject).prototype = r.global.ObjectPrototype
	o._putProp("toString", r.newNativeFunc(r.functionproto_toString, nil, "toString", nil, 0), true, false, true)
	o._putProp("apply", r.newNativeFunc(r.functionproto_apply, nil, "apply", nil, 2), true, false, true)
	o._putProp("call", r.newNativeFunc(r.functionproto_call, nil, "call", nil, 1), true, false, true)
	o._putProp("bind", r.newNativeFunc(r.functionproto_bind, nil, "bind", nil, 1), true, false, true)

	r.global.Function = r.newNativeFuncConstruct(r.builtin_Function, "Function", r.global.FunctionPrototype, 1)
	r.addToGlobal("Function", r.global.Function)
}
