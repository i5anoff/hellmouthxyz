package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/gl/v3.3-core/gl"
	"log"
	"runtime"
	"fmt"
	"image"
	"image/png"
	"image/draw"
	"os"
	"strings"
)

const (
	Width = 640
	Height = 480
	Title = "Rendering meshes with textures"
)

func initGLFW() *glfw.Window {
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

func initOpenGL(window *glfw.Window) {
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

func main() {
	runtime.LockOSThread()

	window := initGLFW()
	initOpenGL(window)

	points := []float32{
		0.0, 0.5, 0.0, 0.5, 1.0,
		0.5, -0.5, 0.0, 1.0, 0.0,
		-0.5, -0.5, 0.0, 0.0, 0.0,
	}

	CheckError()

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(points)*4, gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 20, gl.PtrOffset(0))
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 20, gl.PtrOffset(0))
	gl.BindVertexArray(0)

	CheckError()

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

	res := image.NewRGBA(diffuse.Rect)

	for i := 0; i < len(diffuse.Pix); i += 4 {
		res.Pix[i] = diffuse.Pix[i]
		res.Pix[i+1] = diffuse.Pix[i+1]
		res.Pix[i+2] = diffuse.Pix[i+2]
		res.Pix[i+3] = diffuse.Pix[i+3]
	}

	CheckError()


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

	vertexSourceAsString := `#version 330
in vec3 in_Position;
in vec2 in_Texture;
out vec2 out_Texture;
void main() {
  out_Texture = in_Texture;
  gl_Position = vec4(in_Position, 1.0);
}`

	fragmentSourceAsString := `#version 330
uniform sampler2D diffuse;
in vec2 out_Texture;
out vec4 out_Colour;
void main() {
  out_Colour = vec4(texture(diffuse,out_Texture).rgb , 1.0);
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

	CheckError()


	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.DepthMask(true)
		gl.Disable(gl.BLEND)
		gl.UseProgram(shaderProgram)

		var location int32 = -1
		location = gl.GetUniformLocation(shaderProgram, glStr("diffuse"))
		gl.Uniform1i(location, 0)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texId)

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		CheckError()


		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, 0)

		glfw.PollEvents()
		window.SwapBuffers()
		CheckError()
	}
}

func glStr(str string) *uint8 {
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