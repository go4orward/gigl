package opengl41

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go4orward/gigl"
)

type OpenGLShader struct {
	rc           *OpenGLRenderingContext //
	vshader_code string                  // vertex   shader source code
	fshader_code string                  // fragment shader source code
	program      uint32                  //
	err          error                   //

	uniforms   map[string]map[string]interface{} // shader uniforms to bind
	attributes map[string]map[string]interface{} // shader attributes to bind
}

// ----------------------------------------------------------------------------
// Creating Shader
// ----------------------------------------------------------------------------

func new_opengl_shader(rc *OpenGLRenderingContext, vshader_source string, fshader_source string) (*OpenGLShader, error) {
	// THIS CONSTRUCTOR FUNCTION IS NOT MEANT TO BE CALLED DIRECTLY BY USER.
	// IT SHOULD BE CALLED BY 'OpenGLRenderingContext.CreateShader()'.
	self := OpenGLShader{rc: rc}
	self.uniforms = map[string]map[string]interface{}{}
	self.attributes = map[string]map[string]interface{}{}
	self.CreateShaderProgram(vshader_source, fshader_source)
	return &self, self.err
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
	self.program = gl.CreateProgram()
	gl.AttachShader(self.program, vshader)
	gl.AttachShader(self.program, fshader)
	gl.LinkProgram(self.program)
	gl.GetProgramiv(self.program, gl.LINK_STATUS, &status)
	defer gl.DeleteShader(vshader)
	defer gl.DeleteShader(fshader)
	if status == gl.FALSE {
		gl.GetProgramiv(self.program, gl.INFO_LOG_LENGTH, &logLength)
		logmsg := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(self.program, logLength, nil, gl.Str(logmsg))
		self.err = errors.New("ShaderProgram failed to link")
		fmt.Println("ShaderProgram failed to link : " + strings.TrimSpace(logmsg))
		return
	}
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
// Setting up Shader Bindings
// ----------------------------------------------------------------------------

func (self *OpenGLShader) SetBindingForUniform(name string, dtype string, option interface{}) {
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

func (self *OpenGLShader) SetBindingForAttribute(name string, dtype string, autobinding string) {
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

func (self *OpenGLShader) CheckBindings() {
	// check if the shader was properly built
	if self.err != nil {
		fmt.Printf("ShaderProgram is not ready for CheckBindings (%v)\n", self.err)
		return
	}
	// check uniform locations (type: 'uint32')
	for uname, umap := range self.uniforms {
		location := gl.GetUniformLocation(self.program, gl.Str(uname+"\x00")) // int32
		if location < 0 {
			fmt.Printf("Uniform '%s' cannot be found in the shader program\n", uname)
		} else if umap["dtype"] == nil || (umap["autobinding"] == nil && umap["value"] == nil) {
			fmt.Printf("Invalid binding for uniform '%s' : %v \n", uname, umap)
		} else { // remember the location, since gl.getXXX() is expensive
			umap["location"] = interface{}(location) // save it as interface{}, not uint32
			// fmt.Printf(" Uniform   %-8s (location: %v %T) : %v\n", uname, location, location, umap)
		}
	}
	// check attribute locations (type: 'uint32')
	for aname, amap := range self.attributes {
		location := gl.GetAttribLocation(self.program, gl.Str(aname+"\x00")) // int32
		if location < 0 {
			fmt.Printf("Attribute '%s' cannot be found in the shader program\n", aname)
		} else if amap["dtype"] == nil || (amap["autobinding"] == nil && amap["buffer"] == nil) {
			fmt.Printf("Invalid binding for attribute '%s' : %v \n", aname, amap)
		} else { // remember the location, since gl.getXXX() is expensive
			amap["location"] = interface{}(location) // save it as interface{}, not uint32
			// fmt.Printf(" Attribute %-8s (location: %v %T) : %v\n", aname, location, location, amap)
		}
	}
}

// ----------------------------------------------------------------------------
// Using Shader Bindings
// ----------------------------------------------------------------------------

func (self *OpenGLShader) GetShaderProgram() interface{} {
	return self.program
}

func (self *OpenGLShader) GetUniformBindings() map[string]map[string]interface{} {
	return self.uniforms
}

func (self *OpenGLShader) GetAttributeBindings() map[string]map[string]interface{} {
	return self.attributes
}

// ----------------------------------------------------------------------------
//
// ----------------------------------------------------------------------------

func (self *OpenGLShader) Copy() gigl.GLShader {
	// create a new shader as a copy with empty binding
	shader := OpenGLShader{rc: self.rc, vshader_code: self.vshader_code, fshader_code: self.fshader_code}
	shader.program = self.program
	// initialize shader bindings with empty map
	shader.uniforms = map[string]map[string]interface{}{}
	shader.attributes = map[string]map[string]interface{}{}
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
	return fmt.Sprintf("Shader{V:%s F:%s P:%s Unf:%d Att:%d}", vert, frag, prog, len(self.uniforms), len(self.attributes))
}

func (self *OpenGLShader) ShowInfo() {
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
