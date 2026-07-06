package ebmstest

import (
	"embed"
	"io"
	"net/http"
	"net/http/httptest"
)

//go:embed *.json
var responses embed.FS

const (
	PathGetAllRoute         = "/businfo/getallroute"
	PathGetRouteByID        = "/businfo/getroutebyid/{routeID}"
	PathGetVarsByRoute      = "/businfo/getvarsbyroute/{routeID}"
	PathGetTimetableByRoute = "/businfo/gettimetablebyroute/{routeID}"
	PathGetTripsByTimetable = "/businfo/gettripsbytimetable/{routeID}/{timetableID}"
	PathGetStopsByVar       = "/businfo/getstopsbyvar/{routeID}/{variantID}"
	PathGetPathsByVar       = "/businfo/getpathsbyvar/{routeID}/{variantID}"
)

func NewServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(PathGetAllRoute, getAllRouteHandler)
	mux.HandleFunc(PathGetTimetableByRoute, getTimetableByRouteHandler)
	return httptest.NewServer(mux)
}

func getAllRouteHandler(w http.ResponseWriter, req *http.Request) {
	serve(w, "getallroute.json")
}

func getTimetableByRouteHandler(w http.ResponseWriter, req *http.Request) {
	serve(w, "gettimetablebyroute.json")
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
