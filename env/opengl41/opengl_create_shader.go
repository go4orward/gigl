package opengl41

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go4orward/gigl"
	"github.com/go4orward/gigl/common"
)

type OpenGLShader struct {
	rc             *OpenGLRenderingContext //
	vshader_code   string                  // vertex   shader source code
	fshader_code   string                  // fragment shader source code
	shader_program uint32                  //
	err            error                   //

	gigl.GLShaderBinder
}

// ----------------------------------------------------------------------------
// Creating Shader
// ----------------------------------------------------------------------------

func create_shader(rc *OpenGLRenderingContext, vshader_source string, fshader_source string) (*OpenGLShader, error) {
	// THIS CONSTRUCTOR FUNCTION IS NOT MEANT TO BE CALLED DIRECTLY BY USER.
	// IT SHOULD BE CALLED BY 'OpenGLRenderingContext.CreateShader()'.
	shader := OpenGLShader{rc: rc}
	shader.CreateShaderProgram(vshader_source, fshader_source)
	shader.InitBindings()
	return &shader, shader.err
}

func (self *OpenGLShader) CreateShaderProgram(vshader_source string, fshader_source string) {
	self.vshader_code = vshader_source
	self.fshader_code = fshader_source
	self.err = nil
	c := self.rc.GetConstants()
	var status, logLength int32
	// vertex shader (uint32)
	vshader_source = self.prepare_vshader_source(vshader_source) //
	vshader := gl.CreateShader(c.VERTEX_SHADER)                  // Create a vertex shader object
	vsource, vcfree := gl.Strs(vshader_source)                   // Get C-string of vertex shader source code
	gl.ShaderSource(vshader, 1, vsource, nil)                    // Attach vertex shader source code
	gl.CompileShader(vshader)                                    // Compile the vertex shader
	vcfree()
	gl.GetShaderiv(vshader, c.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		gl.GetShaderiv(vshader, gl.INFO_LOG_LENGTH, &logLength)
		logmsg := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(vshader, logLength, nil, gl.Str(logmsg))
		self.err = errors.New("VShader failed to compile")
		fmt.Println("VShader failed to compile : " + strings.TrimSpace(logmsg))
		return
	}
	// fragment shader (uint32)
	fshader_source = self.prepare_fshader_source(fshader_source) //
	fshader := gl.CreateShader(c.FRAGMENT_SHADER)                // Create fragment shader object
	fsource, fcfree := gl.Strs(fshader_source)                   // Get C-string of fragment shader source code
	gl.ShaderSource(fshader, 1, fsource, nil)                    // Attach fragment shader source code
	gl.CompileShader(fshader)                                    // Compile the fragment shader
	fcfree()
	gl.GetShaderiv(fshader, c.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		gl.GetShaderiv(fshader, gl.INFO_LOG_LENGTH, &logLength)
		logmsg := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(fshader, logLength, nil, gl.Str(logmsg))
		self.err = errors.New("FShader failed to compile")
		fmt.Println("FShader failed to compile : " + strings.TrimSpace(logmsg))
		return
	}
	// shader program (uint32)
	self.shader_program = gl.CreateProgram()
	gl.AttachShader(self.shader_program, vshader)
	gl.AttachShader(self.shader_program, fshader)
	gl.LinkProgram(self.shader_program)
	gl.GetProgramiv(self.shader_program, gl.LINK_STATUS, &status)
	defer gl.DeleteShader(vshader)
	defer gl.DeleteShader(fshader)
	if status == gl.FALSE {
		gl.GetProgramiv(self.shader_program, gl.INFO_LOG_LENGTH, &logLength)
		logmsg := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(self.shader_program, logLength, nil, gl.Str(logmsg))
		self.err = errors.New("ShaderProgram failed to link")
		fmt.Println("ShaderProgram failed to link : " + strings.TrimSpace(logmsg))
		return
	}
}

func (self *OpenGLShader) IsReady() bool {
	return self.shader_program != 0 && self.err == nil
}

func (self *OpenGLShader) GetShaderProgram() any {
	return self.shader_program
}

func (self *OpenGLShader) GetErr() error {
	return self.err
}

func (self *OpenGLShader) prepare_vshader_source(source string) string {
	source = strings.ReplaceAll(source, "attribute", "in")
	source = strings.ReplaceAll(source, "varying", "out")
	return "#version 410\n" + source + "\x00"
}

func (self *OpenGLShader) prepare_fshader_source(source string) string {
	source = strings.ReplaceAll(source, "varying", "in")
	source = strings.Replace(source, "void main", "out vec4 OUTPUT_COLOR;\nvoid main", 1)
	source = strings.ReplaceAll(source, "gl_FragColor", "OUTPUT_COLOR")
	source = strings.ReplaceAll(source, "texture2D", "texture")
	return "#version 410\n" + source + "\x00"
}

// ----------------------------------------------------------------------------
// Shader Bindings
// ----------------------------------------------------------------------------

func (self *OpenGLShader) CheckBindings() {
	// check if the shader was properly built
	if self.err != nil {
		common.Logger.Error("ShaderProgram is not ready for CheckBindings (%v)\n", self.err)
		return
	}
	// check uniform locations (type: 'uint32')
	for uname, utarget := range self.Uniforms {
		location := gl.GetUniformLocation(self.shader_program, gl.Str(uname+"\x00")) // int32
		if location < 0 {
			self.err = fmt.Errorf("Uniform %q cannot be found in the shader program\n", uname)
			common.Logger.Error(self.err.Error())
		} else if utarget.Target == nil {
			self.err = fmt.Errorf("Invalid binding for uniform %q : %v \n", uname, utarget)
			common.Logger.Error(self.err.Error())
		} else { // remember the location, since gl.getXXX() is expensive
			utarget.Loc = location // save it as any
		}
		self.Uniforms[uname] = utarget
	}
	// check attribute locations (type: 'uint32')
	for aname, atarget := range self.Attributes {
		location := gl.GetAttribLocation(self.shader_program, gl.Str(aname+"\x00")) // int32
		if location < 0 {
			self.err = fmt.Errorf("Attribute %q cannot be found in the shader program\n", aname)
			common.Logger.Error(self.err.Error())
		} else if atarget.Target == nil {
			self.err = fmt.Errorf("Invalid binding for attribute %q : %v \n", aname, atarget)
			common.Logger.Error(self.err.Error())
		} else { // remember the location, since gl.getXXX() is expensive
			atarget.Loc = location // save it as any
		}
		self.Attributes[aname] = atarget
	}
}

// ----------------------------------------------------------------------------
//
// ----------------------------------------------------------------------------

func (self *OpenGLShader) Copy() gigl.GLShader {
	// create a new shader as a copy with empty binding
	shader := OpenGLShader{rc: self.rc, vshader_code: self.vshader_code, fshader_code: self.fshader_code}
	shader.shader_program = self.shader_program
	// initialize shader bindings with empty map
	shader.InitBindings()
	return &shader
}

func (self *OpenGLShader) String() string {
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
	return fmt.Sprintf("Shader{V:%s F:%s P:%s Unf:%d Att:%d}", vert, frag, prog, len(self.Uniforms), len(self.Attributes))
}

func (self *OpenGLShader) Summary() string {
	summary := ""
	if self.err == nil && self.shader_program != 0 {
		summary += fmt.Sprintf("Shader  program:Y \n")
	} else if self.err == nil && self.shader_program == 0 {
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
