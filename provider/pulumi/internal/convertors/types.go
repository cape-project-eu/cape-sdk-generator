package convertors

import "time"

func ConvertAny(in any) any {
	return in
}

func ConvertTimeToString(in time.Time) string {
	return in.Format(time.RFC3339)
}

func ConvertStringToTime(in string) time.Time {
	t, _ := time.Parse(time.RFC3339, in)
	return t
}

func ConvertIntToInt64(in int) int64 {
	return int64(in)
}

func ConvertInt64ToInt(in int64) int {
	return int(in)
}
