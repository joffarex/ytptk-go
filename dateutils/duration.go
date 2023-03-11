package dateutils

import (
	"math"
	"strconv"
)

func formatTimePart(part float64) string {
	parsed := strconv.FormatInt(int64(part), 10)

	if parsed == "0" {
		return "00"
	} else if part < 10 {
		return "0" + parsed
	} else {
		return parsed
	}
}

func SecondsToDuration(seconds int) string {
	tempSeconds := float64(seconds)

	hours := math.Floor(tempSeconds / 3600)
	tempSeconds = tempSeconds - hours*3600
	minutes := math.Floor(tempSeconds / 60)
	tempSeconds = tempSeconds - minutes*60

	if minutes == 0 {
		return "00:00:" + formatTimePart(tempSeconds)
	} else if hours == 0 {
		return "00:" + formatTimePart(minutes) + ":" + formatTimePart(tempSeconds)
	} else {
		return formatTimePart(hours) + ":" + formatTimePart(minutes) + ":" + formatTimePart(tempSeconds)
	}

}
