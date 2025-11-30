package leaderboard

import (
	"reflect"
	"testing"
)

func TestLeaderboard(t *testing.T) {
	stats := Map[string, int]{
		"a": 2,
		"b": 3,
		"c": 1,
		"d": 4,
	}
	assertEqual(t, stats.TopN(3),
		[]string{
			"d",
			"b",
			"a",
		},
	)
	assertEqual(t, stats.TopN(5),
		[]string{
			"d",
			"b",
			"a",
			"c",
		},
	)
}

func assertEqual(t *testing.T, actual, expected any) {
	if reflect.DeepEqual(expected, actual) {
		return
	}
	t.Helper()
	t.Errorf("Expected %v, got %v", expected, actual)
}
