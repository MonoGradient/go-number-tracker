package main

import (
	"Go-Tracker/api"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type App struct {
	Router *mux.Router
}

var (
	apiPrefix = "/api/v1"
	port = ":8080"
)
func main() {
	app := App{}
	app.InitializeRoutes()
	app.InitializeMiddleware()
	app.Run(port)
}

func AddJsonContentTypeHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.RequestURI, "api") {
			log.Debugln("Detected API endpoint, adding global json content-type.")
			w.Header().Set("Content-Type", "application/json")
		}
		next.ServeHTTP(w, r)
	})
}

func (app *App) InitializeRoutes() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/storage", api.IdempotentIncrementStorageApi).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/storage/{key}", api.CheckStorageApi).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/storage/{key}", api.IncrementStorageWithKeyApi).Methods(http.MethodPut)
	router.HandleFunc("/api/v1/storage/{key}", api.DeleteStorageApi).Methods(http.MethodDelete)
	router.HandleFunc("/api/v1/storage/{key}/increment", api.IncrementStorageWithKeyApi).Methods(http.MethodPut)
	router.HandleFunc("/api/v1/storage/{key}/decrement", api.DecrementStorageApi).Methods(http.MethodPut)

	app.Router = router
}

func (app *App) InitializeMiddleware() {
	app.Router.Use(mux.CORSMethodMiddleware(app.Router))
	app.Router.Use(LoggingMiddleware)
	app.Router.Use(AddJsonContentTypeHeaderMiddleware)
}
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("API Request to [%v]\n", r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
func (app *App) Run(port string) {
	if err := http.ListenAndServe(port, app.Router); err != nil {
		panic(err)
	}
}

