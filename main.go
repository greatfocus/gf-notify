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

	// background task
	tasks := task.Tasks{}
	tasks.Init(service.DB, service.Config)
	cron.Every(1).Sunday().At("8:00").Do(tasks.RunDatabaseScripts)
	cron.Every(30).Second().Do(tasks.MoveStagedToQueue)
	// cron.Every(10).Second().Do(tasks.SendNewEmails)
	cron.Start()

	// start API service
	service.Start(router.Router(service.DB))
}
