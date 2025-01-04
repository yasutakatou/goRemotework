/*
 * Remote work support tool made in Go language.
 *
 * @author    yasutakatou
 * @copyright 2025 yasutakatou
 * @license   
 */
package main

import (
)

type (
	HANDLE uintptr
	HWND   HANDLE
)

var (
	user32                  = syscall.MustLoadDLL("user32.dll")
	procEnumWindows         = user32.MustFindProc("EnumWindows")
	procGetWindowTextW      = user32.MustFindProc("GetWindowTextW")
	procSetActiveWindow     = user32.MustFindProc("SetActiveWindow")
	procSetForegroundWindow = user32.MustFindProc("SetForegroundWindow")
	procGetForegroundWindow = user32.MustFindProc("GetForegroundWindow")
	procGetWindowRect       = user32.MustFindProc("GetWindowRect")
)

type tasksData struct {
	REGEX  string
    LIMIT  int
    COMMAND string
}

type scheduleData struct {
    DATE   string
    COMMAND string
}

var (
	debug, logging bool
	schedules              []scheduleData
	tasks                  []tasksData
)

func main() {
	_Debug := flag.Bool("debug", false, "[-debug=debug mode (true is enable)]")
	_Logging := flag.Bool("log", false, "[-log=logging mode (true is enable)]")
	_Loop := flag.Int("loop", 60, "[-loop=incident check loop time (Seconds). ]")
    _Verbose := flag.Bool("verbose", false, "[-verbose=incident output verbose (true is enable)]")
    _ScheduleConfig := flag.String("scheduleconfig", "schedule.ini", "[-scheduleconfig=specify the configuration file for scheduled alerts.]")
    _TasksConfig := flag.String("tasksconfig", "tasks.ini", "[-tasksconfig=specify the task aggregation config file.]")
    _OutputConfig := flag.String("outputconfig", "output.txt", "[-outputconfig=specify the output file of the work history.]")

	flag.Parse()

	debug = bool(*_Debug)
	logging = bool(*_Logging)

	if  loadTasksConfig(*_TasksConfig) == false {
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
        taskAlert(*_Verbose)
    }
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

func loadTasksConfig(configFile string) bool {
	debugLog(" -- " + configFile + " --")
	fp, err := os.Open(configFile)
	if err != nil {
		debugLog(configFile + " not found")
		return false
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	for scanner.Scan() {
		buf := scanner.Text()
		debugLog(buf)
		strs := strings.Split(buf, "\t")
		if len(strs) == 3 {
			tasks = append(tasks, tasksData{REGEX: strs[0],LIMIT: strs[1], COMMAND: strs[2]})
		}
	}

	if err = scanner.Err(); err != nil {
		debugLog(configFile + " error")
		return false
	}

	if tasks == nil {
		return false
	}
	return true
}

func loadScheduleConfig(configFile string) bool {
	debugLog(" -- " + configFile + " --")
	fp, err := os.Open(configFile)
	if err != nil {
		debugLog(configFile + " not found")
		return false
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	for scanner.Scan() {
		buf := scanner.Text()
		debugLog(buf)
		if len(strs) == 2 {
			schedules = append(schedules, scheduleData{DATE: strs[0],COMMAND: strs[1]})
		}
	}

	if err = scanner.Err(); err != nil {
		debugLog(configFile + " error")
		return false
	}

	if schedules == nil {
		return false
	}
	return true
}
