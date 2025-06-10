package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parsedStr := strings.Split(data, ",")
	if len(parsedStr) != 3 {
		return 0, "", 0, errors.New("invalid data")
	}

	stepsStr := parsedStr[0]
	activityType := parsedStr[1]
	durationStr := parsedStr[2]

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, err
	}

	if steps <= 0 {
		return 0, "", 0, errors.New("invalid steps")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, err
	}

	if duration <= 0 {
		return 0, "", 0, errors.New("invalid duration")
	}

	return steps, activityType, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distanceInM := stepLength * float64(steps)
	distanceInKm := distanceInM / float64(mInKm)

	return distanceInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distanceInKm := distance(steps, height)

	return distanceInKm / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	var calories float64

	switch activityType {
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	if err != nil {
		return "", err
	}

	speed := meanSpeed(steps, height, duration)
	distanceResult := distance(steps, height)
	durationInHours := duration.Hours()

	var result string

	result += fmt.Sprintf("Тип тренировки: %s\n", activityType)
	result += fmt.Sprintf("Длительность: %.2f ч.\n", durationInHours)
	result += fmt.Sprintf("Дистанция: %.2f км.\n", distanceResult)
	result += fmt.Sprintf("Скорость: %.2f км/ч\n", speed)
	result += fmt.Sprintf("Сожгли калорий: %.2f\n", calories)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("invalid spent calories")
	}

	avgSpeed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()

	calories := (weight * avgSpeed * minutes) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("invalid spent calories")
	}

	calories, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}

	walkingCalories := calories * walkingCaloriesCoefficient

	return walkingCalories, nil
}
