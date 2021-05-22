package gigl

import "fmt"

// ----------------------------------------------------------------------------
// VAO for Rendering 2D/3D SceneObject
// ----------------------------------------------------------------------------

type VAO struct {
	VertBuffer     interface{} // WebGL/OpenGL buffer for geometry's vertex points
	FvtxBuffer     interface{} // WebGL/OpenGL buffer for geometry's face vertex points (points for PER_FACE vertices)
	VertBufferInfo [5]int      // [nverts, stride, coord_size, texture_uv_size, vertex_normal_size]
	FvtxBufferInfo [5]int      // [nverts, stride, coord_size, texture_uv_size, vertex_normal_size]

	EdgeBuffer      interface{} // WebGL/OpenGL buffer for geometry's edge indices
	FaceBuffer      interface{} // WebGL/OpenGL buffer for geometry's face indices
	EdgeBufferCount int         // EdgeBuffer count (number of uint32 vertex indices)
	FaceBufferCount int         // EdgeBuffer count (number of uint32 vertex indices)

	InstanceBuffer     interface{} // WebGL/OpenGL buffer for geometry's instance values
	InstanceBufferInfo [2]int      // [instance_count, instance_stride] of instance poses
}

func (self *VAO) ShowInfo() {
	ox := func(b interface{}) string {
		if b == nil {
			return "X"
		} else {
			return "O"
		}
	}
	fmt.Printf("VAO\n")
	fmt.Printf("  VertBuffer : %s  [ nverts:%d stride:%d coord:%d tuv:%d norm:%d ]\n", ox(self.VertBuffer), self.FvtxBufferInfo[0], self.VertBufferInfo[1], self.VertBufferInfo[2], self.VertBufferInfo[3], self.VertBufferInfo[4])
	fmt.Printf("  FvtxBuffer : %s  [ nverts:%d stride:%d coord:%d tuv:%d norm:%d ]\n", ox(self.FvtxBuffer), self.FvtxBufferInfo[0], self.FvtxBufferInfo[1], self.FvtxBufferInfo[2], self.FvtxBufferInfo[3], self.FvtxBufferInfo[4])
	fmt.Printf("  EdgeBuffer : %s  [ count:%d ]\n", ox(self.EdgeBuffer), self.EdgeBufferCount)
	fmt.Printf("  FaceBuffer : %s  [ count:%d ]\n", ox(self.FaceBuffer), self.FaceBufferCount)
	fmt.Printf("  InstanceBuffer : %s  [ count:%d stride:%d ]\n", ox(self.InstanceBuffer), self.InstanceBufferInfo[0], self.InstanceBufferInfo[1])
}

func (self *VAO) GetVtxBuffer(draw_mode int, xun int) (interface{}, [4]int) {
	if draw_mode == 3 && self.FvtxBuffer != nil {
		pinfo := self.FvtxBufferInfo // use extra vertex buffer (built for FACE drawing)
		switch xun {
		case 0: // vertex position
			nverts, stride, size, offset := pinfo[0], pinfo[1], pinfo[2], 0
			return self.FvtxBuffer, [4]int{nverts, stride, size, offset}
		case 1: // texture uv
			nverts, stride, size, offset := pinfo[0], pinfo[1], pinfo[3], pinfo[2]
			return self.FvtxBuffer, [4]int{nverts, stride, size, offset}
		case 2: // vertex normal
			nverts, stride, size, offset := pinfo[0], pinfo[1], pinfo[4], pinfo[2]+pinfo[3]
			return self.FvtxBuffer, [4]int{nverts, stride, size, offset}
		default:
			fmt.Printf("ERROR: invalid 'xun' (%d) in VAO.GetVtxBuffer()\n", xun)
			return nil, [4]int{0, 0, 0, 0} // nverts, stride, size, offset
		}
	} else {
		pinfo := self.VertBufferInfo // use original vertex buffer
		switch xun {
		case 0: // vertex position
			nverts, stride, size, offset := pinfo[0], pinfo[1], pinfo[2], 0
			return self.VertBuffer, [4]int{nverts, stride, size, offset}
		case 1: // texture uv
			nverts, stride, size, offset := pinfo[0], pinfo[1], pinfo[3], pinfo[2]
			return self.VertBuffer, [4]int{nverts, stride, size, offset}
		case 2: // vertex normal
			nverts, stride, size, offset := pinfo[0], pinfo[1], pinfo[4], pinfo[2]+pinfo[3]
			return self.VertBuffer, [4]int{nverts, stride, size, offset}
		default:
			fmt.Printf("ERROR: invalid 'xun' (%d) in VAO.GetVtxBuffer()\n", xun)
			return nil, [4]int{0, 0, 0, 0} // nverts, stride, size, offset
		}
	}
}

func (self *VAO) GetIdxBuffer(draw_mode int) (interface{}, int) {
	switch draw_mode {
	case 2:
		return self.EdgeBuffer, self.EdgeBufferCount
	case 3:
		return self.FaceBuffer, self.FaceBufferCount
	default:
		fmt.Printf("ERROR: invalid 'draw_mode' (%d) in VAO.GetIdxBuffer()\n", draw_mode)
		return nil, 0
	}
}

func (self *VAO) GetInstanceBuffer() (interface{}, [2]int) {
	count, stride := self.InstanceBufferInfo[0], self.InstanceBufferInfo[1]
	return self.VertBuffer, [2]int{count, stride}
}
