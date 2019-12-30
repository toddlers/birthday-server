package utils

import (
	"regexp"
	"time"
)

func IsUsernameCorrect(name string) bool {
	isAlpha := regexp.MustCompile(`^[A-Za-z]+$`).MatchString
	return isAlpha(name)
}

func IsBirthdayCorrect(birthday string) bool {
	born, _ := time.Parse("2006-01-02", birthday)
	today := time.Now()
	return born.Before(today)
}

func CalculateDays(birthday string) int {
	//yyy-mm-dd
	born, _ := time.Parse("2006-01-02", birthday)
	today := time.Now()
	if (today.Month() == born.Month()) && (today.Day() == born.Day()) {
		return 0
	}
	bd := time.Date(today.Year(), born.Month(), born.Day(), 0, 0, 0, 0, time.UTC)
	if bd.Before(today) {
		bd = time.Date(today.Year()+1, born.Month(), born.Day(), 0, 0, 0, 0, time.UTC)
	}
	days := int(bd.Sub(today).Hours() / 24)
	return days
}
