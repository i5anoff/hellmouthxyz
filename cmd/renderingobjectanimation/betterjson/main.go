package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/hellmouthengine/hellmouthxyz/cmd/renderingobjectanimation"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"runtime"
	"time"
)

type Frames struct {
	Frames []Frame `json:"frames"`
}

type Frame struct {
	Frame int64 `json:"frame"`
	Matrix []float32 `json:"matrix"`
}

func main() {
	runtime.LockOSThread()

	window := common.InitGLFW()
	common.InitOpenGL(window)

	points := []float32{
		0.0, 0.5, 0.0,  	0.5, 1.0,
		-0.5, -0.5, 0.0,  	0.0, 0.0,
		0.5, -0.5, 0.0,  	1.0, 0.0,
	}

	indices := []uint32{
		0,1,2,
	}

	modelMatricesJSON := `{
  "frames": [
    {
      "frame": 1,
      "matrix": [
        1,
        0,
        0,
        -0.5,
        0,
        1,
        0,
        -1,
        0,
        0,
        1,
        0,
        0,
        0,
        0,
        1
      ]
    },
    {
      "frame": 2,
      "matrix": [
        1,
        0,
        0,
        -0.3,
        0,
        1,
        0,
        -1,
        0,
        0,
        1,
        0,
        0,
        0,
        0,
        1
      ]
    },
    {
      "frame": 3,
      "matrix": [
        1,
        0,
        0,
        -0.1,
        0,
        1,
        0,
        -1,
        0,
        0,
        1,
        0,
        0,
        0,
        0,
        1
      ]
    },
    {
      "frame": 4,
      "matrix": [
        1,
        0,
        0,
        0.1,
        0,
        1,
        0,
        -1,
        0,
        0,
        1,
        0,
        0,
        0,
        0,
        1
      ]
    },
    {
      "frame": 5,
      "matrix": [
        1,
        0,
        0,
        0.3,
        0,
        1,
        0,
        -1,
        0,
        0,
        1,
        0,
        0,
        0,
        0,
        1
      ]
    },
    {
      "frame": 6,
      "matrix": [
        1,
        0,
        0,
        0.5,
        0,
        1,
        0,
        -1,
        0,
        0,
        1,
        0,
        0,
        0,
        0,
        1
      ]
    }
  ]
}`

	var frames Frames

	bf := []byte(modelMatricesJSON)
	err := json.Unmarshal(bf, &frames)

	if err != nil {
		log.Fatal(err.Error())
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

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 20, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 20,  gl.PtrOffset(12))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)

	var diffuse *image.RGBA

	ioreader, err := os.Open("grid.png")

	if err != nil {
		log.Fatal("Error opening image grid.png")
	}

	im, err := png.Decode(ioreader)

	if err != nil {
		log.Fatal("Error decoding image grid.png")
	}

	switch trueim := im.(type) {
	case *image.RGBA:
		diffuse = trueim
	default:
		copy := image.NewRGBA(trueim.Bounds())
		draw.Draw(copy, trueim.Bounds(), trueim, image.Pt(0, 0), draw.Src)
		diffuse = copy
	}

	common.CheckError()

	var texId uint32
	gl.GenTextures(1, &texId)
	gl.BindTexture(gl.TEXTURE_2D, texId)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, 1024, 1024, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(diffuse.Pix))

	gl.GenerateMipmap(gl.TEXTURE_2D)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)

	gl.BindTexture(gl.TEXTURE_2D, 0)

	curFrame := 0

	currentModelMatrixElements := []float32{}

	// load all of the model matrices into a single array
	for _, v := range frames.Frames {
		currentModelMatrixElements = append(currentModelMatrixElements, v.Matrix...)
	}

	var modelBuffer uint32

	gl.GenBuffers(1, &modelBuffer)

	gl.BindBuffer(gl.TEXTURE_BUFFER, modelBuffer)
	gl.BufferData(gl.TEXTURE_BUFFER, len(currentModelMatrixElements)*4, gl.Ptr(currentModelMatrixElements), gl.STATIC_DRAW)
	gl.BindBuffer(gl.TEXTURE_BUFFER, 0)

	var modelMatrixId uint32

	gl.GenTextures(1, &modelMatrixId)

	gl.BindTexture(gl.TEXTURE_BUFFER, modelMatrixId)
	gl.TexBuffer(gl.TEXTURE_BUFFER, gl.R32F, modelBuffer)
	gl.BindTexture(gl.TEXTURE_BUFFER, 0)


	vertexSourceAsString := `#version 330

uniform samplerBuffer modelMatrices;
uniform float curFrame;

layout (location = 0) in vec3 in_Position;
layout (location = 1) in vec2 in_Texture;

out vec2 out_Texture;

mat4 getModelMatrix(){
  int index = int(curFrame*16);
  float m00 = texelFetch(modelMatrices, index + 0).r;
  float m01 = texelFetch(modelMatrices, index + 1).r;
  float m02 = texelFetch(modelMatrices, index + 2).r;
  float m03 = texelFetch(modelMatrices, index + 3).r;
  float m10 = texelFetch(modelMatrices, index + 4).r;
  float m11 = texelFetch(modelMatrices, index + 5).r;
  float m12 = texelFetch(modelMatrices, index + 6).r;
  float m13 = texelFetch(modelMatrices, index + 7).r;
  float m20 = texelFetch(modelMatrices, index + 8).r;
  float m21 = texelFetch(modelMatrices, index + 9).r;
  float m22 = texelFetch(modelMatrices, index + 10).r;
  float m23 = texelFetch(modelMatrices, index + 11).r;
  float m30 = texelFetch(modelMatrices, index + 12).r;
  float m31 = texelFetch(modelMatrices, index + 13).r;
  float m32 = texelFetch(modelMatrices, index + 14).r;
  float m33 = texelFetch(modelMatrices, index + 15).r;
  
 return mat4(	m00, m10, m20, m30, 
 				m01, m11, m21, m31, 
 				m02, m12, m22, m32, 
 				m03, m13, m23, m33);
}

void main() {
  out_Texture = in_Texture;
  mat4 mm = getModelMatrix();
  gl_Position = mm * vec4(in_Position, 1.0);
}
`

	fragmentSourceAsString := `#version 330

uniform sampler2D diffuse;

in vec2 out_Texture;

out vec4 out_Colour;

void main() {
  out_Colour = vec4(texture(diffuse,out_Texture).rgb , 1.0);
}
`

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

	animationFrameTime := float64(float64(1000.0)/float64(common.FPS))
	previousTick := time.Now()
	animationCurrentTime := float64(0.0)

	for !window.ShouldClose() {

		timePassed := time.Now().Sub(previousTick)
		animationCurrentTime += float64(float64(timePassed.Nanoseconds()) / float64(1000000.0))
		previousTick = time.Now()

		if animationCurrentTime > animationFrameTime {
			curFrame++
			animationCurrentTime = 0.0

			if curFrame == len(frames.Frames) {
				curFrame = 0
			}
		}

		width, height := window.GetFramebufferSize()
		gl.Viewport(0, 0, int32(width), int32(height))

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.DepthMask(true)
		gl.Disable(gl.BLEND)
		gl.UseProgram(shaderProgram)

		samplerIndex := int32(0)
		var location int32 = -1
		location = gl.GetUniformLocation(shaderProgram, common.GlStr("modelMatrices"))
		gl.Uniform1i(location, samplerIndex)
		samplerIndex++

		location = gl.GetUniformLocation(shaderProgram, common.GlStr("curFrame"))
		gl.Uniform1f(location, float32(curFrame))

		location = gl.GetUniformLocation(shaderProgram, common.GlStr("diffuse"))
		gl.Uniform1i(location, samplerIndex)
		samplerIndex++

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_BUFFER, modelMatrixId)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texId)

		gl.BindVertexArray(vaoId)
		gl.DrawElements(gl.TRIANGLES, 3, gl.UNSIGNED_INT, gl.PtrOffset(0));
		gl.BindVertexArray(0)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_BUFFER, 0)
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, 0)

		glfw.PollEvents()
		window.SwapBuffers()
		common.CheckError()
	}
}

