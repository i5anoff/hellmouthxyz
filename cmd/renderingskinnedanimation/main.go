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
	Width = 640
	Height = 640
	Title = "Rendering skinned animation"
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

	//gl.ClearColor(0, 0.1568627451, 0.2039215686, 1)
	gl.ClearColor(0, 0,0, 1)

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

func main() {
	runtime.LockOSThread()

	window := InitGLFW()
	InitOpenGL(window)

	var boneData Bones

	d := []byte(BoneMatrices)
	err := json.Unmarshal(d, &boneData)

	if err != nil {
		log.Fatal(err.Error())
	}

	var meshes Meshes

	d = []byte(SkinnedVertices)
	err = json.Unmarshal(d, &meshes)

	if err != nil {
		log.Fatal(err.Error())
	}

	var armature *Armature

	d = []byte(ArmatureData)
	err = json.Unmarshal(d, &armature)

	if err != nil {
		log.Fatal(err.Error())
	}

	keyframes := boneData.Data["Cube"]["ArmatureAction"]

	sa := NewSkinnedAnimation(armature, keyframes)

	var vaoId uint32
	gl.GenVertexArrays(1, &vaoId)

	var vboId uint32
	gl.GenBuffers(1, &vboId)

	var vboiId uint32
	gl.GenBuffers(1, &vboiId)

	gl.BindVertexArray(vaoId)


	rawVertices := meshes.Data["Cube"].Vertices

	// X, Y, Z, U, V, Index
	vertices := []float32{}

	for _, v := range rawVertices {
		vertices = append(vertices, v.XYZ...)
		vertices = append(vertices, v.UV...)
		vertices = append(vertices, v.Index)
	}

	indices := meshes.Data["Cube"].Indices

	gl.BindBuffer(gl.ARRAY_BUFFER, vboId)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vboiId)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 24, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 24, gl.PtrOffset(12))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(2, 1, gl.FLOAT, false, 24,  gl.PtrOffset(20))
	gl.EnableVertexAttribArray(2)

	gl.BindVertexArray(0)

	var diffuse *image.RGBA

	ioreader, err := os.Open("../../spritesheet.png")

	if err != nil {
		log.Fatal("Error opening image ../../spritesheet.png")
	}

	im, err := png.Decode(ioreader)

	if err != nil {
		log.Fatal("Error decoding image ../../spritesheet.png")
	}

	switch trueim := im.(type) {
	case *image.RGBA:
		diffuse = trueim
	default:
		copy := image.NewRGBA(trueim.Bounds())
		draw.Draw(copy, trueim.Bounds(), trueim, image.Pt(0, 0), draw.Src)
		diffuse = copy
	}

	var texId uint32
	gl.GenTextures(1, &texId)
	gl.BindTexture(gl.TEXTURE_2D, texId)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, 736, 736, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(diffuse.Pix))

	gl.GenerateMipmap(gl.TEXTURE_2D)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)

	gl.BindTexture(gl.TEXTURE_2D, 0)

	var textureCoordinatesBufffer uint32

	gl.GenBuffers(1, &textureCoordinatesBufffer)

	gl.BindBuffer(gl.TEXTURE_BUFFER, textureCoordinatesBufffer)
	gl.BufferData(gl.TEXTURE_BUFFER, len(animatedTextureCoordinates)*4, gl.Ptr(animatedTextureCoordinates), gl.STATIC_DRAW)
	gl.BindBuffer(gl.TEXTURE_BUFFER, 0)

	var textureCoordinatesId uint32

	gl.GenTextures(1, &textureCoordinatesId)

	gl.BindTexture(gl.TEXTURE_BUFFER, textureCoordinatesId)
	gl.TexBuffer(gl.TEXTURE_BUFFER, gl.R32F, textureCoordinatesBufffer)
	gl.BindTexture(gl.TEXTURE_BUFFER, 0)

	vertexSourceAsString := `#version 330

uniform samplerBuffer textureCoordinates;
uniform float curFrame;

layout (location = 0) in vec3 in_Position;
layout (location = 1) in vec2 in_Texture;
layout (location = 2) in float in_VertexIndex;

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

	animationFrameTime := float64(float64(1000.0)/float64(sa.FPS))
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

			if curFrame > sa.EndFrame {
				curFrame = 1
			}
		}

		width, height := window.GetFramebufferSize()
		gl.Viewport(0, 0, int32(width), int32(height))

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.DepthMask(true)
		gl.Disable(gl.BLEND)

		gl.UseProgram(shaderProgram)

		var location int32 = -1

		location = gl.GetUniformLocation(shaderProgram, GlStr("curFrame"))
		gl.Uniform1f(location, float32(curFrame))

		location = gl.GetUniformLocation(shaderProgram, GlStr("textureCoordinates"))
		gl.Uniform1i(location, 0)

		location = gl.GetUniformLocation(shaderProgram, GlStr("diffuse"))
		gl.Uniform1i(location, 1)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_BUFFER, textureCoordinatesId)
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texId)

		gl.BindVertexArray(vaoId)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0));
		gl.BindVertexArray(0)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_BUFFER, 0)
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, 0)

		glfw.PollEvents()
		window.SwapBuffers()
		CheckError()
	}
}

func NewSkinnedAnimation(armature *Armature, keyframes map[string]*IntToMatrix4fMap) *SkinnedAnimation {
	sa := &SkinnedAnimation{
		map[string]*IntToMatrix4fMap{},
		map[string]*Matrix4f{},
		map[string]*IntToMatrix4fMap{},
		armature,
		1,
		true,
		true,
		1.0 / 1,
		time.Now().UnixNano(),
		0.0,
		"Cube",
		1,
		3,
		1,
	}

	sa.generateFrames(keyframes)
	sa.calculateFinalMatrices(armature)

	return sa
}

func (a *SkinnedAnimation) generateFrames(exportedKeyframes map[string]*IntToMatrix4fMap) {

	for boneName, frameToMatrix := range exportedKeyframes {

		keyframes := []int{}
		bindPoses := []*Matrix4f{}

		for _, key := range frameToMatrix.Keys() {

			m := frameToMatrix.Get(key)

			keyframes = append(keyframes, key)
			bindPoses = append(bindPoses, m)
		}

		bindPoseMatrices := NewIntToMatrix4fMap()

		for i := 0; i < len(keyframes)-1; i++ {
			keyFrameA := keyframes[i]
			keyFrameB := keyframes[i+1]
			bindPoseKeyFrameA := bindPoses[i]
			bindPoseKeyFrameB := bindPoses[i+1]
			if i == 0 {
				bindPoseMatrices.Set(keyFrameA, bindPoseKeyFrameA)
			}
			bindPoseMatrices.Set(keyFrameB, bindPoseKeyFrameB)
		}

		a.AllBindPoseTransformations[boneName] = bindPoseMatrices
	}
}

type SkinnedAnimation struct {
	AllBindPoseTransformations map[string]*IntToMatrix4fMap
	InvertedMatrices           map[string]*Matrix4f
	BindMatrices               map[string]*IntToMatrix4fMap

	Armature        *Armature
	CurrentFrame    int64
	Playing         bool
	ForwardPlay     bool
	FrameTime       float64
	LastTime        int64
	UnprocessedTime float64
	MeshName        string
	StartFrame      int64
	EndFrame        int64
	FPS 			int64
}

func (a *SkinnedAnimation) calculateFinalMatrices(armature *Armature) {

	for _, bone := range armature.Bones {
		a.InvertedMatrices[bone.Name] = bone.MatrixLocalInverted
		bindPoseMatrices := a.AllBindPoseTransformations[bone.Name]
		frames := NewIntToMatrix4fMap()

		if bindPoseMatrices != nil {
			for _, frame := range bindPoseMatrices.Keys() {
				mw := a.matrixWorld(bone, frame)
				frames.Set(frame, mw)
			}
		} else {
			for i := a.StartFrame; i < a.EndFrame; i++ {
				frames.Set(int(i), bone.MatrixLocal)
			}
		}
		a.BindMatrices[bone.Name] = frames
	}
	fmt.Printf("Finished Setting BindMatrices")

}

func (a *SkinnedAnimation) matrixWorld(bone *Bone, frame int) *Matrix4f {
	parent := a.getParent(bone)

	basis := a.AllBindPoseTransformations[bone.Name].Get(frame)
	if parent == nil {
		return bone.MatrixLocal.Mul(basis)
	} else {
		return a.invertedParentChildLocal(parent, bone, frame, basis)
	}
}

func (a *SkinnedAnimation) invertedParentChildLocal(parent, child *Bone, frame int, childBasis *Matrix4f) *Matrix4f {
	parentMatrix := a.matrixWorld(parent, frame)
	invertedParentChildLocal := (parent.MatrixLocalInverted).Mul(child.MatrixLocal)
	invertedParentChildLocalBasis := invertedParentChildLocal.Mul(childBasis)
	return parentMatrix.Mul(invertedParentChildLocalBasis)
}

func (a *SkinnedAnimation) getParent(bone *Bone) *Bone {
	var parent *Bone
	if len(bone.ParentName) > 0 {
		parent = a.Armature.Bones[bone.ParentName]
	}
	return parent
}