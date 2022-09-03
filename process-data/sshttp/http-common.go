package sshttp

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
)

// allowList is a map of allowed endpoints
var allowList = map[string]func(http.ResponseWriter, *http.Request){
	"prometheus":  prometheusHandler,
	"metrics":     metricsHandler,
	"health":      healthHandler,
	"version":     versionHandler,
	"debug/pprof": pprofHandler,
}

// startServer starts the server
func StartServer(serviceName string, routes *mux.Router, port string, adminRoutes *mux.Router, adminPort string) {
	var wg sync.WaitGroup
	wg.Add(2)
	serverStartMessage(serviceName, port, routes)
	serverStartMessage("admin", adminPort, adminRoutes)
	go func() {
		defer wg.Done()
		http.ListenAndServe(fmt.Sprintf(":%s", port), routes)
	}()
	go func() {
		defer wg.Done()
		http.ListenAndServe(fmt.Sprintf(":%s", adminPort), adminRoutes)
	}()
	wg.Wait()
	fmt.Printf("server %s stopped\n", serviceName)

}

// prettyPrint prints the routes
func serverStartMessage(serviceName, port string, r *mux.Router) {
	var routes []string
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, _ := route.GetPathTemplate()
		routes = append(routes, t)
		return nil
	})
	fmt.Printf("Starting %s server on port %s with routes: %v\n", serviceName, port, routes)
}

// prometheusHandler handles prometheus requests
func prometheusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("prometheus"))
}

// metricsHandler handles metrics requests
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("metrics"))
}

// healthHandler handles health requests
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("health"))
}

// versionHandler handles version requests
func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("version"))
}

// pprofHandler handles pprof requests
func pprofHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pprof"))
}

// get allow listed gorrila mux endpoints
func AdminEndpoints(requiredList []string) (*mux.Router, error) {
	// check required in allow list and create router.
	r := mux.NewRouter()
	for _, required := range requiredList {
		if _, ok := allowList[required]; ok {
			r.HandleFunc("/"+required, allowList[required])
		} else {
			return nil, fmt.Errorf("endpoint %s not in allow list", required)
		}
	}
	return r, nil
}

func Ports() (string, string) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	adminPort := os.Getenv("ADMIN_PORT")
	if adminPort == "" {
		adminPort = "8081"
	}
	return port, adminPort
}
