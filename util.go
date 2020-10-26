package bonusly

import "time"

func toStringPtr(s string) *string {
	return &s
}

func fromStringPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func toIntPtr(i int) *int {
	return &i
}

func fromIntPtr(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func toBoolPtr(b bool) *bool {
	return &b
}

func fromBoolPtr(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func toTimePtr(t time.Time) *time.Time {
	return &t
}

func fromTimePtr(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
