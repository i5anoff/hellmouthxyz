package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/gl/v3.3-core/gl"
	"log"
	"runtime"
	"fmt"
)

const (
	Width = 640
	Height = 480
	Title = "Basic Rendering Pipeline"
)

func initGLFW() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	// create a window, requesting OpenGL version 3.3 with the core profile and forward compatibility enabled
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// specify a 24-bit depth buffer, and an 8-bit stencil buffer; 32-bit depth buffers cause problems on old macbooks
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

	// specify a color to wipe the screen with after every frame is drawn; thanks to Ethan Schoonover for Solarized
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
		0.0, 0.5, 0.0,
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
	}

	indices := []uint32{
		0,1,2,
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

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 12, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(0)

	vertexSourceAsString := `#version 330

layout (location = 0) in vec3 in_Position;

void main() {
  gl_Position = vec4(in_Position, 1.0);
}
`

	fragmentSourceAsString := `#version 330

out vec4 out_Colour;

void main() {
  out_Colour = vec4(0.521568, 0.603921, 0.0, 1.0);
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

	for !window.ShouldClose() {
		width, height := window.GetFramebufferSize()
		gl.Viewport(0, 0, int32(width), int32(height))
		// wipe the drawing surface clear
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// bind the shader program to be used in this render pass
		gl.UseProgram(shaderProgram)
		// bind the vao containing the triangle vertex data
		gl.BindVertexArray(vaoId)
		// draw 3 vertices from the currently bound VAO with current in-use shader
		gl.DrawElements(gl.TRIANGLES, 3, gl.UNSIGNED_INT, gl.PtrOffset(0));
		// update other events like input handling
		glfw.PollEvents()
		// put the stuff we've been drawing onto the display
		window.SwapBuffers()
	}
}
