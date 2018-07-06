package portal_http

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/high-value-team/workshop-kubernetes-setup/workshop-ci/interior_interactions"
	"github.com/high-value-team/workshop-kubernetes-setup/workshop-ci/interior_models"
	"github.com/rs/cors"
)

type HTTPPortal struct {
	router *chi.Mux
}

func NewHTTPPortal(interactions *interior_interactions.Interactions, versionNumber string) *HTTPPortal {
	router := chi.NewRouter()
	router.Post("/", NewTriggerPipelineHandler(interactions))
	return &HTTPPortal{router: router}
}

func (portal *HTTPPortal) Run(port int) {
	address := fmt.Sprintf(":%d", port)
	handler := cors.AllowAll().Handler(portal)

	log.Println("Starting Webserver, please go to: ", address)
	err := http.ListenAndServe(address, handler)
	if err != nil {
		log.Println(err)
	}
}

func (portal *HTTPPortal) ServeHTTP(writer http.ResponseWriter, reader *http.Request) {
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	defer handleException(writer)
	portal.router.ServeHTTP(writer, reader)
}

func handleException(writer http.ResponseWriter) {
	r := recover()
	if r != nil {
		switch ex := r.(type) {
		case interior_models.SadException:
			log.Println(ex.Message())
			http.Error(writer, ex.Message(), 404)
		case interior_models.SuprisingException:
			log.Println(ex.Message())
			http.Error(writer, ex.Message(), 500)
		default:
			if err, ok := r.(error); !ok {
				log.Println(err)
				http.Error(writer, err.Error(), 500)
			} else {
				log.Println(err)
				http.Error(writer, fmt.Sprintf("%s", r), 500)
			}
		}
	}
}
