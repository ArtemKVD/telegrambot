package limits

import (
	calc "telegrambot/internal/calculate"
)

type DailyLimits struct {
	Calories int
	Proteins int
	Fats     int
	Carbs    int
}

func Calculate(gender, weightStr, heightStr, program string) (DailyLimits, error) {

	var calories int
	switch program {
	case "lost":
		calories = calc.Kforlost(gender, weightStr, heightStr)
	case "set":
		calories = calc.Kforset(gender, weightStr, heightStr)
	case "get":
		calories = calc.Kforget(gender, weightStr, heightStr)
	default:
		calories = calc.Kforset(gender, weightStr, heightStr)
	}

	var proteins, fats, carbs int
	switch program {
	case "lost":
		proteins, fats, carbs = calc.Lost(calories)
	case "set":
		proteins, fats, carbs = calc.Set(calories)
	case "get":
		proteins, fats, carbs = calc.Get(calories)
	}

	return DailyLimits{
		Calories: calories,
		Proteins: proteins,
		Fats:     fats,
		Carbs:    carbs,
	}, nil
}
