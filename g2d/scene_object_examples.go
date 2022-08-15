package g2d

import (
	"math"

	"github.com/go4orward/gigl"
)

func NewSceneObject_RedTriangle(rc gigl.GLRenderingContext) *SceneObject {
	// This example creates a red triangle with radius 0.5 at (0,0)
	geometry := NewGeometryTriangle(0.5)         // create a triangle with radius 0.5 at (0,0)
	geometry.BuildDataBuffers(true, false, true) // build data buffers for vertices and faces
	material := NewMaterialColors("#aa0000")
	shader := NewShaderForMaterialColors(rc)                    // create shader, and set its bindings
	return NewSceneObject(geometry, material, nil, nil, shader) // set up the scene object (draw FACES only)
}

func NewSceneObject_HexagonWireframe(rc gigl.GLRenderingContext) *SceneObject {
	// This example creates a hexagon with given color and radius 0.5 at (0,0), to be rendered as 'wireframe'
	// (This example demonstrates how 'triangulation of face' works - for faces with more than 3 vertices)
	geometry := NewGeometryPolygon(6, 0.5, 30) // create a hexagon with radius 0.5, with 1st vertex at 30 degree from X axis
	geometry.BuildDataBuffersForWireframe()    // extract wireframe edges from faces
	material := NewMaterialColors("#888888")
	shader := NewShaderForMaterialColors(rc)                    // create shader, and set its bindings
	return NewSceneObject(geometry, material, nil, shader, nil) // set up the scene object (draw EDGES only)
}

func NewSceneObject_RectangleInstancesExample(rc gigl.GLRenderingContext) *SceneObject {
	// This example creates 200*80 instances of a single geometry, each with its own position and color
	geometry := NewGeometryRectangle(0.8)        // create a rectangle of size 1.0
	geometry.BuildDataBuffers(true, false, true) //
	material := NewMaterialColors("#888888")
	shader := NewShaderForInstancePoseColor(rc)                    // create shader, and set its bindings
	scnobj := NewSceneObject(geometry, material, nil, nil, shader) // set up the scene object (draw FACES only)
	scnobj.SetInstanceBuffer(200*80, 5, nil)                       // multiple instances for 200*80 rectangles
	for row := 0; row < 200; row++ {
		for col := 0; col < 80; col++ {
			// r, g, b := float32(1.0), float32(0.1), float32(0.1)
			ii, jj := math.Abs(float64(row)-100)/100, math.Abs(float64(col)-40)/40
			r, g, b := float32(ii), float32(jj), 1-float32((ii+jj)/2)
			// r, g, b := uint8(ii*255), uint8(jj*255), uint8((1-(ii+jj)/2)*255)
			// r, g, b = uint8(255), uint8(255), uint8(255)
			scnobj.SetInstancePoseValues(row*80+col, 0, float32(col), float32(row)) // position (tx,ty)
			scnobj.SetInstancePoseValues(row*80+col, 2, r, g, b)                    // color RGB
		}
	}
	return scnobj
}
