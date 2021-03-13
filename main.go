package main

import (
	"os"

	frame "github.com/greatfocus/gf-frame"
	"github.com/greatfocus/gf-notify/router"
	"github.com/greatfocus/gf-notify/task"
	_ "github.com/greatfocus/pq"
)

// Entry point to the solution
func main() {
	// Get arguments
	var env string
	env = os.Args[1]
	if env == "" {
		panic("Pass the environment")
	}

	// Load params
	frame, server := frame.NewFrame(env)

	// background task
	tasks := task.Tasks{}
	tasks.Init(server.DB, server.Cache, server.Config)
	server.Cron.Every(20).Second().Do(tasks.MoveStagedToQueue)
	server.Cron.Every(10).Second().Do(tasks.ReQueueProcessingEmails) // in case there is panic error some queue may be stack in processing mode
	server.Cron.Every(10).Second().Do(tasks.SendQueuedEmails)
	server.Cron.Every(1).Minute().Do(tasks.MoveOutFailedQueue)
	server.Cron.Every(1).Minute().Do(tasks.MoveOutCompleteQueue)

	server.Mux = router.LoadRouter(server)
	frame.Start(server)
}
