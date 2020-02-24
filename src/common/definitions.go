package common

import (
	"flag"
	"log"
	"os"
	"runtime"
)

//JobSpec : Spec for creating a Job
type JobSpec struct {
	ID      string       `json:"id"`
	Input   []byte       `json:"input"`
	Command string       `json:"command"`
	Filter  MetadataSpec `json:"filters"`
}

//ReturnSpec : Spec for return values to a JobSpec
type ReturnSpec struct {
	ID       string       `json:"id"`
	Worker   string       `json:"worker"`
	Return   []byte       `json:"return"`
	Exit     int          `json:"exit"`
	Metadata MetadataSpec `json:"metadata"`
}

//MetadataSpec : Spec for Filtering and Worker Metadata
type MetadataSpec struct {
	GOARCH string `json:"arch"`
	GOOS   string `json:"os"`
}

//EventSpec : Naive Event Spec for publishing Events
type EventSpec struct {
	ID      string      `json:"id"`
	Type    string      `json:"type"`
	Verb    string      `json:"verb"`
	Payload interface{} `json:"payload"`
}

//JobQueue : Main messaging Queue to distribute Jobs
var JobQueue chan JobSpec

//ReturnQueue : Main messaging Queue to collect Returns
var ReturnQueue chan ReturnSpec

//EventQueue : Main messaging Queue to publish Events
var EventQueue chan EventSpec = make(chan EventSpec)

//Threads : Amount of threads to use
var Threads = flag.Int("threads", runtime.NumCPU(), "Amount of threads to use")

//Delimiter : Input delimiter
var Delimiter = flag.String("delimiter", `\n`, "Input Delimiter, as fmt.Sprintf")

//OutFormat : Output Format
var OutFormat = flag.String("out-format", `%s\n`, "Output Format, as fmt.Printf")

//Command : Command to execute
var Command = flag.String("command", "", "Command to execute")

//DebugFlag : Debug Mode
var DebugFlag = flag.Bool("d", false, "Debug Output")

//Input : STDIN
var Input []string

//JobLength : Counter of all Jobs in Queue
var JobLength int = 0

//Logger : Log Object
var Logger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds)
