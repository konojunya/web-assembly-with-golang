package main

import (
	"encoding/hex"
	"html"
	"html/template"
	"io"
	"log"
	"net/http"
)

type viewData struct {
	Wasm string
}

var wasmAdd = []byte{
	0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00,
	0x01, 0x07, 0x01, 0x60, 0x02, 0x7f, 0x7f, 0x01,
	0x7f, 0x03, 0x02, 0x01, 0x00, 0x07, 0x07, 0x01,
	0x03, 0x61, 0x64, 0x64, 0x00, 0x00, 0x0a, 0x09,
	0x01, 0x07, 0x00, 0x20, 0x00, 0x20, 0x01, 0x6a,
	0x0b,
}

func rootHandle(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/index.html")
	if err != nil {
		log.Fatal(err)
	}

	data := viewData{
		Wasm: html.EscapeString(hex.Dump(wasmAdd)),
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func wasmHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/wasm")
	n, err := w.Write(wasmAdd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	if n != len(wasmAdd) {
		http.Error(w, io.ErrShortWrite.Error(), http.StatusServiceUnavailable)
	}
}

func main() {
	http.HandleFunc("/", rootHandle)
	http.HandleFunc("/wasm", wasmHandle)
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
