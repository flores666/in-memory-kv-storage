package main

import (
	"fmt"
	"in-memory-kv-storage/internal/handlers/storage"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	router := chi.NewRouter()
	fmt.Println("Router Created")

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	storage.NewStorageApi(storage.NewStorageService()).MapRoutes(router)
	fmt.Println("Api Routes Registered")

	fmt.Println("Server started")
	err := http.ListenAndServe(":5001", router)
	if err != nil {
		fmt.Printf("Server error - %s\n", err.Error())
	}

	fmt.Printf("Server error - %s\n", err.Error())
}
