package worker

import (
	"encoding/base64"
	"os"
	"os/exec"
	"runtime"

	eventv0 "github.com/f0o/turbo-pancake/spec/event/v0"
	jobv0 "github.com/f0o/turbo-pancake/spec/job/v0"
	metadatav0 "github.com/f0o/turbo-pancake/spec/metadata/v0"
	returnv0 "github.com/f0o/turbo-pancake/spec/return/v0"
	"github.com/f0o/turbo-pancake/utils"
)

//Worker : Actual Worker function that executes JobSpec
func Worker(worker string) {
	(&eventv0.Spec{
		Type:    eventv0.TypeWorker,
		Verb:    eventv0.VerbCreate,
		Payload: worker,
	}).Publish()
	for job := range jobv0.Queue {
		(&eventv0.Spec{
			Type: eventv0.TypeJob,
			Verb: eventv0.VerbUpdate,
			Payload: ([2]string{
				job.ID,
				worker,
			}),
		}).Publish()
		if Filter(&job) {
			cmd := exec.Command(os.Getenv("SHELL"), "-c", `INPUT=$(echo "$BASEINPUT" | __THREADER_INTERNAL=base64 `+os.Args[0]+`) eval $COMMAND`)
			cmd.Env = append(os.Environ(),
				"COMMAND="+job.Command,
				"BASEINPUT="+base64.StdEncoding.EncodeToString(job.Input),
				"JOBID="+job.ID,
				"WORKERID="+worker,
			)
			exit := 0
			output, err := cmd.CombinedOutput()
			if err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					exit = exitError.ExitCode()
				}
			}
			(&eventv0.Spec{
				Type: eventv0.TypeReturn,
				Verb: eventv0.VerbCreate,
				Payload: ([2]string{
					job.ID,
					worker,
				}),
			}).Publish()
			(&returnv0.Spec{
				Worker: worker,
				Return: output,
				Exit:   exit,
				Metadata: metadatav0.Spec{
					GOARCH: runtime.GOARCH,
					GOOS:   runtime.GOOS,
				},
			}).Publish()
		} else {
			(&eventv0.Spec{
				Type:    eventv0.TypeJob,
				Verb:    eventv0.VerbUpdate,
				Payload: job.ID,
			}).Publish()
			jobv0.Queue <- job
		}
	}
}

//Filter : Filters Jobs on Worker Metadata
func Filter(job *jobv0.Spec) bool {
	ret := false
	switch {
	case job.Filter.GOARCH == runtime.GOARCH:
		ret = true
		fallthrough
	case job.Filter.GOOS == runtime.GOOS:
		ret = true
		fallthrough
	case job.Filter.GOARCH == "any":
		ret = true
		fallthrough
	case job.Filter.GOOS == "any":
		ret = true
		//fallthrough
	}
	utils.Debug(job.ID, "Filter Validated: %+v -> %v", job.Filter, ret)
	return ret
}
