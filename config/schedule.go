package config

import "os"

var theSchedule = "schedules/schedule.json"

func GetTheSchedule() string {
	if os.Getenv("SCHEDULE") != "" {
		theSchedule = os.Getenv("SCHEDULE")
	}
	return theSchedule
}
