package util

import (
	"fmt"
	"time"
)

func AddTime(s string, addDay int64) (string, error) {
	a, err := time.Parse("20060102", s)
	if err != nil {
		return "", err
	}

	a = a.Add(time.Duration(addDay) * 24 * time.Hour)
	result := fmt.Sprintf("%.4d%.2d%.2d", a.Year(), a.Month(), a.Day())

	return result, nil
}

func TimeCompare(i, j string) (bool, error) {
	a, err := time.Parse("20060102", i)
	if err != nil {
		return false, err
	}

	b, err := time.Parse("20060102", j)
	if err != nil {
		return false, err
	}

	return a.After(b), nil
}
