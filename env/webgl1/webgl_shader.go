package webgl1

import (
	"errors"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/go4orward/gigl"
)

type WebGLShader struct {
	rc             *WebGLRenderingContext //
	vshader_code   string                 // vertex   shader source code
	fshader_code   string                 // fragment shader source code
	vert_shader    js.Value               //
	frag_shader    js.Value               //
	shader_program js.Value               //
	err            error                  //

	uniforms   map[string]map[string]interface{} // shader uniforms to bind
	attributes map[string]map[string]interface{} // shader attributes to bind
}

func NewWebGLShader(rc *WebGLRenderingContext, vertex_shader string, fragment_shader string) (*WebGLShader, error) {
	shader := WebGLShader{rc: rc, vshader_code: vertex_shader, fshader_code: fragment_shader}
	shader.vert_shader = rc.context.Call("createShader", rc.constants.VERTEX_SHADER) // Create a vertex shader object
	rc.context.Call("shaderSource", shader.vert_shader, shader.vshader_code)         // Attach vertex shader source code
	rc.context.Call("compileShader", shader.vert_shader)                             // Compile the vertex shader
	if rc.context.Call("getShaderParameter", shader.vert_shader, rc.constants.COMPILE_STATUS).Bool() == false {
		msg := strings.TrimSpace(rc.context.Call("getShaderInfoLog", shader.vert_shader).String())
		shader.err = errors.New("VShader failed to compile : " + msg)
		fmt.Println(shader.err.Error())
	}
	shader.frag_shader = rc.context.Call("createShader", rc.constants.FRAGMENT_SHADER) // Create fragment shader object
	rc.context.Call("shaderSource", shader.frag_shader, shader.fshader_code)           // Attach fragment shader source code
	rc.context.Call("compileShader", shader.frag_shader)                               // Compile the fragmentt shader
	if shader.err == nil && rc.context.Call("getShaderParameter", shader.frag_shader, rc.constants.COMPILE_STATUS).Bool() == false {
		msg := strings.TrimSpace(rc.context.Call("getShaderInfoLog", shader.frag_shader).String())
		shader.err = errors.New("FShader failed to compile : " + msg)
		fmt.Println(shader.err.Error())
	}
	shader.shader_program = rc.context.Call("createProgram")                   // Create a shader program object to store the combined shader program
	rc.context.Call("attachShader", shader.shader_program, shader.vert_shader) // Attach a vertex shader
	rc.context.Call("attachShader", shader.shader_program, shader.frag_shader) // Attach a fragment shader
	rc.context.Call("linkProgram", shader.shader_program)                      // Link both the programs
	if shader.err == nil && rc.context.Call("getProgramParameter", shader.shader_program, rc.constants.LINK_STATUS).Bool() == false {
		msg := strings.TrimSpace(rc.context.Call("getProgramInfoLog", shader.shader_program).String())
		shader.err = errors.New("ShaderProgram failed to link : " + msg)
		fmt.Println(shader.err.Error())
	}
	// initialize shader bindings with empty map
	shader.uniforms = map[string]map[string]interface{}{}
	shader.attributes = map[string]map[string]interface{}{}
	return &shader, shader.err
}

func (self *WebGLShader) String() string {
	vert, frag, prog := "X", "X", "X"
	if self.vshader_code != "" && !self.vert_shader.IsNull() && !self.vert_shader.IsUndefined() {
		vert = "O"
	}
	if self.fshader_code != "" && !self.frag_shader.IsNull() && !self.frag_shader.IsUndefined() {
		vert = "O"
	}
	if !self.shader_program.IsNull() && !self.shader_program.IsUndefined() {
		vert = "O"
	}
	return fmt.Sprintf("Shader{vert:%s frag:%s prog:%s}", vert, frag, prog)
}

func (self *WebGLShader) Copy() gigl.GLShader {
	// create a new shader as a copy with empty binding
	shader := WebGLShader{rc: self.rc, vshader_code: self.vshader_code, fshader_code: self.fshader_code}
	shader.vert_shader = self.vert_shader
	shader.frag_shader = self.frag_shader
	shader.shader_program = self.shader_program
	// initialize shader bindings with empty map
	shader.uniforms = map[string]map[string]interface{}{}
	shader.attributes = map[string]map[string]interface{}{}
	return &shader
}

func (self *WebGLShader) GetShaderProgram() interface{} {
	return self.shader_program
}

func (self *WebGLShader) GetUniformBindings() map[string]map[string]interface{} {
	return self.uniforms
}

func (self *WebGLShader) GetAttributeBindings() map[string]map[string]interface{} {
	return self.attributes
}

func (self *WebGLShader) ShowInfo() {
	if self.err == nil {
		fmt.Printf("Shader  OK\n")
	} else {
		fmt.Printf("Shader  with Error - %s\n", self.err.Error())
	}
	for uname, umap := range self.uniforms {
		fmt.Printf("    Uniform   %-10s: %v\n", uname, umap)
	}
	for aname, amap := range self.attributes {
		fmt.Printf("    Attribute %-10s: %v\n", aname, amap)
	}
}

// ----------------------------------------------------------------------------
// Bindings
// ----------------------------------------------------------------------------

func (self *WebGLShader) SetBindingForUniform(name string, dtype string, option interface{}) {
	// Set uniform binding with its name, data_type, and AUTO/MANUAL option.
	switch option.(type) {
	case string: // let Renderer bind the uniform variable automatically
		autobinding := option.(string)
		autobinding_split := strings.Split(option.(string), ":")
		autobinding0 := autobinding_split[0] // "material.texture:0" (with texture UNIT value)
		switch autobinding0 {
		case "lighting.dlight": // [mat3](3D) directional light information with (direction[3], color[3], ambient[3])
		case "material.color": //  [vec3] uniform color taken from Material
		case "material.texture": // [sampler2D] texture sampler(unit), like "material.texture:0"
		case "renderer.aspect": // AspectRatio of camera, Width : Height
		case "renderer.pvm": //  [mat3](2D) or [mat4](3D) (Proj * View * Model) matrix
		case "renderer.proj": // [mat3](2D) or [mat4](3D) (Projection) matrix
		case "renderer.vwmd": // [mat3](2D) or [mat4](3D) (View * Model) matrix
		default:
			fmt.Printf("Failed to SetBindingForUniform('%s') : unknown autobinding '%s'\n", name, autobinding)
			return
		}
		self.uniforms[name] = map[string]interface{}{"dtype": dtype, "autobinding": autobinding}
	case []float32: // let Renderer set the uniform manually, with the given values
		self.uniforms[name] = map[string]interface{}{"dtype": dtype, "value": option.([]float32)}
	default:
		fmt.Printf("Failed to SetBindingForUniform('%s') : invalid option %v\n", name, option)
		return
	}
}

func (self *WebGLShader) SetBindingForAttribute(name string, dtype string, autobinding string) {
	// Set attribute binding with its name, data_type, with AUTO_BINDING option.
	autobinding_split := strings.Split(autobinding, ":")
	autobinding0 := autobinding_split[0]
	switch autobinding0 {
	case "geometry.coords": // point coordinates
	case "geometry.textuv": // texture UV coordinates
	case "geometry.normal": // (3D only) normal vector
	case "instance.pose": // instance pose, like "instance.pose:<stride>:<offset>"
		if len(autobinding_split) != 3 {
			fmt.Printf("Failed to SetBindingForAttribute('%s') : try 'instance.pose:<stride>:<offset>'\n", name)
			return
		}
	default:
		fmt.Printf("Failed to SetBindingForAttribute('%s') : invalid autobinding '%s'\n", name, autobinding)
		return
	}
	self.attributes[name] = map[string]interface{}{"dtype": dtype, "autobinding": autobinding}
}

func (self *WebGLShader) CheckBindings() {
	// check uniform locations before rendering (since gl.getXXX() is expensive)
	if self.err != nil {
		fmt.Printf("ShaderProgram is not ready for CheckBindings()\n")
		return
	}
	for uname, umap := range self.uniforms {
		location := self.rc.context.Call("getUniformLocation", self.shader_program, uname)
		if location.IsNull() {
			fmt.Printf("Uniform '%s' cannot be found in the shader program\n", uname)
		} else if umap["dtype"] == nil || (umap["autobinding"] == "" && umap["value"] == nil) {
			fmt.Printf("Invalid binding for uniform '%s' : %v \n", uname, umap)
		} else {
			umap["location"] = interface{}(location) // save it as interface{}, not js.Value
		}
	}
	// check attribute locations
	for aname, amap := range self.attributes {
		location := self.rc.context.Call("getAttribLocation", self.shader_program, aname)
		if location.IsNull() {
			fmt.Printf("Attribute '%s' cannot be found in the shader program\n", aname)
		} else if amap["dtype"] == nil || (amap["autobinding"] == "" && amap["buffer"] == nil) {
			fmt.Printf("Invalid binding for attribute '%s' : %v \n", aname, amap)
		} else {
			amap["location"] = interface{}(location) // save it as interface{}, not js.Value
		}
	}
}
