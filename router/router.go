package router

import (
	"log"
	"net/http"

	"github.com/greatfocus/gf-frame/config"

	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-frame/middlewares"
	"github.com/greatfocus/gf-notify/controllers"
)

// Router is exported and used in main.go
func Router(db *database.DB, config *config.Config) *http.ServeMux {
	// create new router
	mux := http.NewServeMux()

	// users
	usersRoute(mux, db, config)

	log.Println("Created routes with controllers")
	return mux
}

// usersRoute created all routes and handlers relating to user controller
func usersRoute(mux *http.ServeMux, db *database.DB, config *config.Config) {
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
	mux.HandleFunc("/notify/channel", middlewares.SetMiddlewareJSON(channelController.Handler, config.Server.AllowedOrigins))
	mux.HandleFunc("/notify/message", middlewares.SetMiddlewareJSON(messageController.Handler, config.Server.AllowedOrigins))
	mux.HandleFunc("/notify/message/bulk", middlewares.SetMiddlewareJSON(messageBulkController.Handler, config.Server.AllowedOrigins))
	mux.HandleFunc("/notify/dashboard", middlewares.SetMiddlewareJSON(dashboardController.Handler, config.Server.AllowedOrigins))
	mux.HandleFunc("/notify/gf-user", middlewares.SetMiddlewareJSON(gfuserController.Handler, config.Server.AllowedOrigins))
	mux.HandleFunc("/notify/template", middlewares.SetMiddlewareJSON(templateController.Handler, config.Server.AllowedOrigins))
	mux.HandleFunc("/notify/template-message", middlewares.SetMiddlewareJSON(templateMessageController.Handler, config.Server.AllowedOrigins))
	mux.HandleFunc("/notify/template-message/bulk", middlewares.SetMiddlewareJSON(templateMessageBulkController.Handler, config.Server.AllowedOrigins))
	mux.HandleFunc("/notify/report", middlewares.SetMiddlewareJSON(reportController.Handler, config.Server.AllowedOrigins))
}
