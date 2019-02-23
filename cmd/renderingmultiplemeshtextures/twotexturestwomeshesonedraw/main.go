package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/hellmouthengine/hellmouthxyz/cmd/renderingmultiplemeshtextures"
	"log"
	"runtime"
	"fmt"
	"image"
	"image/png"
	"image/draw"
	"os"
	"math"
)

func main() {
	runtime.LockOSThread()

	window := common.InitGLFW()
	common.InitOpenGL(window)

	points := []float32{
		0.0, 0.5, 0.0,  	0.5, 1.0,   0.0,
		-0.5, -0.5, 0.0,  	0.0, 0.0,   0.0,
		0.5, -0.5, 0.0,  	1.0, 0.0,   0.0,

		0.5, 0.1, 0.0,  	0.5, 1.0,   1.0,
		0.3, 0.5, 0.0,  	0.0, 0.0,   1.0,
		0.7, 0.5, 0.0,  	1.0, 0.0,   1.0,
	}

	indices := []uint32{
		0,1,2,
		3,4,5,
	}

	var vaoId uint32
	gl.GenVertexArrays(1, &vaoId)

	var vboId uint32
	gl.GenBuffers(1, &vboId)

	var vboiId uint32
	gl.GenBuffers(1, &vboiId)

	gl.BindVertexArray(vaoId)

	gl.BindBuffer(gl.ARRAY_BUFFER, vboId)
	gl.BufferData(gl.ARRAY_BUFFER, len(points)*4, gl.Ptr(points), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vboiId)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 24, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 24,  gl.PtrOffset(12))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(2, 1, gl.FLOAT, false, 24, gl.PtrOffset(20))
	gl.EnableVertexAttribArray(2)

	gl.BindVertexArray(0)

	var diffuse *image.RGBA

	ioreader, err := os.Open("../grid.png")

	if err != nil {
		log.Fatal("Error opening image ../grid.png")
	}

	im, err := png.Decode(ioreader)

	if err != nil {
		log.Fatal("Error decoding image ../grid.png")
	}

	switch trueim := im.(type) {
	case *image.RGBA:
		diffuse = trueim
	default:
		copy := image.NewRGBA(trueim.Bounds())
		draw.Draw(copy, trueim.Bounds(), trueim, image.Pt(0, 0), draw.Src)
		diffuse = copy
	}

	var diffuse2 *image.RGBA

	ioreader, err = os.Open("../grid2.png")

	if err != nil {
		log.Fatal("Error opening image ../grid2.png")
	}

	im, err = png.Decode(ioreader)

	if err != nil {
		log.Fatal("Error decoding image ../grid2.png")
	}

	switch trueim := im.(type) {
	case *image.RGBA:
		diffuse2 = trueim
	default:
		copy := image.NewRGBA(trueim.Bounds())
		draw.Draw(copy, trueim.Bounds(), trueim, image.Pt(0, 0), draw.Src)
		diffuse2 = copy
	}

	var texId uint32
	gl.GenTextures(1, &texId)

	gl.BindTexture(gl.TEXTURE_2D_ARRAY, texId)

	levels := math.Max(math.Log(float64(1024)), math.Log(float64(1024)))

	if levels < 1 {
		levels = 1
	}

	gl.TexStorage3D(gl.TEXTURE_2D_ARRAY, int32(levels), gl.RGBA8, 1024, 1024,2)

	gl.TexSubImage3D(gl.TEXTURE_2D_ARRAY, 0, 0, 0, 0, 1024,1024, 1, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(diffuse.Pix))
	gl.TexSubImage3D(gl.TEXTURE_2D_ARRAY, 0, 0, 0, 1, 1024,1024, 1, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(diffuse2.Pix))

	gl.GenerateMipmap(gl.TEXTURE_2D_ARRAY)

	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)

	gl.BindTexture(gl.TEXTURE_2D_ARRAY, 0)

	vertexSourceAsString := `#version 330

layout (location = 0) in vec3 in_Position;
layout (location = 1) in vec2 in_Texture;
layout (location = 2) in float in_Offset;

out vec2 out_Texture;
out float out_Offset;

void main() {
  out_Texture = in_Texture;
  out_Offset = in_Offset;
  gl_Position = vec4(in_Position, 1.0);
}`

	fragmentSourceAsString := `#version 330

uniform sampler2DArray diffuse;

in vec2 out_Texture;
in float out_Offset;

out vec4 out_Colour;

void main() {
  out_Colour = vec4(texture(diffuse,vec3(out_Texture, out_Offset)).rgb , 255.0);
}`

	vs := gl.CreateShader(gl.VERTEX_SHADER)
	vertexShaderSource, vertexFree := gl.Strs(fmt.Sprintf("%s%s", vertexSourceAsString, "\x00"))
	gl.ShaderSource(vs, 1, vertexShaderSource, nil)
	defer vertexFree()
	gl.CompileShader(vs)

	fs := gl.CreateShader(gl.FRAGMENT_SHADER)
	fragmentShaderSource, fragmentFree := gl.Strs(fmt.Sprintf("%s%s", fragmentSourceAsString, "\x00"))
	gl.ShaderSource(fs, 1, fragmentShaderSource, nil)
	defer fragmentFree()
	gl.CompileShader(fs)

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, fs)
	gl.AttachShader(shaderProgram, vs)

	gl.LinkProgram(shaderProgram)
	gl.ValidateProgram(shaderProgram)

	common.CheckError()

	for !window.ShouldClose() {
		w,h := window.GetFramebufferSize()
		gl.Viewport(0,0,int32(w),int32(h))

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.DepthMask(true)
		gl.Disable(gl.BLEND)
		gl.UseProgram(shaderProgram)

		var location int32 = -1
		location = gl.GetUniformLocation(shaderProgram, common.GlStr("diffuse"))
		gl.Uniform1i(location, 0)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D_ARRAY, texId)

		gl.BindVertexArray(vaoId)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0));
		gl.BindVertexArray(0)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D_ARRAY, 0)

		glfw.PollEvents()
		window.SwapBuffers()

	}
}