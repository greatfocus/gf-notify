package router

import (
	"log"
	"net/http"

	"github.com/greatfocus/gf-frame/server"
	"github.com/greatfocus/gf-notify/controllers"
)

// LoadRouter is exported and used in main.go
func LoadRouter(s *server.MetaData) *http.ServeMux {
	mux := http.NewServeMux()
	loadHandlers(mux, s)
	log.Println("Created routes with controllers")
	return mux
}

// notifyRoute created all routes and handlers relating to controller
func loadHandlers(mux *http.ServeMux, s *server.MetaData) {
	// Initialize controller
	messageController := controllers.MessageController{}
	messageController.Init(s.DB, s.Cache)

	messageBulkController := controllers.MessageBulkController{}
	messageBulkController.Init(s.DB, s.Cache)

	channelController := controllers.ChannelController{}
	channelController.Init(s.DB, s.Cache)

	dashboardController := controllers.DashboardController{}
	dashboardController.Init(s.DB, s.Cache)

	templateController := controllers.TemplateController{}
	templateController.Init(s.DB, s.Cache)

	templateMessageController := controllers.TemplateMessageController{}
	templateMessageController.Init(s.DB, s.Cache)

	templateMessageBulkController := controllers.TemplateMessageBulkController{}
	templateMessageBulkController.Init(s.DB, s.Cache)

	reportController := controllers.ReportController{}
	reportController.Init(s.DB, s.Cache)

	// Initialize routes
	mux.HandleFunc("/notify/channel", server.SetMiddlewareClient(channelController.Handler, s))
	mux.HandleFunc("/notify/message", server.SetMiddlewareClient(messageController.Handler, s))
	mux.HandleFunc("/notify/message/bulk", server.SetMiddlewareClient(messageBulkController.Handler, s))
	mux.HandleFunc("/notify/dashboard", server.SetMiddlewareClient(dashboardController.Handler, s))
	mux.HandleFunc("/notify/template", server.SetMiddlewareClient(templateController.Handler, s))
	mux.HandleFunc("/notify/template-message", server.SetMiddlewareClient(templateMessageController.Handler, s))
	mux.HandleFunc("/notify/template-message/bulk", server.SetMiddlewareClient(templateMessageBulkController.Handler, s))
	mux.HandleFunc("/notify/report", server.SetMiddlewareClient(reportController.Handler, s))
}
