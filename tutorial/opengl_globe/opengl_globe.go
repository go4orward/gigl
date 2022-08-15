package main

import (
	"errors"
	"fmt"
	"log"
	"runtime"

	"github.com/go4orward/gigl/earth"
	opengl "github.com/go4orward/gigl/env/opengl41"
)

func init() { // This is needed to let main() run on the startup thread.
	runtime.LockOSThread() // Ref: https://golang.org/pkg/runtime/#LockOSThread
}

func main() {
	canvas, err := opengl.NewOpenGLCanvas(800, 600, "OpenGL Globe: The Blue Marble", false)
	if err != nil {
		log.Fatal(errors.New("Failed to create OpenGL canvas : " + err.Error()))
	}
	rc := canvas.GetRenderingContext()
	wglobe := earth.NewWorldGlobe(rc, "#000000", "./assets/world.png") // Globe radius is assumed to be 1.0
	wcamera := earth.NewWorldCamera(rc.GetWH(), 15, 1.0)               // camera FOV default is 15° (in degree)
	wcamera.SetPoseByLonLat(0, 0, 10)                                  // longitude 0°, latitude 0°, radius(distance) 10.0
	renderer := earth.NewWorldRenderer(rc)                             // set up the world renderer

	// add user interactions (with mouse)
	canvas.SetEventHandlerForDoubleClick(func(canvasxy [2]int, keystat [4]bool) {
		fmt.Printf("%s\n", wcamera.Summary())
	})
	canvas.SetEventHandlerForMouseDrag(func(canvasxy [2]int, dxy [2]int, keystat [4]bool) {
		wcamera.RotateAroundGlobe(float32(dxy[0])*0.2, float32(dxy[1])*0.2)
	})
	canvas.SetEventHandlerForZoom(func(canvasxy [2]int, scale float32, keystat [4]bool) {
		wcamera.SetZoom(scale) // 'scale' in [ 0.01 ~ 1(default) ~ 100.0 ]
	})
	canvas.SetEventHandlerForWindowResize(func(w int, h int) {
		wcamera.SetAspectRatio(w, h)
	})
	fmt.Println("Try mouse drag & wheel with SHIFT key pressed") // printed in the browser console

	// run UI animation loop
	canvas.Run(func(now float64) {
		renderer.Clear(wglobe)                  // prepare to render (clearing to black background)
		renderer.RenderWorld(wglobe, wcamera)   // render the Globe (and all the layers & glowring)
		wglobe.Rotate([3]float32{0, 0, 1}, 0.1) // rotate the Globe
	})
}
