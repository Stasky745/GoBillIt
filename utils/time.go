package utils

import "time"

func GetLastDayCurrentMonth() time.Time {
	// Get the current date
	now := time.Now()

	// Get the first day of the next month
	firstDayNextMonth := time.Date(
		now.Year(),
		now.Month()+1,
		1,
		0, 0, 0, 0,
		time.UTC,
	)

	// Subtract one day from the first day of the next month to get the last day of the current month
	return firstDayNextMonth.Add(-24 * time.Hour)
}

func GetLastDayNextMonth() time.Time {
	// Get the current date
	now := time.Now()

	// Get the first day of the next 2 month
	firstDayNextMonth := time.Date(
		now.Year(),
		now.Month()+2,
		1,
		0, 0, 0, 0,
		time.UTC,
	)

	// Subtract one day from the first day of the next month to get the last day of the current month
	return firstDayNextMonth.Add(-24 * time.Hour)
}

func FormatDate(d time.Time) string {
	return d.Format("January 2, 2006")
}

func GetCurrentMonthName() string {
	return time.Now().Month().String()
}
