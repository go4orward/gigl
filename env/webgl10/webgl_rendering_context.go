package webgl10

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"syscall/js"
	"unsafe"

	"github.com/go4orward/gigl"
)

type WebGLRenderingContext struct {
	context   js.Value         // WebGLRenderingContext
	constants gigl.GLConstants // WebGL constant values
	ext_uint  js.Value         // extension for "OES_element_index_uint"
	ext_angle js.Value         // extension for "ANGLE_instanced_arrays"
	wh        [2]int           // canvas width & height
}

func NewWebGLRenderingContext(canvas js.Value) (*WebGLRenderingContext, error) {
	context := canvas.Call("getContext", "webgl")
	if context.IsUndefined() {
		context = canvas.Call("getContext", "experimental-webgl")
		if context.IsUndefined() {
			return nil, errors.New("WebGL not supported")
		}
	}
	self := WebGLRenderingContext{context: context}
	self.SetupExtension("UINT32") // extension for UINT32 index
	self.SetupExtension("ANGLE")  // extension for geometry instancing
	self.wh[0] = canvas.Get("clientWidth").Int()
	self.wh[1] = canvas.Get("clientHeight").Int()
	// get WebGL constants
	self.constants.ARRAY_BUFFER = uint32(context.Get("ARRAY_BUFFER").Int())
	self.constants.BLEND = uint32(context.Get("BLEND").Int())
	self.constants.BYTE = uint32(context.Get("BYTE").Int())
	self.constants.CLAMP_TO_EDGE = uint32(context.Get("CLAMP_TO_EDGE").Int())
	self.constants.COLOR_BUFFER_BIT = uint32(context.Get("COLOR_BUFFER_BIT").Int())
	self.constants.COMPILE_STATUS = uint32(context.Get("COMPILE_STATUS").Int())
	self.constants.DEPTH_BUFFER_BIT = uint32(context.Get("DEPTH_BUFFER_BIT").Int())
	self.constants.DEPTH_TEST = uint32(context.Get("DEPTH_TEST").Int())
	self.constants.ELEMENT_ARRAY_BUFFER = uint32(context.Get("ELEMENT_ARRAY_BUFFER").Int())
	self.constants.FLOAT = uint32(context.Get("FLOAT").Int())
	self.constants.FRAGMENT_SHADER = uint32(context.Get("FRAGMENT_SHADER").Int())
	self.constants.LEQUAL = uint32(context.Get("LEQUAL").Int())
	self.constants.LESS = uint32(context.Get("LESS").Int())
	self.constants.LINEAR = uint32(context.Get("LINEAR").Int())
	self.constants.LINES = uint32(context.Get("LINES").Int())
	self.constants.LINK_STATUS = uint32(context.Get("LINK_STATUS").Int())
	self.constants.NEAREST = uint32(context.Get("NEAREST").Int())
	self.constants.ONE = uint32(context.Get("ONE").Int())
	self.constants.ONE_MINUS_SRC_ALPHA = uint32(context.Get("ONE_MINUS_SRC_ALPHA").Int())
	self.constants.POINTS = uint32(context.Get("POINTS").Int())
	self.constants.RGBA = uint32(context.Get("RGBA").Int())
	self.constants.SRC_ALPHA = uint32(context.Get("SRC_ALPHA").Int())
	self.constants.STATIC_DRAW = uint32(context.Get("STATIC_DRAW").Int())
	self.constants.TEXTURE_2D = uint32(context.Get("TEXTURE_2D").Int())
	self.constants.TEXTURE0 = uint32(context.Get("TEXTURE0").Int())
	self.constants.TEXTURE1 = uint32(context.Get("TEXTURE1").Int())
	self.constants.TEXTURE_MIN_FILTER = uint32(context.Get("TEXTURE_MIN_FILTER").Int())
	self.constants.TEXTURE_WRAP_S = uint32(context.Get("TEXTURE_WRAP_S").Int())
	self.constants.TEXTURE_WRAP_T = uint32(context.Get("TEXTURE_WRAP_T").Int())
	self.constants.TRIANGLES = uint32(context.Get("TRIANGLES").Int())
	self.constants.UNSIGNED_BYTE = uint32(context.Get("UNSIGNED_BYTE").Int())
	self.constants.UNSIGNED_INT = uint32(context.Get("UNSIGNED_INT").Int())
	self.constants.UNSIGNED_SHORT = uint32(context.Get("UNSIGNED_SHORT").Int())
	self.constants.VERTEX_SHADER = uint32(context.Get("VERTEX_SHADER").Int())
	return &self, nil
}

func (self *WebGLRenderingContext) GetWH() [2]int {
	return self.wh
}

func (self *WebGLRenderingContext) GetConstants() *gigl.GLConstants {
	return &self.constants
}

func (self *WebGLRenderingContext) GetEnvVariable(vname string, dtype string) interface{} {
	// In WebGL environment, 'EnvVariable' means QueryParameters of the current URL path
	href := js.Global().Get("window").Get("location").Get("href").String()
	url := js.Global().Get("URL").New(href)
	param := url.Get("searchParams").Call("get", vname)
	switch dtype {
	case "int":
		if param.IsNull() {
			return 0
		} else {
			n, _ := strconv.Atoi(param.String())
			return n
		}
	case "bool":
		if param.IsNull() {
			return false
		} else {
			b, _ := strconv.ParseBool(param.String())
			return b
		}
	default:
		if param.IsNull() {
			return ""
		} else {
			return param.String()
		}
	}
}

// ----------------------------------------------------------------------------
// Material & Shader
// ----------------------------------------------------------------------------

func (self *WebGLRenderingContext) LoadMaterial(material gigl.GLMaterial) error {
	return load_material(self, material)
}

func (self *WebGLRenderingContext) SetupMaterial(material gigl.GLMaterial) error {
	return setup_material(self, material)
}

func (self *WebGLRenderingContext) CreateShader(vertex_shader string, fragment_shader string) (gigl.GLShader, error) {
	return create_shader(self, vertex_shader, fragment_shader)
}

// ----------------------------------------------------------------------------
// WebGL DataBuffer
// ----------------------------------------------------------------------------

func (self *WebGLRenderingContext) CreateDataBufferVAO() *gigl.VAO {
	return &gigl.VAO{}
}

func (self *WebGLRenderingContext) CreateVtxDataBuffer(data_slice []float32) interface{} {
	// 'target' : c.ARRAY_BUFFER or c.ELEMENT_ARRAY_BUFFER
	if data_slice == nil {
		return nil
	}
	c := self.GetConstants()
	buffer := self.context.Call("createBuffer")
	self.context.Call("bindBuffer", js.ValueOf(c.ARRAY_BUFFER), buffer)
	var js_typed_array = self.ConvertGoSliceToJsTypedArray(data_slice)
	self.context.Call("bufferData", js.ValueOf(c.ARRAY_BUFFER), js_typed_array, js.ValueOf(c.STATIC_DRAW))
	self.context.Call("bindBuffer", js.ValueOf(c.ARRAY_BUFFER), nil)
	return buffer
}

func (self *WebGLRenderingContext) CreateIdxDataBuffer(data_slice []uint32) interface{} {
	// 'target' : c.ARRAY_BUFFER or c.ELEMENT_ARRAY_BUFFER
	if data_slice == nil {
		return nil
	}
	c := self.GetConstants()
	buffer := self.context.Call("createBuffer")
	self.context.Call("bindBuffer", js.ValueOf(c.ELEMENT_ARRAY_BUFFER), buffer)
	var js_typed_array = self.ConvertGoSliceToJsTypedArray(data_slice)
	self.context.Call("bufferData", js.ValueOf(c.ELEMENT_ARRAY_BUFFER), js_typed_array, js.ValueOf(c.STATIC_DRAW))
	self.context.Call("bindBuffer", js.ValueOf(c.ELEMENT_ARRAY_BUFFER), nil)
	return buffer
}

// ----------------------------------------------------------------------------
// Binding DataBuffer
// ----------------------------------------------------------------------------

func (self *WebGLRenderingContext) GLBindBuffer(target uint32, buffer interface{}) {
	// 'bind_target' : c.ARRAY_BUFFER or c.ELEMENT_ARRAY_BUFFER
	if buffer == nil {
		self.context.Call("bindBuffer", js.ValueOf(target), js.Null())
	} else {
		self.context.Call("bindBuffer", js.ValueOf(target), buffer.(js.Value))
	}
}

// ----------------------------------------------------------------------------
// Binding Texture
// ----------------------------------------------------------------------------

func (self *WebGLRenderingContext) GLActiveTexture(texture_unit int) {
	js_texture_unit := js.ValueOf(js.ValueOf(self.constants.TEXTURE0 + uint32(texture_unit)))
	self.context.Call("activeTexture", js_texture_unit)
}

func (self *WebGLRenderingContext) GLBindTexture(target uint32, texture interface{}) {
	self.context.Call("bindTexture", js.ValueOf(target), texture.(js.Value)) // 'binding_target' : TEXTURE_2D
}

// ----------------------------------------------------------------------------
// Binding Uniforms
// ----------------------------------------------------------------------------

func (self *WebGLRenderingContext) GLUniform1i(location interface{}, v0 int) {
	self.context.Call("uniform1i", location.(js.Value), v0)
}

func (self *WebGLRenderingContext) GLUniform1f(location interface{}, v0 float32) {
	self.context.Call("uniform1f", location.(js.Value), v0)
}

func (self *WebGLRenderingContext) GLUniform2f(location interface{}, v0 float32, v1 float32) {
	self.context.Call("uniform2f", location.(js.Value), v0, v1)
}

func (self *WebGLRenderingContext) GLUniform3f(location interface{}, v0 float32, v1 float32, v2 float32) {
	self.context.Call("uniform3f", location.(js.Value), v0, v1, v2)
}

func (self *WebGLRenderingContext) GLUniform4f(location interface{}, v0 float32, v1 float32, v2 float32, v3 float32) {
	self.context.Call("uniform4f", location.(js.Value), v0, v1, v2, v3)
}

func (self *WebGLRenderingContext) GLUniformMatrix3fv(location interface{}, transpose bool, values []float32) {
	js_typed_array := self.ConvertGoSliceToJsTypedArray(values) // converted to JavaScript 'Float32Array'
	self.context.Call("uniformMatrix3fv", location.(js.Value), transpose, js_typed_array)
}

func (self *WebGLRenderingContext) GLUniformMatrix4fv(location interface{}, transpose bool, values []float32) {
	js_typed_array := self.ConvertGoSliceToJsTypedArray(values) // converted to JavaScript 'Float32Array'
	self.context.Call("uniformMatrix4fv", location.(js.Value), transpose, js_typed_array)
}

// ----------------------------------------------------------------------------
// Binding Attributes
// ----------------------------------------------------------------------------

func (self *WebGLRenderingContext) GLVertexAttribPointer(location interface{}, size int, dtype uint32, normalized bool, stride_in_byte int, offset_in_byte int) {
	self.context.Call("vertexAttribPointer", location.(js.Value), size, js.ValueOf(dtype), normalized, stride_in_byte, offset_in_byte)
}

func (self *WebGLRenderingContext) GLEnableVertexAttribArray(location interface{}) {
	self.context.Call("enableVertexAttribArray", location.(js.Value))
}

func (self *WebGLRenderingContext) GLVertexAttribDivisor(location interface{}, divisor int) {
	// extension for geometry instancing
	if !self.ext_angle.IsNull() {
		self.ext_angle.Call("vertexAttribDivisorANGLE", location.(js.Value), divisor)
	}
}

// ----------------------------------------------------------------------------
// Preparing to Render
// ----------------------------------------------------------------------------

func (self *WebGLRenderingContext) GLClearColor(r float32, g float32, b float32, a float32) {
	self.context.Call("clearColor", r, g, b, a)
}

func (self *WebGLRenderingContext) GLClear(mask uint32) {
	self.context.Call("clear", js.ValueOf(mask))
}

func (self *WebGLRenderingContext) GLEnable(cap uint32) {
	self.context.Call("enable", js.ValueOf(cap))
}

func (self *WebGLRenderingContext) GLDisable(cap uint32) {
	self.context.Call("disable", js.ValueOf(cap))
}

func (self *WebGLRenderingContext) GLDepthFunc(ftn uint32) {
	self.context.Call("depthFunc", js.ValueOf(ftn))
}

func (self *WebGLRenderingContext) GLBlendFunc(sfactor uint32, dfactor uint32) {
	self.context.Call("blendFunc", js.ValueOf(sfactor), js.ValueOf(dfactor))
}

func (self *WebGLRenderingContext) GLUseProgram(shader_program interface{}) {
	self.context.Call("useProgram", shader_program.(js.Value))
}

// ----------------------------------------------------------------------------
// Rendering
// ----------------------------------------------------------------------------

func (self *WebGLRenderingContext) GLDrawArrays(mode uint32, first int, count int) {
	// 'mode' : POINTS
	self.context.Call("drawArrays", js.ValueOf(mode), first, count)
}

func (self *WebGLRenderingContext) GLDrawArraysInstanced(mode uint32, first int, count int, pose_count int) {
	// 'mode' : POINTS
	if !self.ext_angle.IsNull() {
		self.ext_angle.Call("drawArraysInstancedANGLE", js.ValueOf(mode), first, count, pose_count)
	}
}

func (self *WebGLRenderingContext) GLDrawElements(mode uint32, count int, dtype uint32, offset int) {
	// 'mode'  : LINES, TRIANGLES
	// 'dtype' : UNSIGNED_INT
	self.context.Call("drawElements", js.ValueOf(mode), count, js.ValueOf(dtype), offset)
}

func (self *WebGLRenderingContext) GLDrawElementsInstanced(mode uint32, element_count int, dtype uint32, offset int, pose_count int) {
	// 'mode'  : LINES, TRIANGLES
	// 'dtype' : UNSIGNED_INT
	if !self.ext_angle.IsNull() {
		self.ext_angle.Call("drawElementsInstancedANGLE", js.ValueOf(mode), element_count, js.ValueOf(dtype), offset, pose_count)
	}
}

// ----------------------------------------------------------------------------
// WebGL Extensions
// ----------------------------------------------------------------------------

func (self *WebGLRenderingContext) SetupExtension(extname string) {
	switch extname {
	case "UINT32": // extension for UINT32 index, to drawElements() with large number of vertices
		self.ext_uint = self.context.Call("getExtension", "OES_element_index_uint")
	case "ANGLE": // extension for geometry instancing
		self.ext_angle = self.context.Call("getExtension", "ANGLE_instanced_arrays")
	}
}

func (self *WebGLRenderingContext) IsExtensionReady(extname string) bool {
	switch extname {
	case "UINT32": // extension for UINT32 index, to drawElements() with large number of vertices
		return !self.ext_uint.IsNull() && !self.ext_uint.IsUndefined()
	case "ANGLE": // extension for geometry instancing
		return !self.ext_angle.IsNull() && !self.ext_angle.IsUndefined()
	}
	return false
}

// ----------------------------------------------------------------------------
// private functions
// ----------------------------------------------------------------------------

func (self *WebGLRenderingContext) ConvertGoSliceToJsTypedArray(a interface{}) js.Value {
	// Since js.TypedArrayOf() of Go1.11 is no longer supported (due to WASM memory issue),
	// we have to use js.CopyBytesToJS() instead. (Now it runs fine with Go1.15.7, Feb 5 2021)
	//   Ref: syscall/js: replace TypedArrayOf with CopyBytesToGo/CopyBytesToJS
	//   Ref: https://github.com/golang/go/issues/31980  	("js.TypedArrayOf is impossible to use correctly")
	//   Ref: https://go-review.googlesource.com/c/go/+/177537/
	//   Ref: https://github.com/golang/go/issues/32402  	(solution provided by 'hajimehoshi')
	//   Ref: https://github.com/nuberu/webgl				(Golang WebAssembly wrapper for WebGL)
	// Note that this solution sacrifices performance. (WebGL renderer's frame rate will be OK, though)
	// We hope Go/WebAssembly will sort out this issue in the future.
	switch a := a.(type) {
	case []int8:
		b := js.Global().Get("Uint8Array").New(len(a))
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&a))
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		js.CopyBytesToJS(b, byte_slice)
		return js.Global().Get("Int8Array").New(b.Get("buffer"), b.Get("byteOffset"), b.Get("byteLength"))
	case []int16:
		b := js.Global().Get("Uint8Array").New(len(a) * 2)
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&a))
		slice_head.Len *= 2
		slice_head.Cap *= 2
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		js.CopyBytesToJS(b, byte_slice)
		return js.Global().Get("Int16Array").New(b.Get("buffer"), b.Get("byteOffset"), b.Get("byteLength").Int()/2)
	case []int32:
		b := js.Global().Get("Uint8Array").New(len(a) * 4)
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&a))
		slice_head.Len *= 4
		slice_head.Cap *= 4
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		js.CopyBytesToJS(b, byte_slice)
		return js.Global().Get("Int32Array").New(b.Get("buffer"), b.Get("byteOffset"), b.Get("byteLength").Int()/4)
	case []int64:
		b := js.Global().Get("Uint8Array").New(len(a) * 8)
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&a))
		slice_head.Len *= 8
		slice_head.Cap *= 8
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		js.CopyBytesToJS(b, byte_slice)
		return js.Global().Get("BigInt64Array").New(b.Get("buffer"), b.Get("byteOffset"), b.Get("byteLength").Int()/8)
	case []uint8:
		b := js.Global().Get("Uint8Array").New(len(a))
		js.CopyBytesToJS(b, a)
		return b
	case []uint16:
		b := js.Global().Get("Uint8Array").New(len(a) * 2)
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&a))
		slice_head.Len *= 2
		slice_head.Cap *= 2
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		js.CopyBytesToJS(b, byte_slice)
		return js.Global().Get("Uint16Array").New(b.Get("buffer"), b.Get("byteOffset"), b.Get("byteLength").Int()/2)
	case []uint32:
		b := js.Global().Get("Uint8Array").New(len(a) * 4)
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&a))
		slice_head.Len *= 4
		slice_head.Cap *= 4
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		js.CopyBytesToJS(b, byte_slice)
		return js.Global().Get("Uint32Array").New(b.Get("buffer"), b.Get("byteOffset"), b.Get("byteLength").Int()/4)
	case []uint64:
		b := js.Global().Get("Uint8Array").New(len(a) * 4)
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&a))
		slice_head.Len *= 8
		slice_head.Cap *= 8
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		js.CopyBytesToJS(b, byte_slice)
		return js.Global().Get("BigUint64Array").New(b.Get("buffer"), b.Get("byteOffset"), b.Get("byteLength").Int()/8)
	case []float32:
		b := js.Global().Get("Uint8Array").New(len(a) * 4)
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&a))
		slice_head.Len *= 4
		slice_head.Cap *= 4
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		// ShowArrayInfo(byte_slice)
		js.CopyBytesToJS(b, byte_slice)
		return js.Global().Get("Float32Array").New(b.Get("buffer"), b.Get("byteOffset"), b.Get("byteLength").Int()/4)
	case []float64:
		b := js.Global().Get("Uint8Array").New(len(a) * 8)
		slice_head := (*reflect.SliceHeader)(unsafe.Pointer(&a))
		slice_head.Len *= 8
		slice_head.Cap *= 8
		byte_slice := *(*[]byte)(unsafe.Pointer(slice_head))
		js.CopyBytesToJS(b, byte_slice)
		return js.Global().Get("Float64Array").New(b.Get("buffer"), b.Get("byteOffset"), b.Get("byteLength").Int()/8)
	default:
		panic(fmt.Sprintf("Unexpected value at ConvertGoSliceToJsTypedArray(): %T", a))
	}
}
