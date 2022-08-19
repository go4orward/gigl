package gigl

type GLMaterial interface {
	MaterialSummary() string
}

type GLMaterialColors interface {
	MaterialSummary() string

	GetDrawModeColor(draw_mode int) [4]float32
}

type GLMaterialTexture interface {
	MaterialSummary() string

	GetTexturePixbuf() []uint8
	GetTextureWH() [2]int

	GetTexture() any
	SetTexture(texture any)
	GetTextureRGB() [3]float32
	SetTextureRGB(color any)

	IsLoading() bool // Texture is being loaded asynchronously by non-main thread (using Go function).
	IsLoaded() bool  // Texture was successfully loaded, and it needs to be set up by main thread.
	IsReady() bool   // Texture was successfully set up, and it's ready for rendering.
}
