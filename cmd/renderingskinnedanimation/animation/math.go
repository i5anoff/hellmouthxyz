package animation

import (
	"encoding/json"
	"github.com/go-gl/gl/v3.3-core/gl"
	"math"
)

type Matrix4f struct {
	M00 float32
	M01 float32
	M02 float32
	M03 float32

	M10 float32
	M11 float32
	M12 float32
	M13 float32

	M20 float32
	M21 float32
	M22 float32
	M23 float32

	M30 float32
	M31 float32
	M32 float32
	M33 float32
}

func (e *Matrix4f) UnmarshalJSON(b []byte) error {

	c := []float32{}

	if err := json.Unmarshal(b, &c); err == nil {
		x := &Matrix4f{c[0], c[1], c[2], c[3],
			c[4], c[5], c[6], c[7],
			c[8], c[9], c[10], c[11],
			c[12], c[13], c[14], c[15]}
		*e = *x
	} else {
		return err
	}

	return nil
}

func (m *Matrix4f) Mul(r *Matrix4f) *Matrix4f {
	res := new(Matrix4f)

	res.M00 = (m.M00 * r.M00) + (m.M01 * r.M10) + (m.M02 * r.M20) + (m.M03 * r.M30)
	res.M01 = (m.M00 * r.M01) + (m.M01 * r.M11) + (m.M02 * r.M21) + (m.M03 * r.M31)
	res.M02 = (m.M00 * r.M02) + (m.M01 * r.M12) + (m.M02 * r.M22) + (m.M03 * r.M32)
	res.M03 = (m.M00 * r.M03) + (m.M01 * r.M13) + (m.M02 * r.M23) + (m.M03 * r.M33)

	res.M10 = (m.M10 * r.M00) + (m.M11 * r.M10) + (m.M12 * r.M20) + (m.M13 * r.M30)
	res.M11 = (m.M10 * r.M01) + (m.M11 * r.M11) + (m.M12 * r.M21) + (m.M13 * r.M31)
	res.M12 = (m.M10 * r.M02) + (m.M11 * r.M12) + (m.M12 * r.M22) + (m.M13 * r.M32)
	res.M13 = (m.M10 * r.M03) + (m.M11 * r.M13) + (m.M12 * r.M23) + (m.M13 * r.M33)

	res.M20 = (m.M20 * r.M00) + (m.M21 * r.M10) + (m.M22 * r.M20) + (m.M23 * r.M30)
	res.M21 = (m.M20 * r.M01) + (m.M21 * r.M11) + (m.M22 * r.M21) + (m.M23 * r.M31)
	res.M22 = (m.M20 * r.M02) + (m.M21 * r.M12) + (m.M22 * r.M22) + (m.M23 * r.M32)
	res.M23 = (m.M20 * r.M03) + (m.M21 * r.M13) + (m.M22 * r.M23) + (m.M23 * r.M33)

	res.M30 = (m.M30 * r.M00) + (m.M31 * r.M10) + (m.M32 * r.M20) + (m.M33 * r.M30)
	res.M31 = (m.M30 * r.M01) + (m.M31 * r.M11) + (m.M32 * r.M21) + (m.M33 * r.M31)
	res.M32 = (m.M30 * r.M02) + (m.M31 * r.M12) + (m.M32 * r.M22) + (m.M33 * r.M32)
	res.M33 = (m.M30 * r.M03) + (m.M31 * r.M13) + (m.M32 * r.M23) + (m.M33 * r.M33)

	return res
}

func (m *Matrix4f) Get1D() []float32 {
	return []float32{m.M00, m.M01, m.M02, m.M03, m.M10, m.M11, m.M12, m.M13, m.M20, m.M21, m.M22, m.M23, m.M30, m.M31, m.M32, m.M33}
}

func CoTangent(angle float32) float32 {
	return float32(1.0 / math.Tan(float64(angle)))
}

func DegreesToRadians(degrees float32) float32 {
	return degrees * float32(math.Pi/180)
}

func NewProjectionMatrix(width, height int) *Matrix4f {
	projectionMatrix := new(Matrix4f)

	projectionMatrix.M00 = 1
	projectionMatrix.M11 = 1
	projectionMatrix.M22 = 1
	projectionMatrix.M33 = 1

	var fieldOfView = float32(80.0)
	var aspectRatio = float32(float32(width) / float32(height))
	near_plane := 0.1
	far_plane := 300.0

	y_scale := CoTangent(DegreesToRadians(fieldOfView / 2.0))
	x_scale := y_scale / aspectRatio
	frustum_length := far_plane - near_plane

	projectionMatrix.M00 = x_scale
	projectionMatrix.M11 = y_scale
	projectionMatrix.M22 = float32(-((far_plane + near_plane) / frustum_length))
	projectionMatrix.M32 = -1
	projectionMatrix.M23 = float32(-((2 * near_plane * far_plane) / frustum_length))
	projectionMatrix.M33 = 0

	return projectionMatrix
}

func ArrayToTexture(modelMatrixElements []float32) *TextureAndBufferIds {

	var modelMatrixId uint32

	gl.GenTextures(1, &modelMatrixId)

	gl.BindTexture(gl.TEXTURE_BUFFER, modelMatrixId)

	var modelBuffer uint32

	gl.GenBuffers(1, &modelBuffer)

	gl.BindBuffer(gl.TEXTURE_BUFFER, modelBuffer)
	gl.BufferData(gl.TEXTURE_BUFFER, len(modelMatrixElements)*4, gl.Ptr(modelMatrixElements), gl.STATIC_DRAW)

	gl.TexBuffer(gl.TEXTURE_BUFFER, gl.R32F, modelBuffer)
	gl.BindTexture(gl.TEXTURE_BUFFER, 0)

	gl.BindBuffer(gl.TEXTURE_BUFFER, 0)

	return &TextureAndBufferIds{
		modelMatrixId,
		modelBuffer,
	}
}

func UpdateArrayToTexture(modelBuffer uint32, modelMatrixElements []float32) {
	gl.BindBuffer(gl.TEXTURE_BUFFER, modelBuffer)
	gl.BufferData(gl.TEXTURE_BUFFER, len(modelMatrixElements)*4, gl.Ptr(modelMatrixElements), gl.STATIC_DRAW)
	gl.TexBuffer(gl.TEXTURE_BUFFER, gl.R32F, modelBuffer)
	gl.BindBuffer(gl.TEXTURE_BUFFER, 0)
}
