package main

import (
	"log"
	"net/http"

	"github.com/example/iching-app/internal/httpapi"
	"github.com/example/iching-app/internal/service"
	mem "github.com/example/iching-app/internal/storage/memory"
)

func main() {
	repo := mem.NewReadingRepository()
	svc := service.NewReadingService(repo)
	h := httpapi.NewHandler(svc)

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", h.Routes()); err != nil {
		log.Fatal(err)
	}
}
