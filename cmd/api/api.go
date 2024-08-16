package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/NicoHernandezR/Back-end-spotychafa-go/service/mp3"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()

	mp3Store := mp3.NewStore(s.db)
	mp3Handler := mp3.NewHandler(mp3Store)
	mp3Handler.RegisterRouter(router)

	server := http.Server{
		Addr:    s.addr,
		Handler: router,
	}

	log.Printf("server started at %s", s.addr)

	return server.ListenAndServe()

}
