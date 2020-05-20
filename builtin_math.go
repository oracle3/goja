package goja

import (
	"math"
)
//Math.abs(x)
//返回一个数的绝对值。
func (r *Runtime) math_abs(call FunctionCall) Value {
	return floatToValue(math.Abs(call.Argument(0).ToFloat()))
}
//Math.acos(x)
//返回一个数的反余弦值。
func (r *Runtime) math_acos(call FunctionCall) Value {
	return floatToValue(math.Acos(call.Argument(0).ToFloat()))
}
//Math.asin(x)
//返回一个数的反正弦值。
func (r *Runtime) math_asin(call FunctionCall) Value {
	return floatToValue(math.Asin(call.Argument(0).ToFloat()))
}
//Math.atan(x)
//返回一个数的反正切值。
func (r *Runtime) math_atan(call FunctionCall) Value {
	return floatToValue(math.Atan(call.Argument(0).ToFloat()))
}
//Math.atan2(y, x)
//返回 y/x 的反正切值。
func (r *Runtime) math_atan2(call FunctionCall) Value {
	y := call.Argument(0).ToFloat()
	x := call.Argument(1).ToFloat()

	return floatToValue(math.Atan2(y, x))
}
//Math.ceil(x)
//返回大于一个数的最小整数，即一个数向上取整后的值。
func (r *Runtime) math_ceil(call FunctionCall) Value {
	return floatToValue(math.Ceil(call.Argument(0).ToFloat()))
}
//Math.cos(x)
//返回一个数的余弦值。
func (r *Runtime) math_cos(call FunctionCall) Value {
	return floatToValue(math.Cos(call.Argument(0).ToFloat()))
}
//Math.exp(x)
//返回欧拉常数的参数次方，Ex，其中 x 为参数，E 是欧拉常数（2.718...，自然对数的底数）。
func (r *Runtime) math_exp(call FunctionCall) Value {
	return floatToValue(math.Exp(call.Argument(0).ToFloat()))
}
//Math.floor(x)
//返回小于一个数的最大整数，即一个数向下取整后的值。
func (r *Runtime) math_floor(call FunctionCall) Value {
	return floatToValue(math.Floor(call.Argument(0).ToFloat()))
}
//Math.log(x)
//返回一个数的自然对数（㏒e，即 ㏑）。
func (r *Runtime) math_log(call FunctionCall) Value {
	return floatToValue(math.Log(call.Argument(0).ToFloat()))
}
//Math.max([x[, y[, …]]])
//返回零到多个数值中最大值。
func (r *Runtime) math_max(call FunctionCall) Value {
	if len(call.Arguments) == 0 {
		return _negativeInf
	}

	result := call.Arguments[0].ToFloat()
	if math.IsNaN(result) {
		return _NaN
	}
	for _, arg := range call.Arguments[1:] {
		f := arg.ToFloat()
		if math.IsNaN(f) {
			return _NaN
		}
		result = math.Max(result, f)
	}
	return floatToValue(result)
}
//Math.min([x[, y[, …]]])
//返回零到多个数值中最小值。
func (r *Runtime) math_min(call FunctionCall) Value {
	if len(call.Arguments) == 0 {
		return _positiveInf
	}

	result := call.Arguments[0].ToFloat()
	if math.IsNaN(result) {
		return _NaN
	}
	for _, arg := range call.Arguments[1:] {
		f := arg.ToFloat()
		if math.IsNaN(f) {
			return _NaN
		}
		result = math.Min(result, f)
	}
	return floatToValue(result)
}
//Math.pow(x, y)
//返回一个数的 y 次幂。
func (r *Runtime) math_pow(call FunctionCall) Value {
	x := call.Argument(0)
	y := call.Argument(1)
	if x, ok := x.assertInt(); ok {
		if y, ok := y.assertInt(); ok && y >= 0 && y < 64 {
			if y == 0 {
				return intToValue(1)
			}
			if x == 0 {
				return intToValue(0)
			}
			ip := ipow(x, y)
			if ip != 0 {
				return intToValue(ip)
			}
		}
	}

	return floatToValue(math.Pow(x.ToFloat(), y.ToFloat()))
}
//Math.random()
//返回一个 0 到 1 之间的伪随机数。
func (r *Runtime) math_random(call FunctionCall) Value {
	return floatToValue(r.rand())
}
//Math.round(x)
//返回四舍五入后的整数。
func (r *Runtime) math_round(call FunctionCall) Value {
	f := call.Argument(0).ToFloat()
	if math.IsNaN(f) {
		return _NaN
	}

	if f == 0 && math.Signbit(f) {
		return _negativeZero
	}

	t := math.Trunc(f)

	if f >= 0 {
		if f-t >= 0.5 {
			return floatToValue(t + 1)
		}
	} else {
		if t-f > 0.5 {
			return floatToValue(t - 1)
		}
	}

	return floatToValue(t)
}
//Math.sin(x)
//返回一个数的正弦值。
func (r *Runtime) math_sin(call FunctionCall) Value {
	return floatToValue(math.Sin(call.Argument(0).ToFloat()))
}
//Math.sqrt(x)
//返回一个数的平方根。
func (r *Runtime) math_sqrt(call FunctionCall) Value {
	return floatToValue(math.Sqrt(call.Argument(0).ToFloat()))
}
//Math.tan(x)
//返回一个数的正切值。
func (r *Runtime) math_tan(call FunctionCall) Value {
	return floatToValue(math.Tan(call.Argument(0).ToFloat()))
}
// Math库函数注入
func (r *Runtime) createMath(val *Object) objectImpl {
	m := &baseObject{
		class:      "Math",
		val:        val,
		extensible: true,
		prototype:  r.global.ObjectPrototype,
	}
	m.init()
	//欧拉常数，也是自然对数的底数，约等于 2.718。
	m._putProp("E", valueFloat(math.E), false, false, false)
	//10 的自然对数，约等于 2.303。
	m._putProp("LN10", valueFloat(math.Ln10), false, false, false)
	//2 的自然对数，约等于 0.693。
	m._putProp("LN2", valueFloat(math.Ln2), false, false, false)
	//以 2 为底的 E 的对数，约等于 1.443。
	m._putProp("LOG2E", valueFloat(math.Log2E), false, false, false)
	//以 10 为底的 E 的对数，约等于 0.434。
	m._putProp("LOG10E", valueFloat(math.Log10E), false, false, false)
	//圆周率，一个圆的周长和直径之比，约等于 3.14159。
	m._putProp("PI", valueFloat(math.Pi), false, false, false)
	//二分之一 ½ 的平方根，同时也是 2 的平方根的倒数 ，约等于 0.707。
	m._putProp("SQRT1_2", valueFloat(sqrt1_2), false, false, false)
	//2 的平方根，约等于 1.414。
	m._putProp("SQRT2", valueFloat(math.Sqrt2), false, false, false)

	m._putProp("abs", r.newNativeFunc(r.math_abs, nil, "abs", nil, 1), true, false, true)
	m._putProp("acos", r.newNativeFunc(r.math_acos, nil, "acos", nil, 1), true, false, true)
	m._putProp("asin", r.newNativeFunc(r.math_asin, nil, "asin", nil, 1), true, false, true)
	m._putProp("atan", r.newNativeFunc(r.math_atan, nil, "atan", nil, 1), true, false, true)
	m._putProp("atan2", r.newNativeFunc(r.math_atan2, nil, "atan2", nil, 2), true, false, true)
	m._putProp("ceil", r.newNativeFunc(r.math_ceil, nil, "ceil", nil, 1), true, false, true)
	m._putProp("cos", r.newNativeFunc(r.math_cos, nil, "cos", nil, 1), true, false, true)
	m._putProp("exp", r.newNativeFunc(r.math_exp, nil, "exp", nil, 1), true, false, true)
	m._putProp("floor", r.newNativeFunc(r.math_floor, nil, "floor", nil, 1), true, false, true)
	m._putProp("log", r.newNativeFunc(r.math_log, nil, "log", nil, 1), true, false, true)
	m._putProp("max", r.newNativeFunc(r.math_max, nil, "max", nil, 2), true, false, true)
	m._putProp("min", r.newNativeFunc(r.math_min, nil, "min", nil, 2), true, false, true)
	m._putProp("pow", r.newNativeFunc(r.math_pow, nil, "pow", nil, 2), true, false, true)
	m._putProp("random", r.newNativeFunc(r.math_random, nil, "random", nil, 0), true, false, true)
	m._putProp("round", r.newNativeFunc(r.math_round, nil, "round", nil, 1), true, false, true)
	m._putProp("sin", r.newNativeFunc(r.math_sin, nil, "sin", nil, 1), true, false, true)
	m._putProp("sqrt", r.newNativeFunc(r.math_sqrt, nil, "sqrt", nil, 1), true, false, true)
	m._putProp("tan", r.newNativeFunc(r.math_tan, nil, "tan", nil, 1), true, false, true)

	return m
}
// Math库的注入
func (r *Runtime) initMath() {
	r.addToGlobal("Math", r.newLazyObject(r.createMath))
}
