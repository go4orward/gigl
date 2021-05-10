package g2d

import (
	"fmt"
	"math"

	"github.com/go4orward/gigl"
	"github.com/go4orward/gigl/common"
	"github.com/go4orward/gigl/g2d/c2d"
)

// ----------------------------------------------------------------------------
// Geometry
// ----------------------------------------------------------------------------

type Geometry struct {
	verts [][2]float32
	edges [][]uint32
	faces [][]uint32
	tuvs  [][]float32 // texture uv coordinates (PER_FACE [nfaces][6] or PER_VERT [nverts][2])

	data_buffer_vpoints []float32 // serialized data buffer for vertex points : COORD[]
	data_buffer_fpoints []float32 // serialized data buffer for PER_FACE vertex points : COORD[2] + (UV[2])
	data_buffer_lines   []uint32  // serialized data buffer for edge lines : vidx[2]
	data_buffer_faces   []uint32  // serialized data buffer for face triangles : vidx[3]

	// Note that, for PER_FACE texture UV-coordinates, vertices are duplicated for each face
	fpoint_vidx_list  []uint32 // index of vertex_list of each face after PER_FACE data duplication
	fpoint_vert_total int      // total count of vertices after PER_FACE data duplication
	fpoint_info       [4]int   // data size of a point (for triangles) : [ stride, xyz_offset, uv_offset, RESERVED ]
	vpoint_info       [4]int   // data size of a point (for points & lines)

	wbuffer_vpoints interface{} // WebGL data buffer for data_buffer_vpoints (points for vertices)
	wbuffer_fpoints interface{} // WebGL data buffer for data_buffer_fpoints (points for PER_FACE vertices)
	wbuffer_lines   interface{} // WebGL data buffer for data_buffer_lines (indices for lines)
	wbuffer_faces   interface{} // WebGL data buffer for data_buffer_faces (indices for triangles)
}

func NewGeometry() *Geometry {
	var geometry Geometry
	geometry.Clear(true, true, true)
	return &geometry
}

func (self *Geometry) Clear(geom bool, data_buf bool, webgl_buf bool) *Geometry {
	if geom {
		self.verts = [][2]float32{}
		self.edges = [][]uint32{}
		self.faces = [][]uint32{}
		self.tuvs = [][]float32{}
	}
	if data_buf || geom {
		self.data_buffer_vpoints = nil
		self.data_buffer_fpoints = nil
		self.data_buffer_lines = nil
		self.data_buffer_faces = nil
		self.fpoint_vidx_list = nil
		self.fpoint_vert_total = 0
		self.fpoint_info = [4]int{0, 0, 0, 0}
		self.vpoint_info = [4]int{0, 0, 0, 0}
	}
	if webgl_buf || data_buf || geom {
		self.wbuffer_vpoints = nil
		self.wbuffer_fpoints = nil
		self.wbuffer_lines = nil
		self.wbuffer_faces = nil
	}
	return self
}

func (self *Geometry) String() string {
	return fmt.Sprintf("2DGeometry{v:%d e:%d f:%d}\n", len(self.verts), len(self.faces), len(self.faces))
}

func (self *Geometry) ShowInfo() {
	fmt.Printf("2DGeometry with %d verts %d edges %d faces\n", len(self.verts), len(self.edges), len(self.faces))
	if len(self.tuvs) > 0 {
		if self.HasTextureFor("VERTEX") {
			fmt.Printf("    texture UV coords   : [%d][]float32   for each vertex\n", len(self.tuvs))
		} else if self.HasTextureFor("FACE") {
			fmt.Printf("    texture UV coords   : [%d][]float32   for each face\n", len(self.tuvs))
		} else {
			fmt.Printf("    texture UV coords   : [%d][]float32   \n", len(self.tuvs))
		}
	}
	wb := func(wbuffer interface{}) string {
		if wbuffer == nil {
			return "X"
		} else {
			return "O"
		}
	}
	fmt.Printf("    data_buffer_vpoints : %4d (WebGL:%s)  pinfo=%v\n", len(self.data_buffer_vpoints), wb(self.wbuffer_vpoints), self.vpoint_info)
	fmt.Printf("    data_buffer_fpoints : %4d (WebGL:%s)  pinfo=%v\n", len(self.data_buffer_fpoints), wb(self.wbuffer_fpoints), self.fpoint_info)
	fmt.Printf("    data_buffer_lines   : %4d (WebGL:%s)\n", len(self.data_buffer_lines), wb(self.wbuffer_lines))
	fmt.Printf("    data_buffer_faces   : %4d (WebGL:%s)\n", len(self.data_buffer_faces), wb(self.wbuffer_faces))
	// fmt.Printf("    data_buffer_vpoints : %v\n", self.data_buffer_vpoints)
	// fmt.Printf("    data_buffer_fpoints : %v\n", self.data_buffer_fpoints)
}

// ----------------------------------------------------------------------------
//
// ----------------------------------------------------------------------------

func (self *Geometry) SetVertices(vertices [][2]float32) *Geometry {
	self.verts = vertices
	return self
}

func (self *Geometry) SetEdges(edges [][]uint32) *Geometry {
	self.edges = edges
	return self
}

func (self *Geometry) SetFaces(faces [][]uint32) *Geometry {
	self.faces = faces
	return self
}

func (self *Geometry) AddVertex(coords [2]float32) uint32 {
	vidx := len(self.verts)
	self.verts = append(self.verts, coords)
	return uint32(vidx)
}

func (self *Geometry) AddEdge(edge []uint32) uint32 {
	eidx := len(self.edges)
	self.edges = append(self.edges, edge)
	return uint32(eidx)
}

func (self *Geometry) AddFace(face []uint32) uint32 {
	fidx := len(self.faces)
	self.faces = append(self.faces, face)
	return uint32(fidx)
}

// ----------------------------------------------------------------------------
// Transformation of Vertex Coordinates
// ----------------------------------------------------------------------------

func (self *Geometry) Translate(tx float32, ty float32) *Geometry {
	for i := 0; i < len(self.verts); i++ {
		self.verts[i][0] += tx
		self.verts[i][1] += ty
	}
	self.Clear(false, true, true)
	return self
}

func (self *Geometry) Rotate(angle_in_degree float32) *Geometry {
	rad := float64(angle_in_degree * (math.Pi / 180))
	sin, cos := math.Sin(rad), math.Cos(rad)
	for i := 0; i < len(self.verts); i++ {
		x, y := float64(self.verts[i][0]), float64(self.verts[i][1])
		self.verts[i][0] = float32(cos*x - sin*y)
		self.verts[i][1] = float32(sin*x + cos*y)
	}
	self.Clear(false, true, true)
	return self
}

func (self *Geometry) Scale(sx float32, sy float32) *Geometry {
	for i := 0; i < len(self.verts); i++ {
		self.verts[i][0] *= sx
		self.verts[i][1] *= sy
	}
	self.Clear(false, true, true)
	return self
}

func (self *Geometry) AppyMatrix(matrix *common.Matrix3) *Geometry {
	for i := 0; i < len(self.verts); i++ {
		self.verts[i] = matrix.MultiplyVector2(self.verts[i])
	}
	self.Clear(false, true, true)
	return self
}

// ----------------------------------------------------------------------------
// Merge
// ----------------------------------------------------------------------------

func (self *Geometry) Merge(g *Geometry) *Geometry {
	self.Clear(false, true, true)
	vcount := uint32(len(self.verts))
	for _, v := range g.verts {
		self.AddVertex(v)
	}
	for _, e := range g.edges {
		new_edge := make([]uint32, len(e))
		for i := 0; i < len(new_edge); i++ {
			new_edge[i] = e[i] + vcount
		}
		self.AddEdge(new_edge)
	}
	for _, f := range g.faces {
		new_face := make([]uint32, len(f))
		for i := 0; i < len(new_face); i++ {
			new_face[i] = f[i] + vcount
		}
		self.AddFace(new_face)
	}
	// TODO: texture UV coordinates
	return self
}

// ----------------------------------------------------------------------------
// Texture UV coordinates
// ----------------------------------------------------------------------------

func (self *Geometry) HasTextureFor(mode string) bool {
	switch mode {
	case "VERTEX":
		return len(self.tuvs) == len(self.verts) && len(self.tuvs[0]) == 2
	case "FACE":
		return len(self.tuvs) == len(self.faces) && len(self.tuvs[0]) >= 6
	default:
		return self.HasTextureFor("VERTEX") || self.HasTextureFor("FACE")
	}
}

func (self *Geometry) AddTextureUV(tuv []float32) *Geometry {
	if len(tuv) == 0 || len(tuv)%2 == 0 {
		fmt.Printf("Invalid texture coordinates to add : %v\n", tuv)
		return self
	}
	self.tuvs = append(self.tuvs, tuv)
	return self
}

func (self *Geometry) SetTextureUVs(tuvs [][]float32) *Geometry {
	if len(tuvs) == len(self.faces) && len(tuvs[0]) >= 6 { // texture for each face
	} else if len(tuvs) == len(self.verts) && len(tuvs[0]) == 2 { // texture for each vertex
	} else {
		fmt.Printf("Invalid texture UVs : %v\n", tuvs)
		return self
	}
	self.tuvs = tuvs
	return self
}

// ----------------------------------------------------------------------------
// Triangulation
// ----------------------------------------------------------------------------

func (self *Geometry) get_reverse(face_vlist []uint32) []uint32 {
	new_vlist := make([]uint32, len(face_vlist))
	for i := len(face_vlist) - 1; i >= 0; i-- {
		new_vlist[i] = face_vlist[i]
	}
	return new_vlist
}

func (self *Geometry) splice_indices(a []uint32, pos int, delete_count int, new_entries ...uint32) []uint32 {
	head := a[0:pos]
	tail := a[pos+delete_count:]
	return append(append(head, new_entries...), tail...)
}

func (self *Geometry) get_triangulation(face_vlist []uint32) [][]uint32 {
	vindices := make([]uint32, len(face_vlist))
	copy(vindices, face_vlist)
	new_faces := make([][]uint32, 0)
	vidx, vcount := 0, len(vindices)
	iterations, max_iterations := 0, 10*len(vindices)
	for vcount > 3 && iterations < max_iterations {
		i0, i1, i2 := vidx, (vidx+1)%vcount, (vidx+2)%vcount
		v0, v1, v2 := self.verts[vindices[i0]], self.verts[vindices[i1]], self.verts[vindices[i2]]
		if c2d.IsCCW(v0, v1, v2) {
			point_inside := false
			for j := 0; j < vcount; j++ {
				if j != i0 && j != i1 && j != i2 && c2d.IsPointInside(self.verts[vindices[j]], v0, v1, v2) {
					point_inside = true
					break
				}
			}
			if !point_inside {
				new_faces = append(new_faces, []uint32{vindices[i0], vindices[i1], vindices[i2]})
				vindices = self.splice_indices(vindices, i1, 1)
			}
		}
		vcount = len(vindices)
		vidx = (vidx + 1) % vcount
		iterations++
	}
	new_faces = append(new_faces, vindices)
	if iterations == max_iterations {
		fmt.Printf("failed to traiangulate : %v => %v\n", face_vlist, new_faces)
	}
	return new_faces
}

// ----------------------------------------------------------------------------
// Build Data Buffers (serialized)
// ----------------------------------------------------------------------------

func (self *Geometry) IsDataBufferReady() bool {
	return len(self.data_buffer_vpoints) > 0 || len(self.data_buffer_fpoints) > 0
}

func (self *Geometry) count_fpoint_vidx_list() int {
	self.fpoint_vert_total = 0
	self.fpoint_vidx_list = make([]uint32, len(self.faces))
	for i := 0; i < len(self.faces); i++ {
		self.fpoint_vidx_list[i] = uint32(self.fpoint_vert_total)
		self.fpoint_vert_total += len(self.faces[i])
	}
	return self.fpoint_vert_total
}

func (self *Geometry) get_fpoint_new_vidx(fidx int, i int) int {
	if self.fpoint_vidx_list == nil {
		self.count_fpoint_vidx_list()
	}
	return int(self.fpoint_vidx_list[fidx]) + i
}

func (self *Geometry) buffer_copy_xy(buf []float32, pinfo [4]int, new_vidx int, vidx int) {
	stride, offset := pinfo[0], pinfo[1] // XY coordinates
	pos := new_vidx*stride + offset
	buf[pos+0] = self.verts[vidx][0]
	buf[pos+1] = self.verts[vidx][1]
}

func (self *Geometry) buffer_copy_uv(buf []float32, pinfo [4]int, new_vidx int, tuv_idx int, tuv_offset int) {
	stride, offset := pinfo[0], pinfo[2] // UV texture coordinates
	u := uint32(self.tuvs[tuv_idx][tuv_offset+0] * 65535)
	v := uint32(self.tuvs[tuv_idx][tuv_offset+1] * 65535)
	pos := new_vidx*stride + offset
	buf[pos] = math.Float32frombits(u + v<<16) // LittleEndian (lower byte comes first)
}

func (self *Geometry) BuildDataBuffers(for_points bool, for_lines bool, for_faces bool) {
	// create data buffer for vertex points
	self.data_buffer_vpoints, self.vpoint_info = nil, [4]int{0, 0, 0, 0}
	self.data_buffer_fpoints, self.fpoint_info = nil, [4]int{0, 0, 0, 0}
	points_per_face := false
	if for_faces {
		if self.HasTextureFor("FACE") {
			points_per_face = true
			self.count_fpoint_vidx_list()
			self.fpoint_info = [4]int{(2 + 1), 0, 2, 0} // stride, xyz_offset, uv_offset, RESERVED
			self.data_buffer_fpoints = make([]float32, self.fpoint_vert_total*self.fpoint_info[0])
			for fidx, face_vlist := range self.faces {
				for i := 0; i < len(face_vlist); i++ {
					new_vidx, vidx := self.get_fpoint_new_vidx(fidx, i), int(face_vlist[i])
					self.buffer_copy_xy(self.data_buffer_fpoints, self.fpoint_info, new_vidx, vidx)
					self.buffer_copy_uv(self.data_buffer_fpoints, self.fpoint_info, new_vidx, fidx, i*2)
				}
			}
		} else if self.HasTextureFor("VERTEX") {
			self.fpoint_info = [4]int{(2 + 1), 0, 2, 0} // stride, xyz_offset, uv_offset, RESERVED
			self.data_buffer_fpoints = make([]float32, len(self.verts)*self.fpoint_info[0])
			for vidx := 0; vidx < len(self.verts); vidx++ {
				self.buffer_copy_xy(self.data_buffer_fpoints, self.fpoint_info, vidx, vidx)
				self.buffer_copy_uv(self.data_buffer_fpoints, self.fpoint_info, vidx, vidx, 0)
			}
			self.data_buffer_vpoints = self.data_buffer_fpoints
			self.vpoint_info = self.fpoint_info
		} else {
			self.fpoint_info = [4]int{(2 + 0), 0, 0, 0} // stride, xyz_offset, uv_offset, RESERVED
			self.data_buffer_fpoints = make([]float32, len(self.verts)*self.fpoint_info[0])
			for vidx := 0; vidx < len(self.verts); vidx++ {
				self.buffer_copy_xy(self.data_buffer_fpoints, self.fpoint_info, vidx, vidx)
			}
			self.data_buffer_vpoints = self.data_buffer_fpoints
			self.vpoint_info = self.fpoint_info
		}
	} else {
		self.data_buffer_vpoints = nil
	}
	if (for_points || for_lines) && self.data_buffer_vpoints == nil {
		self.vpoint_info = [4]int{2, 0, 0, 0}
		self.data_buffer_vpoints = make([]float32, len(self.verts)*self.vpoint_info[0])
		for vidx := 0; vidx < len(self.verts); vidx++ {
			self.buffer_copy_xy(self.data_buffer_vpoints, self.vpoint_info, vidx, vidx)
		}
	}
	// if self.fpoint_info[0] == 0 {
	// 	self.fpoint_info = [3]int{self.vpoint_info, 0, 0}
	// }
	// create data buffer for line drawings
	if for_lines {
		segment_count := 0
		for _, edge := range self.edges {
			segment_count += len(edge) - 1
		}
		self.data_buffer_lines = make([]uint32, segment_count*2)
		epos := 0
		for _, edge := range self.edges {
			for i := 1; i < len(edge); i++ {
				self.data_buffer_lines[epos+0] = edge[i-1]
				self.data_buffer_lines[epos+1] = edge[i]
				epos += 2
			}
		}
	} else {
		self.data_buffer_lines = nil
	}
	// create data buffer for surface drawings
	if for_faces {
		triangle_count := 0
		for _, face := range self.faces {
			triangle_count += len(face) - 2
		}
		find_index_in_face := func(vidx uint32, face []uint32) int {
			for i := 0; i < len(face); i++ {
				if vidx == face[i] {
					return i
				}
			}
			return 0
		}
		self.data_buffer_faces = make([]uint32, triangle_count*3)
		tpos := 0
		for fidx, face := range self.faces { // []vidx
			triangles := self.get_triangulation(face) // [][3]vidx
			for _, triangle := range triangles {      // [3]vidx
				if points_per_face { // vertex index has been changed due to PER_FACE duplication
					vidx_stt := self.get_fpoint_new_vidx(fidx, 0)
					self.data_buffer_faces[tpos+0] = uint32(vidx_stt + find_index_in_face(triangle[0], face))
					self.data_buffer_faces[tpos+1] = uint32(vidx_stt + find_index_in_face(triangle[1], face))
					self.data_buffer_faces[tpos+2] = uint32(vidx_stt + find_index_in_face(triangle[2], face))
				} else { // vertex index was preserved
					self.data_buffer_faces[tpos+0] = triangle[0]
					self.data_buffer_faces[tpos+1] = triangle[1]
					self.data_buffer_faces[tpos+2] = triangle[2]
				}
				tpos += 3
			}
		}
	} else {
		self.data_buffer_faces = nil
	}
	self.Clear(false, false, true)
}

func (self *Geometry) BuildDataBuffersForWireframe() {
	if self.data_buffer_vpoints == nil {
		// create data buffer for vertex points, only if necessary
		self.data_buffer_vpoints = make([]float32, len(self.verts)*2)
		vpos := 0
		for _, xy := range self.verts {
			self.data_buffer_vpoints[vpos+0] = xy[0]
			self.data_buffer_vpoints[vpos+1] = xy[1]
			vpos += 2
		}
		self.vpoint_info = [4]int{2, 0, 0, 0}
	}
	// create data buffer for edges, by extracting wireframe from faces
	self.data_buffer_lines = make([]uint32, 0)
	for _, face := range self.faces {
		triangles := self.get_triangulation(face)
		for _, t := range triangles {
			self.data_buffer_lines = append(self.data_buffer_lines, t[0], t[1], t[1], t[2], t[2], t[0])
		}
	}
	self.Clear(false, false, true)
}

// ----------------------------------------------------------------------------
// Build WebGL Buffers
// ----------------------------------------------------------------------------

func (self *Geometry) IsWebGLBufferReady() bool {
	return self.wbuffer_vpoints != nil
}

func (self *Geometry) BuildWebGLBuffers(rc gigl.GLRenderingContext, for_points bool, for_lines bool, for_faces bool) {
	// THIS FUCNTION IS MEANT TO BE CALLED BY RENDERER. NO NEED TO BE EXPORTED
	c := rc.GetConstants() // *gigl.Constants
	if for_points {
		self.wbuffer_vpoints = rc.CreateDataBuffer(c.ARRAY_BUFFER, self.data_buffer_vpoints)
	} else {
		self.wbuffer_vpoints = nil
	}
	if for_lines {
		self.wbuffer_lines = rc.CreateDataBuffer(c.ELEMENT_ARRAY_BUFFER, self.data_buffer_lines)
	} else {
		self.wbuffer_lines = nil
	}
	if for_faces && self.data_buffer_faces != nil {
		if self.data_buffer_fpoints != nil {
			self.wbuffer_fpoints = rc.CreateDataBuffer(c.ARRAY_BUFFER, self.data_buffer_fpoints)
		}
		self.wbuffer_faces = rc.CreateDataBuffer(c.ELEMENT_ARRAY_BUFFER, self.data_buffer_faces)
	} else {
		self.wbuffer_faces = nil
	}
}

func (self *Geometry) GetWebGLBuffer(draw_mode int) (interface{}, int, [4]int) {
	switch draw_mode {
	case 1: // "POINTS", "VERTICES"
		if self.data_buffer_fpoints == nil {
			return self.wbuffer_vpoints, len(self.data_buffer_vpoints), self.vpoint_info
		} else if self.data_buffer_vpoints == nil {
			return self.wbuffer_fpoints, len(self.data_buffer_fpoints), self.fpoint_info
		} else {
			return self.wbuffer_fpoints, len(self.data_buffer_fpoints), self.fpoint_info
		}
	case 2: // "LINES", "EDGES"
		return self.wbuffer_lines, len(self.data_buffer_lines), [4]int{0, 0, 0, 0}
	case 3: // "TRIANGLES", "FACES"
		return self.wbuffer_faces, len(self.data_buffer_faces), [4]int{0, 0, 0, 0}
	default:
		fmt.Printf("Invalid mode (%d) for GetWebGLBuffer()\n", draw_mode)
		return nil, 0, [4]int{0, 0, 0, 0}
	}
}
