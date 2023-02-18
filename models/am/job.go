package am

import (
	"github.com/arjendevos/gosql/models/dm"
)

type Job struct {
	*dm.Job

	InstagramAccounts dm.InstagramAccountSlice `boil:"InstagramAccounts" json:"InstagramAccounts" toml:"InstagramAccounts" yaml:"InstagramAccounts"`
}

type JobSlice []*Job

func SqlBoilerJobsToApiJobs(a dm.JobSlice) JobSlice {
	var s JobSlice
	for _, d := range a {
		s = append(s, SqlBoilerJobToApiJob(d))
	}
	return s
}

func SqlBoilerJobToApiJob(a *dm.Job) *Job {
	return &Job{
		Job: a,

		InstagramAccounts: a.R.InstagramAccounts,
	}
}
