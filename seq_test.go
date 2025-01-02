package seq_test

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"strconv"
	"testing"

	"github.com/heliorosa/seq"
	"github.com/stretchr/testify/require"
)

var (
	S     = []int{1, 2, 42, 99}
	sSeq  = slices.Values(S)
	sSum  = 1 + 2 + 42 + 99
	sProd = 1 * 2 * 42 * 99
	sStr  = []string{"1", "2", "42", "99"}
	SMap  = map[int]int{0: 1, 1: 2, 2: 42, 3: 99}
	M     = map[int]string{1: "one", 2: "two", 7: "seven"}
	mSeq  = maps.All(M)
	MStr  = map[string]string{"1": "one", "2": "two", "7": "seven"}
)

func TestFilter(t *testing.T) {
	require.Equal(
		t,
		S[:2],
		slices.Collect(seq.Filter(sSeq, func(v int) bool { return v < 10 })),
	)
}

func TestFilter2(t *testing.T) {
	require.Equal(
		t,
		M,
		maps.Collect(seq.Filter2(mSeq, func(k int, v string) bool { return k < 10 })),
	)
}

func TestMap(t *testing.T) {
	require.Equal(
		t,
		sStr,
		slices.Collect(seq.Map(sSeq, func(v int) string { return strconv.Itoa(v) })),
	)
}

func TestMap2(t *testing.T) {
	require.Equal(
		t,
		MStr,
		maps.Collect(seq.Map2(mSeq, func(k int, v string) (string, string) { return strconv.Itoa(k), v })),
	)
}

func TestReduce(t *testing.T) {
	require.Equal(t, sSum, seq.Reduce(sSeq, 0, func(ac int, v int) int { return ac + v }))
}

func TestReduce2(t *testing.T) {
	toString := func(k int, v string) string { return fmt.Sprintf("%d - %s", k, v) }
	r := seq.Reduce2(mSeq, []string{}, func(ac []string, k int, v string) []string {
		return append(ac, toString(k, v))
	})
	slices.Sort(r)
	e := make([]string, 0, len(M))
	for k, v := range M {
		e = append(e, toString(k, v))
	}
	slices.Sort(e)
	require.Equal(t, e, r)
}

func TestFold(t *testing.T) {
	require.Equal(t, sSum, seq.Fold(sSeq, func(ac int, v int) int { return ac + v }))
}

func TestKeys(t *testing.T) {
	e := slices.Collect(maps.Keys(M))
	slices.Sort(e)
	got := slices.Collect(seq.Keys(mSeq))
	slices.Sort(got)
	require.Equal(t, e, got)
}

func TestValues(t *testing.T) {
	e := slices.Collect(maps.Values(M))
	slices.Sort(e)
	got := slices.Collect(seq.Values(mSeq))
	slices.Sort(got)
	require.Equal(t, e, got)
}

func TestEnumerate(t *testing.T) {
	require.Equal(t, SMap, maps.Collect(seq.Enumerate(sSeq)))
}

func TestCollectSlice(t *testing.T) {
	require.Equal(t, S, seq.CollectSlice(sSeq, 16))
}

func TestCollectMap(t *testing.T) {
	require.Equal(t, SMap, seq.CollectMap(maps.All(SMap), 16))
}

func TestMin(t *testing.T) {
	m, ok := seq.Min(sSeq)
	require.True(t, ok)
	require.Equal(t, 1, m)
}

func TestMax(t *testing.T) {
	m, ok := seq.Max(sSeq)
	require.True(t, ok)
	require.Equal(t, 99, m)
}

func TestSum(t *testing.T) {
	require.Equal(t, sSum, seq.Sum(sSeq))
}

func TestProduct(t *testing.T) {
	require.Equal(t, sProd, seq.Product(sSeq))
}

func iterLen[T any](seq iter.Seq[T]) int {
	r := 0
	for range seq {
		r++
	}
	return r
}

func iterLen2[K, V any](seq iter.Seq2[K, V]) int {
	r := 0
	for range seq {
		r++
	}
	return r
}

func TestSkip(t *testing.T) {
	require.Equal(t, len(S)-2, iterLen(seq.Skip(sSeq, 2)))
}

func TestSkip2(t *testing.T) {
	require.Equal(t, len(M)-2, iterLen2(seq.Skip2(mSeq, 2)))
}

func TestLimit(t *testing.T) {
	require.Equal(t, 2, iterLen(seq.Limit(sSeq, 2)))
}

func TestLimit2(t *testing.T) {
	require.Equal(t, 2, iterLen2(seq.Limit2(mSeq, 2)))
}
