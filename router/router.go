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

	gfuserController := controllers.GFUserController{}
	gfuserController.Init(s.DB)

	templateController := controllers.TemplateController{}
	templateController.Init(s.DB)

	templateMessageController := controllers.TemplateMessageController{}
	templateMessageController.Init(s.DB)

	templateMessageBulkController := controllers.TemplateMessageBulkController{}
	templateMessageBulkController.Init(s.DB)

	reportController := controllers.ReportController{}
	reportController.Init(s.DB)

	// Initialize routes
	mux.HandleFunc("/notify/channel", middlewares.SetMiddlewareJSON(channelController.Handler, s))
	mux.HandleFunc("/notify/message", middlewares.SetMiddlewareJSON(messageController.Handler, s))
	mux.HandleFunc("/notify/message/bulk", middlewares.SetMiddlewareJSON(messageBulkController.Handler, s))
	mux.HandleFunc("/notify/dashboard", middlewares.SetMiddlewareJSON(dashboardController.Handler, s))
	mux.HandleFunc("/notify/gf-user", middlewares.SetMiddlewareJSON(gfuserController.Handler, s))
	mux.HandleFunc("/notify/template", middlewares.SetMiddlewareJSON(templateController.Handler, s))
	mux.HandleFunc("/notify/template-message", middlewares.SetMiddlewareJSON(templateMessageController.Handler, s))
	mux.HandleFunc("/notify/template-message/bulk", middlewares.SetMiddlewareJSON(templateMessageBulkController.Handler, s))
	mux.HandleFunc("/notify/report", middlewares.SetMiddlewareJSON(reportController.Handler, s))
}
