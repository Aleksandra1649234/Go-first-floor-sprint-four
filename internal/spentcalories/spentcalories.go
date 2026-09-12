package spentcalories

import (
	"errors"
	"fmt"
	"log"
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
	dataSl := strings.Split(data, ",")
	if len(dataSl) != 3 {
		return 0, "", 0, errors.New("Длина слайса != 3")
	}
	stepsCount, err := strconv.Atoi(dataSl[0])
	if err != nil {
		return 0, "", 0, err
	}
	if stepsCount <= 0 {
		return 0, "", 0, errors.New("неверные шаги")
	}
	dur, err := time.ParseDuration(dataSl[2])
	if err != nil {
		return 0, "", 0, err
	}
	if dur <= 0 {
		return 0, "", 0, errors.New("неверная продолжительность")
	}
	return stepsCount, dataSl[1], dur, nil
}

func distance(steps int, height float64) float64 {
	return (height * stepLengthCoefficient * float64(steps)) / float64(mInKm)
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0.0
	}
	avgSpeed := distance(steps, height) / duration.Hours()
	return avgSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	stepsCount, typeOfTraining, duration, err := parseTraining(data)
	if err == nil {
		var calories float64
		switch typeOfTraining {
		case "Ходьба":
			calories, _ = WalkingSpentCalories(stepsCount, weight, height, duration)
			return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", typeOfTraining, duration.Hours(), distance(stepsCount, height), meanSpeed(stepsCount, height, duration), calories), nil
		case "Бег":
			calories, _ = RunningSpentCalories(stepsCount, weight, height, duration)
			return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", typeOfTraining, duration.Hours(), distance(stepsCount, height), meanSpeed(stepsCount, height, duration), calories), nil
		default:
			return "", errors.New("неизвестный тип тренировки")
		}
	} else {
		log.Println(err)
		return "", err
	}
}
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps < 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные параметры тренировки")
	}
	speed := meanSpeed(steps, height, duration)
	if speed <= 0 {
		return 0, errors.New("некорректная скорость")
	}
	calories := (weight * speed * float64(duration.Minutes())) / minInH
	return calories, nil
}
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps < 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные параметры тренировки")
	}
	speed := meanSpeed(steps, height, duration)
	if speed <= 0 {
		return 0, errors.New("некорректная скорость")
	}
	calories := (weight * speed * float64(duration.Minutes())) / minInH * walkingCaloriesCoefficient
	return calories, nil
}