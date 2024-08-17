package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/NicoHernandezR/Back-end-spotychafa-go/service/mp3"
	awsS3 "github.com/NicoHernandezR/Back-end-spotychafa-go/service/s3"
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

	awsS3.S3.Region = "us-east-1"
	awsS3.S3.NewSession(awsS3.S3.Region)

	router := http.NewServeMux()

	mp3Store := mp3.NewStore(s.db)
	mp3Handler := mp3.NewHandler(mp3Store, awsS3.S3)
	mp3Handler.RegisterRouter(router)

	server := http.Server{
		Addr:    s.addr,
		Handler: router,
	}

	log.Printf("server started at %s", s.addr)

	return server.ListenAndServe()

}
