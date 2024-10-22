package util

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
)

func IntArrayToString(arr []int) string {
	return strings.Trim(strings.Join(strings.Split(fmt.Sprint(arr), " "), ","), "[]")
}

func ValidateActiveDays(activeDays []int) error {
	if len(activeDays) == 0 {
		return apperror.ErrEmptyActiveDays
	}
	for _, dayId := range activeDays {
		if dayId < appconst.MinDayId || dayId > appconst.MaxDayId {
			return apperror.ErrInvalidDayId
		}
	}
	return nil
}

func ValidateHourFormat(hour string) error {
	const regexPattern = "^([0-1]?[0-9]|2[0-3]):[0-5][0-9][+-][0-9]?[0-9]$"
	ok, _ := regexp.MatchString(regexPattern, hour)
	if !ok {
		return apperror.ErrInvalidHourFormat
	}
	return nil
}

func ValidateOperationalHourFormat(startHour, stopHour string) error {
	if err := ValidateHourFormat(startHour); err != nil {
		return err
	}
	if err := ValidateHourFormat(stopHour); err != nil {
		return err
	}
	return nil
}

func ValidatePartnerRequest(activeDays []int, startHour, stopHour string) error {
	if err := ValidateActiveDays(activeDays); err != nil {
		return err
	}
	if err := ValidateOperationalHourFormat(startHour, stopHour); err != nil {
		return err
	}
	return nil
}
