package g3d

import (
	"fmt"
	"math"

	"github.com/go4orward/gigl"
	"github.com/go4orward/gigl/common"
)

// ----------------------------------------------------------------------------
// SceneObject
// ----------------------------------------------------------------------------

type SceneObject struct {
	Geometry    gigl.GLGeometry // geometry interface
	Material    gigl.GLMaterial // material
	VShader     gigl.GLShader   // vert shader and its bindings
	EShader     gigl.GLShader   // edge shader and its bindings
	FShader     gigl.GLShader   // face shader and its bindings
	modelmatrix common.Matrix4  //
	UseDepth    bool            // depth test flag (default is true)
	UseBlend    bool            // blending flag with alpha (default is false)
	children    []*SceneObject  //
	// multiple instance poses
	instance_count  int       // number of instances
	instance_stride int       // number of values of a single pose
	instance_buffer []float32 //
	// VAO (set of RenderingContext buffers)
	vao *gigl.VAO //
	//
	err error
}

func NewSceneObject(geometry gigl.GLGeometry, material gigl.GLMaterial,
	vshader gigl.GLShader, eshader gigl.GLShader, fshader gigl.GLShader) *SceneObject {
	// 'geometry' : geometric shape (vertices, edges, faces) to be rendered
	// 'material' : color, texture, or other material properties	: OPTIONAL (can be 'nil')
	// 'vshader' : shader for VERTICES (POINTS) 					: OPTIONAL (can be 'nil')
	// 'eshader' : shader for EDGES (LINES) 						: OPTIONAL (can be 'nil')
	// 'fshader' : shader for FACES (TRIANGLES) 					: OPTIONAL (can be 'nil')
	// Note that geometry & material & shader can be shared among different SceneObjects.
	if geometry == nil {
		return nil
	}
	// Note that 'material' & 'shader' can be nil, in which case its parent's 'material' & 'shader' will be used to render.
	sobj := SceneObject{Geometry: geometry, Material: material, VShader: vshader, EShader: eshader, FShader: fshader}
	sobj.modelmatrix.SetIdentity()
	sobj.UseDepth = true  // depth test is turned on by default
	sobj.UseBlend = false // alpha blending is turned off by default
	sobj.children = nil
	return &sobj
}

func (self *SceneObject) IsReady() bool {
	return self.Geometry != nil && self.err == nil
}

func (self *SceneObject) Summary() string {
	summary := "SceneObject " + self.Geometry.Summary()
	if self.instance_buffer != nil {
		summary += fmt.Sprintf("  Instancess : count=%d stride=%d\n", self.instance_count, self.instance_stride)
	}
	if self.Material != nil {
		summary += fmt.Sprintf("  %s\n", self.Material.MaterialSummary())
	}
	if self.VShader != nil {
		summary += fmt.Sprintf("  VERT %s\n", self.VShader.Summary())
	}
	if self.EShader != nil {
		summary += fmt.Sprintf("  EDGE %s\n", self.EShader.Summary())
	}
	if self.FShader != nil {
		summary += fmt.Sprintf("  FACE %s\n", self.FShader.Summary())
	}
	summary += fmt.Sprintf("  Flags    : UseDepth=%t  UseBlend=%t\n", self.UseDepth, self.UseBlend)
	summary += fmt.Sprintf("  Children : %d", len(self.children))
	return summary
}

// ----------------------------------------------------------------------------
// Basic Access
// ----------------------------------------------------------------------------

func (self *SceneObject) AddChild(child *SceneObject) *SceneObject {
	if self.children == nil {
		self.children = make([]*SceneObject, 0)
	}
	self.children = append(self.children, child)
	return self
}

func (self *SceneObject) GetModelMatrix() *common.Matrix4 {
	return &self.modelmatrix
}

func (self *SceneObject) GetChildren() []*SceneObject {
	return self.children
}

// ----------------------------------------------------------------------------
// Multiple Instance Poses
// ----------------------------------------------------------------------------

func (self *SceneObject) ClearInstanceBuffer() {
	self.instance_buffer = nil
	self.instance_count = 0
	self.instance_stride = 0
}

func (self *SceneObject) SetInstanceBuffer(instance_count int, instance_stride int, data []float32) *SceneObject {
	// This function is OPTIONAL (only if multiple instances of the geometry are rendered)
	self.instance_buffer = make([]float32, instance_count*instance_stride)
	self.instance_count = instance_count
	self.instance_stride = instance_stride
	if data != nil {
		for i := 0; i < len(self.instance_buffer) && i < len(data); i++ {
			self.instance_buffer[i] = data[i]
		}
	}
	return self
}

func (self *SceneObject) SetInstancePoseValues(instance_index int, offset int, values ...float32) {
	// This function is OPTIONAL (only if multiple instances of the geometry are rendered)
	if (offset + len(values)) > self.instance_stride {
		common.Logger.Error("SetInstancePoseValues() failed : invalid offset (%d) and value count (%d) for the stride (%d)\n", offset, len(values), self.instance_stride)
		return
	}
	pos := instance_index * self.instance_stride
	for i := 0; i < len(values); i++ {
		self.instance_buffer[pos+offset+i] = values[i]
	}
}

func (self *SceneObject) SetInstanceColorValues(instance_index int, offset int, v0 uint8, v1 uint8, v2 uint8, v3 uint8) {
	// This function is OPTIONAL (only if multiple instances of the geometry are rendered)
	if (offset + 1) > self.instance_stride {
		common.Logger.Error("SetInstanceColorValue() failed : invalid offset (%d) for the stride (%d)\n", offset, self.instance_stride)
		return
	}
	pos := instance_index * self.instance_stride
	b0, b1, b2, b3 := uint32(0), uint32(0), uint32(0), uint32(0)
	self.instance_buffer[pos+offset] = math.Float32frombits(b0 + b1<<8 + b2<<16 + b3<<24) // LittleEndian (lower byte comes first)
}

// ----------------------------------------------------------------------------
// Translation, Rotation, Scaling (by manipulating MODEL matrix)
// ----------------------------------------------------------------------------

func (self *SceneObject) SetTransformation(txyz [3]float32, axis [3]float32, angle_in_degree float32, sxyz [3]float32) *SceneObject {
	translation := common.NewMatrix4().SetTranslation(txyz[0], txyz[1], txyz[2])
	rotation := common.NewMatrix4().SetRotationByAxis(axis, angle_in_degree)
	scaling := common.NewMatrix4().SetScaling(sxyz[0], sxyz[1], sxyz[2])
	self.modelmatrix.SetMultiplyMatrices(translation, rotation, scaling)
	return self
}

func (self *SceneObject) Translate(tx float32, ty float32, tz float32) *SceneObject {
	translation := common.NewMatrix4().SetTranslation(tx, ty, tz)
	self.modelmatrix.SetMultiplyMatrices(translation, &self.modelmatrix)
	return self
}

func (self *SceneObject) Rotate(axis [3]float32, angle_in_degree float32) *SceneObject {
	rotation := common.NewMatrix4().SetRotationByAxis(axis, angle_in_degree)
	self.modelmatrix.SetMultiplyMatrices(rotation, &self.modelmatrix)
	return self
}

func (self *SceneObject) Scale(sx float32, sy float32, sz float32) *SceneObject {
	scaling := common.NewMatrix4().SetScaling(sx, sy, sz)
	self.modelmatrix.SetMultiplyMatrices(scaling, &self.modelmatrix)
	return self
}
