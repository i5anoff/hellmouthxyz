package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"runtime"
	"strings"
)

const (
	Width = 480
	Height = 480
	Title = "Skinned animation"
	FPS = 1.0
)

func InitGLFW() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	glfw.WindowHint(glfw.DepthBits, 24)
	glfw.WindowHint(glfw.StencilBits, 8)

	window, err := glfw.CreateWindow(Width, Height, Title, nil, nil)

	if err != nil {
		panic(err)
	}

	window.SetCursorPos(0, 0)
	window.MakeContextCurrent()

	return window
}

func InitOpenGL(window *glfw.Window) {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	gl.ClearColor(0, 0.1568627451, 0.2039215686, 1)

	width, height := window.GetFramebufferSize()
	//width, height := window.GetSize()
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.FrontFace(gl.CCW)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
}

func GlStr(str string) *uint8 {
	if !strings.HasSuffix(str, "\x00") {
		str = str + "\x00"
	}
	return gl.Str(str)
}


func CheckError() {

	err := gl.GetError()

	if err != gl.NO_ERROR {
		log.Print("CheckError")
		log.Fatal(err)
	}
}

type Mesh struct {
	Indices []uint32 `json:"indices"`
	Coordinates []Coordinate `json:"coordinates"`
}

type Coordinate struct {
	Index int `json:"index"`
	Vertices []float32 `json:"xyz"`
	Textures []float32 `json:"uvs"`
	Skin map[string]float32 `json:"skin"`
}

func main() {
	runtime.LockOSThread()

	window := InitGLFW()
	InitOpenGL(window)

	meshes := `{"Cube": {"indices": [1, 2, 0, 3, 9, 2, 7, 4, 6, 5, 12, 13, 14, 0, 2, 15, 16, 17, 8, 10, 18, 6, 11, 7, 2, 19, 14, 17, 8, 15, 1, 3, 2, 3, 20, 9, 7, 21, 4, 5, 22, 12, 14, 23, 0, 15, 24, 16, 8, 25, 10, 6, 26, 11, 2, 9, 19, 17, 25, 8], "coordinates": [{"xyz": [-0.1, -0.0, -0.1], "index": 0, "uvs": [0.33333, 0.5], "skin": {"Bone": 0.90976, "Bone.001": 0.09024}}, {"xyz": [-0.1, -0.0, 0.1], "index": 1, "uvs": [0.0, 0.5], "skin": {"Bone": 0.91735, "Bone.001": 0.08265}}, {"xyz": [-0.1, 0.2, -0.1], "index": 2, "uvs": [0.33333, 0.25], "skin": {"Bone": 0.49966, "Bone.001": 0.50034}}, {"xyz": [-0.1, 0.2, 0.1], "index": 3, "uvs": [0.0, 0.25], "skin": {"Bone": 0.50034, "Bone.001": 0.49966}}, {"xyz": [0.1, -0.0, -0.1], "index": 4, "uvs": [0.33333, 0.5], "skin": {"Bone": 0.91735, "Bone.001": 0.08265}}, {"xyz": [0.1, -0.0, 0.1], "index": 5, "uvs": [0.33333, 0.5], "skin": {"Bone": 0.90976, "Bone.001": 0.09024}}, {"xyz": [0.1, 0.2, -0.1], "index": 6, "uvs": [0.33333, 0.75], "skin": {"Bone": 0.50034, "Bone.001": 0.49966}}, {"xyz": [0.1, 0.2, 0.1], "index": 7, "uvs": [0.0, 0.75], "skin": {"Bone": 0.49966, "Bone.001": 0.50034}}, {"xyz": [-0.1, 0.4, 0.1], "index": 8, "uvs": [0.66667, 0.5], "skin": {"Bone": 0.09024, "Bone.001": 0.90976}}, {"xyz": [-0.1, 0.4, -0.1], "index": 9, "uvs": [0.33333, 0.0], "skin": {"Bone": 0.08265, "Bone.001": 0.91735}}, {"xyz": [0.1, 0.4, -0.1], "index": 10, "uvs": [1.0, 0.75], "skin": {"Bone": 0.09024, "Bone.001": 0.90976}}, {"xyz": [0.1, 0.4, 0.1], "index": 11, "uvs": [0.0, 1.0], "skin": {"Bone": 0.08265, "Bone.001": 0.91735}}, {"xyz": [-0.1, -0.0, -0.1], "index": 12, "uvs": [0.66667, 0.75], "skin": {"Bone": 0.90976, "Bone.001": 0.09024}}, {"xyz": [0.1, -0.0, -0.1], "index": 13, "uvs": [0.33333, 0.75], "skin": {"Bone": 0.91735, "Bone.001": 0.08265}}, {"xyz": [0.1, 0.2, -0.1], "index": 14, "uvs": [0.66667, 0.25], "skin": {"Bone": 0.50034, "Bone.001": 0.49966}}, {"xyz": [-0.1, 0.2, 0.1], "index": 15, "uvs": [0.66667, 0.25], "skin": {"Bone": 0.50034, "Bone.001": 0.49966}}, {"xyz": [0.1, -0.0, 0.1], "index": 16, "uvs": [1.0, 0.0], "skin": {"Bone": 0.90976, "Bone.001": 0.09024}}, {"xyz": [0.1, 0.2, 0.1], "index": 17, "uvs": [1.0, 0.25], "skin": {"Bone": 0.49966, "Bone.001": 0.50034}}, {"xyz": [-0.1, 0.4, -0.1], "index": 18, "uvs": [0.66667, 0.75], "skin": {"Bone": 0.08265, "Bone.001": 0.91735}}, {"xyz": [0.1, 0.4, -0.1], "index": 19, "uvs": [0.66667, 0.0], "skin": {"Bone": 0.09024, "Bone.001": 0.90976}}, {"xyz": [-0.1, 0.4, 0.1], "index": 20, "uvs": [0.0, 0.0], "skin": {"Bone": 0.09024, "Bone.001": 0.90976}}, {"xyz": [0.1, -0.0, 0.1], "index": 21, "uvs": [0.0, 0.5], "skin": {"Bone": 0.90976, "Bone.001": 0.09024}}, {"xyz": [-0.1, -0.0, 0.1], "index": 22, "uvs": [0.66667, 0.5], "skin": {"Bone": 0.91735, "Bone.001": 0.08265}}, {"xyz": [0.1, -0.0, -0.1], "index": 23, "uvs": [0.66667, 0.5], "skin": {"Bone": 0.91735, "Bone.001": 0.08265}}, {"xyz": [-0.1, -0.0, 0.1], "index": 24, "uvs": [0.66667, 0.0], "skin": {"Bone": 0.91735, "Bone.001": 0.08265}}, {"xyz": [0.1, 0.4, 0.1], "index": 25, "uvs": [1.0, 0.5], "skin": {"Bone": 0.08265, "Bone.001": 0.91735}}, {"xyz": [0.1, 0.4, -0.1], "index": 26, "uvs": [0.33333, 1.0], "skin": {"Bone": 0.09024, "Bone.001": 0.90976}}]}}`

	//armatureData := `{"Armature": {"name": "Armature", "bones": {"Bone": {"name": "Bone", "matrixLocal": [1.0, 0.0, 0.0, 0.0, 0.0, 0.0, -1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0], "matrixLocalInverted": [1.0, -0.0, 0.0, -0.0, -0.0, 0.0, 1.0, 0.0, 0.0, -1.0, 0.0, -0.0, -0.0, 0.0, -0.0, 1.0]}, "Bone.001": {"name": "Bone.001", "matrixLocal": [1.0, 0.0, 0.0, 0.0, 0.0, 0.0, -1.0, 2.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0], "parentName": "Bone", "matrixLocalInverted": [1.0, -0.0, 0.0, -0.0, -0.0, 0.0, 1.0, 0.0, 0.0, -1.0, 0.0, 2.0, -0.0, 0.0, -0.0, 1.0]}}, "matrixWorld": [1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0]}}`
	//
	//boneMatrices := `{"Cube": {"ArmatureAction": {"Bone": {"1": [1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0], "2": [1.0, 0.0, 0.0, 0.0, 0.0, 0.93969, -0.34202, 0.0, 0.0, 0.34202, 0.93969, 0.0, 0.0, 0.0, 0.0, 1.0], "3": [1.0, 0.0, 0.0, 0.0, 0.0, 0.76604, -0.64279, 0.0, 0.0, 0.64279, 0.76604, 0.0, 0.0, 0.0, 0.0, 1.0]}, "Bone.001": {"1": [1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0], "2": [1.0, 0.0, 0.0, 0.0, 0.0, 0.93969, -0.34202, 0.0, 0.0, 0.34202, 0.93969, 0.0, 0.0, 0.0, 0.0, 1.0], "3": [1.0, 0.0, 0.0, 0.0, 0.0, 0.76604, -0.64279, 0.0, 0.0, 0.64279, 0.76604, 0.0, 0.0, 0.0, 0.0, 1.0]}}}}`

	// load vertices into buffer
	// load indices into buffer
	// set up simple render pass with simple mesh first

	var mesh map[string]Mesh

	bf := []byte(meshes)
	err := json.Unmarshal(bf, &mesh)

	if err != nil {
		log.Fatal(err.Error())
	}

	cube := mesh["Cube"]

	points := make([]float32, len(cube.Coordinates) * 8) // 8 = 3 location coordinates + 2 texture coordinates + 1 mesh offset + 1 skin length + 1 skin offset

	i := 0
	skinOffset := float32(0.0)
	skin := []float32{}
	for _, c := range cube.Coordinates{
		points[i] = c.Vertices[0]
		i++
		points[i] = c.Vertices[1]
		i++
		points[i] = c.Vertices[2]
		i++
		points[i] = c.Textures[0]
		i++
		points[i] = 1 - c.Textures[1]
		i++
		points[i] = 0 // hardcode the mesh offset to 0 since we only have one mesh to render
		i++
		points[i] = float32(len(c.Skin))
		i++
		points[i] = skinOffset
		i++

		skinOffset += float32(len(c.Skin))

		for bn, bv := range c.Skin {

		}

		//points = append(points, c.Vertices[0])
		//points = append(points, c.Vertices[1])
		//points = append(points, c.Vertices[2])
		//
		//points = append(points, c.Textures[0])
		//points = append(points, 1 - c.Textures[1])
		//
		//fmt.Println("added vertices and textures...")
		//fmt.Println(len(points))
	}




	// try rendering jum verts / 5?

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
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(cube.Indices)*4, gl.Ptr(cube.Indices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 32, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 32, gl.PtrOffset(12))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(2, 1, gl.FLOAT, false, 32, gl.PtrOffset(20))
	gl.EnableVertexAttribArray(2)

	gl.VertexAttribPointer(3, 1, gl.FLOAT, false, 32, gl.PtrOffset(24))
	gl.EnableVertexAttribArray(3)

	gl.VertexAttribPointer(4, 1, gl.FLOAT, false, 32, gl.PtrOffset(28))
	gl.EnableVertexAttribArray(4)


	gl.BindVertexArray(0)

	var diffuse *image.RGBA

	ioreader, err := os.Open("../../grid.png")

	if err != nil {
		log.Fatal("Error opening image ../../grid.png")
	}

	im, err := png.Decode(ioreader)

	if err != nil {
		log.Fatal("Error decoding image ../../grid.png")
	}

	switch trueim := im.(type) {
	case *image.RGBA:
		diffuse = trueim
	default:
		copy := image.NewRGBA(trueim.Bounds())
		draw.Draw(copy, trueim.Bounds(), trueim, image.Pt(0, 0), draw.Src)
		diffuse = copy
	}

	CheckError()

	var texId uint32
	gl.GenTextures(1, &texId)
	gl.BindTexture(gl.TEXTURE_2D, texId)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 1024, 1024, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(diffuse.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.BindTexture(gl.TEXTURE_2D, 0)

	//curFrame := 0
	//
	//currentModelMatrixElements := []float32{}
	//
	//// load all of the model matrices into a single array
	//for _, v := range frames.Frames {
	//	currentModelMatrixElements = append(currentModelMatrixElements, v.Matrix...)
	//}

	//var modelBuffer uint32
	//
	//gl.GenBuffers(1, &modelBuffer)
	//
	//gl.BindBuffer(gl.TEXTURE_BUFFER, modelBuffer)
	//gl.BufferData(gl.TEXTURE_BUFFER, len(currentModelMatrixElements)*4, gl.Ptr(currentModelMatrixElements), gl.STATIC_DRAW)
	//gl.BindBuffer(gl.TEXTURE_BUFFER, 0)
	//
	//var modelMatrixId uint32
	//
	//gl.GenTextures(1, &modelMatrixId)
	//
	//gl.BindTexture(gl.TEXTURE_BUFFER, modelMatrixId)
	//gl.TexBuffer(gl.TEXTURE_BUFFER, gl.R32F, modelBuffer)
	//gl.BindTexture(gl.TEXTURE_BUFFER, 0)


	vertexSourceAsString := `#version 330

layout (location = 0) in vec3 in_Position;
layout (location = 1) in vec2 in_Texture;
layout (location = 2) in float in_MeshOffset;
layout (location = 3) in float in_NumBones;
layout (location = 4) in float in_SkinOffset;

out vec2 out_Texture;

void main() {
  out_Texture = in_Texture;
  gl_Position = vec4(in_Position, 1.0);
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

	CheckError()

	//animationFrameTime := float64(float64(1000.0)/float64(FPS))
	//previousTick := time.Now()
	//animationCurrentTime := float64(0.0)

	for !window.ShouldClose() {

		//timePassed := time.Now().Sub(previousTick)
		//animationCurrentTime += float64(float64(timePassed.Nanoseconds()) / float64(1000000.0))
		//previousTick = time.Now()
		//
		//if animationCurrentTime > animationFrameTime {
		//	curFrame++
		//	animationCurrentTime = 0.0
		//
		//	if curFrame == len(frames.Frames) {
		//		curFrame = 0
		//	}
		//}

		width, height := window.GetFramebufferSize()
		gl.Viewport(0, 0, int32(width), int32(height))

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.DepthMask(true)
		gl.Disable(gl.BLEND)
		gl.UseProgram(shaderProgram)

		samplerIndex := int32(0)
		var location int32 = -1
		//location = gl.GetUniformLocation(shaderProgram, GlStr("modelMatrices"))
		//gl.Uniform1i(location, samplerIndex)
		//samplerIndex++
		//
		//location = gl.GetUniformLocation(shaderProgram, GlStr("curFrame"))
		//gl.Uniform1f(location, float32(curFrame))

		location = gl.GetUniformLocation(shaderProgram, GlStr("diffuse"))
		gl.Uniform1i(location, samplerIndex)
		samplerIndex++

		fmt.Println(location)
		fmt.Println(samplerIndex)

		//gl.ActiveTexture(gl.TEXTURE0)
		//gl.BindTexture(gl.TEXTURE_BUFFER, modelMatrixId)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texId)

		gl.BindVertexArray(vaoId)
		//gl.DrawElements(gl.TRIANGLES, int32(len(cube.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0));
		gl.DrawElements(gl.TRIANGLES, int32(len(cube.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0));

		gl.BindVertexArray(0)

		//gl.ActiveTexture(gl.TEXTURE0)
		//gl.BindTexture(gl.TEXTURE_BUFFER, 0)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, 0)

		glfw.PollEvents()
		window.SwapBuffers()
		CheckError()
	}
}

