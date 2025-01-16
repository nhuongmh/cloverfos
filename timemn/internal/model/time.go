package model

import "time"

type HProcessWeight int

const (
	LEISURE HProcessWeight = iota
	LEVEL_1
	LEVEL_2
	LEVEL_3
	LEVEL_4
	LEVEL_5
	EMERGENCY
)

type TimeSlot struct {
	Base
	Start time.Time
	End   time.Time
}

type Deadline struct {
	Base
	DeadlineType int
	HardDeadline time.Time
	// WithinDuration time.Duration
}

type RepetitiveType int

const (
	Daily RepetitiveType = iota
	Weekdays
	Weekends
	Weekly
	Biweekly
	Monthly
	Yearly
	Custom
)

type Repetitive struct {
}

// Human Process
type HProcess struct {
	Base
	Name              string
	Description       string
	Weight            HProcessWeight
	Deadline          Deadline
	EstimatedTimeCost time.Duration
	State             string
	Repeat            bool
}

type TScheduler interface {
	GetFreetime(from, to time.Time) *[]TimeSlot
	GetTasksOfCustomTime(from, to time.Time) *[]TimeSlot
	GetWhatNow() *TimeSlot
	GetTodayTasks() *[]TimeSlot
	GetThisWeekTasks() *[]TimeSlot
	SetActually()
	ScheduleTask(hp *HProcess) (*[]TimeSlot, error)
}

type Day interface {
}
