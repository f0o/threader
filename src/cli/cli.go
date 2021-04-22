package cli

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/f0o/turbo-pancake/common"
	jobv0 "github.com/f0o/turbo-pancake/spec/job/v0"
	metadatav0 "github.com/f0o/turbo-pancake/spec/metadata/v0"
	returnv0 "github.com/f0o/turbo-pancake/spec/return/v0"
	"github.com/f0o/turbo-pancake/utils"
	"github.com/f0o/turbo-pancake/worker"
)

//spawnWorkers : Wrapper that creates Worker-Goroutines
func spawnWorkers() {
	for w := 0; w < *common.Threads; w++ {
		uuid, _ := utils.UUIDv4()
		utils.Debug("test", "%v", uuid)
		go worker.Worker(uuid[:])
	}
}

//getReturns : Main Loop to wait and process ReturnSpec from the ReturnQueue
func getReturns() {
	for ret := range returnv0.Queue {
		common.JobLength--
		if ret.Exit > 0 {
			utils.Debug(ret.ID, "Worker '%s' Exited with Code '%d'", ret.Worker, ret.Exit)
		}
		utils.Debug(ret.ID, "Returned: %v", ret.Return)
		fmt.Printf(*common.OutFormat, ret.Return)
		if common.JobLength == 0 {
			break
		}
	}
}

//createQueues : Creates Queues from Input
func createQueues() {
	common.JobQueue = make(chan common.JobSpec, len(common.Input))
	common.ReturnQueue = make(chan common.ReturnSpec, len(common.Input))
}

//createJobs : Creates JobSpecs
func createJobs() {
	for i := 0; i < len(common.Input); i++ {
		uuid, _ := utils.UUIDv4()
		select {
		case jobv0.Queue <- jobv0.Spec{
			ID:      uuid[:],
			Input:   []byte(common.Input[i]),
			Command: *common.Command,
			Filter: metadatav0.Spec{
				GOARCH: "any",
				GOOS:   runtime.GOOS,
			},
		}:
			utils.Debug(uuid[:], "Created")
			common.JobLength++
		}
	}
}

//internalExec : Execute Internal Routines to remove dependencies on external tools
func internalExec() {
	switch os.Getenv("__THREADER_INTERNAL") {
	case "base64":
		output, _ := base64.StdEncoding.DecodeString(utils.GetSTDIN())
		fmt.Printf("%s", output)
	default:
		os.Exit(128)
	}
	os.Exit(0)
}

//init : See Golang's init()
func init() {
	// Decide if we're running as "internal" Command or not
	if os.Getenv("__THREADER_INTERNAL") != "" {
		internalExec()
	} else {
		flag.Parse()
		if *common.Command == "" {
			flag.PrintDefaults()
			os.Exit(2)
		}
		dtmp, _ := strconv.Unquote(`"` + *common.Delimiter + `"`)
		common.Delimiter = &dtmp
		otmp, _ := strconv.Unquote(`"` + *common.OutFormat + `"`)
		common.OutFormat = &otmp
		common.Input = strings.Split(utils.GetSTDIN(), *common.Delimiter)
	}
}

//Main : See Golang's main()
func Main() {
	createQueues()
	spawnWorkers()
	go createJobs()
	getReturns()
	defer close(common.JobQueue)
	defer close(common.ReturnQueue)
}
