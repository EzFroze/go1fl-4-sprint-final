package daysteps

import (
	"errors"
	"fmt"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parsedData := strings.Split(data, ",")
	if len(parsedData) != 2 {
		return 0, 0, errors.New(fmt.Sprint("invalid data:", data))
	}
	stepsSrt := parsedData[0]
	dateStr := parsedData[1]

	steps, err := strconv.Atoi(stepsSrt)
	if err != nil || steps == 0 {
		return 0, 0, err
	}

	duration, err := time.ParseDuration(dateStr)

	if err != nil {
		return 0, 0, err
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)

	if err != nil {
		return ""
	}

	lengthInM := stepLength * float64(steps)

	lengthInKm := lengthInM / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	if err != nil {
		return ""
	}

	var result string
	result += fmt.Sprintf("Количество шагов: %d\n", steps)
	result += fmt.Sprintf("Дистанция составила: %.2f км.\n", lengthInKm)
	result += fmt.Sprintf("Вы сожгли: %.2f ккал.\n", calories)

	return result
}
