package main

import (
	frame "github.com/greatfocus/gf-frame"
	"github.com/greatfocus/gf-frame/cron"
	"github.com/greatfocus/gf-notify/router"
	"github.com/greatfocus/gf-notify/task"
	_ "github.com/lib/pq"
)

// Entry point to the solution
func main() {
	// Load configurations
	service := frame.Create("dev.json")

	tasks := task.Tasks{}
	tasks.Init(service.DB, service.Config)
	cron.Every(1).Minute().Do(tasks.SendNewEmails)
	<-cron.Start()

	// start API service
	service.Start(router.Router(service.DB))
}
