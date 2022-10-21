package main

import (
	// "fmt"
	// "io/ioutil"
	// "log"
	"net/http"
	// "html/template"
	// "github.com/gorilla/mux"
)

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>This is the about page</h1>"))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello!</h1>"))
}

func rootHandlers(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../../assets/index.html")
}

func main() {
	// handle func is convienance method on http. registers function to a path
	// on default serve mux.
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/about", aboutHandler)
	mux.HandleFunc("/test", rootHandlers)
	// staticHandler := http.FileServer(http.Dir("/assets"))
	http.ListenAndServe(":9090", mux)

	// the default serve mux is an http handler
	// everything related to the server in go is an http handler
	// serve mux basically just maps (routes) a function to a pattern
	// this below --v uses the default serve mux

}
