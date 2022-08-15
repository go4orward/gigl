package opengl41

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go4orward/gigl"
)

type OpenGLCanvas struct {
	wh     [2]int                  //
	window *glfw.Window            //
	rc     *OpenGLRenderingContext //
	paused bool                    //
}

var glfw_initialized bool = false
var gl_initialized bool = false

func NewOpenGLCanvas(width int, height int, title string, resizable bool) (*OpenGLCanvas, error) {
	if !glfw_initialized {
		if err := glfw.Init(); err != nil {
			log.Fatalln("failed to initialize glfw:", err)
		}
		glfw_initialized = true
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, err
	}
	window.MakeContextCurrent()
	if !gl_initialized {
		if err := gl.Init(); err != nil { // Initialize Glow
			log.Fatalln("failed to initialize gl: ", err)
		}
		gl_initialized = true
	}

	// create WebGL context
	self := OpenGLCanvas{window: window, wh: [2]int{width, height}}
	self.rc = NewOpenGLRenderingContext(width, height)

	return &self, nil
}

func (self *OpenGLCanvas) GetWH() [2]int {
	return self.wh
}

func (self *OpenGLCanvas) GetRenderingContext() gigl.GLRenderingContext {
	return self.rc
}

func (self *OpenGLCanvas) String() string {
	return fmt.Sprintf("OpenGLCanvas{size:%dx%d}\n", self.wh[0], self.wh[1])
}

func (self *OpenGLCanvas) SwapBuffers() {
	self.window.SwapBuffers()
}

// ----------------------------------------------------------------------------
// User Interactions (Event Handling)
// ----------------------------------------------------------------------------

func (self *OpenGLCanvas) SetEventHandlerForClick(handler func(canvasxy [2]int, keystat [4]bool)) {
}

func (self *OpenGLCanvas) SetEventHandlerForDoubleClick(handler func(canvasxy [2]int, keystat [4]bool)) {
}

func (self *OpenGLCanvas) SetEventHandlerForMouseOver(handler func(canvasxy [2]int, keystat [4]bool)) {
}

func (self *OpenGLCanvas) SetEventHandlerForMouseDrag(handler func(canvasxy [2]int, dxy [2]int, keystat [4]bool)) {
}

func (self *OpenGLCanvas) SetEventHandlerForZoom(handler func(canvasxy [2]int, scale float32, keystat [4]bool)) {
}

func (self *OpenGLCanvas) SetEventHandlerForScroll(handler func(canvasxy [2]int, dx int, dy int, keystat [4]bool)) {
}

func (self *OpenGLCanvas) SetEventHandlerForKeyPress(handler func(key string, code string, keystat [4]bool)) {
}

func (self *OpenGLCanvas) SetEventHandlerForWindowResize(handler func(w int, h int)) {
}

// ----------------------------------------------------------------------------
// Animating with DrawHandler
// ----------------------------------------------------------------------------

func (self *OpenGLCanvas) Run(draw_handler func(now float64)) {
	// run UI animation loop forever, with the given 'draw_handler'
	self.paused = false
	for !self.window.ShouldClose() {
		if draw_handler != nil && !self.paused {
			now := glfw.GetTime()
			draw_handler(now)
			self.window.SwapBuffers()
		}
		glfw.PollEvents()
	}
}

func (self *OpenGLCanvas) RunOnce(draw_handler func(now float64)) {
	// run UI animation loop only once, with the given 'draw_handler'
	var first_time bool = true
	for !self.window.ShouldClose() {
		if draw_handler != nil && first_time {
			now := glfw.GetTime()
			draw_handler(now)
			self.window.SwapBuffers()
			first_time = false
		}
		glfw.PollEvents()
	}
}

func (self *OpenGLCanvas) Pause() {
	self.paused = true
}

func (self *OpenGLCanvas) Resume() {
	self.paused = false
}
