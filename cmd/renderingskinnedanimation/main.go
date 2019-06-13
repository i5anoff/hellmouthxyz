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
	"time"
)

const (
	Width = 480
	Height = 480
	Title = "Skinned animation"
    VertexData = `{"Cube": {"indices": [1, 3, 0, 5, 11, 6, 4, 12, 0, 5, 2, 13, 14, 7, 15, 16, 17, 18, 10, 9, 8, 4, 19, 20, 21, 22, 7, 17, 23, 18, 1, 24, 3, 5, 25, 11, 4, 20, 12, 5, 6, 2, 14, 21, 7, 16, 26, 17, 10, 27, 9, 4, 28, 19, 21, 29, 22, 17, 30, 23], "coordinates": [{"xyz": [1.0, 0.0, -1.0], "index": 0, "uvs": [0.33333, 0.5], "totalWeight": 1.0, "skin": {"Bone": 0.9098, "Bone.001": 0.0902}}, {"xyz": [1.0, 0.0, 1.0], "index": 1, "uvs": [0.66667, 0.5], "totalWeight": 1.0, "skin": {"Bone": 0.91731, "Bone.001": 0.08269}}, {"xyz": [-1.0, 0.0, 1.0], "index": 2, "uvs": [0.66667, 0.5], "totalWeight": 1.0, "skin": {"Bone": 0.9098, "Bone.001": 0.0902}}, {"xyz": [-1.0, 0.0, -1.0], "index": 3, "uvs": [0.33333, 0.75], "totalWeight": 1.0, "skin": {"Bone": 0.91731, "Bone.001": 0.08269}}, {"xyz": [1.0, 2.0, -1.0], "index": 4, "uvs": [0.33333, 0.75], "totalWeight": 1.0, "skin": {"Bone": 0.5, "Bone.001": 0.5}}, {"xyz": [1.0, 2.0, 1.0], "index": 5, "uvs": [0.33333, 0.25], "totalWeight": 1.0, "skin": {"Bone": 0.5, "Bone.001": 0.5}}, {"xyz": [-1.0, 2.0, 1.0], "index": 6, "uvs": [0.66667, 0.25], "totalWeight": 1.0, "skin": {"Bone": 0.5, "Bone.001": 0.5}}, {"xyz": [-1.0, 2.0, -1.0], "index": 7, "uvs": [0.33333, 0.25], "totalWeight": 1.0, "skin": {"Bone": 0.5, "Bone.001": 0.5}}, {"xyz": [1.0, 4.0, -1.0], "index": 8, "uvs": [1.0, 0.5], "totalWeight": 1.0, "skin": {"Bone": 0.0902, "Bone.001": 0.9098}}, {"xyz": [1.0, 4.0, 1.0], "index": 9, "uvs": [0.66667, 0.5], "totalWeight": 1.0, "skin": {"Bone": 0.08269, "Bone.001": 0.91731}}, {"xyz": [-1.0, 4.0, -1.0], "index": 10, "uvs": [1.0, 0.75], "totalWeight": 1.0, "skin": {"Bone": 0.08269, "Bone.001": 0.91731}}, {"xyz": [-1.0, 4.0, 1.0], "index": 11, "uvs": [0.66667, 0.0], "totalWeight": 1.0, "skin": {"Bone": 0.0902, "Bone.001": 0.9098}}, {"xyz": [1.0, 0.0, 1.0], "index": 12, "uvs": [0.0, 0.5], "totalWeight": 1.0, "skin": {"Bone": 0.91731, "Bone.001": 0.08269}}, {"xyz": [1.0, 0.0, 1.0], "index": 13, "uvs": [0.33333, 0.5], "totalWeight": 1.0, "skin": {"Bone": 0.91731, "Bone.001": 0.08269}}, {"xyz": [-1.0, 0.0, 1.0], "index": 14, "uvs": [0.0, 0.5], "totalWeight": 1.0, "skin": {"Bone": 0.9098, "Bone.001": 0.0902}}, {"xyz": [-1.0, 0.0, -1.0], "index": 15, "uvs": [0.33333, 0.5], "totalWeight": 1.0, "skin": {"Bone": 0.91731, "Bone.001": 0.08269}}, {"xyz": [1.0, 0.0, -1.0], "index": 16, "uvs": [0.66667, 0.0], "totalWeight": 1.0, "skin": {"Bone": 0.9098, "Bone.001": 0.0902}}, {"xyz": [-1.0, 2.0, -1.0], "index": 17, "uvs": [1.0, 0.25], "totalWeight": 1.0, "skin": {"Bone": 0.5, "Bone.001": 0.5}}, {"xyz": [1.0, 2.0, -1.0], "index": 18, "uvs": [0.66667, 0.25], "totalWeight": 1.0, "skin": {"Bone": 0.5, "Bone.001": 0.5}}, {"xyz": [1.0, 4.0, 1.0], "index": 19, "uvs": [0.0, 1.0], "totalWeight": 1.0, "skin": {"Bone": 0.08269, "Bone.001": 0.91731}}, {"xyz": [1.0, 2.0, 1.0], "index": 20, "uvs": [0.0, 0.75], "totalWeight": 1.0, "skin": {"Bone": 0.5, "Bone.001": 0.5}}, {"xyz": [-1.0, 2.0, 1.0], "index": 21, "uvs": [0.0, 0.25], "totalWeight": 1.0, "skin": {"Bone": 0.5, "Bone.001": 0.5}}, {"xyz": [-1.0, 4.0, -1.0], "index": 22, "uvs": [0.33333, 0.0], "totalWeight": 1.0, "skin": {"Bone": 0.08269, "Bone.001": 0.91731}}, {"xyz": [1.0, 4.0, -1.0], "index": 23, "uvs": [0.66667, 0.5], "totalWeight": 1.0, "skin": {"Bone": 0.0902, "Bone.001": 0.9098}}, {"xyz": [-1.0, 0.0, 1.0], "index": 24, "uvs": [0.66667, 0.75], "totalWeight": 1.0, "skin": {"Bone": 0.9098, "Bone.001": 0.0902}}, {"xyz": [1.0, 4.0, 1.0], "index": 25, "uvs": [0.33333, 0.0], "totalWeight": 1.0, "skin": {"Bone": 0.08269, "Bone.001": 0.91731}}, {"xyz": [-1.0, 0.0, -1.0], "index": 26, "uvs": [1.0, 0.0], "totalWeight": 1.0, "skin": {"Bone": 0.91731, "Bone.001": 0.08269}}, {"xyz": [-1.0, 4.0, 1.0], "index": 27, "uvs": [0.66667, 0.75], "totalWeight": 1.0, "skin": {"Bone": 0.0902, "Bone.001": 0.9098}}, {"xyz": [1.0, 4.0, -1.0], "index": 28, "uvs": [0.33333, 1.0], "totalWeight": 1.0, "skin": {"Bone": 0.0902, "Bone.001": 0.9098}}, {"xyz": [-1.0, 4.0, 1.0], "index": 29, "uvs": [0.0, 0.0], "totalWeight": 1.0, "skin": {"Bone": 0.0902, "Bone.001": 0.9098}}, {"xyz": [-1.0, 4.0, -1.0], "index": 30, "uvs": [1.0, 0.5], "totalWeight": 1.0, "skin": {"Bone": 0.08269, "Bone.001": 0.91731}}]}}`
    ArmatureData = `{"Armature": {"bones": {"Bone": {"matrix_local": [1.0, 0.0, 0.0, 0.0, 0.0, 0.0, -1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0], "name": "Bone", "matrix_local_inverted": [1.0, -0.0, 0.0, -0.0, -0.0, 0.0, 1.0, 0.0, 0.0, -1.0, 0.0, -0.0, -0.0, 0.0, -0.0, 1.0]}, "Bone.001": {"parentName": "Bone", "matrix_local": [1.0, 0.0, 0.0, 0.0, 0.0, 0.0, -1.0, 2.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0], "name": "Bone.001", "matrix_local_inverted": [1.0, -0.0, 0.0, -0.0, -0.0, 0.0, 1.0, 0.0, 0.0, -1.0, 0.0, 2.0, -0.0, 0.0, -0.0, 1.0]}}, "matrix_world": [1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0], "name": "Armature"}}`
    AnimationMatrices = `{"Cube": {"ArmatureAction": {"Bone": {"1": [1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0], "2": [0.9763, 0.0, -0.21644, 0.0, 0.0, 1.0, 0.0, 0.0, 0.21644, 0.0, 0.9763, 0.0, 0.0, 0.0, 0.0, 1.0], "3": [0.90631, 0.0, -0.42262, 0.0, 0.0, 1.0, 0.0, 0.0, 0.42262, 0.0, 0.90631, 0.0, 0.0, 0.0, 0.0, 1.0]}, "Bone.001": {"1": [1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0], "2": [0.9763, 0.0, -0.21644, 0.0, 0.0, 1.0, 0.0, 0.0, 0.21644, 0.0, 0.9763, 0.0, 0.0, 0.0, 0.0, 1.0], "3": [0.90631, 0.0, -0.42262, 0.0, 0.0, 1.0, 0.0, 0.0, 0.42262, 0.0, 0.90631, 0.0, 0.0, 0.0, 0.0, 1.0]}}}}`
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

func main() {
	runtime.LockOSThread()

	// initialise glfw and opengl
	window := InitGLFW()
	InitOpenGL(window)


	// unmarshal vertex data
	var vertexData map[string]Mesh

	vertexByteArray := []byte(VertexData)
	err := json.Unmarshal(vertexByteArray, &vertexData)

	if err != nil {
		log.Fatal(err.Error())
	}

	// unmarshal animation data
	var animationData map[string]map[string]map[string]*IntToMatrix4fMap

	animationByteArray := []byte(AnimationMatrices)
	err = json.Unmarshal(animationByteArray, &animationData)

	if err != nil {
		log.Fatal(err.Error())
	}

	// unmarshal armature data
	var armatureData map[string]*Armature

	armatureByteArray := []byte(ArmatureData)
	err = json.Unmarshal(armatureByteArray, &armatureData)

	if err != nil {
		log.Fatal(err.Error())
	}


	// initialise what will eventually be texture buffers bound to our shader program
	boneMatrixBuffer := []float32{}
	invertedBoneMatrixBuffer := []float32{}
	skinBuffer := []float32{}
	offsetBuffer := []float32{}


	armatureActionFrames := animationData["Cube"]["ArmatureAction"]
	armature := armatureData["Armature"]
	// create a new skinned animation instance using the armature data, and the animation data
	skinnedAnimation := NewSkinnedAnimation(armature, armatureActionFrames)
	// fetch the bind matrices that were generated in the skinned animation instance
	bindMatrices := skinnedAnimation.BindMatrices

	boneNameToIndex := NewStringToFloat32Map()
	boneIndex := float32(0)

	// loop over the bones of the armature that the mesh is parented to
	for boneName := range armature.Bones {

		// get the per-frame matrices of the current bone in the armature
		frameMatrices := bindMatrices[boneName]

		for _, animationFrame := range frameMatrices.Keys() {
			frameMatrix := frameMatrices.Get(animationFrame)
			boneMatrixBuffer = append(boneMatrixBuffer, frameMatrix.Get1D()...)
		}

		// keep track of the name of the bone whose matrices we just appended to a list, along with the index with which it was iterated on
		boneNameToIndex.Set(boneName, boneIndex)

		// incremenet the bone index
		boneIndex++
	}

	// collect the inverted bone matrices for the bones in the armature based on the order they were encountered when collecting their pose matrices
	for _, boneName := range boneNameToIndex.Keys() {
		matrixLocalInverted := armature.Bones[boneName].MatrixLocalInverted
		invertedBoneMatrixBuffer = append(invertedBoneMatrixBuffer, matrixLocalInverted.Get1D()...)
	}

	cubeVertexData := vertexData["Cube"]

	for _, v := range cubeVertexData.Coordinates {
		for sk, sv := range v.Skin {
			skinBuffer = append(skinBuffer, boneNameToIndex.Get(sk))
			skinBuffer = append(skinBuffer, sv/v.TotalWeight)
		}
	}

	offsetBuffer = append(offsetBuffer, []float32{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}...)

	skinBufferID := ArrayToTexture(skinBuffer)
	boneMatrixsBufferID := ArrayToTexture(boneMatrixBuffer)
	invertedBoneMatrixBufferID := ArrayToTexture(invertedBoneMatrixBuffer)
	offsetBufferID := ArrayToTexture(offsetBuffer)



	// load all the vertex attributes for all vertices into a float array, which will later become a vertex buffer
	points := make([]float32, len(cubeVertexData.Coordinates) * 8) // 8 = 3 location coordinates + 2 texture coordinates + 1 mesh offset + 1 skin length + 1 skin offset

	currentPointElement := 0

	skinOffset := float32(0.0)

	for _, coordinate := range cubeVertexData.Coordinates{
		points[currentPointElement] = coordinate.Vertices[0]
		currentPointElement++
		points[currentPointElement] = coordinate.Vertices[1]
		currentPointElement++
		points[currentPointElement] = coordinate.Vertices[2]
		currentPointElement++
		points[currentPointElement] = coordinate.Textures[0]
		currentPointElement++
		points[currentPointElement] = 1 - coordinate.Textures[1]
		currentPointElement++
		points[currentPointElement] = 0 // hardcode the mesh offset to 0 since we only have one mesh to render
		currentPointElement++
		points[currentPointElement] = float32(len(coordinate.Skin))
		currentPointElement++
		points[currentPointElement] = skinOffset
		currentPointElement++

		skinOffset += float32(len(coordinate.Skin)) * 2
	}


	var vaoId uint32
	gl.GenVertexArrays(1, &vaoId)

	var vboId uint32
	gl.GenBuffers(1, &vboId)

	var vboiId uint32
	gl.GenBuffers(1, &vboiId)

	gl.BindVertexArray(vaoId)

	// pass the vertex attribute float array to an array buffer
	gl.BindBuffer(gl.ARRAY_BUFFER, vboId)
	gl.BufferData(gl.ARRAY_BUFFER, len(points)*4, gl.Ptr(points), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vboiId)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(cubeVertexData.Indices)*4, gl.Ptr(cubeVertexData.Indices), gl.STATIC_DRAW)

	// define the layout of a single vertex; define the first 3 floats as one attribute (xyz location)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 32, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// define the next 2 floats as one attribute (uv texture coordinates)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 32, gl.PtrOffset(12))
	gl.EnableVertexAttribArray(1)

	// define the next float as one attribute (mesh offset)
	gl.VertexAttribPointer(2, 1, gl.FLOAT, false, 32, gl.PtrOffset(20))
	gl.EnableVertexAttribArray(2)

	// define the next float as one attribute (number of bones)
	gl.VertexAttribPointer(3, 1, gl.FLOAT, false, 32, gl.PtrOffset(24))
	gl.EnableVertexAttribArray(3)

	// define the next float as one attribute (skin offset)
	gl.VertexAttribPointer(4, 1, gl.FLOAT, false, 32, gl.PtrOffset(28))
	gl.EnableVertexAttribArray(4)

	gl.BindVertexArray(0)


	// load in a texture to apply to the mesh
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


	// create an identity matrix
	modelMatrixElements := []float32{
		1,0,0,0,
		0,1,0,0,
		0,0,1,0,
		0,0,0,1,
	}

	var modelMatrixBuffer uint32

	gl.GenBuffers(1, &modelMatrixBuffer)

	gl.BindBuffer(gl.TEXTURE_BUFFER, modelMatrixBuffer)
	gl.BufferData(gl.TEXTURE_BUFFER, len(modelMatrixElements)*4, gl.Ptr(modelMatrixElements), gl.STATIC_DRAW)
	gl.BindBuffer(gl.TEXTURE_BUFFER, 0)

	var modelMatrixId uint32

	gl.GenTextures(1, &modelMatrixId)

	gl.BindTexture(gl.TEXTURE_BUFFER, modelMatrixId)
	gl.TexBuffer(gl.TEXTURE_BUFFER, gl.R32F, modelMatrixBuffer)
	gl.BindTexture(gl.TEXTURE_BUFFER, 0)

	vertexSourceAsString := `#version 330

uniform mat4 projectionMatrix;
uniform mat4 viewMatrix;

uniform samplerBuffer modelMatrices;
uniform samplerBuffer offsets; // curFrame|numFrames|meshSkinOffset|animationOffset|meshOffset|invertedMatrixOffset
uniform samplerBuffer skin;   // vertex -> { boneIndex : boneInfluence } | (skinOffset + in_VertexIndex) -> for i in in_NumberOfBones -> influence (passed by engine)
uniform samplerBuffer boneMatrices; // mesh -> animation -> frame -> { boneIndex : mat4 } | in_ModelOffset -> { animation : frame : boneIndex : mat4 }
uniform samplerBuffer invertedMatrices;

layout (location = 0) in vec3 in_Position;
layout (location = 1) in vec2 in_Texture;
layout (location = 2) in float in_ModelOffset;
layout (location = 3) in float in_NumberOfBones;
layout (location = 4) in float in_SkinOffset;

out vec2 out_Texture;

mat4 getMatrix(int index, samplerBuffer fpgbuffer){
  float m00 = texelFetch(fpgbuffer, index + 0).r;
  float m01 = texelFetch(fpgbuffer, index + 1).r;
  float m02 = texelFetch(fpgbuffer, index + 2).r;
  float m03 = texelFetch(fpgbuffer, index + 3).r;
  float m10 = texelFetch(fpgbuffer, index + 4).r;
  float m11 = texelFetch(fpgbuffer, index + 5).r;
  float m12 = texelFetch(fpgbuffer, index + 6).r;
  float m13 = texelFetch(fpgbuffer, index + 7).r;
  float m20 = texelFetch(fpgbuffer, index + 8).r;
  float m21 = texelFetch(fpgbuffer, index + 9).r;
  float m22 = texelFetch(fpgbuffer, index + 10).r;
  float m23 = texelFetch(fpgbuffer, index + 11).r;
  float m30 = texelFetch(fpgbuffer, index + 12).r;
  float m31 = texelFetch(fpgbuffer, index + 13).r;
  float m32 = texelFetch(fpgbuffer, index + 14).r;
  float m33 = texelFetch(fpgbuffer, index + 15).r;

	return mat4(m00, m10, m20, m30,
  				m01, m11, m21, m31,
  				m02, m12, m22, m32,
  				m03, m13, m23, m33);
}

mat4 getModelMatrix(){
  int index = int(in_ModelOffset*16);
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
  mat4 modelMatrix = getModelMatrix();

  int offset = int(in_ModelOffset * 6);
  float curFrame = texelFetch(offsets, offset).r;   				// current frame of animation currently playing
  vec3 mod_position = in_Position;

   if (curFrame != -1){
	  float numFramesInAnimation = texelFetch(offsets, offset + 1).r;   	// offset for how many bones are in the current meshes armature SHOULD BE NUMBER OF FRAMES IN CURRENT ANIMATION
	  float skinOffset = texelFetch(offsets, offset + 2).r; 				// how far into the universal skin object we have to go to get the skin for the current mesh
	  float animationOffset = texelFetch(offsets, offset + 3).r;			// how far into the matrices does the current animation start?
	  float meshAnimationOffset = texelFetch(offsets, offset + 4).r;		// how far into the matrices does the current mesh start?
	  float invertedMatrixOffset = texelFetch(offsets, offset + 5).r;		// how far into the inverted bone matrices does this meshes matrices start?

	  float vertexSkinOffset = skinOffset + in_SkinOffset;      			// get the starting point of the current vertices skin

	  mod_position = vec3(0,0,0);

	  for(float i=0;i<in_NumberOfBones;++i) {
	  	int vOffset = int(vertexSkinOffset + i*2);        	// for each bone that affects this vertex..

	  	float boneIndex = texelFetch(skin, vOffset).r;			// get the bones index (for fetching its mat4)
	  	float boneInfluence = texelFetch(skin, vOffset + 1).r;	// get the bones weight

	  	float matIndex = meshAnimationOffset + animationOffset +  (boneIndex * numFramesInAnimation * 16) + ((curFrame - 1) * 16);

	  	mat4 boneMat = getMatrix(int(matIndex), boneMatrices);
	  	mat4 invertedMat = getMatrix(int(invertedMatrixOffset + (boneIndex * 16)), invertedMatrices);

	    mod_position = mod_position + ((boneMat * (invertedMat * vec4(in_Position, 1.0))) *  boneInfluence ).xyz;
	  }
  }


  vec3 worldPos = (modelMatrix * vec4(mod_position, 1.0)).xyz;

  gl_Position = projectionMatrix * viewMatrix * vec4(worldPos, 1.0);
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

	projectionMatrix := NewProjectionMatrix()

	// initialise the view matrix as an identity matrix
	viewMatrix := new(Matrix4f)
	viewMatrix.M00 = 1
	viewMatrix.M11 = 1
	viewMatrix.M22 = 1
	viewMatrix.M33 = 1

	// the view matrix is effectively the inverted camera position
	viewMatrix.M03 = 0
	viewMatrix.M13 = -2
	viewMatrix.M23 = -10


	animationFrameTime := float64(float64(1000.0)/float64(skinnedAnimation.FPS))
	previousTick := time.Now()
	animationCurrentTime := float64(0.0)
	curFrame := int64(1)


	for !window.ShouldClose() {

		timePassed := time.Now().Sub(previousTick)
		animationCurrentTime += float64(float64(timePassed.Nanoseconds()) / float64(1000000.0))
		previousTick = time.Now()

		if animationCurrentTime > animationFrameTime {
			curFrame++
			animationCurrentTime = 0.0

			if curFrame > skinnedAnimation.EndFrame {
				curFrame = 1
			}
		}

		width, height := window.GetFramebufferSize()
		gl.Viewport(0, 0, int32(width), int32(height))

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.DepthMask(true)
		gl.Disable(gl.BLEND)

		complete := []float32{float32(curFrame), float32(skinnedAnimation.EndFrame), 0.0, 0.0, 0.0, 0.0}

		UpdateArrayToTexture(offsetBufferID.BufferID, complete)

		gl.UseProgram(shaderProgram)
		var location int32 = -1

		// set uniform values in the shader program
		location = gl.GetUniformLocation(shaderProgram, GlStr("projectionMatrix"))
		gl.UniformMatrix4fv(location, 1, true, &projectionMatrix.Get1D()[0])

		location = gl.GetUniformLocation(shaderProgram, GlStr("viewMatrix"))
		gl.UniformMatrix4fv(location, 1, true, &viewMatrix.Get1D()[0])

		location = gl.GetUniformLocation(shaderProgram, GlStr("modelMatrices"))
		gl.Uniform1i(location, 0)

		location = gl.GetUniformLocation(shaderProgram, GlStr("offsets"))
		gl.Uniform1i(location, 1)

		location = gl.GetUniformLocation(shaderProgram, GlStr("skin"))
		gl.Uniform1i(location, 2)

		location = gl.GetUniformLocation(shaderProgram, GlStr("boneMatrices"))
		gl.Uniform1i(location, 3)

		location = gl.GetUniformLocation(shaderProgram, GlStr("invertedMatrices"))
		gl.Uniform1i(location, 4)

		location = gl.GetUniformLocation(shaderProgram, GlStr("diffuse"))
		gl.Uniform1i(location, 5)


		// bind the buffers at the appropriate texture slots
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_BUFFER, modelMatrixId)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_BUFFER, offsetBufferID.TextureID)

		gl.ActiveTexture(gl.TEXTURE2)
		gl.BindTexture(gl.TEXTURE_BUFFER, skinBufferID.TextureID)

		gl.ActiveTexture(gl.TEXTURE3)
		gl.BindTexture(gl.TEXTURE_BUFFER, boneMatrixsBufferID.TextureID)

		gl.ActiveTexture(gl.TEXTURE4)
		gl.BindTexture(gl.TEXTURE_BUFFER, invertedBoneMatrixBufferID.TextureID)

		gl.ActiveTexture(gl.TEXTURE5)
		gl.BindTexture(gl.TEXTURE_2D, texId)

		gl.BindVertexArray(vaoId)
		gl.DrawElements(gl.TRIANGLES, int32(len(cubeVertexData.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0));
		gl.BindVertexArray(0)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_BUFFER, 0)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_BUFFER, 0)

		gl.ActiveTexture(gl.TEXTURE2)
		gl.BindTexture(gl.TEXTURE_BUFFER, 0)

		gl.ActiveTexture(gl.TEXTURE3)
		gl.BindTexture(gl.TEXTURE_BUFFER, 0)

		gl.ActiveTexture(gl.TEXTURE4)
		gl.BindTexture(gl.TEXTURE_BUFFER, 0)

		gl.ActiveTexture(gl.TEXTURE5)
		gl.BindTexture(gl.TEXTURE_2D, 0)

		glfw.PollEvents()
		window.SwapBuffers()
		CheckError()
	}
}

