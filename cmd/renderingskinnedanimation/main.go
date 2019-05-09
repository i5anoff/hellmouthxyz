package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"log"
	"runtime"
	"strings"
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


func CheckError(msg string) {

	err := gl.GetError()

	if err != gl.NO_ERROR {
		log.Printf("CheckError: %s\n", msg)
		log.Fatal(err)
	}
}

func main() {
	runtime.LockOSThread()

	window := InitGLFW()
	InitOpenGL(window)

	var meshes map[string]Mesh

	d := []byte(SkinnedVertices)
	err := json.Unmarshal(d, &meshes)

	if err != nil {
		log.Fatal(err.Error())
	}

	var vaoId uint32
	gl.GenVertexArrays(1, &vaoId)

	var vboId uint32
	gl.GenBuffers(1, &vboId)

	var vboiId uint32
	gl.GenBuffers(1, &vboiId)


	rawVertices := meshes["Cube"].Vertices

	// X, Y, Z, U, V, model offset, number of bones, skin offset
	vertices := []float32{}

	for _, v := range rawVertices {
		vertices = append(vertices, v.XYZ...)
	}

	indices := meshes["Cube"].Indices

	gl.BindVertexArray(vaoId)

	gl.BindBuffer(gl.ARRAY_BUFFER, vboId)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vboiId)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 12, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)


	vertexSourceAsString := `#version 330

layout (location = 0) in vec3 in_Position;

void main() { 
  gl_Position = vec4(in_Position, 1.0);
}
`
	fragmentSourceAsString := `#version 330

out vec4 out_Colour;

void main() {
  
  out_Colour = vec4(1.0, 0.0, 1.0 , 1.0);
}
`

	vs := gl.CreateShader(gl.VERTEX_SHADER)
	vertexShaderSource, vertexFree := gl.Strs(fmt.Sprintf("%s%s", vertexSourceAsString, "\x00"))
	gl.ShaderSource(vs, 1, vertexShaderSource, nil)
	defer vertexFree()
	gl.CompileShader(vs)

	var status int32
	gl.GetShaderiv(vs, gl.COMPILE_STATUS, &status)

	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(vs, gl.INFO_LOG_LENGTH, &logLength)

		loge := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(vs, logLength, nil, gl.Str(loge))
		log.Fatal("%s SHADER LOG: %s", "vs", loge)
	}

	fs := gl.CreateShader(gl.FRAGMENT_SHADER)
	fragmentShaderSource, fragmentFree := gl.Strs(fmt.Sprintf("%s%s", fragmentSourceAsString, "\x00"))
	gl.ShaderSource(fs, 1, fragmentShaderSource, nil)
	defer fragmentFree()
	gl.CompileShader(fs)

	gl.GetShaderiv(fs, gl.COMPILE_STATUS, &status)

	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(fs, gl.INFO_LOG_LENGTH, &logLength)

		loge := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(fs, logLength, nil, gl.Str(loge))
		log.Fatal("%s SHADER LOG: %s", "vs", loge)
	}

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, fs)
	gl.AttachShader(shaderProgram, vs)

	gl.LinkProgram(shaderProgram)

	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shaderProgram, logLength, nil, GlStr(log))
		fmt.Println("LINK STATUS")
		fmt.Println(log)
	}


	gl.ValidateProgram(shaderProgram)
	var vStatus int32
	gl.GetProgramiv(shaderProgram, gl.VALIDATE_STATUS, &vStatus)

	if vStatus == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shaderProgram, logLength, nil, GlStr(log))
		fmt.Println("VALIDATE STATUS")
		fmt.Println(log)
	}

	for !window.ShouldClose() {

		width, height := window.GetFramebufferSize()
		gl.Viewport(0, 0, int32(width), int32(height))

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(shaderProgram)

		gl.EnableVertexAttribArray(0)

		gl.BindVertexArray(vaoId)
		gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, gl.PtrOffset(0));
		gl.BindVertexArray(0)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

