package main

import (
	"encoding/json"
	"fmt"
	// "io/ioutil"
	// "log"
	"net/http"
	"strings"

	// "html/template"
	"strconv"

	"github.com/gorilla/mux"
	"test.com/m/internal/database"
)

func parseParams(req *http.Request, prefix string, num int) ([]string, error) {
	url := strings.TrimPrefix(req.URL.Path, prefix)
	params := strings.Split(url, "/")
	if len(params) != num || len(params[0]) == 0 || len(params[1]) == 0 {
		return nil, fmt.Errorf("Bad format. Expecting exactly %d params", num)
	}
	return params, nil
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>This is the about page</h1>"))
}

func listenStreamHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print(mux.Vars(r))
	// streamer := mux.Vars(r)["stream"]
	// livestream.StartStream(streamer)
}

func upsertStreamerHandler(w http.ResponseWriter, r *http.Request) {
	streamer := mux.Vars(r)["stream"]
	is_active, _ := strconv.ParseBool(mux.Vars(r)["disable"])
	err := database.InsertStreamer(streamer, is_active)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func listStreamersHandler(w http.ResponseWriter, r *http.Request) {
	x := database.GetStreamerData()
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(x)
	w.Write([]byte("<h1>This is the about page</h1>"))
}

func disableStreamerHandler(w http.ResponseWriter, r *http.Request) {
	streamer := mux.Vars(r)["stream"]
	database.InsertStreamer(streamer, false)
	w.WriteHeader(200)
	w.Write([]byte(streamer))
}

func rootHandlers(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../../assets/index.html")
}

func main() {
	// handle func is convienance method on http. registers function to a path
	// on default serve mux.
	mux := mux.NewRouter().StrictSlash(true)
	
	mux.HandleFunc("/stream", listenStreamHandler).Queries("stream", "{stream}").Methods("POST")
	mux.HandleFunc("/about", aboutHandler)
	mux.HandleFunc("/stream/upsert", upsertStreamerHandler).Queries("stream", "{stream}", "disable", "{disable}").Methods("POST")
	mux.HandleFunc("/stream/list", listStreamersHandler).Methods("GET")
	http.ListenAndServe(":8080", mux)
	

}
