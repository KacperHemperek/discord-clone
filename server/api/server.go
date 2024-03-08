package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kacperhemperek/discord-go/store"
	"log"
	"net/http"
)

type Server struct {
	port int
}

func NewApiServer(port int) *Server {
	return &Server{port: port}
}

func (s *Server) Start() {
	router := mux.NewRouter()
	db := store.NewDB()

	fmt.Println(db)

	router.HandleFunc("/api/users", HandlerFunc(func(w http.ResponseWriter, r *http.Request, c *Context) error {
		return WriteJson(w, http.StatusOK, c)
	})).Methods("GET")

	portStr := fmt.Sprintf(":%d", s.port)

	fmt.Printf("Server is running on port %d\n", s.port)
	log.Fatal(http.ListenAndServe(portStr, router))
}
