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

	reportController := controllers.ReportController{}
	reportController.Init(db)

	// Initialize routes
	mux.HandleFunc("/channel", middlewares.SetMiddlewareJSON(channelController.Handler))
	mux.HandleFunc("/message", middlewares.SetMiddlewareJSON(messageController.Handler))
	mux.HandleFunc("/message/bulk", middlewares.SetMiddlewareJSON(messageBulkController.Handler))
	mux.HandleFunc("/dashboard", middlewares.SetMiddlewareJSON(dashboardController.Handler))
	mux.HandleFunc("/gf-user", middlewares.SetMiddlewareJSON(gfuserController.Handler))
	mux.HandleFunc("/template", middlewares.SetMiddlewareJSON(templateController.Handler))
	mux.HandleFunc("/template-message", middlewares.SetMiddlewareJSON(templateMessageController.Handler))
	mux.HandleFunc("/template-message/bulk", middlewares.SetMiddlewareJSON(templateMessageBulkController.Handler))
	mux.HandleFunc("/report", middlewares.SetMiddlewareJSON(reportController.Handler))
}
