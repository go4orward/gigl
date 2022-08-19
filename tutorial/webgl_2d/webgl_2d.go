package main

import (
	"github.com/go4orward/gigl/common"
	webgl "github.com/go4orward/gigl/env/webgl10"
	"github.com/go4orward/gigl/g2d"
)

type Config struct {
	loglevel  string //
	logfilter string //
}

func main() {
	cfg := Config{loglevel: "info", logfilter: ""}
	if cfg.loglevel != "" {
		common.SetLogger(common.NewConsoleLogger(cfg.loglevel)).SetTraceFilter(cfg.logfilter).SetOption("", false)
	}
	// THIS CODE IS SUPPOSED TO BE BUILT AS WEBASSEMBLY AND RUN INSIDE A BROWSER.
	// BUILD IT LIKE 'GOOS=js GOARCH=wasm go build -o example.wasm examples/example.go'.
	canvas, err := webgl.NewWebGLCanvas("wasmcanvas") // ID of canvas element
	if err != nil {
		common.Logger.Error("Failed to start WebGL : %v\n", err)
		return
	}
	rc := canvas.GetRenderingContext()
	geometry := g2d.NewGeometryTriangle(0.5)                          // create geometry (a triangle with radius 0.5)
	geometry.BuildDataBuffers(true, true, true)                       // build data buffers for vertices, edges and faces
	mcolors := g2d.NewMaterialColors("#bbbbff", "#bbbbff", "#0000ff") // create material (with light-blue color)
	shader := g2d.NewShaderForMaterialColors(rc)                      // shader with auto-binded color & PVM matrix
	scnobj := g2d.NewSceneObject(geometry, mcolors, nil, shader, shader).Rotate(40)
	scene := g2d.NewScene("#ffffff").Add(scnobj) // scene holds all the SceneObjects to be rendered
	camera := g2d.NewCamera(rc.GetWH(), 2, 1)    // FOV 2 means range of [-1,+1] in X, ZoomLevel is 1.0
	renderer := g2d.NewRenderer(rc)              // set up the renderer

	canvas.RunOnce(func(now float64) {
		renderer.Clear(scene)               // prepare to render (clearing to white background)
		renderer.RenderScene(scene, camera) // render the scene (iterating over all the SceneObjects in it)
		renderer.RenderAxes(camera, 1.0)    // render the axes (just for visual reference)
	})
}
