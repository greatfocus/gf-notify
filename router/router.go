package router

import (
	"log"
	"net/http"

	"github.com/greatfocus/gf-notify/handler"
	"github.com/greatfocus/gf-sframe/server"
)

// Router is exported and used in main.go
func LoadRouter(s *server.Meta) *http.ServeMux {
	mux := http.NewServeMux()
	loadHandlers(mux, s)
	log.Println("Created routes with handler")
	return mux
}

// notifyRoute created all routes and handlers relating to controller
func loadHandlers(mux *http.ServeMux, s *server.Meta) {
	// Initialize controller
	messageHandler := handler.Message{}
	messageHandler.Init(s)
	mux.Handle("/notify/message", server.Use(messageHandler,
		server.SetHeaders(),
		server.CheckLimitsRates(),
		server.CheckCors(s),
		server.CheckAllowedIPRange(s),
		server.CheckProcessTimeout(s),
		server.WithoutAuth()))

	dashboardHandler := handler.Dashboard{}
	dashboardHandler.Init(s)
	mux.Handle("/notify/dashboard", server.Use(dashboardHandler,
		server.SetHeaders(),
		server.CheckLimitsRates(),
		server.CheckCors(s),
		server.CheckAllowedIPRange(s),
		server.CheckProcessTimeout(s),
		server.WithoutAuth()))

	templateHandler := handler.Template{}
	templateHandler.Init(s)
	mux.Handle("/notify/template", server.Use(templateHandler,
		server.SetHeaders(),
		server.CheckLimitsRates(),
		server.CheckCors(s),
		server.CheckAllowedIPRange(s),
		server.CheckProcessTimeout(s),
		server.WithoutAuth()))

	templateMessageHandler := handler.TemplateMessage{}
	templateMessageHandler.Init(s)
	mux.Handle("/notify/template/message", server.Use(templateMessageHandler,
		server.SetHeaders(),
		server.CheckLimitsRates(),
		server.CheckCors(s),
		server.CheckAllowedIPRange(s),
		server.CheckProcessTimeout(s),
		server.WithoutAuth()))

	reportHandler := handler.Report{}
	reportHandler.Init(s)
	mux.Handle("/notify/report", server.Use(reportHandler,
		server.SetHeaders(),
		server.CheckLimitsRates(),
		server.CheckCors(s),
		server.CheckAllowedIPRange(s),
		server.CheckProcessTimeout(s),
		server.WithoutAuth()))

	// Initialize routes
	channelHandler := handler.Channel{}
	channelHandler.Init(s)
	mux.Handle("/notify/channel", server.Use(channelHandler,
		server.SetHeaders(),
		server.CheckLimitsRates(),
		server.CheckCors(s),
		server.CheckAllowedIPRange(s),
		server.CheckProcessTimeout(s),
		server.WithoutAuth()))
}
