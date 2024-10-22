package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
)

func TimeToCron(timeStr string) (string, error) {
	timeStr = strings.ReplaceAll(timeStr, "+", " +")
	parts := strings.Split(timeStr, " ")
	if len(parts) != 2 {
		return "", apperror.ErrInvalidHourFormat
	}

	t, err := time.Parse("15:04", parts[0])
	if err != nil {
		return "", apperror.ErrInternalServerError
	}

	cronSchedule := fmt.Sprintf("%d %d * * *", t.Minute(), t.Hour())
	return cronSchedule, nil
}
