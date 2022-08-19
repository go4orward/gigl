package main

import (
	"fmt"

	"github.com/go4orward/gigl/common"
	"github.com/go4orward/gigl/env/webgl10"
)

type Config struct {
	loglevel  string //
	logfilter string //
}

func main() {
	cfg := Config{loglevel: "info", logfilter: ""}
	if cfg.loglevel != "" {
		common.SetLogger(common.NewConsoleLogger(cfg.loglevel)).SetTraceFilter(cfg.logfilter).SetOption("", false)
	}
	// THIS CODE IS SUPPOSED TO BE BUILT AS WEBASSEMBLY AND RUN INSIDE A BROWSER.
	// BUILD IT LIKE 'GOOS=js GOARCH=wasm go build -o example.wasm examples/example.go'.
	fmt.Println("Hello WebGL 1.0")                      // printed in the browser console
	canvas, err := webgl10.NewWebGLCanvas("wasmcanvas") // ID of canvas element
	if err != nil {
		common.Logger.Error("Failed to start WebGL : %v\n", err)
		return
	}

	webgl10.DrawSimplestTriangle(canvas)

	canvas.RunOnce(nil)
}
