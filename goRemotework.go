/*
 * Remote work support tool made in Go language.
 *
 * @author    yasutakatou
 * @copyright 2025 yasutakatou
 * @license   BSD 3-Clause
 */
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"syscall"
	"time"
	"unsafe"

	"github.com/ChromeTemp/Popup"
)

type taskListsData struct {
	NAME  string
	TIME  int
	LIMIT int
}

type tasksData struct {
	NAME    string
	REGEX   string
	LIMIT   int
	COMMAND string
	MESSAGE string
}

type scheduleData struct {
	DATE    string
	COMMAND string
	MESSAGE string
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
	_Loop := flag.Int("loop", 5, "[-loop=incident check loop time (Seconds). ]")
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
				scheduleAlert()
			}
		}()
	}

	tasklists = append(tasklists, taskListsData{NAME: "OTHER", TIME: 0, LIMIT: 0})
	for {
		time.Sleep(time.Second * time.Duration(*_Loop))
		taskAlert(*_OutputConfig, *_Loop)
	}
	os.Exit(0)
}

func scheduleAlert() {
	wdays := [...]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

	const layout = "2006/01/02 15:04:05"
	t := time.Now()
	nowDate := t.Format(layout) + " " + wdays[t.Weekday()]
	debugLog("nowDate: " + nowDate)

	for _, rule := range schedules {
		dateRegex := regexp.MustCompile(rule.DATE)
		if dateRegex.MatchString(nowDate) == true {
			debugLog("dateRegex: " + rule.DATE)
			if rule.COMMAND == "POPUP" {
				Popup.Alert("popup", rule.MESSAGE)
			} else {
				go execCommand(rule.COMMAND, rule.MESSAGE)
			}
		}

	}
}

func taskAlert(filename string, duration int) {
	title := getCurrentWindow()
	debugLog("title: " + title)

	mFlag := false
	for _, rule := range tasks {
		debugLog("ruleRegex: " + rule.REGEX)
		ruleRegex := regexp.MustCompile(rule.REGEX)
		if ruleRegex.MatchString(title) == true {
			mFlag = true
			debugLog("matchRegex: " + rule.REGEX)
			tFlag := false
			for i := 0; i < len(tasklists); i++ {
				if tasklists[i].NAME == rule.NAME {
					tasklists[i].TIME = tasklists[i].TIME + duration
					tFlag = true

					if rule.LIMIT != 0 && len(rule.COMMAND) != 0 && len(rule.MESSAGE) != 0 {
						if tasklists[i].LIMIT-duration <= 0 {
							tasklists[i].LIMIT = rule.LIMIT
							if rule.COMMAND == "POPUP" {
								Popup.Alert("popup", rule.MESSAGE)
							} else {
								go execCommand(rule.COMMAND, rule.MESSAGE)
							}
						} else {
							tasklists[i].LIMIT = tasklists[i].LIMIT - duration
						}
					}
				}
			}
			if tFlag == false {
				if rule.LIMIT != 0 && len(rule.COMMAND) != 0 && len(rule.MESSAGE) != 0 {
					tasklists = append(tasklists, taskListsData{NAME: rule.NAME, TIME: duration, LIMIT: rule.LIMIT})
				} else {
					tasklists = append(tasklists, taskListsData{NAME: rule.NAME, TIME: duration, LIMIT: 0})
				}
			}
		}
	}

	if mFlag == false {
		for i := 0; i < len(tasklists); i++ {
			if tasklists[i].NAME == "OTHER" {
				tasklists[i].TIME = tasklists[i].TIME + duration
			}
		}
	}

	sort.Slice(tasklists, func(i, j int) bool {
		return tasklists[i].TIME > tasklists[j].TIME
	})

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	fmt.Println(" - - - tasklists - - - ")
	for _, rule := range tasklists {
		fmt.Println(rule.NAME + " " + strconv.Itoa(rule.TIME) + " " + strconv.Itoa(rule.LIMIT))
		fmt.Fprintln(file, rule.NAME+" "+strconv.Itoa(rule.TIME))
	}
	fmt.Println(" - - - - - - - - - - ")

}

func execCommand(command, message string) {
	debugLog("command: " + command + " " + message)
	cmd := exec.Command("cmd", "/c", command+" "+message)
	output, err := cmd.CombinedOutput()
	if err != nil {
		debugLog(fmt.Sprint(err) + ": " + string(output))
	} else {
		debugLog(string(output))
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
		if len(record) == 5 {
			i, err := strconv.Atoi(record[2])
			if err == nil {
				tasks = append(tasks, tasksData{NAME: record[0], REGEX: record[1], LIMIT: i, COMMAND: record[3], MESSAGE: record[4]})
				fmt.Println(record)
			} else if record[2] == "NO" && record[3] == "NO" {
				tasks = append(tasks, tasksData{NAME: record[0], REGEX: record[1], LIMIT: 0, COMMAND: "", MESSAGE: ""})
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
		if len(record) == 3 {
			schedules = append(schedules, scheduleData{DATE: record[0], COMMAND: record[1], MESSAGE: record[2]})
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
