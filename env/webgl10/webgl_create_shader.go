package webgl10

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
	shader_program js.Value               //
	err            error                  //

	uniforms   map[string]map[string]interface{} // shader uniforms to bind
	attributes map[string]map[string]interface{} // shader attributes to bind
}

// ----------------------------------------------------------------------------
// Creating Shader
// ----------------------------------------------------------------------------

func create_shader(rc *WebGLRenderingContext, vshader_source string, fshader_source string) (gigl.GLShader, error) {
	// THIS CONSTRUCTOR FUNCTION IS NOT MEANT TO BE CALLED DIRECTLY BY USER.
	// IT SHOULD BE CALLED BY 'WebGLRenderingContext.CreateShader()'.
	shader := WebGLShader{rc: rc}
	shader.uniforms = map[string]map[string]interface{}{}
	shader.attributes = map[string]map[string]interface{}{}
	shader.CreateShaderProgram(vshader_source, fshader_source)
	return &shader, shader.err
}

func (self *WebGLShader) CreateShaderProgram(vshader_source string, fshader_source string) {
	self.vshader_code = vshader_source
	self.fshader_code = fshader_source
	rc, c := self.rc, self.rc.constants
	vshader := rc.context.Call("createShader", c.VERTEX_SHADER) // Create a vertex shader object
	rc.context.Call("shaderSource", vshader, self.vshader_code) // Attach vertex shader source code
	rc.context.Call("compileShader", vshader)                   // Compile the vertex shader
	if rc.context.Call("getShaderParameter", vshader, c.COMPILE_STATUS).Bool() == false {
		msg := strings.TrimSpace(rc.context.Call("getShaderInfoLog", vshader).String())
		self.err = errors.New("VShader failed to compile")
		fmt.Println("VShader failed to compile : " + msg)
		return
	}
	defer rc.context.Call("deleteShader", vshader)
	fshader := rc.context.Call("createShader", c.FRAGMENT_SHADER) // Create fragment shader object
	rc.context.Call("shaderSource", fshader, self.fshader_code)   // Attach fragment shader source code
	rc.context.Call("compileShader", fshader)                     // Compile the fragmentt shader
	if self.err == nil && rc.context.Call("getShaderParameter", fshader, c.COMPILE_STATUS).Bool() == false {
		msg := strings.TrimSpace(rc.context.Call("getShaderInfoLog", fshader).String())
		self.err = errors.New("FShader failed to compile")
		fmt.Println("FShader failed to compile : " + msg)
		return
	}
	defer rc.context.Call("deleteShader", fshader)
	self.shader_program = rc.context.Call("createProgram")        // Create a shader program object to store the combined shader program
	rc.context.Call("attachShader", self.shader_program, vshader) // Attach a vertex shader
	rc.context.Call("attachShader", self.shader_program, fshader) // Attach a fragment shader
	rc.context.Call("linkProgram", self.shader_program)           // Link both the programs
	if self.err == nil && rc.context.Call("getProgramParameter", self.shader_program, c.LINK_STATUS).Bool() == false {
		msg := strings.TrimSpace(rc.context.Call("getProgramInfoLog", self.shader_program).String())
		self.err = errors.New("ShaderProgram failed to link")
		fmt.Println("ShaderProgram failed to link : " + msg)
		return
	}
}

func (self *WebGLShader) GetErr() error {
	return self.err
}

// ----------------------------------------------------------------------------
// Setting up Shader Bindings
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
	case "instance.pose", "instance.color": // instance pose or color, like "instance.pose:<stride>:<offset>"
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
	// check if the shader was properly built
	if self.err != nil {
		fmt.Printf("ShaderProgram is not ready for CheckBindings()\n")
		return
	}
	// check uniform locations (type: 'object')
	for uname, umap := range self.uniforms {
		location := self.rc.context.Call("getUniformLocation", self.shader_program, uname)
		if location.IsNull() {
			fmt.Printf("Uniform '%s' cannot be found in the shader program\n", uname)
		} else if umap["dtype"] == nil || (umap["autobinding"] == "" && umap["value"] == nil) {
			fmt.Printf("Invalid binding for uniform '%s' : %v \n", uname, umap)
		} else { // remember the location, since gl.getXXX() is expensive
			umap["location"] = interface{}(location) // save it as interface{}, not js.Value
		}
	}
	// check attribute locations (type: 'object')
	for aname, amap := range self.attributes {
		location := self.rc.context.Call("getAttribLocation", self.shader_program, aname)
		if location.IsNull() {
			fmt.Printf("Attribute '%s' cannot be found in the shader program\n", aname)
		} else if amap["dtype"] == nil || (amap["autobinding"] == "" && amap["buffer"] == nil) {
			fmt.Printf("Invalid binding for attribute '%s' : %v \n", aname, amap)
		} else { // remember the location, since gl.getXXX() is expensive
			amap["location"] = interface{}(location) // save it as interface{}, not js.Value
		}
	}
}

// ----------------------------------------------------------------------------
// Using Shader Bindings
// ----------------------------------------------------------------------------

func (self *WebGLShader) GetShaderProgram() interface{} {
	return self.shader_program
}

func (self *WebGLShader) GetUniformBindings() map[string]map[string]interface{} {
	return self.uniforms
}

func (self *WebGLShader) GetAttributeBindings() map[string]map[string]interface{} {
	return self.attributes
}

// ----------------------------------------------------------------------------
//
// ----------------------------------------------------------------------------

func (self *WebGLShader) Copy() gigl.GLShader {
	// create a new shader as a copy with empty binding
	// (so that the same 'shader_program' can be shared among different rendering targets)
	shader := WebGLShader{rc: self.rc, vshader_code: self.vshader_code, fshader_code: self.fshader_code}
	shader.shader_program = self.shader_program
	// initialize shader bindings with empty map
	shader.uniforms = map[string]map[string]interface{}{}
	shader.attributes = map[string]map[string]interface{}{}
	return &shader
}

func (self *WebGLShader) String() string {
	vert, frag, prog := "X", "X", "X"
	if self.err == nil || !strings.HasPrefix(self.err.Error(), "VShader") {
		vert = "O"
		if self.err == nil || !strings.HasPrefix(self.err.Error(), "FShader") {
			frag = "O"
			if self.err == nil {
				prog = "O"
			}
		}
	}
	return fmt.Sprintf("Shader{V:%s F:%s P:%s}", vert, frag, prog)
}

func (self *WebGLShader) Summary() string {
	summary := ""
	if self.err == nil {
		summary += fmt.Sprintf("Shader  program:%v(%T)\n", self.shader_program, self.shader_program)
	} else {
		summary += fmt.Sprintf("Shader  with Error - %s\n", self.err.Error())
	}
	for uname, umap := range self.uniforms {
		summary += fmt.Sprintf("    Uniform   %-10s: %v\n", uname, umap)
	}
	for aname, amap := range self.attributes {
		summary += fmt.Sprintf("    Attribute %-10s: %v\n", aname, amap)
	}
	return summary
}
