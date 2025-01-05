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
	"regexp"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)

type taskListsData struct {
	NAME string
	TIME int
}

type tasksData struct {
	NAME    string
	REGEX   string
	LIMIT   int
	CURRENT int
	COMMAND string
}

type scheduleData struct {
	DATE    string
	COMMAND string
}

var (
	debug, logging bool
	schedules      []scheduleData
	tasks          []tasksData
	tasklists      []taskListsData
)

func main() {
	_Debug := flag.Bool("debug", false, "[-debug=debug mode (true is enable)]")
	_Logging := flag.Bool("log", false, "[-log=logging mode (true is enable)]")
	_Loop := flag.Int("loop", 60, "[-loop=incident check loop time (Seconds). ]")
	_ScheduleConfig := flag.String("scheduleconfig", "schedule.ini", "[-scheduleconfig=specify the configuration file for scheduled alerts.]")
	_TasksConfig := flag.String("tasksconfig", "tasks.ini", "[-tasksconfig=specify the task aggregation config file.]")
	_OutputConfig := flag.String("outputconfig", "output.txt", "[-outputconfig=specify the output file of the work history.]")

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

	for {
		time.Sleep(time.Second * time.Duration(*_Loop))
		taskAlert(*_OutputConfig, *_Loop)
	}
	os.Exit(0)
}

func taskAlert(filename string, duration int) {
	title := getCurrentWindow()

	for _, rule := range tasks {
		debugLog("ruleRegex: " + rule.REGEX)
		ruleRegex := regexp.MustCompile(rule.REGEX)
		if ruleRegex.MatchString(title) == true {
			tFlag := false
			for i := 0; i < len(tasklists); i++ {
				if tasklists[i].NAME == rule.NAME {
					tasklists[i].TIME = tasklists[i].TIME + duration
					tFlag = true
				}
			}
			if tFlag == false {
				tasklists = append(tasklists, taskListsData{NAME: rule.NAME, TIME: duration})
			}
		}
	}
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
		if len(record) == 4 {
			i, err := strconv.Atoi(record[2])
			if err == nil {
				tasks = append(tasks, tasksData{NAME: record[0], REGEX: record[1], LIMIT: i, CURRENT: i, COMMAND: record[3]})
				fmt.Println(record)
			} else if record[1] == "NO" && record[2] == "NO" {
				tasks = append(tasks, tasksData{NAME: record[0], REGEX: record[1], LIMIT: 0, CURRENT: 0, COMMAND: ""})
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

func getCurrentWindow() string {
	//https://stakiran.hatenablog.com/entry/2019/08/19/191433

	user32, err := syscall.LoadDLL("user32.dll")
	if err != nil {
		panic(err)
	}
	defer user32.Release()

	procGetForegroundWindow, err := user32.FindProc("GetForegroundWindow")
	if err != nil {
		panic(err)
	}
	hwnd, _, _ := procGetForegroundWindow.Call()

	procGetWindowTextLength, err := user32.FindProc("GetWindowTextLengthW")
	if err != nil {
		panic(err)
	}
	textLength, _, _ := procGetWindowTextLength.Call(hwnd)
	textLength = textLength + 1

	procGetWindowText, err := user32.FindProc("GetWindowTextW")
	if err != nil {
		panic(err)
	}

	buf := make([]uint16, textLength)
	procGetWindowText.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), textLength)

	text := syscall.UTF16ToString(buf)
	return text
}
