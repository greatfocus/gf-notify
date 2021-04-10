package main

import (
	"log"
	"os"

	frame "github.com/greatfocus/gf-frame"
	"github.com/greatfocus/gf-notify/router"
	"github.com/greatfocus/gf-notify/task"
	_ "github.com/greatfocus/pq"
)

// Entry point to the solution
func main() {
	// Get arguments
	var env = os.Args[1]
	if env == "" {
		panic("pass the environment")
	}

	// Load params
	frame, server := frame.NewFrame(env)

	// background task
	tasks := task.Tasks{}
	tasks.Init(server.DB, server.Cache, server.Config)
	err := server.Cron.Every(20).Second().Do(tasks.MoveStagedToQueue)
	if err != nil {
		log.Fatalf("Cron Job failed: MoveStagedToQueue: %v", err)
	}
	err = server.Cron.Every(10).Second().Do(tasks.ReQueueProcessingEmails) // in case there is panic error some queue may be stack in processing mode
	if err != nil {
		log.Fatalf("Cron Job failed: ReQueueProcessingEmails: %v", err)
	}
	err = server.Cron.Every(10).Second().Do(tasks.SendQueuedEmails)
	if err != nil {
		log.Fatalf("Cron Job failed: SendQueuedEmails: %v", err)
	}
	err = server.Cron.Every(1).Minute().Do(tasks.MoveOutFailedQueue)
	if err != nil {
		log.Fatalf("Cron Job failed: MoveOutFailedQueue: %v", err)
	}
	err = server.Cron.Every(1).Minute().Do(tasks.MoveOutCompleteQueue)
	if err != nil {
		log.Fatalf("Cron Job failed: MoveOutCompleteQueue: %v", err)
	}

	server.Mux = router.LoadRouter(server)
	frame.Start(server)
}
