package g2d

import (
	"math"

	"github.com/go4orward/gigl/common"
)

func NewSceneObject_2DAxes(rc common.GLRenderingContext, length float32) *SceneObject {
	// This example creates two lines for X (red) and Y (green) axes, with origin at (0,0)
	geometry := NewGeometry()                                            // create an empty geometry
	geometry.SetVertices([][2]float32{{0, 0}, {length, 0}, {0, length}}) // add three vertices
	geometry.SetEdges([][]uint32{{0, 1}, {0, 2}})                        // add two edges
	geometry.BuildDataBuffers(true, true, false)                         // build data buffers for vertices and edges
	shader := NewShader_2DAxes(rc)                                       // create shader, and set its bindings
	return NewSceneObject(geometry, nil, nil, shader, nil)               // set up the scene object (draw LINES)
}

func NewSceneObject_RedTriangle(rc common.GLRenderingContext) *SceneObject {
	// This example creates a red triangle with radius 0.5 at (0,0)
	geometry := NewGeometry_Triangle(0.5)                  // create a triangle with radius 0.5 at (0,0)
	geometry.BuildDataBuffers(true, false, true)           // build data buffers for vertices and faces
	shader := NewShader_MaterialColor(rc)                  // create shader, and set its bindings
	return NewSceneObject(geometry, nil, nil, nil, shader) // set up the scene object (draw FACES only)
}

func NewSceneObject_HexagonWireframe(rc common.GLRenderingContext) *SceneObject {
	// This example creates a hexagon with given color and radius 0.5 at (0,0), to be rendered as 'wireframe'
	// (This example demonstrates how 'triangulation of face' works - for faces with more than 3 vertices)
	geometry := NewGeometry_Polygon(6, 0.5, 30)                 // create a hexagon with radius 0.5, with 1st vertex at 30 degree from X axis
	geometry.BuildDataBuffersForWireframe()                     // extract wireframe edges from faces
	material, _ := rc.CreateMaterial("#888888")                 // create material
	shader := NewShader_MaterialColor(rc)                       // create shader, and set its bindings
	return NewSceneObject(geometry, material, nil, shader, nil) // set up the scene object (draw EDGES only)
}

func NewSceneObject_RectangleInstancesExample(rc common.GLRenderingContext) *SceneObject {
	// This example creates 200*80 instances of a single geometry, each with its own position and color
	geometry := NewGeometry_Rectangle(0.8)                         // create a rectangle of size 1.0
	geometry.BuildDataBuffers(true, false, true)                   //
	material, _ := rc.CreateMaterial("#888888")                    // create material
	shader := NewShader_InstancePoseColor(rc)                      // create shader, and set its bindings
	scnobj := NewSceneObject(geometry, material, nil, nil, shader) // set up the scene object (draw FACES only)
	scnobj.SetupPoses(5, 200*80, nil)                              // multiple poses for 200*80 rectangle instances
	for row := 0; row < 200; row++ {
		for col := 0; col < 80; col++ {
			ii, jj := math.Abs(float64(row)-100)/100, math.Abs(float64(col)-40)/40
			r, g, b := float32(ii), float32(jj), 1-float32((ii+jj)/2)
			scnobj.SetPoseValues(row*80+col, 0, float32(col), float32(row)) // position (tx,ty)
			scnobj.SetPoseValues(row*80+col, 2, r, g, b)                    // color    (r,g,b)
		}
	}
	return scnobj
}
