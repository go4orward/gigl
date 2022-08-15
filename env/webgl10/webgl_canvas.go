package webgl10

import (
	"errors"
	"fmt"
	"math"
	"syscall/js"

	"github.com/go4orward/gigl"
)

type WebGLCanvas struct {
	id     string                 // canvas DOM element's ID
	canvas js.Value               // canvas DOM element
	wh     [2]int                 //
	rc     *WebGLRenderingContext //

	mouse_event_common_handler js.Func //
	mouse_wheel_common_handler js.Func //
	mouse_dragging             bool
	mouse_sxy                  [2]int
	mouse_wheel_scale          float64 // in the range of [0 ~ 500(default) ~ 1000]
	evthandler_for_click       func(canvasxy [2]int, keystat [4]bool)
	evthandler_for_dblclick    func(canvasxy [2]int, keystat [4]bool)
	evthandler_for_mouse_over  func(canvasxy [2]int, keystat [4]bool)
	evthandler_for_mouse_drag  func(canvasxy [2]int, dxy [2]int, keystat [4]bool)
	evthandler_for_zoom        func(canvasxy [2]int, scale float32, keystat [4]bool)
	evthandler_for_scroll      func(canvasxy [2]int, dx int, dy int, keystat [4]bool)
	wasm_handler_for_draw      js.Func
	user_handler_for_draw      func(now float64)
	paused                     bool
}

func NewWebGLCanvas(canvas_id string) (*WebGLCanvas, error) {
	self := WebGLCanvas{id: canvas_id}
	var err error
	// initialize the canvas
	doc := js.Global().Get("document")
	self.canvas = doc.Call("getElementById", canvas_id)
	if self.canvas.IsNull() {
		err := errors.New("Canvas not found (ID:'" + canvas_id + "')")
		js.Global().Call("alert", "Failed to start WebGL : "+err.Error())
		return nil, err
	}
	self.wh[0] = self.canvas.Get("clientWidth").Int()
	self.wh[1] = self.canvas.Get("clientHeight").Int()
	// self.width = doc.Get("body").Get("clientWidth").Int()
	// self.height = doc.Get("body").Get("clientHeight").Int()
	// Contrary to the usual html elements, a Canvas element needs it's width and height attributes for logical size.
	// (CSS width and height you set in HTML only stretches the result, and it may cause blurry image)
	// Ref: https://stackoverflow.com/questions/4938346/canvas-width-and-height-in-html5
	self.canvas.Set("width", self.wh[0])  // IMPORTANT!
	self.canvas.Set("height", self.wh[1]) // IMPORTANT!
	// context.Call("viewport", 0, 0, camera.wh[0], camera.wh[1]) // (LowerLeft.x, LowerLeft.y, width, height)
	// (if 'viewport' is not updated, rendering may blur after window.resize)
	self.mouse_wheel_scale = 500 // in the range of [0 ~ 500(default) ~ 1000]
	// create WebGL context
	self.rc, err = NewWebGLRenderingContext(self.canvas)
	if err != nil {
		js.Global().Call("alert", "Failed to start WebGL : "+err.Error())
		return nil, err
	}
	return &self, nil
}

func (self *WebGLCanvas) GetRenderingContext() gigl.GLRenderingContext {
	return (self.rc)
}

func (self *WebGLCanvas) String() string {
	return fmt.Sprintf("WebGLCanvas{id:'%s' size:%dx%d}\n", self.id, self.wh[0], self.wh[1])
}

func (self *WebGLCanvas) GetWebGLRenderingContext() js.Value {
	return self.rc.context
}

func (self *WebGLCanvas) GetWebGLConstants() gigl.GLConstants {
	return self.rc.constants
}

func (self *WebGLCanvas) ConvertGoSliceToJsTypedArray(a interface{}) js.Value {
	return self.rc.ConvertGoSliceToJsTypedArray(a)
}

// ----------------------------------------------------------------------------
// User Interactions (Event Handling)
// ----------------------------------------------------------------------------

func (self *WebGLCanvas) setup_mouse_event_common_handler() {
	self.mouse_event_common_handler = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			fmt.Println("Invalid GoCallback call (for EventHandling) from Javascript")
			return nil
		}
		event := args[0]                    // js.Value (event object)
		etype := event.Get("type").String() // canvas := event.Get("srcElement")
		switch etype {
		case "click":
			cxy := [2]int{event.Get("clientX").Int(), event.Get("clientY").Int()}
			dx, dy := (cxy[0] - self.mouse_sxy[0]), (cxy[1] - self.mouse_sxy[1])
			keystat := [4]bool{event.Get("altKey").Bool(), event.Get("ctrlKey").Bool(), event.Get("metaKey").Bool(), event.Get("shiftKey").Bool()}
			if dx < -3 || dx > +3 || dy < -3 || dy > +3 {
				// ignore
			} else if self.evthandler_for_click != nil {
				self.evthandler_for_click(cxy, keystat)
			} else {
				fmt.Printf("%s (%d %d) %v\n", etype, cxy[0], cxy[1], keystat)
			}
		case "dblclick":
			cxy := [2]int{event.Get("clientX").Int(), event.Get("clientY").Int()}
			keystat := [4]bool{event.Get("altKey").Bool(), event.Get("ctrlKey").Bool(), event.Get("metaKey").Bool(), event.Get("shiftKey").Bool()}
			if self.evthandler_for_dblclick != nil {
				self.evthandler_for_dblclick(cxy, keystat)
			} else {
				fmt.Printf("%s (%d %d) %v\n", etype, cxy[0], cxy[1], keystat)
			}
		case "mousemove":
			if self.mouse_dragging {
				cxy := [2]int{event.Get("clientX").Int(), event.Get("clientY").Int()}
				dxy := [2]int{event.Get("movementX").Int(), event.Get("movementY").Int()}
				keystat := [4]bool{event.Get("altKey").Bool(), event.Get("ctrlKey").Bool(), event.Get("metaKey").Bool(), event.Get("shiftKey").Bool()}
				if self.evthandler_for_mouse_drag != nil {
					self.evthandler_for_mouse_drag(cxy, dxy, keystat)
				} else {
					fmt.Printf("%s (%d %d) with %v\n", etype, dxy[0], dxy[1], keystat)
				}
			} else {
				if self.evthandler_for_mouse_over != nil {
					cxy := [2]int{event.Get("clientX").Int(), event.Get("clientY").Int()}
					keystat := [4]bool{event.Get("altKey").Bool(), event.Get("ctrlKey").Bool(), event.Get("metaKey").Bool(), event.Get("shiftKey").Bool()}
					self.evthandler_for_mouse_over(cxy, keystat)
				}
			}
		case "mousedown":
			self.mouse_dragging = true
			self.mouse_sxy = [2]int{event.Get("clientX").Int(), event.Get("clientY").Int()}
		case "mouseup":
			self.mouse_dragging = false
		case "mouseleave":
			self.mouse_dragging = false
		default:
			fmt.Println(etype)
		}
		return nil
	})
	self.canvas.Call("addEventListener", "click", self.mouse_event_common_handler)
	self.canvas.Call("addEventListener", "dblclick", self.mouse_event_common_handler)
	self.canvas.Call("addEventListener", "mousemove", self.mouse_event_common_handler)
	self.canvas.Call("addEventListener", "mousedown", self.mouse_event_common_handler)
	self.canvas.Call("addEventListener", "mouseup", self.mouse_event_common_handler)
	self.canvas.Call("addEventListener", "mouseleave", self.mouse_event_common_handler)
}

func (self *WebGLCanvas) SetEventHandlerForClick(handler func(canvasxy [2]int, keystat [4]bool)) {
	self.evthandler_for_click = handler
	if self.mouse_event_common_handler.IsUndefined() {
		self.setup_mouse_event_common_handler()
	}
}

func (self *WebGLCanvas) SetEventHandlerForDoubleClick(handler func(canvasxy [2]int, keystat [4]bool)) {
	self.evthandler_for_dblclick = handler
	if self.mouse_event_common_handler.IsUndefined() {
		self.setup_mouse_event_common_handler()
	}
}

func (self *WebGLCanvas) SetEventHandlerForMouseOver(handler func(canvasxy [2]int, keystat [4]bool)) {
	self.evthandler_for_mouse_over = handler
	if self.mouse_event_common_handler.IsUndefined() {
		self.setup_mouse_event_common_handler()
	}
}

func (self *WebGLCanvas) SetEventHandlerForMouseDrag(handler func(canvasxy [2]int, dxy [2]int, keystat [4]bool)) {
	self.evthandler_for_mouse_drag = handler
	if self.mouse_event_common_handler.IsUndefined() {
		self.setup_mouse_event_common_handler()
	}
}

func (self *WebGLCanvas) setup_mouse_wheel_common_handler() {
	// For zooming,   'handler()' is given 2nd argument of 'scale' in the range of [ 0.01 ~ 1(default) ~ 100.0 ]
	// For scrolling, 'handler()' is given 2nd argument of 'delta' in the range of [ -200 ~ 0 ~ +200 ]
	js_handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0] // js.Value (event object), event.Get("type"), event.Get("srcElement")
		keystat := [4]bool{event.Get("altKey").Bool(), event.Get("ctrlKey").Bool(), event.Get("metaKey").Bool(), event.Get("shiftKey").Bool()}
		if keystat[3] { // ZOOM, if SHIFT is was pressed
			if self.evthandler_for_zoom != nil {
				cxy := [2]int{event.Get("clientX").Int(), event.Get("clientY").Int()}
				delta := float64(event.Get("deltaY").Int())
				if math.Abs(delta) > 100 { // on Windows, mouse wheel delta is too big (+/-125)
					delta = delta * 0.1
				}
				self.mouse_wheel_scale += delta // [ 0 ~ 500(default) ~ 1000 ]
				self.mouse_wheel_scale = float64(math.Max(0, math.Min(self.mouse_wheel_scale, 1000)))
				scale_exp := (self.mouse_wheel_scale - 500.0) / 250.0 // [ -2 ~ 0(default) ~ +2 ]
				scale := math.Pow(10, scale_exp)                      // [ 0.01 ~ 1(default) ~ 100.0 ]
				self.evthandler_for_zoom(cxy, float32(scale), keystat)
			}
		} else { // SCROLL
			if self.evthandler_for_scroll != nil {
				cxy := [2]int{event.Get("clientX").Int(), event.Get("clientY").Int()}
				dx, dy := event.Get("deltaX").Int(), event.Get("deltaY").Int()
				self.evthandler_for_scroll(cxy, dx, dy, keystat)
			}
		}
		return nil
	})
	self.canvas.Call("addEventListener", "wheel", js_handler)
}

func (self *WebGLCanvas) SetEventHandlerForZoom(handler func(canvasxy [2]int, scale float32, keystat [4]bool)) {
	// 'scale' in the range of [ 0.01 ~ 1(default) ~ 100.0 ]
	self.evthandler_for_zoom = handler
	if self.mouse_wheel_common_handler.IsUndefined() {
		self.setup_mouse_wheel_common_handler()
	}
}

func (self *WebGLCanvas) SetEventHandlerForScroll(handler func(canvasxy [2]int, dx int, dy int, keystat [4]bool)) {
	// 'scroll' in the range of [ -200 ~ 0 ~ +200 ] 	// (-): swipe_down, (+): swipe_up
	self.evthandler_for_scroll = handler
	if self.mouse_wheel_common_handler.IsUndefined() {
		self.setup_mouse_wheel_common_handler()
	}
}

func (self *WebGLCanvas) SetEventHandlerForKeyPress(handler func(key string, code string, keystat [4]bool)) {
	js_handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0] // js.Value (event object), event.Get("type"), event.Get("srcElement")
		if handler != nil {
			key, code := event.Get("key").String(), event.Get("code").String()
			keystat := [4]bool{event.Get("altKey").Bool(), event.Get("ctrlKey").Bool(), event.Get("metaKey").Bool(), event.Get("shiftKey").Bool()}
			handler(key, code, keystat)
		}
		return nil
	})
	// js.Global().Get("document").Call("addEventListener", "keypress", js_handler)
	js.Global().Get("document").Call("addEventListener", "keydown", js_handler) // ARROW keys are captured by 'keydown' only
}

func (self *WebGLCanvas) SetEventHandlerForWindowResize(handler func(w int, h int)) {
	js_handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		w := js.Global().Get("window").Get("innerWidth").Int()
		h := js.Global().Get("window").Get("innerHeight").Int()
		if handler != nil {
			handler(w, h)
		} else {
			fmt.Printf("window.resize %d %d\n", w, h)
		}
		return nil
	})
	js.Global().Get("window").Call("addEventListener", "resize", js_handler)
}

// ----------------------------------------------------------------------------
// Animating with DrawHandler
// ----------------------------------------------------------------------------

func (self *WebGLCanvas) Run(draw_handler func(now float64)) {
	// run UI animation loop forever, with the given 'draw_handler'
	self.paused = false
	if draw_handler != nil {
		self.user_handler_for_draw = draw_handler
		self.wasm_handler_for_draw = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if self.user_handler_for_draw != nil && !self.paused {
				now := args[0].Float() // DOMHighResTimeStamp similar to performance.now()
				self.user_handler_for_draw(now)
			}
			js.Global().Call("requestAnimationFrame", self.wasm_handler_for_draw)
			return nil
		})
		js.Global().Call("requestAnimationFrame", self.wasm_handler_for_draw)
		// What it actually does is like:
		//   requestAnimationFrame(drawHandlerForAnimationFrame);
		//   function drawHandlerForAnimationFrame() {
		//     if draw_handler != nil {
		//         draw_handler();   // draw the scene by calling Go renderer function
		//     }
		//     requestAnimationFrame(drawHandlerForAnimationFrame); // call itself again for the next frame
		//   }
	}
	<-make(chan bool) // wait for events (without exiting)
}

func (self *WebGLCanvas) RunOnce(draw_handler func(now float64)) {
	// run UI animation loop only once, with the given 'draw_handler'
	if draw_handler != nil {
		self.user_handler_for_draw = draw_handler
		self.wasm_handler_for_draw = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if self.user_handler_for_draw != nil {
				now := args[0].Float() // DOMHighResTimeStamp similar to performance.now()
				self.user_handler_for_draw(now)
			}
			return nil
		})
		js.Global().Call("requestAnimationFrame", self.wasm_handler_for_draw)
	}
	<-make(chan bool) // wait for events (without exiting)
}

func (self *WebGLCanvas) Pause() {
	self.paused = true
}

func (self *WebGLCanvas) Resume() {
	self.paused = false
}
