package webgl10

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/go4orward/gigl"
	"github.com/go4orward/gigl/common"
)

type WebGLShader struct {
	rc             *WebGLRenderingContext //
	vshader_code   string                 // vertex   shader source code
	fshader_code   string                 // fragment shader source code
	shader_program js.Value               //
	err            error                  //

	gigl.GLShaderBinder
}

// ----------------------------------------------------------------------------
// Creating Shader
// ----------------------------------------------------------------------------

func create_shader(rc *WebGLRenderingContext, vshader_source string, fshader_source string) (gigl.GLShader, error) {
	// THIS CONSTRUCTOR FUNCTION IS NOT MEANT TO BE CALLED DIRECTLY BY USER.
	// IT SHOULD BE CALLED BY 'WebGLRenderingContext.CreateShader()'.
	shader := WebGLShader{rc: rc}
	shader.CreateShaderProgram(vshader_source, fshader_source)
	shader.InitBindings()
	return &shader, shader.err
}

func (self *WebGLShader) CreateShaderProgram(vshader_source string, fshader_source string) {
	self.vshader_code = vshader_source
	self.fshader_code = fshader_source
	self.err = nil
	rc, c := self.rc, self.rc.constants
	vshader := rc.context.Call("createShader", c.VERTEX_SHADER) // Create a vertex shader object
	rc.context.Call("shaderSource", vshader, self.vshader_code) // Attach vertex shader source code
	rc.context.Call("compileShader", vshader)                   // Compile the vertex shader
	if rc.context.Call("getShaderParameter", vshader, c.COMPILE_STATUS).Bool() == false {
		msg := strings.TrimSpace(rc.context.Call("getShaderInfoLog", vshader).String())
		self.err = fmt.Errorf("VShader failed to compile (%s)", msg)
		common.Logger.Error(self.err.Error())
		return
	}
	defer rc.context.Call("deleteShader", vshader)
	fshader := rc.context.Call("createShader", c.FRAGMENT_SHADER) // Create fragment shader object
	rc.context.Call("shaderSource", fshader, self.fshader_code)   // Attach fragment shader source code
	rc.context.Call("compileShader", fshader)                     // Compile the fragmentt shader
	if self.err == nil && rc.context.Call("getShaderParameter", fshader, c.COMPILE_STATUS).Bool() == false {
		msg := strings.TrimSpace(rc.context.Call("getShaderInfoLog", fshader).String())
		self.err = fmt.Errorf("FShader failed to compile (%s)", msg)
		common.Logger.Error(self.err.Error())
		return
	}
	defer rc.context.Call("deleteShader", fshader)
	self.shader_program = rc.context.Call("createProgram")        // Create a shader program object to store the combined shader program
	rc.context.Call("attachShader", self.shader_program, vshader) // Attach a vertex shader
	rc.context.Call("attachShader", self.shader_program, fshader) // Attach a fragment shader
	rc.context.Call("linkProgram", self.shader_program)           // Link both the programs
	if self.err == nil && rc.context.Call("getProgramParameter", self.shader_program, c.LINK_STATUS).Bool() == false {
		msg := strings.TrimSpace(rc.context.Call("getProgramInfoLog", self.shader_program).String())
		self.err = fmt.Errorf("ShaderProgram failed to link (%s)", msg)
		common.Logger.Error(self.err.Error())
		return
	}
}

func (self *WebGLShader) IsReady() bool {
	return !self.shader_program.IsNull() && self.err == nil
}

func (self *WebGLShader) GetShaderProgram() any {
	return self.shader_program
}

func (self *WebGLShader) GetErr() error {
	return self.err
}

// ----------------------------------------------------------------------------
// Shader Bindings
// ----------------------------------------------------------------------------

func (self *WebGLShader) CheckBindings() {
	// check if the shader was properly built
	if self.err != nil {
		common.Logger.Error("ShaderProgram is not ready for CheckBindings()\n")
		return
	}
	// check uniform locations (type: 'object')
	for uname, utarget := range self.Uniforms {
		location := self.rc.context.Call("getUniformLocation", self.shader_program, uname)
		if location.IsNull() {
			self.err = fmt.Errorf("Uniform %q cannot be found in the shader program\n", uname)
			common.Logger.Error(self.err.Error())
		} else if utarget.Target == nil {
			self.err = fmt.Errorf("Invalid binding for uniform %q : %v \n", uname, utarget)
			common.Logger.Error(self.err.Error())
		} else { // remember the location, since gl.getXXX() is expensive
			utarget.Loc = location // save it as interface{}, not js.Value
		}
		self.Uniforms[uname] = utarget
	}
	// check attribute locations (type: 'object')
	for aname, atarget := range self.Attributes {
		location := self.rc.context.Call("getAttribLocation", self.shader_program, aname)
		if location.IsNull() {
			self.err = fmt.Errorf("Attribute %q cannot be found in the shader program\n", aname)
			common.Logger.Error(self.err.Error())
		} else if atarget.Target == nil {
			self.err = fmt.Errorf("Invalid binding for attribute %q : %v \n", aname, atarget)
			common.Logger.Error(self.err.Error())
		} else { // remember the location, since gl.getXXX() is expensive
			atarget.Loc = location // save it as interface{}, not js.Value
		}
		self.Attributes[aname] = atarget
	}
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
	shader.InitBindings()
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
	if self.err == nil && !self.shader_program.IsNull() {
		summary += fmt.Sprintf("Shader  program:Y \n")
	} else if self.err == nil && self.shader_program.IsNull() {
		summary += fmt.Sprintf("Shader  program:N \n")
	} else {
		summary += fmt.Sprintf("Shader  with Error (%s)\n", self.err.Error())
	}
	for uname, ut := range self.Uniforms {
		summary += fmt.Sprintf("    Uniform   %-10s: %s\n", uname, ut.String())
	}
	for aname, at := range self.Attributes {
		summary += fmt.Sprintf("    Attribute %-10s: %s\n", aname, at.String())
	}
	return strings.TrimSuffix(summary, "\n")
}
