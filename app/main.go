package main

import (
	"carro-ideal/app/db"
	"carro-ideal/app/internal/admin"
	"carro-ideal/app/internal/api"
	"carro-ideal/app/internal/web"
	"carro-ideal/app/repository"
	"carro-ideal/app/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	db.Connect()

	userRepo := repository.NewUserRepository(db.GetDB())
	userService := service.NewUserService(userRepo)

	webHandler := web.NewHandler(userService)
	apiHandler := api.NewHandler(userService)
	adminHandler := admin.NewHandler(userService)

	r := mux.NewRouter()
	web.RegisterRoutes(r, webHandler)
	api.RegisterRoutes(r, apiHandler)
	admin.RegisterRoutes(r, adminHandler)

	log.Println("Servidor iniciado em http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Erro ao iniciar servidor:", err)
	}
}
