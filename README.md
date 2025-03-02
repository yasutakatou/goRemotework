# goRemotework
**Remote work support tool made in Go language. Bridges the gap between junior engineers who are unfamiliar with remote work and managers who want to manage it properly.**

## Challenges of onboarding junior engineers in remote work
## Evidence.
### Junior engineer's point of view
- (1) Although daily work is summarized in daily reports, etc., it is difficult to visualize what was done because it is hard to grasp exactly how much time was actually consumed on what. In addition, the time to write the daily report itself is time-consuming and tends to be loose.
- (2) When they are stuck in a task, they keep getting stuck in trying to solve it on their own because they have not developed the judgment of how much research they should do to report back.
- (3) They have not developed the ability to manage their time, and although they try to manage their time by putting appointments on a calendar tool, they cannot receive time alerts properly due to conflicts with meeting schedules.

As a result, junior engineers who have no (or little) experience working in an office and who are not accustomed to time management work remotely, which results in stagnation in improving their work skills and, in turn, requires more effort from managers.

### Manager's Perspective
As a counterpart to the junior engineer's viewpoint, the following efforts will be required of the manager

- (1) Skills planning is difficult even after reading the daily report; it cannot be thought of in terms of KPI. Also, there is no other way but to base it on the goodness of character in terms of the need to manage whether there is any slacking off or time not spent on work.
- (2) As in the office, it is not possible to look over one's back and see if one is stuck and call out at the right time, so the work takes longer than expected.
- (3) Increased time and effort to check things that could have been prevented by a simple voice call in the office. It is also necessary to check whether there are any omissions in the tasks instructed

## Claim
This is not an engineering technology problem of the human resource development system that both sides are facing, **but a potential underlying problem is that remote work is based on office work, which takes time for junior engineers to get used to and makes it difficult for senior engineers to provide proper guidance and training opportunities**.<br>
Therefore, it seems that apps, SaaS, and tools with mechanisms to minimize the disadvantages are not available for junior engineers, with reference to the following anti-patterns that have occurred

### anti-pattern
- Implement remote monitoring tool to save keystrokes and record work status from in-camera to record turnover, etc. Provide appropriate assistance to junior engineers from the output of the tool.

This method is based on a sexist approach. Management practices that are deprecated even for office work, such as reduced psychological safety, excessive monitoring, micromanagement, etc.

- All-day office work ※**Focuses on the disadvantageous aspects based on the claims, as the advantages are also significant**

It seems to be effective from a training perspective, as it is incorporated into the objective of increasing attendance for the purpose of improving productivity. However, there may be situations where senior engineers and managers may want to apply commuting time to work time in a post-optimization situation where their work is optimized for remote work

## Guarantee
### Desired state
- (1) Visualize and quantify the results as transparently and evenly as possible without individual differences.
- (2) Catch up with managers and senior engineers when they are in a state of excessive maze of thoughts.
- (3) Tasks can be started at the timing agreed with the manager and senior engineers.

**because support and training for junior engineers is a kind of input and output using interpersonal actions as an interface, and is important for gaining engineering experience**.

# Functions
- (1) **Collect information on applications used in the background and measure the time spent on each.**
  - (ollecting at the application information level, which is a layer below keystroke measurement monitoring, prevents excessive monitoring.
- (2) **Aggregate the time spent on a particular application, and if it exceeds a certain amount of time, run a defined Runbook and automatically trigger an alert.**
  - Enable automatic catch-up on situations that boil over. Also, management can set the alert time to increase psychological safety for reporting
- (3) **notify at a specified time and date so that people are aware of what to do at what time**.
  - Prevent task forgetting by separating it from the meeting calendar. This also fosters a sense of commitment to the task on both sides by setting differences with the management side

# Specific actions
**Using the assumption that business PCs are always started up and shut down, loop the operations (1) to (3) at regular intervals to make them into a tool**.

- (1) Collect information on applications used in the background and measure the time consumed for each.
  - Obtain application titles for windows in the foreground, and write the analogous time to a file on a regular basis
  - Shrink the same application that has different titles by using regular expressions and aggregate them as the same item so that they can be treated as a single task *For example, a browser whose title changes depending on the open screen.
- (2) Aggregate the operation time of specific applications, and if the time exceeds a certain amount of time, execute the defined Runbook and automatically issue an alert.
  - Similarly to (1) above, when the time that a specific application is running in the foreground exceeds a certain time, a defined command is executed.
  - In the same way as (1) above, set up a mechanism that allows the same items to be aggregated by shrinking them using regular expressions.
- (3) Notification at a specified time or date and time allows the user to notice what to do at what point in time.
  - Execute commands defined at a specific date, time, or hour by regular expressions

## Usage
The system operates by the rules defined in each configuration file. Therefore, the following team operation would be desirable

- PMs and senior engineers prepare the “Task Aggregation Definition” and “Schedule Notification” configuration files and share them with the team.
- Set up the system to start automatically when Windows starts up [example of setting](https://note.com/bright_clover112/n/ncd35e325b202)
- Upload the output work time file periodically.
- PM evaluates the work time file and tunes it for proper management.

# Operation screen
- (1) Collect information on applications used in the background and measure the time consumed for each

Aggregate regular expression matches in windows that were operated in the foreground.

![image](https://github.com/user-attachments/assets/99e76027-42f7-4da7-ad1f-90607ef72baa)

- (2) Aggregate the operation time of specific applications, and if it exceeds a certain amount of time, run the defined Runbook and automatically fire an alert.

Example of an alert fired because of continuous search on stackoverflow.

![image](https://github.com/user-attachments/assets/3737c0ae-4e53-4d22-b437-d96529958674)

- (3) Notification at a specified time and date so that you are aware of what to do at what time.

Let's take a break using the Pomodoro Technique at 45 minutes every hour, example.

![image](https://github.com/user-attachments/assets/6fae4eee-fb2b-4451-a07f-9d65e9a2909e)

# How to install
Get the tools from this repository

```
go get github.com/yasutakatou/goRemotework
```

Or clone and build it.

```
git clone https://github.com/yasutakatou/goRemotework
cd caplint
go build .
```

If it's too much trouble, [you can extract the binary files at hand](https://github.com/yasutakatou/goRemotework/releases).

# How to uninstall

Since it is Go language, you can delete the binary file!

# Configuration file

**Configuration file is in tsv (tab split value) format. Define tab-separated values**.

## Task Aggregation Definition (and Alert Runbook)

This configuration file defines task names to be aggregated and (if set) alerts for time exceeded

```
(1) Definition name (2) Regular expression of the window name to be aggregated (3) Time exceeded (in seconds) (4) Command to be executed when time exceeded (5) Arguments to be given to the command in “5)”.
```

For example, if you want to measure the time spent looking at the AWS management console and be alerted when it exceeds 30 minutes, set the following
(PMs and senior engineers can set up alerts, which lowers the barrier to escalation)

```
AWS . *Amazon Web Services.* 1800 popup.bat AWS_USe
```

To measure only window name and typed time without alerting, write “NO” in (3) to (5) as follows

```
Notepad . *Notebook . * NO NO NO NO	NO
```

All window names that do not fit in the regular expression will be grouped into “OTHER”. (It is designed not to over-monitor even if you have opened Youtube with background music for work.

### Added pop-up window functionality (v0.2)

Specify “POPUP” as the command name, followed by a message to display an OS-native popup

```
Chrome	.*Chrome.*	100	POPUP	ポップアップ テスト ！！
```

![image](https://github.com/user-attachments/assets/c7768c79-d3d9-431c-ac37-0632ab061edc)

note) In previous versions, the only way to notify messages was to pass parameters to the command, so it was not possible to notify messages with spaces or non-English messages such as Japanese.

## Schedule notifications

This configuration file will implement alerts at set times
The reason for separating the task aggregation definitions and files is to allow for team-shared alert settings, plus the ability to add your own alerts for individual alerts.

```
(1) Regular expressions to generate alerts (2) Alerts for each team member
```

The date and time are recorded in the following format.

```
2025/01/06 13:45:38 Mon
```

The year, month, day, time, and day of the week. For example, if you want to be alerted every day from 10am to 1pm, you would set the following

```
. */. */. * 1[0-3]:. *:. * . * popup.bat 10-13HourAlert
```

If it is Monday at 9:00, then

```
. */. */. * 9:00:00 Mon popup.bat MondayAlert
```

Note that if you extend the loop interval options too far (see below), you may not get an alert.

### Added pop-up window functionality (v0.2)

Specify “POPUP” as the command name, followed by a message to display an OS-native popup

```
Chrome	.*Chrome.*	100	POPUP	ポップアップ テスト ！！
```

![image](https://github.com/user-attachments/assets/c7768c79-d3d9-431c-ac37-0632ab061edc)

note) In previous versions, the only way to notify messages was to pass parameters to the command, so it was not possible to notify messages with spaces or non-English messages such as Japanese.

# Option

```
  -debug
        [-debug=debug mode (true is enable)]
  -log
        [-log=logging mode (true is enable)]
  -loop int
        [-loop=incident check loop time (Seconds). ] (default 60)
  -outputconfig string
        [-outputconfig=specify the output file of the work history.] (default “output.txt”)
  -scheduleconfig string
        [-scheduleconfig=specify the configuration file for scheduled alerts.] (default “schedule.ini”)
  -tasksconfig string
        [-tasksconfig=specify the task aggregation config file.] (default “tasks.ini”)
```

## -debug

Runs in debug mode. Various output is printed.

## -log

Option to output the log from debug mode.

## -loop int

Interval (in seconds) to loop the entire operation.<br>
If it is too long, it cannot be measured accurately, so it is better to keep it as short as possible and operate in increments of a few seconds.
Also, **if the schedule notification setting is crossed by the loop interval, alerts will not be issued, so be careful about this**.

## -outputconfig string

If you want to change the file to measure only window name and typed time from the default output.txt, specify the name of the config file with this.<br>
The measurement file is always overwritten because it is assumed to be uploaded and synchronized at OS shutdown. Therefore, if an unexpected OS reboot is possible, it is better to change the name of this file to the date and time, etc.

## -scheduleconfig string

If you want to change the schedule notification configuration file from the default schedule.ini, specify the configuration file name with this.

## -tasksconfig string

Specify the name of the configuration file when you want to change the task aggregation definition file from the default tasks.ini.

# Lisence
BSD 3-Clause
