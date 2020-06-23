package router

import (
	"log"
	"net/http"

	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-frame/middlewares"
	"github.com/greatfocus/gf-notify/controllers"
)

// Router is exported and used in main.go
func Router(db *database.DB) *http.ServeMux {
	// create new router
	mux := http.NewServeMux()

	// users
	usersRoute(mux, db)

	log.Println("Created routes with controllers")
	return mux
}

// usersRoute created all routes and handlers relating to user controller
func usersRoute(mux *http.ServeMux, db *database.DB) {
	// Initialize controller
	notifyController := controllers.NotifyController{}
	notifyController.Init(db)

	// Initialize routes
	mux.HandleFunc("/api/messages", middlewares.SetMiddlewareJwt(notifyController.Handler))
	mux.HandleFunc("/api/dashboard", middlewares.SetMiddlewareJwt(notifyController.Handler))
}
