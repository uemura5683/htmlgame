package main

import (
	"syscall/js"
	"fmt"
)


func main() {
	c := make(chan struct{})
	window := js.Global()
	document := window.Get("document")
	canvasEl := document.Call("getElementById", "canvas")

	bodyW := window.Get("innerWidth").Float() * 0.44
	bodyH := window.Get("innerHeight").Float()
	fmt.Println(bodyW, bodyH)
	canvasEl.Set("width", bodyW)
	canvasEl.Set("height", bodyH)

	ctx := canvasEl.Call("getContext", "2d")

	drawImage := func(x, y float64) {
		img := window.Get("Image").New()
		img.Set("src", "images/go.svg")
		imageWidth := img.Get("width").Float()
		imageHeight := img.Get("height").Float()
		ctx.Call("drawImage", img, x-imageWidth/2, y-imageHeight/2)
		img.Call("addEventListener", "load", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return nil
		}))
	}

	clearCanvas := func() {
		ctx := canvasEl.Call("getContext", "2d")
		ctx.Call("clearRect", 0, 0, bodyW, bodyH)
	}

	e := goids.CreateEnv(bodyW, bodyH, 30, 4, 2, 100)
	var animation js.Func
	animation = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		clearCanvas()
		e.SetHeight(bodyH)
		e.SetWidth(bodyW)
		e.Update()
		for _, goid := range e.Goids() {
			drawImage(goid.Position().X, goid.Position().Y, goid.ImageType())
		}

		window.Call("requestAnimationFrame", animation)
		return nil
	})
	defer animation.Release()

	window.Call("requestAnimationFrame", animation)

	<-c
}