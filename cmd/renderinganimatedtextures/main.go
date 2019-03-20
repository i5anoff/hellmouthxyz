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
	"strings"
	"time"
)

const (
	Width = 640
	Height = 640
	Title = "Rendering animated textures"
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

	points := []float32{
		-0.15, 0.15, 0.0, 0.0,  // top left
		0.15, 0.15, 0.0, 1.0,   // top right
		-0.15, -0.15, 0.0, 2.0, // bottom left
		0.15, -0.15, 0.0, 3.0,  // bottom right
	}

	indices := []uint32{
		0,1,2,
		1,2,3,
	}

	var textureAnimation TextureAnimation

	bf := []byte(AnimatedTextureFrames)
	err := json.Unmarshal(bf, &textureAnimation)

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

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 16, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 1, gl.FLOAT, false, 16,  gl.PtrOffset(12))
	gl.EnableVertexAttribArray(1)

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

	animatedTextureCoordinates := []float32{}

	// load all of the model matrices into a single array
	for i := textureAnimation.StartFrame; i <= textureAnimation.EndFrame; i++ {
		for j := 0; j < 4; j++ {
			uvs := textureAnimation.TextureFrames[fmt.Sprintf("%d", i)][fmt.Sprintf("%d", j)] // terrible
			animatedTextureCoordinates = append(animatedTextureCoordinates, float32(uvs[0]), float32(uvs[1]))
		}
	}

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
layout (location = 1) in float in_VertexIndex;


out vec2 out_Texture;

void main() {
  int index = int(curFrame-1.0) * 8;

  int jindex = index + int(int(in_VertexIndex) * 2);

  float uvx = texelFetch(textureCoordinates, jindex + 0).r;
  float uvy = texelFetch(textureCoordinates, jindex + 1).r;
 
  out_Texture = vec2(uvx, uvy);
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

	animationFrameTime := float64(float64(1000.0)/float64(textureAnimation.FPS))
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

			if curFrame > textureAnimation.EndFrame {
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

		location = gl.GetUniformLocation(shaderProgram, common.GlStr("curFrame"))
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

