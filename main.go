package main

import (
	frame "github.com/greatfocus/gf-frame"
	"github.com/greatfocus/gf-frame/scheduler"
	"github.com/greatfocus/gf-notify/router"
	"github.com/greatfocus/gf-notify/task"
	_ "github.com/lib/pq"
)

// Entry point to the solution
func main() {
	// Load configurations
	service := frame.Create("dev.json")

	// configure scheduled jobs
	job := task.MessageOut{}
	s := scheduler.Scheduler{}
	s.Every(3).Minute().Do(job.Start, &service.DB)

	// start API service
	service.Start(router.Router(service.DB))
}
