package main

import (
	"errors"
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	opengl "github.com/go4orward/gigl/env/opengl41"
)

func init() { // This is needed to let main() run on the startup thread.
	runtime.LockOSThread() // Ref: https://golang.org/pkg/runtime/#LockOSThread
}

func main() {
	canvas, err := opengl.NewOpenGLCanvas(800, 600, "OpenGL1st: Triangle in Clip Space", false)
	if err != nil {
		log.Fatal(errors.New("Failed to create OpenGL canvas : " + err.Error()))
	}

	// Geometry
	vertices := []float32{-0.5, 0.5, 0, -0.5, -0.5, 0, 0.5, -0.5, 0}
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Shaders
	var status, logLength int32
	vshader_source := `#version 410
		in vec3 xyz;
		void main(void) {
			gl_Position = vec4(xyz, 1.0);		// coordinates in the CLIP space [-1.0 ~ +1.0]
		}` + "\x00"
	fshader_source := `#version 410
		out vec4 OUTPUT_COLOR;
		void main(void) {
			OUTPUT_COLOR = vec4(0.0, 0.0, 1.0, 1.0);
		}` + "\x00"
	vshader := gl.CreateShader(gl.VERTEX_SHADER) // Create a vertex shader object (uint32)
	vsource, vcfree := gl.Strs(vshader_source)   // Get C-string of vertex shader source code
	gl.ShaderSource(vshader, 1, vsource, nil)    // Attach vertex shader source code
	gl.CompileShader(vshader)                    // Compile the vertex shader
	vcfree()
	gl.GetShaderiv(vshader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		gl.GetShaderiv(vshader, gl.INFO_LOG_LENGTH, &logLength)
		logmsg := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(vshader, logLength, nil, gl.Str(logmsg))
		log.Fatal(errors.New("VShader failed to compile : " + logmsg))
	}
	fshader := gl.CreateShader(gl.FRAGMENT_SHADER) // Create fragment shader object (uint32)
	fsource, fcfree := gl.Strs(fshader_source)     // Get C-string of fragment shader source code
	gl.ShaderSource(fshader, 1, fsource, nil)      // Attach fragment shader source code
	gl.CompileShader(fshader)                      // Compile the fragment shader
	fcfree()
	gl.GetShaderiv(fshader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		gl.GetShaderiv(fshader, gl.INFO_LOG_LENGTH, &logLength)
		logmsg := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(fshader, logLength, nil, gl.Str(logmsg))
		log.Fatal(errors.New("FShader failed to compile : " + logmsg))
	}
	shader_program := gl.CreateProgram() // shader program (uint32)
	gl.AttachShader(shader_program, vshader)
	gl.AttachShader(shader_program, fshader)
	gl.LinkProgram(shader_program)
	gl.GetProgramiv(shader_program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		gl.GetProgramiv(shader_program, gl.INFO_LOG_LENGTH, &logLength)
		logmsg := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shader_program, logLength, nil, gl.Str(logmsg))
		log.Fatal(errors.New("ShaderProgram failed to link : " + logmsg))
	}
	gl.UseProgram(shader_program) // Let the completed shader program to be used
	gl.DeleteShader(vshader)
	gl.DeleteShader(fshader)

	// Attributes
	loc := uint32(gl.GetAttribLocation(shader_program, gl.Str("xyz\x00"))) // Get the location of attribute 'xyz' in the shader
	gl.VertexAttribPointerWithOffset(loc, 3, gl.FLOAT, false, 3*4, 0)      // Point 'xyz' location to the positions of ARRAY_BUFFER
	gl.EnableVertexAttribArray(loc)                                        // Enable the use of attribute 'xyz' from ARRAY_BUFFER

	// Rendering
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)

	canvas.SwapBuffers()
	canvas.RunOnce(nil)
}
