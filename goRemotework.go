/*
 * Remote work support tool made in Go language.
 *
 * @author    yasutakatou
 * @copyright 2025 yasutakatou
 * @license   
 */
package main

import (
	"syscall"
	"flag"
	"fmt"
	"os"
	"time"
	"io"
	"log"
	"strings"
	"unsafe"
	"bufio"
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

type RECTdata struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

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
    //_Verbose := flag.Bool("verbose", false, "[-verbose=incident output verbose (true is enable)]")
    _ScheduleConfig := flag.String("scheduleconfig", "schedule.ini", "[-scheduleconfig=specify the configuration file for scheduled alerts.]")
    _TasksConfig := flag.String("tasksconfig", "tasks.ini", "[-tasksconfig=specify the task aggregation config file.]")
    //_OutputConfig := flag.String("outputconfig", "output.txt", "[-outputconfig=specify the output file of the work history.]")

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

    // for {
    //     time.Sleep(time.Second * time.Duration(*_Loop))
    //     taskAlert(*_Verbose)
    // }
	ListWindow()
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
	reader := bufio.NewReader(strings.NewReader(configFile))

	for {
		buf, err := reader.ReadString('\n')
		debugLog(buf)
		// strs := strings.Split(buf, "\t")
		// if len(strs) == 3 {
		// 	tasks = append(tasks, tasksData{REGEX: strs[0],LIMIT: strs[1], COMMAND: strs[2]})
		// }

		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
	}

	// if tasks == nil {
	// 	return false
	// }
	return true
}

func loadScheduleConfig(configFile string) bool {
	debugLog(" -- " + configFile + " --")
	reader := bufio.NewReader(strings.NewReader(configFile))

	for {
		buf, err := reader.ReadString('\n')
		debugLog(buf)
		// strs := strings.Split(buf, "\t")
		// if len(strs) == 2 {
		// 	schedules = append(schedules, scheduleData{DATE: strs[0],COMMAND: strs[1]})
		// }

		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
	}

	// if schedules == nil {
	// 	return false
	// }
	return true
}

func GetWindowRect(hwnd HWND, rect *RECTdata) (err error) {
	r1, _, e1 := syscall.Syscall(procGetWindowRect.Addr(), 7, uintptr(hwnd), uintptr(unsafe.Pointer(rect)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
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

func EnumWindows(enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procEnumWindows.Addr(), 2, uintptr(enumFunc), uintptr(lparam), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}


func ListWindow() []string {
	var rect RECTdata

	ret := []string{}

	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		b := make([]uint16, 200)
		_, err := GetWindowText(h, &b[0], int32(len(b)))
		if err != nil {
			return 1
		}

		GetWindowRect(HWND(h), &rect)
		if rect.Left != 0 || rect.Top != 0 || rect.Right != 0 || rect.Bottom != 0 {
			if debug == true {
				fmt.Printf("Window Title '%s' window: handle=0x%x\n", syscall.UTF16ToString(b), h)
				if rect.Left != 0 || rect.Top != 0 || rect.Right != 0 || rect.Bottom != 0 {
					fmt.Printf("window rect: ")
					fmt.Println(rect)
				}
			}
			ret = append(ret, fmt.Sprintf("%s : %x", syscall.UTF16ToString(b), h))
		}
		return 1
	})
	EnumWindows(cb, 0)
	return ret
}
