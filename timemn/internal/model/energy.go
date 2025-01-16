package model

import (
	"context"
	"time"
)

const (
	SLEEPING_WEIGHT = 0.4
	NUTS_WEIGHT     = 0.2
	EXERCISE_WEIGHT = 0.1
	HS_WEIGHT       = 0.1
	FEELING_WEIGHT  = 0.2
)

type DailySleepMetric struct {
	StartSleepingTime time.Time
	EndSleepingTime   time.Time
}

type DailyExercise map[string]float64

type DailyPhysicalInput struct {
	Row        int
	Date       time.Time
	Sleep      DailySleepMetric
	Nuts       int
	Sxs        int
	Exercise   DailyExercise
	Feeling    float64
	SleepScore float64 // calculated before if any
	Etf        float64 // calculated before if any
	Etd        float64 // calculated before if any
	Note       string
}

type EnergyMetric struct {
	Name   string
	Weight float64
}

type EnergyMngService interface {
	// EvaluateEnergyToday(input *DailyPhysicalInput) (float32, error)
	EvaluateAllFromSheet(ctx context.Context, forceOverwrite bool) error
}
