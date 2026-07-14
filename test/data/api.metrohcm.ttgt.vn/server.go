package ttgttest

import (
	"embed"
	"io"
	"net/http"
	"net/http/httptest"
)

//go:embed *.json
var responses embed.FS

const (
	PathTransitVehicles = "/transit/vehicles"
)

func NewServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(PathTransitVehicles, transitVehiclesHandler)
	return httptest.NewServer(mux)
}

func transitVehiclesHandler(w http.ResponseWriter, req *http.Request) {
	routeId := req.URL.Query().Get("routeId")
	varId := req.URL.Query().Get("varId")

	if routeId == "384" {
		switch varId {
		case "1":
			serve(w, "transitvehicles-1.json")
			return
		case "2":
			serve(w, "transitvehicles-2.json")
			return
		}
	}

	w.WriteHeader(http.StatusInternalServerError)
}

func serve(w http.ResponseWriter, path string) {
	file, err := responses.Open(path)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(w, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
