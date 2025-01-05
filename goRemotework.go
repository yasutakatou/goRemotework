/*
 * Remote work support tool made in Go language.
 *
 * @author    yasutakatou
 * @copyright 2025 yasutakatou
 * @license
 */
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)

type (
	HANDLE uintptr
	HWND   HANDLE
)

var (
	user32                  = syscall.MustLoadDLL("user32.dll")
	procGetWindowTextW      = user32.MustFindProc("GetWindowTextW")
	procGetForegroundWindow = user32.MustFindProc("GetForegroundWindow")
	debug, logging          bool
	schedules               []scheduleData
	tasks                   []tasksData
)

type tasksData struct {
	REGEX   string
	LIMIT   int
	COMMAND string
}

type scheduleData struct {
	DATE    string
	COMMAND string
}

func main() {
	_Debug := flag.Bool("debug", false, "[-debug=debug mode (true is enable)]")
	_Logging := flag.Bool("log", false, "[-log=logging mode (true is enable)]")
	_Loop := flag.Int("loop", 60, "[-loop=incident check loop time (Seconds). ]")
	//_Verbose := flag.Bool("verbose", false, "[-verbose=incident output verbose (true is enable)]")
	_ScheduleConfig := flag.String("scheduleconfig", "schedule.ini", "[-scheduleconfig=specify the configuration file for scheduled alerts.]")
	_TasksConfig := flag.String("tasksconfig", "tasks.ini", "[-tasksconfig=specify the task aggregation config file.]")
	//_OutputConfig := flag.String("outputconfig", "output.txt", "[-outputconfig=specify the output file of the work history.]")

	flag.Parse()

	debug = bool(*_Debug)
	logging = bool(*_Logging)

	if loadTasksConfig(*_TasksConfig) == false {
		fmt.Printf("Fail to read task config file: %v\n", *_TasksConfig)
		os.Exit(1)
	}

	if loadScheduleConfig(*_ScheduleConfig) == true {
		go func() {
			for {
				time.Sleep(time.Second * time.Duration(*_Loop))
				//scheduleAlert(*_Verbose)
			}
		}()
	}

	// for {
	//     time.Sleep(time.Second * time.Duration(*_Loop))
	//     taskAlert(*_Verbose)
	// }
	fmt.Println(getWindow())
	os.Exit(0)
}

func debugLog(message string) {
	var file *os.File
	var err error

	if debug == true {
		fmt.Println(message)
	}

	if logging == false {
		return
	}

	const layout = "2006-01-02_15"
	const layout2 = "2006/01/02 15:04:05"
	t := time.Now()
	filename := t.Format(layout) + ".log"
	logHead := "[" + t.Format(layout2) + "] "

	if Exists(filename) == true {
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)
	} else {
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	}

	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()
	fmt.Fprintln(file, logHead+message)
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func loadTasksConfig(configFile string) bool {
	debugLog(" -- " + configFile + " --")
	var fp *os.File
	var err error
	fp, err = os.Open(configFile)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader := csv.NewReader(fp)
	reader.Comma = '\t'
	reader.LazyQuotes = true
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if len(record) == 3 {
			i, err := strconv.Atoi(record[1])
			if err == nil {
				tasks = append(tasks, tasksData{REGEX: record[0], LIMIT: i, COMMAND: record[2]})
				fmt.Println(record)
			} else if record[1] == "NO" && record[2] == "NO" {
				tasks = append(tasks, tasksData{REGEX: record[0], LIMIT: 0, COMMAND: ""})
				fmt.Println(record)
			}
		}
	}
	if tasks == nil {
		return false
	}
	return true
}

func loadScheduleConfig(configFile string) bool {
	debugLog(" -- " + configFile + " --")
	var fp *os.File
	var err error
	fp, err = os.Open(configFile)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader := csv.NewReader(fp)
	reader.Comma = '\t'
	reader.LazyQuotes = true
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if len(record) == 2 {
			schedules = append(schedules, scheduleData{DATE: record[0], COMMAND: record[1]})
			fmt.Println(record)
		}
	}
	if schedules == nil {
		return false
	}
	return true
}

func GetWindowText(hwnd syscall.Handle, str *uint16, maxCount int32) (len int32, err error) {
	r0, _, e1 := syscall.Syscall(procGetWindowTextW.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(str)), uintptr(maxCount))
	if len = int32(r0); len == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func getWindow() string {
	hwnd, _, _ := syscall.Syscall(procGetForegroundWindow.Addr(), 6, 0, 0, 0)
	if debug == true {
		fmt.Printf("currentWindow: handle=0x%x\n", hwnd)
	}
	return ""
}
