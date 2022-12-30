package main

import (
	"encoding/json"
	"fmt"
	// "io/ioutil"
	// "log"
	"net/http"
	// s "strings"

	// "html/template"
	"strconv"

	"github.com/gorilla/mux"
	"test.com/m/internal/database"
)

type Upsert struct {
	Streamer    string `json:"streamer"`
	Is_Active   bool   `json:"is_active"`
	Status_Code int    `json:"status_code"`
}

func upsertResponse(streamer string, is_active bool, status_code int) Upsert {
	var response Upsert
	response.Streamer = streamer
	response.Is_Active = is_active
	response.Status_Code = status_code

	return response
}

func upsertStreamerHandler(w http.ResponseWriter, r *http.Request, c database.InsertStreamerAPI) {
	w.Header().Set("Content-Type", "application/json")
	streamer := mux.Vars(r)["stream"]
	is_active, _ := strconv.ParseBool(mux.Vars(r)["disable"])
	err := c(streamer, is_active)
	if err != nil {
		response := upsertResponse(streamer, is_active, 500)
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(400)
		return
	} else {
		response := upsertResponse(streamer, is_active, 200)
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(300)
	}

}

func listStreamersHandler(w http.ResponseWriter, r *http.Request) {
	x, err := database.GetStreamerData()
	if err != nil {
		panic(err)
	}
	fmt.Print(x)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(x)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/index.html")

}

func main() {
	// handle func is convienance method on http. registers function to a path
	// on default serve mux.
	mux := mux.NewRouter().StrictSlash(true)
	mux.HandleFunc("/about", homeHandler)
	mux.HandleFunc("/stream/upsert", func(w http.ResponseWriter, r *http.Request) {
		upsertStreamerHandler(w, r, database.InsertStreamer)
	}).Queries("stream", "{stream}", "disable", "{disable}").Methods("POST")
	mux.HandleFunc("/stream/list", listStreamersHandler).Methods("GET")
	http.ListenAndServe(":1334", mux)

}
