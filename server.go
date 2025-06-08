package main

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"log"
	"os"
)

//go:embed index.html
var hometext []byte

func main() {
	http.HandleFunc("POST /api", put)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/", home)

	fmt.Println("Listening on :5000...")

	http.ListenAndServe(":5000", nil)
}

func put(w http.ResponseWriter, req *http.Request) {
	url := req.FormValue("input")
	log.Println("input: ", url)
	if url == "" {
		fmt.Fprintf(w, "Key is required!\n")
		return
	}
	out, err := Makesite(context.Background(), url)
	if err != nil {
		// log and then....
		w.Write([]byte("error when making call to LLM, probably out of quota"))
		return
	}
	w.Write([]byte(out))
}

func hello(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("hello"+os.Getenv("WORLD")))
}

func home(w http.ResponseWriter, req *http.Request) {
	w.Write(hometext)
}
