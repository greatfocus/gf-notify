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

	messageBulkController := controllers.MessageBulkController{}
	messageBulkController.Init(db)

	channelController := controllers.ChannelController{}
	channelController.Init(db)

	dashboardController := controllers.DashboardController{}
	dashboardController.Init(db)

	gfuserController := controllers.GFUserController{}
	gfuserController.Init(db)

	templateController := controllers.TemplateController{}
	templateController.Init(db)

	templateMessageController := controllers.TemplateMessageController{}
	templateMessageController.Init(db)

	templateMessageBulkController := controllers.TemplateMessageBulkController{}
	templateMessageBulkController.Init(db)

	// Initialize routes
	mux.HandleFunc("/api/channel", middlewares.SetMiddlewareJSON(channelController.Handler))
	mux.HandleFunc("/api/message", middlewares.SetMiddlewareJSON(messageController.Handler))
	mux.HandleFunc("/api/message/bulk", middlewares.SetMiddlewareJSON(messageBulkController.Handler))
	mux.HandleFunc("/api/dashboard", middlewares.SetMiddlewareJSON(dashboardController.Handler))
	mux.HandleFunc("/api/gf-user", middlewares.SetMiddlewareJSON(gfuserController.Handler))
	mux.HandleFunc("/api/template", middlewares.SetMiddlewareJSON(templateController.Handler))
	mux.HandleFunc("/api/template-message", middlewares.SetMiddlewareJSON(templateMessageController.Handler))
	mux.HandleFunc("/api/template-message/bulk", middlewares.SetMiddlewareJSON(templateMessageBulkController.Handler))
}
