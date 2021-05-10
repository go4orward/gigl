package webgl1

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
	self.constants.ARRAY_BUFFER = context.Get("ARRAY_BUFFER")
	self.constants.BLEND = context.Get("BLEND")
	self.constants.BYTE = context.Get("BYTE")
	self.constants.CLAMP_TO_EDGE = context.Get("CLAMP_TO_EDGE")
	self.constants.COLOR_BUFFER_BIT = context.Get("COLOR_BUFFER_BIT")
	self.constants.COMPILE_STATUS = context.Get("COMPILE_STATUS")
	self.constants.DEPTH_BUFFER_BIT = context.Get("DEPTH_BUFFER_BIT")
	self.constants.DEPTH_TEST = context.Get("DEPTH_TEST")
	self.constants.ELEMENT_ARRAY_BUFFER = context.Get("ELEMENT_ARRAY_BUFFER")
	self.constants.FLOAT = context.Get("FLOAT")
	self.constants.FRAGMENT_SHADER = context.Get("FRAGMENT_SHADER")
	self.constants.LEQUAL = context.Get("LEQUAL")
	self.constants.LINEAR = context.Get("LINEAR")
	self.constants.LINES = context.Get("LINES")
	self.constants.LINK_STATUS = context.Get("LINK_STATUS")
	self.constants.NEAREST = context.Get("NEAREST")
	self.constants.ONE = context.Get("ONE")
	self.constants.ONE_MINUS_SRC_ALPHA = context.Get("ONE_MINUS_SRC_ALPHA")
	self.constants.POINTS = context.Get("POINTS")
	self.constants.RGBA = context.Get("RGBA")
	self.constants.SRC_ALPHA = context.Get("SRC_ALPHA")
	self.constants.STATIC_DRAW = context.Get("STATIC_DRAW")
	self.constants.TEXTURE_2D = context.Get("TEXTURE_2D")
	self.constants.TEXTURE0 = context.Get("TEXTURE0")
	self.constants.TEXTURE1 = context.Get("TEXTURE1")
	self.constants.TEXTURE_MIN_FILTER = context.Get("TEXTURE_MIN_FILTER")
	self.constants.TEXTURE_WRAP_S = context.Get("TEXTURE_WRAP_S")
	self.constants.TEXTURE_WRAP_T = context.Get("TEXTURE_WRAP_T")
	self.constants.TRIANGLES = context.Get("TRIANGLES")
	self.constants.UNSIGNED_BYTE = context.Get("UNSIGNED_BYTE")
	self.constants.UNSIGNED_INT = context.Get("UNSIGNED_INT")
	self.constants.UNSIGNED_SHORT = context.Get("UNSIGNED_SHORT")
	self.constants.VERTEX_SHADER = context.Get("VERTEX_SHADER")
	return &self, nil
}

func (self *WebGLRenderingContext) GetWH() [2]int {
	return self.wh
}

func (self *WebGLRenderingContext) GetContext() js.Value {
	return self.context
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
// Create Material & Shader
// ----------------------------------------------------------------------------

func (self *WebGLRenderingContext) CreateMaterial(source string, options ...interface{}) (gigl.GLMaterial, error) {
	return NewWebGLMaterial(self, source, options...)
}

func (self *WebGLRenderingContext) CreateShader(vertex_shader string, fragment_shader string) (gigl.GLShader, error) {
	return NewWebGLShader(self, vertex_shader, fragment_shader)
}

// ----------------------------------------------------------------------------
// WebGL Data Buffer
// ----------------------------------------------------------------------------

func (self *WebGLRenderingContext) CreateDataBuffer(target interface{}, data_slice interface{}) interface{} {
	// 'target' : c.ARRAY_BUFFER or c.ELEMENT_ARRAY_BUFFER
	if data_slice == nil {
		return nil
	}
	c := self.GetConstants()
	buffer := self.context.Call("createBuffer", c.ARRAY_BUFFER.(js.Value)) // TODO: no argument needed
	self.context.Call("bindBuffer", target.(js.Value), buffer)
	var js_typed_array = self.ConvertGoSliceToJsTypedArray(data_slice)
	self.context.Call("bufferData", target.(js.Value), js_typed_array, c.STATIC_DRAW.(js.Value))
	self.context.Call("bindBuffer", target.(js.Value), nil)
	return buffer
}

func (self *WebGLRenderingContext) GLBindBuffer(target interface{}, buffer interface{}) {
	// 'bind_target' : c.ARRAY_BUFFER or c.ELEMENT_ARRAY_BUFFER
	if buffer == nil {
		self.context.Call("bindBuffer", target.(js.Value), js.Null())
	} else {
		self.context.Call("bindBuffer", target.(js.Value), buffer.(js.Value))
	}
}

// ----------------------------------------------------------------------------
// Binding Texture
// ----------------------------------------------------------------------------

func (self *WebGLRenderingContext) GLCreateTexture() interface{} {
	return self.context.Call("createTexture")
}

func (self *WebGLRenderingContext) GLActiveTexture(texture_unit int) {
	js_texture_unit := js.ValueOf(self.constants.TEXTURE0.(js.Value).Int() + texture_unit)
	self.context.Call("activeTexture", js_texture_unit)
}

func (self *WebGLRenderingContext) GLBindTexture(target interface{}, texture interface{}) {
	// 'binding_target' : TEXTURE_2D
	self.context.Call("bindTexture", target.(js.Value), texture.(js.Value))
}

func (self *WebGLRenderingContext) GLTexImage2DFromPixelBuffer(target interface{}, level int, internalformat interface{}, width int, height int, border int, format interface{}, dtype interface{}, pixbuf []uint8) {
	js_buffer := self.ConvertGoSliceToJsTypedArray(pixbuf)
	// rc.Context.Call("texImage2D", c.TEXTURE_2D, 0, c.RGBA, size.X, size.Y, 0, c.RGBA, c.UNSIGNED_BYTE, js_buffer)
	self.context.Call("texImage2D", target.(js.Value), level, internalformat.(js.Value), width, height, border, format.(js.Value), dtype.(js.Value), js_buffer)

}

func (self *WebGLRenderingContext) GLTexImage2DFromImgObject(target interface{}, level int, internalformat interface{}, format interface{}, dtype interface{}, imgobj interface{}) {
	// rc.Context.Call("texImage2D", c.TEXTURE_2D, 0, c.RGBA, c.RGBA, c.UNSIGNED_BYTE, canvas)
	// 'imgobj' can be ImageData, ImageBitmap, HTMLImageElement, HTMLCanvasElement, or HTMLVideoElement
	self.context.Call("texImage2D", target.(js.Value), level, internalformat.(js.Value), format.(js.Value), dtype.(js.Value), imgobj.(js.Value))
}

func (self *WebGLRenderingContext) GLGenerateMipmap(target interface{}) {
	// rc.Context.Call("generateMipmap", c.TEXTURE_2D)
	self.context.Call("generateMipmap", target.(js.Value))

}

func (self *WebGLRenderingContext) GLTexParameteri(target interface{}, pname interface{}, pvalue interface{}) {
	self.context.Call("texParameteri", target.(js.Value), pname.(js.Value), pvalue.(js.Value))
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

func (self *WebGLRenderingContext) GLVertexAttribPointer(location interface{}, size int, dtype interface{}, normalized bool, stride_in_byte int, offset_in_byte int) {
	self.context.Call("vertexAttribPointer", location.(js.Value), size, dtype.(js.Value), normalized, stride_in_byte, offset_in_byte)
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

func (self *WebGLRenderingContext) GLClear(mask interface{}) {
	self.context.Call("clear", mask.(js.Value))
}

func (self *WebGLRenderingContext) GLEnable(cap interface{}) {
	self.context.Call("enable", cap.(js.Value))
}

func (self *WebGLRenderingContext) GLDisable(cap interface{}) {
	self.context.Call("disable", cap.(js.Value))
}

func (self *WebGLRenderingContext) GLDepthFunc(ftn interface{}) {
	self.context.Call("depthFunc", ftn.(js.Value))
}

func (self *WebGLRenderingContext) GLBlendFunc(sfactor interface{}, dfactor interface{}) {
	self.context.Call("blendFunc", sfactor.(js.Value), dfactor.(js.Value))
}

func (self *WebGLRenderingContext) GLUseProgram(shader_program interface{}) {
	self.context.Call("useProgram", shader_program.(js.Value))
}

// ----------------------------------------------------------------------------
// Rendering
// ----------------------------------------------------------------------------

func (self *WebGLRenderingContext) GLDrawArrays(mode interface{}, first int, count int) {
	// 'mode' : POINTS
	self.context.Call("drawArrays", mode.(js.Value), first, count)
}

func (self *WebGLRenderingContext) GLDrawArraysInstanced(mode interface{}, first int, count int, pose_count int) {
	// 'mode' : POINTS
	if !self.ext_angle.IsNull() {
		self.ext_angle.Call("drawArraysInstancedANGLE", mode.(js.Value), first, count, pose_count)
	}
}

func (self *WebGLRenderingContext) GLDrawElements(mode interface{}, count int, dtype interface{}, offset int) {
	// 'mode'  : LINES, TRIANGLES
	// 'dtype' : UNSIGNED_INT
	self.context.Call("drawElements", mode.(js.Value), count, dtype.(js.Value), offset)
}

func (self *WebGLRenderingContext) GLDrawElementsInstanced(mode interface{}, element_count int, dtype interface{}, offset int, pose_count int) {
	// 'mode'  : LINES, TRIANGLES
	// 'dtype' : UNSIGNED_INT
	if !self.ext_angle.IsNull() {
		self.ext_angle.Call("drawElementsInstancedANGLE", mode.(js.Value), element_count, dtype.(js.Value), offset, pose_count)
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