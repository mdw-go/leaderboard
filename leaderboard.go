package leaderboard

import (
	"cmp"
	"maps"
	"slices"
)

type Map[K comparable, V cmp.Ordered] map[K]V

func (this Map[K, V]) TopN(n int) []K {
	return slices.SortedStableFunc(maps.Keys(this),
		this.compare)[:min(n, len(this))]
}
func (this Map[K, V]) compare(i, j K) int {
	return -cmp.Compare(this[i], this[j])
}
