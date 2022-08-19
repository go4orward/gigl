package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go4orward/gigl/common"
	opengl "github.com/go4orward/gigl/env/opengl41"
	"github.com/go4orward/gigl/g2d"
	"github.com/go4orward/gigl/g3d"
)

type Config struct {
	loglevel  string //
	logfilter string //
}

func main() {
	cfg := Config{loglevel: "info", logfilter: ""}
	flag.StringVar(&cfg.loglevel, "L", cfg.loglevel, "log level (error/warn/info/debug/trace)")
	flag.StringVar(&cfg.logfilter, "F", cfg.logfilter, "log filter for trace log messages")
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Printf("Usage:  geometry_viewer  INPUT_FILE\n")
		flag.PrintDefaults()
		os.Exit(0)
	}
	if cfg.loglevel != "" {
		common.SetLogger(common.NewConsoleLogger(cfg.loglevel)).SetTraceFilter(cfg.logfilter).SetOption("", false)
	}
	input_geometry := flag.Arg(0)
	var geometry *g3d.Geometry = nil
	switch input_geometry {
	case "cube":
		geometry = g3d.NewGeometryCube(1, 1, 1)
	case "sphere":
		geometry = g3d.NewGeometrySphere(1.0, 36, 18)
	case "cylinder":
		geometry = g3d.NewGeometryCylinder(8, 1.0, 3.0, 0, true)
	default:
		geometry = g3d.NewGeometryCylinder(8, 1.0, 3.0, 0, true)
		// TODO(go4orward) - load models from OBJ or STL file format
	}

	// TODO(go4orward) - crashes on OpenGL on Mac

	canvas, err := opengl.NewOpenGLCanvas(800, 600, "Geometry Viewer", false)
	if err != nil {
		log.Fatal(errors.New("Failed to create OpenGL canvas : " + err.Error()))
	}
	rc := canvas.GetRenderingContext()

	geometry.BuildNormalsForFace()               // calculate normal vectors for each face
	geometry.BuildDataBuffers(true, false, true) // build data buffers for vertices and faces
	material := g2d.NewMaterialColors("#ffffff") // create material (with texture image)
	shader := g3d.NewShader_NormalColor(rc)      // use the standard NORMAL+Color shader
	scene := g3d.NewScene("#000000")
	scene.Add(g3d.NewSceneObject(geometry, material, nil, nil, shader))
	// llayer := g3d.NewOverlayLabelLayer(rc, 20, true).AddLabelsForTest()
	// mlayer := g3d.NewOverlayMarkerLayer(rc).AddMarkersForTest()
	// scene.AddOverlay(mlayer)
	cam_ip := g3d.CamInternalParams{WH: rc.GetWH(), Fov: 15, Zoom: 1.0, NearFar: [2]float32{1, 100}}
	cam_ep := g3d.CamExternalPose{From: [3]float32{0, 0, 10}, At: [3]float32{0, 0, 0}, Up: [3]float32{0, 1, 0}}
	camera := g3d.NewCamera(true, &cam_ip, &cam_ep)
	renderer := g3d.NewRenderer(rc) // set up the renderer

	canvas.Run(func(now float64) {
		renderer.Clear(scene)               // prepare to render (clearing to white background)
		renderer.RenderScene(scene, camera) // render the scene (iterating over all the SceneObjects in it)
		renderer.RenderAxes(camera, 1.0)    // render the axes (just for visual reference)
		scene.Get(0).Rotate([3]float32{0, 1, 1}, 1.0)
	})
}
