// This code launches a simple HTTP server to test the pre-built WASM bundle
// (This code was inspired by https://github.com/bobcob7/wasm-basic-triangle)
package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
)

func main() {
	// HTML file
	html, err := ioutil.ReadFile("./webgl_test.html")
	if err != nil {
		log.Fatalf("Could not read webgl_test.html file: %s\n", err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(html)
	})

	// texture images
	http.Handle("/assets/", http.FileServer(http.Dir("../..")))

	// 'wasm_exec.js'
	exjs, err := ioutil.ReadFile(runtime.GOROOT() + "/misc/wasm/wasm_exec.js")
	if err != nil {
		log.Fatalf("Could not read wasm_exec.js file: %s\n", err)
	}
	http.HandleFunc("/wasm_exec.js", func(w http.ResponseWriter, r *http.Request) {
		w.Write(exjs)
	})

	// 'webgl_test.wasm'
	wasm, err := ioutil.ReadFile("./webgl_test.wasm")
	if err != nil {
		log.Fatalf("Could not read wasm file: %s\n", err)
	}
	http.HandleFunc("/webgl_test.wasm", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/wasm")
		w.WriteHeader(http.StatusOK)
		w.Write(wasm)
	})

	// start the server
	log.Printf("WebGL test sever is running.\n")
	log.Printf("Open http://localhost:8080 in your browser.\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
