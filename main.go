package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type Args struct {
	Delimiter    string
	Input        string
	ThreadAmount int
	Run          string
	Runs         int
}

var args Args
var logger = log.New(os.Stdout, "", 0)

func init() {
	// presets
}

func main() {
	// first we gonne parse the args
	parseArgs()

	// one simple check
	cpuAmount := runtime.NumCPU()
	if args.ThreadAmount > cpuAmount {
		fmt.Println("You set '", args.ThreadAmount, "' threads but only have '", cpuAmount, "' cores. It`s your choice...")
	}

	// check what type of input we get by following priority
	// 1. a input-command string to be executed, the return
	//    used as index
	// 2. arg -input as string input
	var input map[int]string
	// is there a command given that provides input?
	if args.Runs != -1 {
		input = make(map[int]string)
	} else {
		args.Input = getPipedInput()
		input = splitInput()
	}

	// decide if its smart minions or stupid minions
	if 0 == len(input) {
		runHeadlessMinions()
	} else {
		runSmartMinions(input)
	}

}

func getPipedInput() string {
	reader := bufio.NewReader(os.Stdin)
	var output []rune

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}

	return string(output)
}

func runSmartMinions(input map[int]string) {
	// create an intercom map
	intercom := make(map[int]chan string)

	// some number juggling , probably smarter ways
	// out there but for now its fine
	specialCase := -1
	base := float64(len(input)) / float64(args.ThreadAmount)
	checkBase := math.Floor(float64(base))
	if base != checkBase {
		checkAmount := int(checkBase) * args.ThreadAmount
		specialCase = len(input) - checkAmount
	}

	// start/stop id
	start := 0
	stop := 0

	// run the minions
	for i := 0; i < args.ThreadAmount; i++ {
		// create the channel
		intercom[i] = make(chan string)

		// check how many entries we got to push
		// and calc the stop
		subAmount := int(checkBase)
		if specialCase != -1 && i == 0 {
			subAmount = int(checkBase) + specialCase
		}
		stop = stop + subAmount

		// dont judge me i want to finish this
		subInput := make([]string, subAmount)

		for x := start; x < stop; x++ {
			if val, ok := input[x]; ok {
				subInput[x-start] = val
			} else {
				logger.Println("Could not find ", x, " in ", input)
			}

		}

		// set the next start to the current stop
		start = stop
		// start the minion
		go createSmartMinion(i, subInput, intercom[i])
	}

	// check for the minions to be done
	done := 0
	for {

		// check all thread channels for return
		for i := 0; i < args.ThreadAmount; i++ {
			ret := <-intercom[i]
			if "done" == ret {
				done++
			}
		}

		// all threads are done
		if args.ThreadAmount == done {
			break
		}
	}
	//fmt.Println("All work is done. Exiting")

}

func createSmartMinion(id int, input []string, intercom chan string) {
	//fmt.Println("Minion nr", id, "spawned. ")
	for iid, singleIn := range input {
		// provide some vars for the command
		// run the command
		execString := prepareExecString(singleIn, id, iid)
		ret, _ := custExec(execString)
		logger.Println(ret)
	}
	//fmt.Println("Minion nr", id, "finished its job.")
	intercom <- "done"
}

func prepareExecString(input string, threadID int, inputID int) string {
	args.Run = strings.Replace(args.Run, "\"", "\\\"", -1)
	commandString := "INPUTSTR=\"" + input + "\";THREADID=\"" + strconv.Itoa(threadID) + "\";INPUTID=\"" + strconv.Itoa(inputID) + "\";" + args.Run + ";"
	return commandString
}

func runHeadlessMinions() {
	// create an intercom map
	intercom := make(map[int]chan string)

	// check some things
	if args.ThreadAmount > args.Runs {
		fmt.Println("You did set '", args.ThreadAmount, "' but set only '", args.Runs, "' runs. Resetting the Threadamount")
		args.ThreadAmount = args.Runs
	}

	// some number juggling , probably smarter ways
	// out there but for now its fine
	specialCase := -1
	base := float64(args.Runs) / float64(args.ThreadAmount)
	checkBase := math.Floor(float64(base))
	if base != checkBase {
		checkAmount := int(checkBase) * args.ThreadAmount
		specialCase = args.Runs - checkAmount
	}

	// run the minions
	for i := 0; i < args.ThreadAmount; i++ {
		intercom[i] = make(chan string)
		runs := int(checkBase)
		if specialCase != -1 && i == 0 {
			runs = int(checkBase) + specialCase
		}
		go createStupidMinion(i, runs, intercom[i])
	}

	// check for the minions to be done
	done := 0
	for {
		// check all thread channels for return
		for i := 0; i < args.ThreadAmount; i++ {
			ret := <-intercom[i]
			if "done" == ret {
				done++
			}
		}

		// all threads are done
		if args.ThreadAmount == done {
			break
		}
	}
	//fmt.Println("All work is done. Exiting")
}

func createStupidMinion(id int, runs int, intercom chan string) {
	//fmt.Println("Minion nr", id, "spawned. ")
	for i := 0; i < runs; i++ {
		prepared := prepareExecString("", id, i)
		ret, _ := custExec(prepared)
		if "" != ret {
			logger.Println(ret)
		}
	}
	//fmt.Println("Minion nr", id, "finished its job.")
	intercom <- "done"
}

func splitInput() map[int]string {
	// predefine return map
	data := make(map[int]string)

	// split the input string
	inputSlice := strings.Split(args.Input, args.Delimiter)

	// check if the input amount is less than the given thread amount
	if args.ThreadAmount > len(inputSlice) {
		fmt.Println("You set to execute more threads than inputs provided. Resetting threads amount to input amount '", len(inputSlice), "'.")
		args.ThreadAmount = len(inputSlice)
	}

	// transform the input to map
	i := 0
	for _, val := range inputSlice {
		data[i] = val
		i++
	}

	return data
}

func custExec(cmd string) (string, error) {
	output, err := exec.Command("/bin/bash", "-c", cmd).Output()
	if nil != err {
		return "", errors.New("Error executin command")
	}
	return strings.TrimRight(string(output), "\n"), nil
}

func parseArgs() {
	// delimiter to be used for custom generator input
	var delimiter string
	flag.StringVar(&delimiter, "delimiter", "\n", "Cutoff color value")

	// thread amount to be used
	var threadAmount int
	flag.IntVar(&threadAmount, "threads", 1, "Amount of threads to be used")

	// amount of runs in case you dont provide input
	var runs int
	flag.IntVar(&runs, "runs", -1, "Amount of threads to be used")

	// input by string
	var run string
	flag.StringVar(&run, "run", "", "Job to be executed")

	// parse the flags
	flag.Parse()

	args = Args{
		Delimiter:    delimiter,
		ThreadAmount: threadAmount,
		Run:          run,
		Input:        "",
		Runs:         runs,
	}

}