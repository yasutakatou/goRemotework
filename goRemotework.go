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

	if Exists(*_TasksConfig) == true {
		loadConfig(0)
	} else {
		fmt.Printf("Fail to read task config file: %v\n", *_TasksConfig)
		os.Exit(1)
	}

    if loadConfig(1) == true {
        go func() {
			for {
				time.Sleep(time.Second * time.Duration(*_Loop))
				//scheduleAlert()
			}
		}()
    }

    for {
        time.Sleep(time.Second * time.Duration(*_Loop))
        //scheduleAlert()
    }
    os.Exit(0)
}
