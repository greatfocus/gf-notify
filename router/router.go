package router

import (
	"log"
	"net/http"

	"github.com/greatfocus/gf-frame/server"

	"github.com/greatfocus/gf-frame/middlewares"
	"github.com/greatfocus/gf-notify/controllers"
)

// Router is exported and used in main.go
func Router(s *server.Server) *http.ServeMux {
	// create new router
	mux := http.NewServeMux()

	// users
	usersRoute(mux, s)

	log.Println("Created routes with controllers")
	return mux
}

// usersRoute created all routes and handlers relating to user controller
func usersRoute(mux *http.ServeMux, s *server.Server) {
	// Initialize controller
	messageController := controllers.MessageController{}
	messageController.Init(s.DB)

	messageBulkController := controllers.MessageBulkController{}
	messageBulkController.Init(s.DB)

	channelController := controllers.ChannelController{}
	channelController.Init(s.DB)

	dashboardController := controllers.DashboardController{}
	dashboardController.Init(s.DB)

	templateController := controllers.TemplateController{}
	templateController.Init(s.DB)

	templateMessageController := controllers.TemplateMessageController{}
	templateMessageController.Init(s.DB)

	templateMessageBulkController := controllers.TemplateMessageBulkController{}
	templateMessageBulkController.Init(s.DB)

	reportController := controllers.ReportController{}
	reportController.Init(s.DB)

	// Initialize routes
	mux.HandleFunc("/notify/channel", middlewares.SetMiddlewareClient(channelController.Handler, s))
	mux.HandleFunc("/notify/message", middlewares.SetMiddlewareClient(messageController.Handler, s))
	mux.HandleFunc("/notify/message/bulk", middlewares.SetMiddlewareClient(messageBulkController.Handler, s))
	mux.HandleFunc("/notify/dashboard", middlewares.SetMiddlewareClient(dashboardController.Handler, s))
	mux.HandleFunc("/notify/template", middlewares.SetMiddlewareClient(templateController.Handler, s))
	mux.HandleFunc("/notify/template-message", middlewares.SetMiddlewareClient(templateMessageController.Handler, s))
	mux.HandleFunc("/notify/template-message/bulk", middlewares.SetMiddlewareClient(templateMessageBulkController.Handler, s))
	mux.HandleFunc("/notify/report", middlewares.SetMiddlewareClient(reportController.Handler, s))
}
