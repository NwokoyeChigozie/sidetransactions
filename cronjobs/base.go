package cronjobs

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/vesicash/transactions-ms/external/request"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/utility"
)

var (
	cronJobs = map[string]CronJobObject{
		"transaction-inspection-period": {CronJob: HandleTransactionInspectionPeriod, Interval: time.Hour * 24},
		"transaction-auto-mark":         {CronJob: HandleTransactionAutoMark, Interval: time.Hour * 24},
		"transactions-auto-close":       {CronJob: HandleTransactionAutoClose, Interval: time.Hour * 24},
		"transaction-close":             {CronJob: HandleTransactionClose, Interval: time.Minute * 10},
		"update-status":                 {CronJob: HandleUpdateStatus, Interval: time.Minute * 10},
	}
)

type CronJob func(extReq request.ExternalRequest, db postgresql.Databases)

type CronJobObject struct {
	CronJob  CronJob
	Interval time.Duration
}

func Scheduler(extReq request.ExternalRequest, db postgresql.Databases, mutex *sync.Mutex, cronJob CronJob, interval time.Duration) {
	for {
		mutex.Lock()
		cronJob(extReq, db)
		mutex.Unlock()
		time.Sleep(interval)
	}
}

func SetupCronJobs(extReq request.ExternalRequest, db postgresql.Databases, selectedJobs []string) {
	mutex := &sync.Mutex{}
	for _, v := range selectedJobs {
		jobName := strings.ToLower(v)
		cronJob, ok := cronJobs[jobName]

		if ok {
			utility.LogAndPrint(extReq.Logger, fmt.Sprintf("starting cronjob: %s", jobName))
			go Scheduler(extReq, db, mutex, cronJob.CronJob, cronJob.Interval)
		} else {
			utility.LogAndPrint(extReq.Logger, fmt.Sprintf("Cronjob not found: %s", jobName))
		}

	}

	select {}
}
