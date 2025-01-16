package energy

import (
	"strconv"
	"strings"
	"time"

	"github.com/nhuongmh/cvfs/timemn/internal/model"
)

func parseDate(date string) (time.Time, error) {
	const layout = "Mon, 1/2/06"
	parsedTime, err := time.Parse(layout, date)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime, nil
}

func parseSleepTime(sleepStart string, sleepEnd string) (model.DailySleepMetric, error) {
	const layout = "3:04 PM"
	startTime, err := time.Parse(layout, sleepStart)
	if err != nil {
		return model.DailySleepMetric{}, err
	}

	endTime, err := time.Parse(layout, sleepEnd)
	if err != nil {
		return model.DailySleepMetric{}, err
	}

	return model.DailySleepMetric{
		StartSleepingTime: startTime,
		EndSleepingTime:   endTime,
	}, nil
}

func parseExercise(exercise string) (*model.DailyExercise, error) {
	exercises := make(model.DailyExercise)
	if strings.TrimSpace(exercise) == "" {
		return &exercises, nil
	}
	lines := strings.Split(exercise, "\n")
	for _, line := range lines {
		parts := strings.Split(line, "-")
		if len(parts) != 2 {
			return nil, model.ErrInvalidInput
		}
		activity := strings.TrimSpace(parts[0])
		durationStr := strings.TrimSpace(parts[1])

		var duration float64
		if strings.HasSuffix(durationStr, "min") {
			minutes, err := strconv.Atoi(strings.TrimSuffix(durationStr, "min"))
			if err != nil {
				return nil, err
			}
			duration = float64(minutes) / 60
		} else if strings.HasSuffix(durationStr, "hr") {
			hours, err := strconv.Atoi(strings.TrimSuffix(durationStr, "hr"))
			if err != nil {
				return nil, err
			}
			duration = float64(hours)
		} else {
			return nil, model.ErrInvalidInput
		}

		exercises[activity] = duration
	}

	return &exercises, nil
}
