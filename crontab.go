package cronv

import (
	"fmt"
	"strconv"
	"strings"
)

type Schedule struct {
	Minute     string
	Hour       string
	DayOfMonth string
	Month      string
	DayOfWeek  string
	Year       string
}

func (self *Schedule) ToCrontab() string {
	dest := strings.Join([]string{self.Minute, self.Hour, self.DayOfMonth,
		self.Month, self.DayOfWeek, self.Year}, " ")
	return strings.Trim(dest, " ")
}

type Crontab struct {
	Line     string
	Schedule *Schedule
	Job      string
}

func (c *Crontab) IsRunningEveryMinutes() bool {
	for i, v := range strings.Split(c.Schedule.ToCrontab(), " ") {
		if v != "*" && (i > 0 || v != "*/1") {
			return false
		}
	}
	return true
}

type InvalidTaskError struct {
	Line string
}

func (e *InvalidTaskError) Error() string {
	return fmt.Sprintf("Invalid task: '%s'", e.Line)
}

func ParseCrontab(line string) (*Crontab, error) {
	// TODO use regrex to parse: https://gist.github.com/istvanp/310203
	parts := strings.Split(line, " ")
	if len(parts) < 5 {
		return nil, &InvalidTaskError{line}
	}

	// https://en.wikipedia.org/wiki/Cron#Predefined_scheduling_definitions
	schedule := &Schedule{}
	job := []string{}
	c := 0
	for _, v := range parts {
		if len(v) == 0 {
			continue
		}
		switch c {
		case 0:
			schedule.Minute = v
		case 1:
			schedule.Hour = v
		case 2:
			schedule.DayOfMonth = v
		case 3:
			schedule.Month = v
		case 4:
			schedule.DayOfWeek = v
		default:
			job = append(job, v)
		}
		c++
	}
	crontab := &Crontab{
		Line:     line,
		Schedule: schedule,
		Job:      strings.Join(job, " "),
	}
	return crontab, nil
}

func ParseJobDuration(line string) (float64, error) {
	durationMinutes := 1.0
	parts := strings.Split(line, "#")
	if len(parts) > 1 {
		minutes, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return durationMinutes, err
		}
		durationMinutes = minutes
	}
	return durationMinutes, nil
}
