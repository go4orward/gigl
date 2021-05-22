package main

import (
	"errors"
	"fmt"
	"log"
	"runtime"

	"github.com/go4orward/gigl/env/opengl41"
	"github.com/go4orward/gigl/world"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	window, err := opengl41.NewOpenGLWindow(800, 600, "OpenGL Globe: The Blue Marble", false)
	if err != nil {
		log.Fatal(errors.New("Failed to create OpenGL window : " + err.Error()))
	}
	rc := window.GetGLRenderingContext()
	wglobe := world.NewWorldGlobe(rc, "#000000", "./assets/gopher.png") // Globe radius is assumed to be 1.0
	wcamera := world.NewWorldCamera(rc.GetWH(), 15, 1.0)                // camera FOV default is 15° (in degree)
	wcamera.SetPoseByLonLat(0, 0, 10)                                   // longitude 0°, latitude 0°, radius(distance) 10.0
	renderer := world.NewWorldRenderer(rc)                              // set up the world renderer

	// add user interactions (with mouse)
	window.SetEventHandlerForDoubleClick(func(canvasxy [2]int, keystat [4]bool) {
		wcamera.ShowInfo()
	})
	window.SetEventHandlerForMouseDrag(func(canvasxy [2]int, dxy [2]int, keystat [4]bool) {
		wcamera.RotateAroundGlobe(float32(dxy[0])*0.2, float32(dxy[1])*0.2)
	})
	window.SetEventHandlerForZoom(func(canvasxy [2]int, scale float32, keystat [4]bool) {
		wcamera.SetZoom(scale) // 'scale' in [ 0.01 ~ 1(default) ~ 100.0 ]
	})
	window.SetEventHandlerForWindowResize(func(w int, h int) {
		wcamera.SetAspectRatio(w, h)
	})
	fmt.Println("Try mouse drag & wheel with SHIFT key pressed") // printed in the browser console

	// run UI animation loop
	window.Run(func(now float64) {
		renderer.Clear(wglobe)                  // prepare to render (clearing to black background)
		renderer.RenderWorld(wglobe, wcamera)   // render the Globe (and all the layers & glowring)
		wglobe.Rotate([3]float32{0, 0, 1}, 0.1) // rotate the Globe
	})
}
