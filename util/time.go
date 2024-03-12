package util

import "time"

func TimeNowAsString() string {
	return time.Now().Format(time.RFC3339Nano)
}

func Time2String(t uint64) string {
	return time.Unix(int64(t/1000), 0).Format(time.RFC3339Nano)
}

func TimeNowInMillis() int64 {
	return time.Now().UnixNano() / 1e6
}
