package main

import (
	"encoding/json"
	"fmt"
	"go_web_server/homepage"
	"go_web_server/util"
	"log"
	"net/http"
)

type TestData struct {
	Name string
	Age  int
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./favicon_io/favicon.ico")
}

func cssHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	fileName := "./css/" + r.URL.Path
	http.ServeFile(w, r, fileName)
}

func Router(mux *http.ServeMux) {
	mux.HandleFunc("/favicon.ico", faviconHandler)
	mux.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data := TestData{Name: "Ajay", Age: 24}
		json.NewEncoder(w).Encode(data)
	})
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Test Route!\n")
	})
	// Serve CSS files
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))

	mux.HandleFunc("/api", homepage.HomePage(util.JSON, mux))
	mux.HandleFunc("/", homepage.HomePage(util.HTML, mux))
}

func main() {
	log.Println("Starting server....")
	mux := http.NewServeMux()
	Router(mux)
	log.Fatal(http.ListenAndServe(":9000", mux))
}
