package main

import (
	"github.com/greatfocus/gf-notify/router"
	"github.com/greatfocus/gf-notify/task"
	_ "github.com/greatfocus/gf-pq"
	frame "github.com/greatfocus/gf-sframe"
)

// Entry point to the solution
func main() {
	frame := frame.NewFrame("gf-notify", "notify")
	mux := router.LoadRouter(frame.Server)

	// background task
	tasks := task.Tasks{}
	tasks.Init(frame.Server)
	frame.Server.Cron.MustAddJob("* * * * *", tasks.MoveStagedToQueue)       // 20 Seconds
	frame.Server.Cron.MustAddJob("* * * * *", tasks.ReQueueProcessingEmails) // 10 Seconds
	frame.Server.Cron.MustAddJob("* * * * *", tasks.SendQueuedEmails)        // 10 Seconds
	frame.Server.Cron.MustAddJob("* * * * *", tasks.MoveOutFailedQueue)      // 1 Minute
	frame.Server.Cron.MustAddJob("* * * * *", tasks.MoveOutCompleteQueue)    // 1 Minute
	frame.Start(mux)
}
