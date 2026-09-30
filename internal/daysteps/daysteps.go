package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	dataSl := strings.Split(data, ",")
	if len(dataSl) != 2 {
		return 0, 0, errors.New("slice length != 2")
	}
	stepsCount, err := strconv.Atoi(dataSl[0])
	if err != nil {
		return 0, 0, err
	}
	if stepsCount <= 0 {
		return 0, 0, errors.New("steps count <= 0")
	}
	dur, err := time.ParseDuration(dataSl[1])
	if err != nil {
		return 0, 0, err
	}
	if dur <= 0 {
		return 0, 0, errors.New("неверная продолжительность")
	}
	return stepsCount, dur, nil
}

func DayActionInfo(data string, weight, height float64) string {
	stepsCount, duration, err := parsePackage(data)
	if err != nil {
		// Не логируем здесь. Если нужно, логирует вызывающий код.
		return ""
	}

	calories, _ := spentcalories.WalkingSpentCalories(stepsCount, weight, height, duration)
	distance := (float64(stepsCount) * stepLength) / float64(mInKm)
	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		stepsCount, distance, calories,
	)
}
