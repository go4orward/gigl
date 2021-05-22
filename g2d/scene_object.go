package g2d

import (
	"fmt"
	"math"

	"github.com/go4orward/gigl"
	"github.com/go4orward/gigl/common"
	"github.com/go4orward/gigl/g2d/c2d"
)

type SceneObject struct {
	Geometry    *Geometry       // geometry interface
	Material    gigl.GLMaterial // material
	VShader     gigl.GLShader   // vert shader and its bindings
	EShader     gigl.GLShader   // edge shader and its bindings
	FShader     gigl.GLShader   // face shader and its bindings
	modelmatrix common.Matrix3  // model transformation matrix of this SceneObject
	UseDepth    bool            // depth test flag (default is true)
	UseBlend    bool            // blending flag with alpha (default is false)
	children    []*SceneObject  // OPTIONAL, children of this SceneObject (to be rendered recursively)
	bbox        [2][2]float32   // bounding box
	// multiple instance poses
	instance_count  int       // number of instances
	instance_stride int       // number of values of a single pose
	instance_buffer []float32 //
	// VAO (set of RenderingContext buffers)
	vao *gigl.VAO //
}

func NewSceneObject(geometry *Geometry, material gigl.GLMaterial,
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
	sobj := SceneObject{Geometry: geometry, Material: material, VShader: vshader, EShader: eshader, FShader: fshader}
	sobj.modelmatrix.SetIdentity()
	sobj.UseDepth = false // new drawings will overwrite old ones by default
	sobj.UseBlend = false // alpha blending is turned off by default
	sobj.children = nil   // OPTIONAL, only if current SceneObject has any child SceneObjects
	sobj.bbox = c2d.BBoxInit()
	return &sobj
}

func (self *SceneObject) ShowInfo() {
	fmt.Printf("SceneObject ")
	self.Geometry.ShowInfo()
	if self.instance_buffer != nil {
		fmt.Printf("  Instance Poses : nposes=%d stride=%d\n", self.instance_count, self.instance_stride)
	}
	if self.Material != nil {
		fmt.Printf("  ")
		self.Material.ShowInfo()
	}
	if self.VShader != nil {
		fmt.Printf("  VERT ")
		self.VShader.ShowInfo()
	}
	if self.EShader != nil {
		fmt.Printf("  EDGE ")
		self.EShader.ShowInfo()
	}
	if self.FShader != nil {
		fmt.Printf("  FACE ")
		self.FShader.ShowInfo()
	}
	fmt.Printf("  Flags    : UseDepth=%t  UseBlend=%t\n", self.UseDepth, self.UseBlend)
	fmt.Printf("  Children : %d\n", len(self.children))
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
		fmt.Printf("WARNING: SetInstancePoseValues() failed : invalid offset (%d) and value count (%d) for the stride (%d)\n", offset, len(values), self.instance_stride)
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
		fmt.Printf("WARNING: SetInstanceColorValues() failed : invalid offset (%d) for the stride (%d)\n", offset, self.instance_stride)
		return
	}
	pos := instance_index * self.instance_stride
	b0, b1, b2, b3 := uint32(v0), uint32(v1), uint32(v2), uint32(v3)
	self.instance_buffer[pos+offset] = math.Float32frombits(b0 + b1<<8 + b2<<16 + b3<<24) // LittleEndian (lower byte comes first)
}

// ----------------------------------------------------------------------------
// Translation, Rotation, Scaling (by manipulating MODEL matrix)
// ----------------------------------------------------------------------------

func (self *SceneObject) SetTransformation(txy [2]float32, angle_in_degree float32, sxy [2]float32) *SceneObject {
	translation := common.NewMatrix3().SetTranslation(txy[0], txy[1])
	rotation := common.NewMatrix3().SetRotation(angle_in_degree)
	scaling := common.NewMatrix3().SetScaling(sxy[0], sxy[1])
	self.modelmatrix.SetMultiplyMatrices(translation, rotation, scaling)
	return self
}

func (self *SceneObject) Rotate(angle_in_degree float32) *SceneObject {
	rotation := common.NewMatrix3().SetRotation(angle_in_degree)
	self.modelmatrix.SetMultiplyMatrices(rotation, &self.modelmatrix)
	return self
}

func (self *SceneObject) Translate(tx float32, ty float32) *SceneObject {
	translation := common.NewMatrix3().SetTranslation(tx, ty)
	self.modelmatrix.SetMultiplyMatrices(translation, &self.modelmatrix)
	return self
}

func (self *SceneObject) Scale(sx float32, sy float32) *SceneObject {
	scaling := common.NewMatrix3().SetScaling(sx, sy)
	self.modelmatrix.SetMultiplyMatrices(scaling, &self.modelmatrix)
	return self
}

// ----------------------------------------------------------------------------
// Bounding Box
// ----------------------------------------------------------------------------

func (self *SceneObject) GetBoundingBox(m *common.Matrix3, renew bool) [2][2]float32 {
	if !c2d.BBoxIsSet(self.bbox) || renew {
		bbox := c2d.BBoxInit()
		// apply the transformation matrx
		var mm *common.Matrix3 = nil
		if m != nil {
			mm = m.MultiplyToTheRight(&self.modelmatrix)
		} else {
			mm = self.modelmatrix.Copy() // new matrix
		}
		// add all the vertices of the geometry
		for _, v := range self.Geometry.verts {
			xy := mm.MultiplyVector2(v)
			c2d.BBoxAddPoint(&bbox, xy)
		}
		if self.instance_count > 0 {
			bbox_posed := c2d.BBoxInit()
			for i := 0; i < self.instance_count; i++ {
				idx := i * self.instance_stride
				txy := self.instance_buffer[idx : idx+2]
				for k := 0; k < 2; k++ {
					bboxp := [2]float32{bbox[k][0] + txy[0], bbox[k][1] + txy[1]}
					bboxp = mm.MultiplyVector2(bboxp)
					c2d.BBoxAddPoint(&bbox_posed, bboxp)
				}
			}
			bbox = bbox_posed
		}
		for _, sobj := range self.children {
			bbox = c2d.BBoxMerge(bbox, sobj.GetBoundingBox(mm, renew))
		}
		self.bbox = bbox
	}
	return self.bbox
}
