package opengl41

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go4orward/gigl"
)

type OpenGLWindow struct {
	wh     [2]int                  //
	window *glfw.Window            //
	rc     *OpenGLRenderingContext //

}

func NewOpenGLWindow(width int, height int, title string, resizable bool) (*OpenGLWindow, error) {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
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
	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	// create WebGL context
	self := OpenGLWindow{window: window, wh: [2]int{width, height}}
	self.rc = NewOpenGLRenderingContext(width, height)
	return &self, nil
}

func (self *OpenGLWindow) GetGLRenderingContext() gigl.GLRenderingContext {
	return self.rc
}

func (self *OpenGLWindow) String() string {
	return fmt.Sprintf("OpenGLWindow{size:%dx%d}\n", self.wh[0], self.wh[1])
}

// ----------------------------------------------------------------------------
// User Interactions (Event Handling)
// ----------------------------------------------------------------------------

func (self *OpenGLWindow) SetEventHandlerForClick(handler func(canvasxy [2]int, keystat [4]bool)) {
}

func (self *OpenGLWindow) SetEventHandlerForDoubleClick(handler func(canvasxy [2]int, keystat [4]bool)) {
}

func (self *OpenGLWindow) SetEventHandlerForMouseOver(handler func(canvasxy [2]int, keystat [4]bool)) {
}

func (self *OpenGLWindow) SetEventHandlerForMouseDrag(handler func(canvasxy [2]int, dxy [2]int, keystat [4]bool)) {
}

func (self *OpenGLWindow) SetEventHandlerForZoom(handler func(canvasxy [2]int, scale float32, keystat [4]bool)) {
}

func (self *OpenGLWindow) SetEventHandlerForScroll(handler func(canvasxy [2]int, dx int, dy int, keystat [4]bool)) {
}

func (self *OpenGLWindow) SetEventHandlerForKeyPress(handler func(key string, code string, keystat [4]bool)) {
}

func (self *OpenGLWindow) SetEventHandlerForWindowResize(handler func(w int, h int)) {
}

// ----------------------------------------------------------------------------
// Animating with DrawHandler
// ----------------------------------------------------------------------------

func (self *OpenGLWindow) Run(draw_handler func(now float64)) {
	// run UI animation loop forever, with the given 'draw_handler'
	for !self.window.ShouldClose() {
		if draw_handler != nil {
			now := glfw.GetTime()
			draw_handler(now)
			self.window.SwapBuffers()
		}
		glfw.PollEvents()
	}
}

func (self *OpenGLWindow) RunOnce(draw_handler func(now float64)) {
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
