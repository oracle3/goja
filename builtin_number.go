package goja

import (
	"math"
	"strconv"
)
//Number.prototype.valueOf()
//valueOf() 方法返回一个被 Number 对象包装的原始值。
func (r *Runtime) numberproto_valueOf(call FunctionCall) Value {
	this := call.This
	if !isNumber(this) {
		r.typeErrorResult(true, "Value is not a number")
	}
	if _, ok := this.assertInt(); ok {
		return this
	}

	if _, ok := this.assertFloat(); ok {
		return this
	}

	if obj, ok := this.(*Object); ok {
		if v, ok := obj.self.(*primitiveValueObject); ok {
			return v.pValue
		}
	}

	r.typeErrorResult(true, "Number.prototype.valueOf is not generic")
	return nil
}

func isNumber(v Value) bool {
	switch t := v.(type) {
	case valueFloat, valueInt:
		return true
	case *Object:
		switch t := t.self.(type) {
		case *primitiveValueObject:
			return isNumber(t.pValue)
		}
	}
	return false
}
//toString() 方法返回指定 Number 对象的字符串表示形式。
func (r *Runtime) numberproto_toString(call FunctionCall) Value {
	if !isNumber(call.This) {
		r.typeErrorResult(true, "Value is not a number")
	}
	var radix int
	if arg := call.Argument(0); arg != _undefined {
		radix = int(arg.ToInteger())
	} else {
		radix = 10
	}

	if radix < 2 || radix > 36 {
		panic(r.newError(r.global.RangeError, "toString() radix argument must be between 2 and 36"))
	}

	num := call.This.ToFloat()

	if math.IsNaN(num) {
		return stringNaN
	}

	if math.IsInf(num, 1) {
		return stringInfinity
	}

	if math.IsInf(num, -1) {
		return stringNegInfinity
	}

	if radix == 10 {
		var fmt byte
		if math.Abs(num) >= 1e21 {
			fmt = 'e'
		} else {
			fmt = 'f'
		}
		return asciiString(strconv.FormatFloat(num, fmt, -1, 64))
	}

	return asciiString(dtobasestr(num, radix))
}
//Number.prototype.toFixed()
//toFixed() 方法使用定点表示法来格式化一个数值。
func (r *Runtime) numberproto_toFixed(call FunctionCall) Value {
	prec := call.Argument(0).ToInteger()
	if prec < 0 || prec > 20 {
		panic(r.newError(r.global.RangeError, "toFixed() precision must be between 0 and 20"))
	}

	num := call.This.ToFloat()
	if math.IsNaN(num) {
		return stringNaN
	}
	if math.Abs(num) >= 1e21 {
		return asciiString(strconv.FormatFloat(num, 'g', -1, 64))
	}
	return asciiString(strconv.FormatFloat(num, 'f', int(prec), 64))
}
//Number.prototype.toExponential()
//toExponential() 方法以指数表示法返回该数值字符串表示形式。
func (r *Runtime) numberproto_toExponential(call FunctionCall) Value {
	prec := call.Argument(0).ToInteger()
	if prec < 0 || prec > 20 {
		panic(r.newError(r.global.RangeError, "toExponential() precision must be between 0 and 20"))
	}

	num := call.This.ToFloat()
	if math.IsNaN(num) {
		return stringNaN
	}
	if math.Abs(num) >= 1e21 {
		return asciiString(strconv.FormatFloat(num, 'g', -1, 64))
	}
	return asciiString(strconv.FormatFloat(num, 'e', int(prec), 64))
}
//Number.prototype.toPrecision()
//toPrecision() 方法以指定的精度返回该数值对象的字符串表示。
func (r *Runtime) numberproto_toPrecision(call FunctionCall) Value {
	prec := call.Argument(0).ToInteger()
	if prec < 0 || prec > 20 {
		panic(r.newError(r.global.RangeError, "toPrecision() precision must be between 0 and 20"))
	}

	num := call.This.ToFloat()
	if math.IsNaN(num) {
		return stringNaN
	}
	if math.Abs(num) >= 1e21 {
		return asciiString(strconv.FormatFloat(num, 'g', -1, 64))
	}
	return asciiString(strconv.FormatFloat(num, 'g', int(prec), 64))
}
// 构造Number类
func (r *Runtime) initNumber() {
	r.global.NumberPrototype = r.newPrimitiveObject(valueInt(0), r.global.ObjectPrototype, classNumber)
	o := r.global.NumberPrototype.self
	o._putProp("valueOf", r.newNativeFunc(r.numberproto_valueOf, nil, "valueOf", nil, 0), true, false, true)
	o._putProp("toString", r.newNativeFunc(r.numberproto_toString, nil, "toString", nil, 0), true, false, true)
	o._putProp("toLocaleString", r.newNativeFunc(r.numberproto_toString, nil, "toLocaleString", nil, 0), true, false, true)
	o._putProp("toFixed", r.newNativeFunc(r.numberproto_toFixed, nil, "toFixed", nil, 1), true, false, true)
	o._putProp("toExponential", r.newNativeFunc(r.numberproto_toExponential, nil, "toExponential", nil, 1), true, false, true)
	o._putProp("toPrecision", r.newNativeFunc(r.numberproto_toPrecision, nil, "toPrecision", nil, 1), true, false, true)

	r.global.Number = r.newNativeFunc(r.builtin_Number, r.builtin_newNumber, "Number", r.global.NumberPrototype, 1)
	o = r.global.Number.self
	//Number.MAX_VALUE
	//能表示的最大正数。最小的负数是 -MAX_VALUE。
	o._putProp("MAX_VALUE", valueFloat(math.MaxFloat64), false, false, false)
	o._putProp("MIN_VALUE", valueFloat(math.SmallestNonzeroFloat64), false, false, false)
	//Number.NaN
	//特殊的“非数字”值。
	o._putProp("NaN", _NaN, false, false, false)
	//Number.NEGATIVE_INFINITY
	//特殊的负无穷大值，在溢出时返回该值。
	o._putProp("NEGATIVE_INFINITY", _negativeInf, false, false, false)
	//Number.POSITIVE_INFINITY
	//特殊的正无穷大值，在溢出时返回该值。
	o._putProp("POSITIVE_INFINITY", _positiveInf, false, false, false)
	//Number.EPSILON
	//两个可表示(representable)数之间的最小间隔。
	o._putProp("EPSILON", _epsilon, false, false, false)
	r.addToGlobal("Number", r.global.Number)

}
