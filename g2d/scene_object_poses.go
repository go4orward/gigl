package g2d

import (
	"fmt"

	"github.com/go4orward/gigl"
)

type SceneObjectPoses struct {
	Size       int         // number of values for a single pose
	Count      int         // number of poses
	DataBuffer []float32   //
	rcbuffer   interface{} // DataBuffer for the RenderingContext
}

func NewSceneObjectPoses(size int, count int, data []float32) *SceneObjectPoses {
	poses := SceneObjectPoses{Size: size, Count: count}
	poses.DataBuffer = make([]float32, size*count)
	if data != nil {
		for i := 0; i < len(poses.DataBuffer) && i < len(data); i++ {
			poses.DataBuffer[i] = data[i]
		}
	}
	poses.rcbuffer = nil
	return &poses
}

func (self *SceneObjectPoses) ShowInfo() {
	fmt.Printf("Instance Poses : size = %d & count = %d\n", self.Size, self.Count)
}

// ------------------------------------------------------------------------
// Setting Instance Pose
// ------------------------------------------------------------------------

func (self *SceneObjectPoses) SetPose(index int, offset int, values ...float32) bool {
	if (offset + len(values)) > self.Size {
		return false
	}
	pos := index * self.Size
	for i := 0; i < len(values); i++ {
		self.DataBuffer[pos+offset+i] = values[i]
	}
	return true
}

// ----------------------------------------------------------------------------
// Build Buffers of RenderingContext
// ----------------------------------------------------------------------------

func (self *SceneObjectPoses) IsRcBufferReady() bool {
	return self.rcbuffer != nil
}

func (self *SceneObjectPoses) BuildRcBuffer(rc gigl.GLRenderingContext) {
	// THIS FUNCTION IS MEANT TO BE CALLED BY RENDERER.
	c := rc.GetConstants()
	if self.DataBuffer != nil {
		self.rcbuffer = rc.CreateDataBuffer(c.ARRAY_BUFFER, self.DataBuffer)
	} else {
		self.rcbuffer = nil
	}
}

func (self *SceneObjectPoses) GetRcBuffer() (interface{}, int, [4]int) {
	return self.rcbuffer, len(self.DataBuffer), [4]int{0, 0, 0, 0}
}
