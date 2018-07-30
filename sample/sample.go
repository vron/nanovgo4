// +build !js

package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/vron/nanovgo4"
	"github.com/vron/nanovgo4/perfgraph"
	"github.com/vron/nanovgo4/sample/demo"
)

var blowup bool
var premult bool

func key(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	} else if key == glfw.KeySpace && action == glfw.Press {
		blowup = !blowup
	} else if key == glfw.KeyP && action == glfw.Press {
		premult = !premult
	}
}

func init() {
}

func main() {
	runtime.LockOSThread()

	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	// demo MSAA
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Samples, 16)

	window, err := glfw.CreateWindow(1000, 600, "nanovgo4", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	window.SetKeyCallback(key)

	ctx, err := nanovgo4.NewContext(0) // | nanovgo4.Debug /*nanovgo4.AntiAlias | nanovgo4.StencilStrokes | nanovgo4.Debug*/)
	defer ctx.Delete()

	if err != nil {
		panic(err)
	}

	demoData := LoadDemo(ctx)

	glfw.SwapInterval(0)

	fps := perfgraph.NewPerfGraph("Frame Time", "sans")

	for !window.ShouldClose() {
		t, _ := fps.UpdateGraph()

		fbWidth, fbHeight := window.GetFramebufferSize()
		winWidth, winHeight := window.GetSize()
		mx, my := window.GetCursorPos()

		pixelRatio := float32(fbWidth) / float32(winWidth)
		gl.Viewport(0, 0, int32(fbWidth), int32(fbHeight))
		if premult {
			gl.ClearColor(0, 0, 0, 0)
		} else {
			gl.ClearColor(0.3, 0.3, 0.32, 1.0)
		}
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
		gl.Enable(gl.CULL_FACE)
		gl.Disable(gl.DEPTH_TEST)

		ctx.BeginFrame(winWidth, winHeight, pixelRatio)

		demo.RenderDemo(ctx, float32(mx), float32(my), float32(winWidth), float32(winHeight), t, blowup, demoData)
		fps.RenderGraph(ctx, 5, 5)
		ctx.EndFrame()

		gl.Enable(gl.DEPTH_TEST)
		window.SwapBuffers()
		glfw.PollEvents()
	}

	demoData.FreeData(ctx)
}

func LoadDemo(ctx *nanovgo4.Context) *demo.DemoData {
	d := &demo.DemoData{}
	for i := 0; i < 12; i++ {
		path := fmt.Sprintf("images/image%d.jpg", i+1)
		d.Images = append(d.Images, ctx.CreateImage(path, 0))
		if d.Images[i] == 0 {
			log.Fatalf("Could not load %s", path)
		}
	}

	d.FontIcons = ctx.CreateFont("icons", "entypo.ttf")
	if d.FontIcons == -1 {
		log.Fatalln("Could not add font icons.")
	}
	d.FontNormal = ctx.CreateFont("sans", "Roboto-Regular.ttf")
	if d.FontNormal == -1 {
		log.Fatalln("Could not add font italic.")
	}
	d.FontBold = ctx.CreateFont("sans-bold", "Roboto-Bold.ttf")
	if d.FontBold == -1 {
		log.Fatalln("Could not add font bold.")
	}
	return d
}
