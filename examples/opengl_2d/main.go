package main

import (
	"errors"
	"log"
	"runtime"

	"github.com/go4orward/gigl/env/opengl41"
	"github.com/go4orward/gigl/g2d"
)

func init() { // This is needed to let main() run on the startup thread.
	runtime.LockOSThread() // Ref: https://golang.org/pkg/runtime/#LockOSThread
}

func main() {
	window, err := opengl41.NewOpenGLWindow(800, 600, "OpenGL2D: Triangle in World Space", false)
	if err != nil {
		log.Fatal(errors.New("Failed to create OpenGL window : " + err.Error()))
	}

	rc := window.GetGLRenderingContext()
	geometry := g2d.NewGeometry_Triangle(0.5)   // create geometry (a triangle with radius 0.5)
	geometry.BuildDataBuffers(true, true, true) // build data buffers for vertices, edges and faces
	material, _ := rc.CreateMaterial("#bbbbff") // create material (with light-blue color)
	material.SetColorForDrawMode(2, "#0000ff")  // set edge color for the material
	shader := g2d.NewShader_MaterialColor(rc)   // shader with auto-binded color & PVM matrix
	scnobj := g2d.NewSceneObject(geometry, material, nil, shader, shader).Rotate(40)
	scene := g2d.NewScene("#ffffff").Add(scnobj) // scene holds all the SceneObjects to be rendered
	camera := g2d.NewCamera(rc.GetWH(), 2, 1)    // FOV 2 means range of [-1,+1] in X, ZoomLevel is 1.0
	renderer := g2d.NewRenderer(rc)              // set up the renderer

	window.Run(func(now float64) {
		renderer.Clear(scene)               // prepare to render (clearing with the bkg color of the scene)
		renderer.RenderScene(scene, camera) // render the scene (iterating over all the SceneObjects in it)
		renderer.RenderAxes(camera, 1.0)    // render the axes (just for visual reference)
	})
}
