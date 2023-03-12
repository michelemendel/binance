package utils

import "time"

const TimestampFormatLayout = "2006-01-02T15:04:05-07:00 MST"

// type Timestamp string

func StampTimeNow() string {
	// return time.Now().Format(TimestampFormatLayout)
	return time.Now().Format(time.RFC3339Nano)
}

func Now() time.Time {
	return time.Now()
}

func ToTime(t int64) string {
	return time.Unix(int64(t/1000), 0).Format(time.RFC3339Nano)
}

func ToTimeMillis() int64 {
	return time.Now().UnixNano() / 1e6
}
