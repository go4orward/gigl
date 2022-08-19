package main

import (
	"errors"
	"log"
	"runtime"

	opengl "github.com/go4orward/gigl/env/opengl41"
)

func init() { // This is needed to let main() run on the startup thread.
	runtime.LockOSThread() // Ref: https://golang.org/pkg/runtime/#LockOSThread
}

func main() {
	canvas, err := opengl.NewOpenGLCanvas(1200, 900, "OpenGL1st: Triangle in Clip Space", false)
	if err != nil {
		log.Fatal(errors.New("Failed to create OpenGL canvas : " + err.Error()))
	}

	opengl.DrawSimplestTriangle(canvas)

	canvas.SwapBuffers()
	canvas.RunOnce(nil)
}
