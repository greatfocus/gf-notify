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
	messageController := controllers.MessageController{}
	messageController.Init(db)

	// Initialize controller
	channelController := controllers.ChannelController{}
	channelController.Init(db)

	// Initialize routes
	mux.HandleFunc("/api/channels", middlewares.SetMiddlewareJSON(channelController.Handler))
	mux.HandleFunc("/api/messages", middlewares.SetMiddlewareJSON(messageController.Handler))
	//mux.HandleFunc("/api/dashboard", middlewares.SetMiddlewareJwt(messageController.Handler))
}
