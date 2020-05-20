package goja

func (r *Runtime) initErrors() {
	//通过Error的构造器可以创建一个错误对象。当运行时错误产生时，
	//Error的实例对象会被抛出。Error对象也可用于用户自定义的异常的基础对象
	r.global.ErrorPrototype = r.NewObject()
	o := r.global.ErrorPrototype.self
	o._putProp("message", stringEmpty, true, false, true)
	o._putProp("name", stringError, true, false, true)
	o._putProp("toString", r.newNativeFunc(r.error_toString, nil, "toString", nil, 0), true, false, true)

	r.global.Error = r.newNativeFuncConstruct(r.builtin_Error, "Error", r.global.ErrorPrototype, 1)
	o = r.global.Error.self
	r.addToGlobal("Error", r.global.Error)

	r.global.TypeErrorPrototype = r.builtin_new(r.global.Error, []Value{})
	o = r.global.TypeErrorPrototype.self
	o._putProp("name", stringTypeError, true, false, true)
	//TypeError
	//创建一个error实例，表示错误的原因：变量或参数不属于有效类型。
	r.global.TypeError = r.newNativeFuncConstructProto(r.builtin_Error, "TypeError", r.global.TypeErrorPrototype, r.global.Error, 1)
	r.addToGlobal("TypeError", r.global.TypeError)

	r.global.ReferenceErrorPrototype = r.builtin_new(r.global.Error, []Value{})
	o = r.global.ReferenceErrorPrototype.self
	o._putProp("name", stringReferenceError, true, false, true)
	//ReferenceError
	//创建一个error实例，表示错误的原因：无效引用。
	r.global.ReferenceError = r.newNativeFuncConstructProto(r.builtin_Error, "ReferenceError", r.global.ReferenceErrorPrototype, r.global.Error, 1)
	r.addToGlobal("ReferenceError", r.global.ReferenceError)

	r.global.SyntaxErrorPrototype = r.builtin_new(r.global.Error, []Value{})
	o = r.global.SyntaxErrorPrototype.self
	o._putProp("name", stringSyntaxError, true, false, true)
	//SyntaxError
	//创建一个error实例，表示错误的原因：eval()在解析代码的过程中发生的语法错误。
	r.global.SyntaxError = r.newNativeFuncConstructProto(r.builtin_Error, "SyntaxError", r.global.SyntaxErrorPrototype, r.global.Error, 1)
	r.addToGlobal("SyntaxError", r.global.SyntaxError)

	r.global.RangeErrorPrototype = r.builtin_new(r.global.Error, []Value{})
	o = r.global.RangeErrorPrototype.self
	o._putProp("name", stringRangeError, true, false, true)
	//RangeError
	//创建一个error实例，表示错误的原因：数值变量或参数超出其有效范围。
	r.global.RangeError = r.newNativeFuncConstructProto(r.builtin_Error, "RangeError", r.global.RangeErrorPrototype, r.global.Error, 1)
	r.addToGlobal("RangeError", r.global.RangeError)

	r.global.EvalErrorPrototype = r.builtin_new(r.global.Error, []Value{})
	o = r.global.EvalErrorPrototype.self
	o._putProp("name", stringEvalError, true, false, true)
	//EvalError
	//创建一个error实例，表示错误的原因：与 eval() 有关。
	r.global.EvalError = r.newNativeFuncConstructProto(r.builtin_Error, "EvalError", r.global.EvalErrorPrototype, r.global.Error, 1)
	r.addToGlobal("EvalError", r.global.EvalError)

	r.global.URIErrorPrototype = r.builtin_new(r.global.Error, []Value{})
	o = r.global.URIErrorPrototype.self
	o._putProp("name", stringURIError, true, false, true)
	//URIError
	//创建一个error实例，表示错误的原因：给 encodeURI()或  decodeURl()传递的参数无效。
	r.global.URIError = r.newNativeFuncConstructProto(r.builtin_Error, "URIError", r.global.URIErrorPrototype, r.global.Error, 1)
	r.addToGlobal("URIError", r.global.URIError)

	r.global.GoErrorPrototype = r.builtin_new(r.global.Error, []Value{})
	o = r.global.GoErrorPrototype.self
	o._putProp("name", stringGoError, true, false, true)
	//GoError是个自建的错误类型，目前不知道干啥用，elikong
	r.global.GoError = r.newNativeFuncConstructProto(r.builtin_Error, "GoError", r.global.GoErrorPrototype, r.global.Error, 1)
	r.addToGlobal("GoError", r.global.GoError)
}
