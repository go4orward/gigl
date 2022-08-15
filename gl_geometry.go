package gigl

type GLGeometry interface {
	// This interface defines a set of functions that both 2D and 3D Geometry have to provide.
	// NewSceneObject() function requires this GLGeometry interface instead of a 2D or 3D Geometry,
	// since 3D SceneObject should be able to use both 2D and 3D Geometry (as in g3d.OverlayMarkerLayer).
	IsDataBufferReady() bool
	IsVtxBufferRebuiltForFaces() bool
	GetVtxBuffer(draw_mode int) []float32  // data buffer of vertices (mode 0:original_verts, 1:face_verts_only)
	GetIdxBuffer(draw_mode int) []uint32   // data buffer of indices  (mode 2:for_edges, 3:for_faces)
	GetVtxBufferInfo(draw_mode int) [5]int // data buffer info : [nverts, stride, xyz_size, uv_size, normal_size]
	GetIdxBufferCount(draw_mode int) int   // data buffer count : number of vertex indices
	Summary() string                       //
}
