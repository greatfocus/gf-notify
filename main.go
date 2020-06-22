package main

import (
	"github.com/greatfocus/gf-frame/scheduler"
	"github.com/greatfocus/gf-notify/router"
	"github.com/greatfocus/gf-notify/tasks"
	frame "github.com/greatfocus/gt-frame"
	_ "github.com/lib/pq"
)

// Entry point to the solution
func main() {
	// Load configurations
	service := frame.Create("dev.json")

	// configure scheduled jobs
	scheduler.Every(10).Minute().Do(tasks.MessageOut, service.DB)
	service.Start(router.Router(service.DB))
}
