package goja

type objectArrayBuffer struct {
	baseObject
	data []byte
}

func (o *objectArrayBuffer) export() interface{} {
	return o.data
}

func (r *Runtime) _newArrayBuffer(proto *Object, o *Object) *objectArrayBuffer {
	if o == nil {
		o = &Object{runtime: r}
	}
	b := &objectArrayBuffer{
		baseObject: baseObject{
			class:      classObject,
			val:        o,
			prototype:  proto,
			extensible: true,
		},
	}
	o.self = b
	b.init()
	return b
}

func (r *Runtime) builtin_ArrayBuffer(args []Value, proto *Object) *Object {
	b := r._newArrayBuffer(proto, nil)
	if len(args) > 0 {
		b.data = make([]byte, toLength(args[0]))
	}
	return b.val
}
//ArrayBuffer.prototype.byteLength
//byteLength访问器属性表示一个ArrayBuffer 对象的字节长度。
func (r *Runtime) arrayBufferProto_getByteLength(call FunctionCall) Value {
	o := r.toObject(call.This)
	if b, ok := o.self.(*objectArrayBuffer); ok {
		return intToValue(int64(len(b.data)))
	}
	r.typeErrorResult(true, "Object is not ArrayBuffer: %s", o)
	panic("unreachable")
}
//ArrayBuffer.prototype.slice()
//slice()方法返回一个新的 ArrayBuffer ，它的内容是这个ArrayBuffer的字节副本，从begin（包括），到end（不包括）
func (r *Runtime) arrayBufferProto_slice(call FunctionCall) Value {
	o := r.toObject(call.This)
	if b, ok := o.self.(*objectArrayBuffer); ok {
		l := int64(len(b.data))
		start := toLength(call.Argument(0))
		if start < 0 {
			start = l + start
		}
		if start < 0 {
			start = 0
		} else if start > l {
			start = l
		}
		var stop int64
		if arg := call.Argument(1); arg != _undefined {
			stop = toLength(arg)
			if stop < 0 {
				stop = int64(len(b.data)) + stop
			}
			if stop < 0 {
				stop = 0
			} else if stop > l {
				stop = l
			}

		} else {
			stop = l
		}

		ret := r._newArrayBuffer(r.global.ArrayBufferPrototype, nil)

		if stop > start {
			ret.data = b.data[start:stop]
		}

		return ret.val
	}
	r.typeErrorResult(true, "Object is not ArrayBuffer: %s", o)
	panic("unreachable")
}
// 注入ArrayBuffer类
func (r *Runtime) createArrayBufferProto(val *Object) objectImpl {
	b := r._newArrayBuffer(r.global.Object, val)
	byteLengthProp := &valueProperty{
		accessor:     true,
		configurable: true,
		getterFunc:   r.newNativeFunc(r.arrayBufferProto_getByteLength, nil, "get byteLength", nil, 0),
	}
	b._put("byteLength", byteLengthProp)
	b._putProp("slice", r.newNativeFunc(r.arrayBufferProto_slice, nil, "slice", nil, 2), true, false, true)
	return b
}
//ArrayBuffer 对象用来表示通用的、固定长度的原始二进制数据缓冲区。
//它是一个字节数组，通常在其他语言中称为“byte array”。
func (r *Runtime) initTypedArrays() {

	r.global.ArrayBufferPrototype = r.newLazyObject(r.createArrayBufferProto)

	r.global.ArrayBuffer = r.newNativeFuncConstruct(r.builtin_ArrayBuffer, "ArrayBuffer", r.global.ArrayBufferPrototype, 1)
	r.addToGlobal("ArrayBuffer", r.global.ArrayBuffer)
}
