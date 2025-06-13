package svc

import (
	"fmt"
	"time"

	"looklook/app/mqueue/cmd/scheduler/internal/config"

	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
)

// create scheduler
func newScheduler(c config.Config) *asynq.Scheduler {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		logx.Errorf("LoadLocation err %+v", err)
		return nil
	}
	return asynq.NewScheduler(
		asynq.RedisClientOpt{
			Addr:     c.Redis.Host,
			Password: c.Redis.Pass,
		}, &asynq.SchedulerOpts{
			Location: location,
			EnqueueErrorHandler: func(task *asynq.Task, opts []asynq.Option, err error) {
				fmt.Printf("Scheduler EnqueueErrorHandler <<<<<<<===>>>>> err : %+v , task : %+v", err, task)
			},
		})
}
