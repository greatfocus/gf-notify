package main

import (
	frame "github.com/greatfocus/gf-frame"
	"github.com/greatfocus/gf-notify/router"
	"github.com/greatfocus/gf-notify/task"
	_ "github.com/lib/pq"
)

// Entry point to the solution
func main() {
	// Load configurations
	server := frame.Create("dev.json")

	// background task
	tasks := task.Tasks{}
	tasks.Init(server.DB, server.Config)
	server.Cron.Every(1).Sunday().At("8:00").Do(tasks.RunDatabaseScripts)
	// cron.Every(1).Day().At("3:00").Do(tasks.RunDatabaseScripts) // consider making API calls to users to validate gf users e.g if they have been disabled
	server.Cron.Every(20).Second().Do(tasks.MoveStagedToQueue)
	server.Cron.Every(10).Second().Do(tasks.ReQueueProcessingEmails) // in case there is panic error some queue may be stack in processing mode
	server.Cron.Every(10).Second().Do(tasks.SendQueuedEmails)
	server.Cron.Every(1).Minute().Do(tasks.MoveOutFailedQueue)
	server.Cron.Every(1).Minute().Do(tasks.MoveOutCompleteQueue)
	server.Cron.Start()

	// start API service
	server.Start(router.Router(server.DB, server.Config))
}
