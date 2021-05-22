package g3d

import (
	"fmt"
	"math"

	"github.com/go4orward/gigl/common"
	"github.com/go4orward/gigl/g3d/c3d"
)

// ----------------------------------------------------------------------------
// Geometry
// ----------------------------------------------------------------------------

type Geometry struct {
	verts [][3]float32 // vertices
	edges [][]uint32   // edges
	faces [][]uint32   // faces
	tuvs  [][]float32  // texture uv coordinates (PER_FACE [nfaces][6] or PER_VERT [nverts][2])
	norms [][3]float32 // normal vectors (PER_FACE [nfaces][3] or PER_VERT [nverts][3])

	dbuffer_vpoint      []float32 // data buffer for vertex points : COORD[] + (UV[2]) + (NORMAL[3])
	dbuffer_fpoint      []float32 // data buffer for PER_FACE vertex points : COORD[3] + (UV[2]) + (NORMAL[3])
	dbuffer_line        []uint32  // data buffer for edge lines     : list of vertex indices
	dbuffer_face        []uint32  // data buffer for face triangles : list of vertex indices
	dbuffer_vpoint_info [4]int    // data buffer info : [ stride, xyz_size, uv_size, normal_size ]
	dbuffer_fpoint_info [4]int    // data buffer info : [ stride, xyz_size, uv_size, normal_size ]

	// Note that, for PER_FACE texture UV-coordinates and normal vectors, vertices are duplicated for each face
	fpoint_vidx_list  []uint32 // index of vertex_list of each face after PER_FACE data duplication
	fpoint_vert_total int      // total count of vertices after PER_FACE data duplication
}

func NewGeometry() *Geometry {
	var geometry Geometry
	geometry.Clear(true, true, true)
	return &geometry
}

func (self *Geometry) Clear(geom bool, data_buf bool, webgl_buf bool) *Geometry {
	if geom {
		self.verts = [][3]float32{}
		self.edges = [][]uint32{}
		self.faces = [][]uint32{}
		self.tuvs = [][]float32{}
		self.norms = [][3]float32{}
	}
	if geom || data_buf {
		self.dbuffer_vpoint = nil
		self.dbuffer_fpoint = nil
		self.dbuffer_line = nil
		self.dbuffer_face = nil
		self.dbuffer_fpoint_info = [4]int{0, 0, 0, 0}
		self.dbuffer_vpoint_info = [4]int{0, 0, 0, 0}
		self.fpoint_vidx_list = nil
		self.fpoint_vert_total = 0
	}
	return self
}

func (self *Geometry) String() string {
	return fmt.Sprintf("3DGeometry{v:%d e:%d f:%d}\n", len(self.verts), len(self.faces), len(self.faces))
}

func (self *Geometry) ShowInfo() {
	fmt.Printf("3DGeometry with %d verts %d edges %d faces\n", len(self.verts), len(self.edges), len(self.faces))
	if len(self.tuvs) > 0 {
		if self.HasTextureFor("VERTEX") {
			fmt.Printf("    texture coords : [%d][]float32   for each vertex\n", len(self.tuvs))
		} else if self.HasTextureFor("FACE") {
			fmt.Printf("    texture coords : [%d][]float32   for each face\n", len(self.tuvs))
		} else {
			fmt.Printf("    texture coords : [%d][]float32   incomplete\n", len(self.tuvs))
		}
	}
	if len(self.norms) > 0 {
		if self.HasNormalFor("VERTEX") {
			fmt.Printf("    normal vectors : [%d][3]float32   for each vertex\n", len(self.norms))
		} else if self.HasNormalFor("FACE") {
			fmt.Printf("    normal vectors : [%d][3]float32   for each face\n", len(self.norms))
		} else {
			fmt.Printf("    normal vectors : [%d][3]float32   incomplete\n", len(self.norms))
		}
	}
	fmt.Printf("    dbuffer_vpoint : %4d  pinfo=%v\n", len(self.dbuffer_vpoint)/self.dbuffer_vpoint_info[0], self.dbuffer_vpoint_info)
	fmt.Printf("    dbuffer_fpoint : %4d  pinfo=%v\n", len(self.dbuffer_fpoint)/self.dbuffer_fpoint_info[0], self.dbuffer_fpoint_info)
	fmt.Printf("    dbuffer_line   : %4d  \n", len(self.dbuffer_line))
	fmt.Printf("    dbuffer_face   : %4d  \n", len(self.dbuffer_face))
	// fmt.Printf("dbuffer_vpoint : %v\n", self.dbuffer_vpoint)
	// fmt.Printf("dbuffer_fpoint : %v\n", self.dbuffer_fpoint)
}

// ----------------------------------------------------------------------------
//
// ----------------------------------------------------------------------------

func (self *Geometry) SetVertices(vertices [][3]float32) *Geometry {
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

func (self *Geometry) AddVertex(coords [3]float32) uint32 {
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

func (self *Geometry) Translate(tx float32, ty float32, tz float32) *Geometry {
	for i := 0; i < len(self.verts); i++ {
		self.verts[i][0] += tx
		self.verts[i][1] += ty
		self.verts[i][2] += tz
	}
	self.Clear(false, true, true)
	return self
}

func (self *Geometry) Rotate(axis [3]float32, angle_in_degree float32) *Geometry {
	rM := common.NewMatrix4()
	rM.SetRotationByAxis(axis, angle_in_degree)
	for i := 0; i < len(self.verts); i++ {
		self.verts[i] = rM.MultiplyVector3(self.verts[i])
	}
	self.Clear(false, true, true)
	return self
}

func (self *Geometry) Scale(sx float32, sy float32, sz float32) *Geometry {
	for i := 0; i < len(self.verts); i++ {
		self.verts[i][0] *= sx
		self.verts[i][1] *= sy
		self.verts[i][2] *= sz
	}
	self.Clear(false, true, true)
	return self
}

func (self *Geometry) AppyMatrix4(m *common.Matrix4) *Geometry {
	for i := 0; i < len(self.verts); i++ {
		self.verts[i] = m.MultiplyVector3(self.verts[i])
	}
	self.Clear(false, true, true)
	return self
}

// ----------------------------------------------------------------------------
// Merge
// ----------------------------------------------------------------------------

func (self *Geometry) Merge(g *Geometry) *Geometry {
	self.Clear(false, true, true)
	nverts := uint32(len(self.verts))
	for _, v := range g.verts {
		self.AddVertex(v)
	}
	for _, e := range g.edges {
		new_edge := make([]uint32, len(e))
		for i := 0; i < len(e); i++ {
			new_edge[i] = nverts + e[i]
		}
		self.AddEdge(new_edge)
	}
	for _, f := range g.faces {
		new_face := make([]uint32, len(f))
		for i := 0; i < len(f); i++ {
			new_face[i] = nverts + f[i]
		}
		self.AddFace(new_face)
	}
	// TODO: texture UV coordinates and normal vectors
	return self
}

// ----------------------------------------------------------------------------
// Texture UV coordinates
// ----------------------------------------------------------------------------

func (self *Geometry) HasTextureFor(mode string) bool {
	switch mode {
	case "VERTEX":
		return len(self.tuvs) > 0 && len(self.tuvs[0]) == 2
	case "FACE":
		return len(self.tuvs) > 0 && len(self.tuvs[0]) >= 6
	default:
		return self.HasTextureFor("VERTEX") || self.HasTextureFor("FACE")
	}
}

func (self *Geometry) AddTextureUV(tuv []float32) *Geometry {
	if len(tuv) == 0 || len(tuv)%2 != 0 {
		fmt.Printf("Invalid texture coordinates to add : %v\n", tuv)
		return self
	}
	self.tuvs = append(self.tuvs, tuv)
	return self
}

func (self *Geometry) SetTextureUVs(tuvs [][]float32) *Geometry {
	self.tuvs = tuvs
	return self
}

// ----------------------------------------------------------------------------
// Normal Vectors
// ----------------------------------------------------------------------------

func (self *Geometry) HasNormalFor(mode string) bool {
	switch mode {
	case "VERTEX":
		return len(self.norms) > 0 && len(self.norms) == len(self.verts)
	case "FACE":
		return len(self.norms) > 0 && len(self.norms) == len(self.faces)
	default:
		return self.HasNormalFor("VERTEX") || self.HasNormalFor("FACE")
	}
}

func (self *Geometry) AddNormal(normal_vector [3]float32) *Geometry {
	self.norms = append(self.norms, normal_vector)
	return self
}

func (self *Geometry) SetNormals(normals [][3]float32) *Geometry {
	self.norms = normals
	return self
}

func (self *Geometry) BuildNormalsForVertex() {
	// self.norms == [nverts][3]float32
	self.norms = make([][3]float32, len(self.verts)) // self.norms == [nverts][3]float32
	for i := 0; i < len(self.verts); i++ {
		self.norms[i] = self.GetVertexNormal(i)
	}
}

func (self *Geometry) BuildNormalsForFace() {
	self.norms = make([][3]float32, len(self.faces)) // self.norms == [nfaces][3]float32
	for i := 0; i < len(self.faces); i++ {
		self.norms[i] = self.GetFaceNormal(i)
	}
}

func (self *Geometry) GetFaceNormal(fidx int) [3]float32 {
	findices := self.faces[fidx]
	i0, i1, i2 := findices[0], findices[1], findices[len(findices)-1]
	v0, v1, v2 := self.verts[i0], self.verts[i1], self.verts[i2]
	v01 := c3d.Normalize(c3d.SubAB(v1, v0))
	v02 := c3d.Normalize(c3d.SubAB(v2, v0))
	return c3d.Normalize(c3d.CrossAB(v01, v02))
}

func (self *Geometry) GetVertexNormal(vidx int) [3]float32 {
	if true {
		normals := [][3]float32{}
		for fidx, face := range self.faces {
			for _, face_vidx := range face {
				if face_vidx == uint32(vidx) {
					face_normal := self.GetFaceNormal(fidx)
					normals = append(normals, face_normal)
					break
				}
			}
		}
		return c3d.Normalize(c3d.AverageAll(normals))
	} else {
		neighbors := [][2]uint32{}
		for _, face_vlist := range self.faces {
			for i, face_vidx := range face_vlist {
				if face_vidx == uint32(vidx) {
					face_len := len(face_vlist)
					neighbor_next := uint32((i + 1) % face_len)
					neighbor_prev := uint32((i - 1 + face_len) % face_len)
					neighbors = append(neighbors, [2]uint32{neighbor_next, neighbor_prev})
					break
				}
			}
		}
		cross_sum := [3]float32{0, 0, 0}
		for _, nv := range neighbors {
			vnext := c3d.SubAB(self.verts[nv[0]], self.verts[vidx])
			vprev := c3d.SubAB(self.verts[nv[1]], self.verts[vidx])
			cross_sum = c3d.AddAB(cross_sum, c3d.CrossAB(vnext, vprev))
		}
		return c3d.Normalize(cross_sum)
	}
}

func (self *Geometry) ChangeNormal(idx int, normal_vector [3]float32) *Geometry {
	self.norms[idx] = normal_vector
	return self
}

// ----------------------------------------------------------------------------
// Trianulation
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

func (self *Geometry) get_triangulation(face_vlist []uint32, face_normal [3]float32) [][]uint32 {
	vindices := make([]uint32, len(face_vlist))
	copy(vindices, face_vlist)
	new_faces := make([][]uint32, 0)
	vidx, vcount := 0, len(vindices)
	iterations, max_iterations := 0, 10*len(vindices)
	for vcount > 3 && iterations < max_iterations {
		i0, i1, i2 := vidx, (vidx+1)%vcount, (vidx+2)%vcount
		v0, v1, v2 := self.verts[vindices[i0]], self.verts[vindices[i1]], self.verts[vindices[i2]]
		if c3d.DotAB(face_normal, c3d.CrossAB(c3d.SubAB(v1, v0), c3d.SubAB(v2, v0))) > 0 {
			point_inside := false
			for j := 0; j < vcount; j++ {
				if j != i0 && j != i1 && j != i2 && c3d.IsPointInside(self.verts[vindices[j]], v0, v1, v2) {
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
	// fmt.Printf("%v => %v\n", face_vlist, new_faces)
	return new_faces
}

// ----------------------------------------------------------------------------
// Build Data Buffers
// ----------------------------------------------------------------------------

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

func (self *Geometry) buffer_copy_xyz(buf []float32, pinfo [4]int, new_vidx int, vidx int) {
	stride, offset := pinfo[0], 0 // XY coordinates in 3 bytes
	pos := new_vidx*stride + offset
	buf[pos+0] = self.verts[vidx][0]
	buf[pos+1] = self.verts[vidx][1]
	buf[pos+2] = self.verts[vidx][2]
}

func (self *Geometry) buffer_copy_tuv(buf []float32, pinfo [4]int, new_vidx int, tuv_idx int, tuv_offset int) {
	stride, offset := pinfo[0], pinfo[1] // UV texture coordinates in 1 byte
	u := uint32(self.tuvs[tuv_idx][tuv_offset+0] * 65535)
	v := uint32(self.tuvs[tuv_idx][tuv_offset+1] * 65535)
	pos := new_vidx*stride + offset
	buf[pos] = math.Float32frombits(u + v<<16) // LittleEndian (lower byte comes first)
}

func (self *Geometry) buffer_copy_nor(buf []float32, pinfo [4]int, new_vidx int, nor_idx int) {
	stride, offset := pinfo[0], pinfo[1]+pinfo[2] // normal vector in 1 byte
	nx := uint32(self.norms[nor_idx][0] * 127)
	ny := uint32(self.norms[nor_idx][1] * 127)
	nz := uint32(self.norms[nor_idx][2] * 127)
	pos := new_vidx*stride + offset
	buf[pos] = math.Float32frombits(nx + ny<<8 + nz<<16) // LittleEndian (lower byte comes first)
}

func (self *Geometry) BuildDataBuffers(for_points bool, for_lines bool, for_faces bool) {
	// create data buffer for vertex points
	self.dbuffer_vpoint, self.dbuffer_vpoint_info = nil, [4]int{0, 0, 0, 0}
	self.dbuffer_fpoint, self.dbuffer_fpoint_info = nil, [4]int{0, 0, 0, 0}
	points_per_face := false
	if for_faces {
		points_per_face = self.HasNormalFor("FACE") || self.HasTextureFor("FACE")
		if points_per_face {
			self.count_fpoint_vidx_list()
		}
		if self.HasNormalFor("FACE") && self.HasTextureFor("FACE") {
			self.dbuffer_fpoint_info = [4]int{(3 + 1 + 1), 3, 1, 1} // size, xyz_size, uv_size, normal_size
			self.dbuffer_fpoint = make([]float32, self.fpoint_vert_total*self.dbuffer_fpoint_info[0])
			for fidx, face_vlist := range self.faces {
				for i := 0; i < len(face_vlist); i++ {
					new_vidx, vidx := self.get_fpoint_new_vidx(fidx, i), int(face_vlist[i])
					self.buffer_copy_xyz(self.dbuffer_fpoint, self.dbuffer_fpoint_info, new_vidx, vidx)
					self.buffer_copy_tuv(self.dbuffer_fpoint, self.dbuffer_fpoint_info, new_vidx, fidx, i*2)
					self.buffer_copy_nor(self.dbuffer_fpoint, self.dbuffer_fpoint_info, new_vidx, fidx)
				}
			}
		} else if self.HasNormalFor("FACE") && self.HasTextureFor("VERTEX") {
			self.dbuffer_fpoint_info = [4]int{(3 + 1 + 1), 3, 1, 1} // size, xyz_size, uv_size, normal_size
			self.dbuffer_fpoint = make([]float32, self.fpoint_vert_total*self.dbuffer_fpoint_info[0])
			for fidx, face_vlist := range self.faces {
				for i := 0; i < len(face_vlist); i++ {
					new_vidx, vidx := self.get_fpoint_new_vidx(fidx, i), int(face_vlist[i])
					self.buffer_copy_xyz(self.dbuffer_fpoint, self.dbuffer_fpoint_info, new_vidx, vidx)
					self.buffer_copy_tuv(self.dbuffer_fpoint, self.dbuffer_fpoint_info, new_vidx, vidx, 0)
					self.buffer_copy_nor(self.dbuffer_fpoint, self.dbuffer_fpoint_info, new_vidx, fidx)
				}
			}
		} else if self.HasNormalFor("FACE") && !self.HasTextureFor("") {
			self.dbuffer_fpoint_info = [4]int{(3 + 0 + 1), 3, 0, 1} // size, xyz_size, uv_size, normal_size
			self.dbuffer_fpoint = make([]float32, self.fpoint_vert_total*self.dbuffer_fpoint_info[0])
			for fidx, face_vlist := range self.faces {
				for i := 0; i < len(face_vlist); i++ {
					new_vidx, vidx := self.get_fpoint_new_vidx(fidx, i), int(face_vlist[i])
					self.buffer_copy_xyz(self.dbuffer_fpoint, self.dbuffer_fpoint_info, new_vidx, vidx)
					self.buffer_copy_nor(self.dbuffer_fpoint, self.dbuffer_fpoint_info, new_vidx, fidx)
				}
			}
		} else if self.HasNormalFor("VERTEX") && self.HasTextureFor("FACE") {
			self.dbuffer_fpoint_info = [4]int{(3 + 1 + 1), 3, 1, 1} // size, xyz_size, uv_size, normal_size
			self.dbuffer_fpoint = make([]float32, self.fpoint_vert_total*self.dbuffer_fpoint_info[0])
			for fidx, face_vlist := range self.faces {
				for i := 0; i < len(face_vlist); i++ {
					new_vidx, vidx := self.get_fpoint_new_vidx(fidx, i), int(face_vlist[i])
					self.buffer_copy_xyz(self.dbuffer_fpoint, self.dbuffer_fpoint_info, new_vidx, vidx)
					self.buffer_copy_tuv(self.dbuffer_fpoint, self.dbuffer_fpoint_info, new_vidx, fidx, i*2)
					self.buffer_copy_nor(self.dbuffer_fpoint, self.dbuffer_fpoint_info, new_vidx, vidx)
				}
			}
		} else if self.HasNormalFor("VERTEX") && self.HasTextureFor("VERTEX") {
			self.dbuffer_fpoint_info = [4]int{(3 + 1 + 1), 3, 1, 1} // size, xyz_size, uv_size, normal_size
			self.dbuffer_fpoint = make([]float32, len(self.verts)*self.dbuffer_fpoint_info[0])
			for vidx := 0; vidx < len(self.verts); vidx++ {
				self.buffer_copy_xyz(self.dbuffer_fpoint, self.dbuffer_fpoint_info, vidx, vidx)
				self.buffer_copy_tuv(self.dbuffer_fpoint, self.dbuffer_fpoint_info, vidx, vidx, 0)
				self.buffer_copy_nor(self.dbuffer_fpoint, self.dbuffer_fpoint_info, vidx, vidx)
			}
			self.dbuffer_vpoint = self.dbuffer_fpoint
			self.dbuffer_vpoint_info = self.dbuffer_fpoint_info
		} else if self.HasNormalFor("VERTEX") && !self.HasTextureFor("") {
			self.dbuffer_fpoint_info = [4]int{(3 + 0 + 1), 3, 0, 1} // size, xyz_size, uv_size, normal_size
			self.dbuffer_fpoint = make([]float32, len(self.verts)*self.dbuffer_fpoint_info[0])
			for vidx := 0; vidx < len(self.verts); vidx++ {
				self.buffer_copy_xyz(self.dbuffer_fpoint, self.dbuffer_fpoint_info, vidx, vidx)
				self.buffer_copy_nor(self.dbuffer_fpoint, self.dbuffer_fpoint_info, vidx, vidx)
			}
			self.dbuffer_vpoint = self.dbuffer_fpoint
			self.dbuffer_vpoint_info = self.dbuffer_fpoint_info
		} else if !self.HasNormalFor("") && self.HasTextureFor("FACE") {
			self.dbuffer_fpoint_info = [4]int{(3 + 1 + 0), 3, 1, 0} // size, xyz_size, uv_size, normal_size
			self.dbuffer_fpoint = make([]float32, self.fpoint_vert_total*self.dbuffer_fpoint_info[0])
			for fidx, face_vlist := range self.faces {
				for i := 0; i < len(face_vlist); i++ {
					new_vidx, vidx := self.get_fpoint_new_vidx(fidx, i), int(face_vlist[i])
					self.buffer_copy_xyz(self.dbuffer_fpoint, self.dbuffer_fpoint_info, new_vidx, vidx)
					self.buffer_copy_tuv(self.dbuffer_fpoint, self.dbuffer_fpoint_info, new_vidx, fidx, i*2)
				}
			}
		} else if !self.HasNormalFor("") && self.HasTextureFor("VERTEX") {
			self.dbuffer_fpoint_info = [4]int{(3 + 1 + 0), 3, 1, 0} // size, xyz_size, uv_size, normal_size
			self.dbuffer_fpoint = make([]float32, len(self.verts)*self.dbuffer_fpoint_info[0])
			for vidx := 0; vidx < len(self.verts); vidx++ {
				self.buffer_copy_xyz(self.dbuffer_fpoint, self.dbuffer_fpoint_info, vidx, vidx)
				self.buffer_copy_tuv(self.dbuffer_fpoint, self.dbuffer_fpoint_info, vidx, vidx, 0)
			}
			self.dbuffer_vpoint = self.dbuffer_fpoint
			self.dbuffer_vpoint_info = self.dbuffer_fpoint_info
		} else if !self.HasNormalFor("") && !self.HasTextureFor("") {
			self.dbuffer_fpoint_info = [4]int{(3 + 0 + 0), 3, 0, 0} // size, xyz_size, uv_size, normal_size
			self.dbuffer_fpoint = make([]float32, len(self.verts)*self.dbuffer_fpoint_info[0])
			for vidx := 0; vidx < len(self.verts); vidx++ {
				self.buffer_copy_xyz(self.dbuffer_fpoint, self.dbuffer_fpoint_info, vidx, vidx)
			}
			self.dbuffer_vpoint = self.dbuffer_fpoint
			self.dbuffer_vpoint_info = self.dbuffer_fpoint_info
		}
	} else {
		self.dbuffer_fpoint = nil
	}
	if (for_points || for_lines) && self.dbuffer_vpoint == nil {
		self.dbuffer_vpoint_info = [4]int{3, 3, 0, 0}
		self.dbuffer_vpoint = make([]float32, len(self.verts)*self.dbuffer_vpoint_info[0])
		for vidx := 0; vidx < len(self.verts); vidx++ {
			self.buffer_copy_xyz(self.dbuffer_vpoint, self.dbuffer_vpoint_info, vidx, vidx)
		}
	}
	// create data buffer for edge lines
	if for_lines {
		segment_count := 0
		for _, edge := range self.edges {
			segment_count += len(edge) - 1
		}
		self.dbuffer_line = make([]uint32, segment_count*2)
		epos := 0
		for _, edge := range self.edges {
			for i := 1; i < len(edge); i++ {
				self.dbuffer_line[epos+0] = edge[i-1]
				self.dbuffer_line[epos+1] = edge[i]
				epos += 2
			}
		}
	} else {
		self.dbuffer_line = nil
	}
	// create data buffer for faces
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
		self.dbuffer_face = make([]uint32, triangle_count*3)
		tpos := 0
		for fidx, face := range self.faces { // []vidx
			nv := self.GetFaceNormal(fidx)
			triangles := self.get_triangulation(face, nv) // [][3]vidx
			for _, triangle := range triangles {          // [3]vidx
				if points_per_face { // vertex index has been changed due to PER_FACE duplication
					vidx_stt := self.get_fpoint_new_vidx(fidx, 0)
					self.dbuffer_face[tpos+0] = uint32(vidx_stt + find_index_in_face(triangle[0], face))
					self.dbuffer_face[tpos+1] = uint32(vidx_stt + find_index_in_face(triangle[1], face))
					self.dbuffer_face[tpos+2] = uint32(vidx_stt + find_index_in_face(triangle[2], face))
				} else { // vertex index was preserved
					self.dbuffer_face[tpos+0] = triangle[0]
					self.dbuffer_face[tpos+1] = triangle[1]
					self.dbuffer_face[tpos+2] = triangle[2]
				}
				tpos += 3
			}
		}
	} else {
		self.dbuffer_face = nil
	}
	self.Clear(false, false, true)
}

func (self *Geometry) BuildDataBuffersForWireframe() {
	if self.dbuffer_vpoint == nil {
		// create data buffer for vertex points, only if necessary
		self.dbuffer_vpoint = make([]float32, len(self.verts)*3)
		vpos := 0
		for _, xyz := range self.verts {
			self.dbuffer_vpoint[vpos+0] = xyz[0]
			self.dbuffer_vpoint[vpos+1] = xyz[1]
			self.dbuffer_vpoint[vpos+2] = xyz[2]
			vpos += 3
		}
		self.dbuffer_vpoint_info = [4]int{3, 3, 0, 0}
	}
	// create data buffer for edges, by extracting wireframe from faces
	self.dbuffer_line = make([]uint32, 0)
	for fidx, face := range self.faces {
		normal := self.GetFaceNormal(fidx)
		triangles := self.get_triangulation(face, normal)
		for _, t := range triangles {
			self.dbuffer_line = append(self.dbuffer_line, t[0], t[1], t[1], t[2], t[2], t[0])
		}
	}
	self.Clear(false, false, true)
}

// ----------------------------------------------------------------------------
// Get Data Buffer (GLGeometry Interface)
// ----------------------------------------------------------------------------

func (self *Geometry) IsDataBufferReady() bool {
	return len(self.dbuffer_vpoint) > 0 || len(self.dbuffer_fpoint) > 0
}

func (self *Geometry) IsVtxBufferRebuiltForFaces() bool {
	return self.dbuffer_fpoint != nil
}

func (self *Geometry) GetVtxBuffer(draw_mode int) []float32 {
	if draw_mode == 3 && self.dbuffer_fpoint != nil {
		return self.dbuffer_fpoint // use extra vertex buffer (built for FACE drawing)
	} else {
		return self.dbuffer_vpoint // use original vertex buffer
	}
}

func (self *Geometry) GetVtxBufferInfo(draw_mode int) [5]int {
	if draw_mode == 3 && self.dbuffer_fpoint != nil {
		pinfo := self.dbuffer_fpoint_info // use extra vertex buffer (built for FACE drawing)s
		return [5]int{(len(self.dbuffer_fpoint) / pinfo[0]), pinfo[0], pinfo[1], pinfo[2], pinfo[3]}
	} else {
		pinfo := self.dbuffer_vpoint_info // use original vertex buffer
		return [5]int{(len(self.dbuffer_vpoint) / pinfo[0]), pinfo[0], pinfo[1], pinfo[2], pinfo[3]}
	}
}

func (self *Geometry) GetIdxBuffer(draw_mode int) []uint32 {
	switch draw_mode {
	case 2:
		return self.dbuffer_line
	case 3:
		return self.dbuffer_face
	default:
		fmt.Printf("ERROR: invalid 'draw_mode' (%d) in geom.GetIdxBuffer()\n", draw_mode)
		return nil
	}
}

func (self *Geometry) GetIdxBufferCount(draw_mode int) int {
	switch draw_mode {
	case 2:
		return len(self.dbuffer_line)
	case 3:
		return len(self.dbuffer_face)
	default:
		fmt.Printf("ERROR: invalid 'draw_mode' (%d) in geom.GetIdxBufferCount()\n", draw_mode)
		return 0 // n
	}
}
