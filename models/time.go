package models

import "time"

type MinuteRoundedTime struct {
	time.Time
}

func NewMinuteRoundedTime(t time.Time) MinuteRoundedTime {
	return MinuteRoundedTime{time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())}
}
func (t MinuteRoundedTime) MicroSecond() int64 {
	return t.Unix() * 1000
}
