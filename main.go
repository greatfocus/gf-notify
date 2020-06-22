package main

import (
	frame "github.com/greatfocus/gf-frame"
	"github.com/greatfocus/gf-frame/scheduler"
	"github.com/greatfocus/gf-notify/router"
	"github.com/greatfocus/gf-notify/tasks"
	_ "github.com/lib/pq"
)

// Entry point to the solution
func main() {
	// Load configurations
	service := frame.Create("dev.json")

	// configure scheduled jobs
	s := scheduler.Scheduler{}
	s.Every(10).Minute().Do(tasks.MessageOut, service.DB)

	// start API service
	service.Start(router.Router(service.DB))
}
