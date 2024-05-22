package core

import (
	"fmt"
	"strconv"
	"time"
)

func Time2Chunk(t time.Time) string {
	// chunk by 10 minutes
	return fmt.Sprintf("%d", (t.Unix()/600)*600)
}

func Chunk2RecentTime(chunk string) time.Time {
	i, _ := strconv.ParseInt(chunk, 10, 64)
	return time.Unix(i+600, 0)
}

func Chunk2ImmediateTime(chunk string) time.Time {
	i, _ := strconv.ParseInt(chunk, 10, 64)
	return time.Unix(i, 0)
}

func TypedIDToType(id string) string {
	if len(id) != 27 {
		return ""
	}
	prefix := id[0]
	switch prefix {
	case 'a':
		return "association"
	case 'm':
		return "message"
	default:
		return ""
	}
}

func hasChar(s string, c byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return true
		}
	}
	return false
}

func IsCKID(keyID string) bool {
	return len(keyID) == 42 && keyID[:3] == "cck" && !hasChar(keyID, '.')
}

func IsCCID(keyID string) bool {
	return len(keyID) == 42 && keyID[:3] == "con" && !hasChar(keyID, '.')
}