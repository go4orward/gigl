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

	IsLoading() bool
	IsLoaded() bool
	IsReady() bool
}
