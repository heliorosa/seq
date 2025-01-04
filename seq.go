package seq

import (
	"cmp"
	"iter"
	"slices"
)

func Filter[T any](seq iter.Seq[T], filterFunc func(v T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range seq {
			if filterFunc(v) && !yield(v) {
				return
			}
		}
	}
}

func Filter2[K, V any](seq iter.Seq2[K, V], filterFunc func(k K, V V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			if filterFunc(k, v) && !yield(k, v) {
				return
			}
		}
	}
}

func Map[S, D any](seq iter.Seq[S], mapFunc func(v S) D) iter.Seq[D] {
	return func(yield func(D) bool) {
		for v := range seq {
			if !yield(mapFunc(v)) {
				return
			}
		}
	}
}

func Map2[SK, SV, DK, DV any](seq iter.Seq2[SK, SV], mapFunc func(k SK, v SV) (DK, DV)) iter.Seq2[DK, DV] {
	return func(yield func(DK, DV) bool) {
		for k, v := range seq {
			if !yield(mapFunc(k, v)) {
				return
			}
		}
	}
}

func Reduce[T, D any](seq iter.Seq[T], ac D, reduceFunc func(ac D, v T) D) D {
	r := ac
	for v := range seq {
		r = reduceFunc(r, v)
	}
	return r
}

func Reduce2[K, V, D any](seq iter.Seq2[K, V], ac D, reduceFunc func(ac D, k K, v V) D) D {
	r := ac
	for k, v := range seq {
		r = reduceFunc(r, k, v)
	}
	return r
}

func Fold[T any](seq iter.Seq[T], foldFunc func(ac T, v T) T) T {
	next, stop := iter.Pull(seq)
	defer stop()
	r, ok := next()
	if !ok {
		return r
	}
	for {
		v, ok := next()
		if !ok {
			return r
		}
		r = foldFunc(r, v)
	}
}

func Keys[K, V any](seq iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range seq {
			if !yield(k) {
				return
			}
		}
	}
}

func Values[K, V any](seq iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range seq {
			if !yield(v) {
				return
			}
		}
	}
}

func EnumerateFunc[K, V any](seq iter.Seq[V], enumFunc func(v V) K) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for v := range seq {
			if !yield(enumFunc(v), v) {
				return
			}
		}
	}
}

func Enumerate[T any](seq iter.Seq[T]) iter.Seq2[int, T] {
	counter := 0
	return EnumerateFunc(seq, func(v T) int {
		r := counter
		counter++
		return r
	})
}

func CollectSlice[T any](seq iter.Seq[T], size int) []T {
	return Reduce(
		seq,
		make([]T, 0, size),
		func(ac []T, v T) []T { return append(ac, v) },
	)
}

func CollectMap[K comparable, V any](seq iter.Seq2[K, V], size int) map[K]V {
	return Reduce2(
		seq,
		make(map[K]V, size),
		func(ac map[K]V, k K, v V) map[K]V {
			ac[k] = v
			return ac
		},
	)
}

func Min[T cmp.Ordered](seq iter.Seq[T]) T {
	return Fold(seq, func(ac, v T) T { return when(ac < v, ac, v) })
}

func Max[T cmp.Ordered](seq iter.Seq[T]) T {
	return Fold(seq, func(ac, v T) T { return when(ac > v, ac, v) })
}

type Numeric interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | uint |
		~int8 | ~int16 | ~int32 | ~int64 | int |
		~float32 | ~float64 | ~complex64 | ~complex128
}

type Addable interface{ Numeric | ~string }

func Sum[T Addable](seq iter.Seq[T]) T {
	return Fold(seq, func(ac, v T) T { return ac + v })
}

func Product[T Numeric](seq iter.Seq[T]) T {
	return Fold(seq, func(ac, v T) T { return ac * v })
}

func Skip[T any](seq iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		skipped := 0
		for v := range seq {
			if skipped < n {
				skipped++
				continue
			}
			if !yield(v) {
				return
			}
		}
	}
}

func Skip2[K, V any](seq iter.Seq2[K, V], n int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		skipped := 0
		for k, v := range seq {
			if skipped < n {
				skipped++
				continue
			}
			if !yield(k, v) {
				return
			}
		}
	}
}

func Limit[T any](seq iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		yielded := 0
		for v := range seq {
			if yielded == n || !yield(v) {
				return
			}
			yielded++
		}
	}
}

func Limit2[K, V any](seq iter.Seq2[K, V], n int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		yielded := 0
		for k, v := range seq {
			if yielded == n || !yield(k, v) {
				return
			}
			yielded++
		}
	}
}

func Find[T any](seq iter.Seq[T], f func(T) bool) (r T, ok bool) {
	var z T
	for v := range seq {
		if f(v) {
			return v, true
		}
	}
	return z, false
}

func Find2[K, V any](seq iter.Seq2[K, V], f func(k K, v V) bool) (k K, v V, ok bool) {
	var (
		zk K
		zv V
	)
	for k, v := range seq {
		if f(k, v) {
			return k, v, true
		}
	}
	return zk, zv, false
}

func Contains[T comparable](seq iter.Seq[T], v T) bool {
	_, ok := Find(seq, func(vv T) bool { return v == vv })
	return ok
}

func ContainsKey[K comparable, V any](seq iter.Seq2[K, V], key K) bool {
	_, _, ok := Find2(seq, func(k K, _ V) bool { return k == key })
	return ok
}

func ContainsValue[K any, V comparable](seq iter.Seq2[K, V], val V) bool {
	_, _, ok := Find2(seq, func(_ K, v V) bool { return v == val })
	return ok
}

func Flatten[T any](seq iter.Seq[iter.Seq[T]]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for s := range seq {
			for v := range s {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func Flatten2[K, V any](seq iter.Seq[iter.Seq2[K, V]]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for s := range seq {
			for k, v := range s {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

func Concat[T any](seqs ...iter.Seq[T]) iter.Seq[T] {
	return Flatten(slices.Values(seqs))
}

func Concat2[K, V any](seqs ...iter.Seq2[K, V]) iter.Seq2[K, V] {
	return Flatten2(slices.Values(seqs))
}

func Drain[T any](seq iter.Seq[T]) {
	for range seq {
	}
}

func Drain2[K, V any](seq iter.Seq2[K, V]) {
	for range seq {
	}
}

func Dedup[T comparable](seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		seen := make(map[T]struct{}, 16)
		for v := range seq {
			if _, ok := seen[v]; !ok {
				if !yield(v) {
					return
				}
				seen[v] = struct{}{}
			}
		}
	}
}

func GenerateFunc[T any](start T, nextFunc func(last T) T, continueFunc func(v T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		cur := start
		for continueFunc(cur) {
			if !yield(cur) {
				return
			}
			cur = nextFunc(cur)
		}
	}
}

func Generate(start int, stop int, step int) iter.Seq[int] {
	return GenerateFunc(
		start,
		func(last int) int { return last + step },
		func(v int) bool { return v < stop },
	)
}

func Repeat[T any](seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for {
			for v := range seq {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func Repeat2[K, V any](seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for {
			for k, v := range seq {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

func CountFunc[T any](seq iter.Seq[T], countFunc func(v T) bool) int {
	r := 0
	for v := range seq {
		if countFunc(v) {
			r++
		}
	}
	return r
}

func Count[T comparable](seq iter.Seq[T], value T) int {
	return CountFunc(seq, func(v T) bool { return value == v })
}

func CountFunc2[K, V any](seq iter.Seq2[K, V], countFunc func(k K, v V) bool) int {
	r := 0
	for k, v := range seq {
		if countFunc(k, v) {
			r++
		}
	}
	return r
}

func Len[T any](seq iter.Seq[T]) int {
	return CountFunc(seq, func(_ T) bool { return true })
}

func Len2[K, V any](seq iter.Seq2[K, V]) int {
	return CountFunc2(seq, func(_ K, _ V) bool { return true })
}

func CompareFunc[T any](seq1 iter.Seq[T], seq2 iter.Seq[T], compareFunc func(v1 T, v2 T) bool) (equal int, total int) {
	next, stop := iter.Pull(seq2)
	defer stop()
	for v1 := range seq1 {
		v2, ok := next()
		if !ok {
			return 0, 0
		}
		total++
		if compareFunc(v1, v2) {
			equal++
		}
	}
	return equal, total
}

func Compare[T comparable](seq1 iter.Seq[T], seq2 iter.Seq[T]) (equal int, total int) {
	return CompareFunc(seq1, seq2, func(v1 T, v2 T) bool { return v1 == v2 })
}

func ComparePercent[T comparable](seq1 iter.Seq[T], seq2 iter.Seq[T]) float64 {
	equal, total := Compare(seq1, seq2)
	return float64(equal) / float64(total)
}

func Equal[T comparable](seq1 iter.Seq[T], seq2 iter.Seq[T]) bool {
	equal, total := Compare(seq1, seq2)
	return equal == total
}

func SortFunc[T any](seq iter.Seq[T], cmpFunc func(a T, b T) int) iter.Seq[T] {
	return func(yield func(T) bool) {
		s := slices.Collect(seq)
		slices.SortFunc(s, cmpFunc)
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
	}
}

func Sort[T cmp.Ordered](seq iter.Seq[T]) iter.Seq[T] {
	return SortFunc(seq, func(a, b T) int {
		if a < b {
			return -1
		}
		return when(a > b, 1, 0)
	})
}

func AnyFunc[T any](seq iter.Seq[T], cmpFunc func(a T) bool) bool {
	for v := range seq {
		if cmpFunc(v) {
			return true
		}
	}
	return false
}

func Any[T comparable](seq iter.Seq[T], value T) bool {
	return AnyFunc(seq, func(v T) bool { return value == v })
}

func AnyFunc2[K, V any](seq iter.Seq2[K, V], cmpFunc func(k K, v V) bool) bool {
	for k, v := range seq {
		if cmpFunc(k, v) {
			return true
		}
	}
	return false
}

func AllFunc[T any](seq iter.Seq[T], cmpFunc func(a T) bool) bool {
	for v := range seq {
		if !cmpFunc(v) {
			return false
		}
	}
	return true
}

func All[T comparable](seq iter.Seq[T], value T) bool {
	return AnyFunc(seq, func(v T) bool { return value == v })
}

func AllFunc2[K, V any](seq iter.Seq2[K, V], cmpFunc func(k K, v V) bool) bool {
	for k, v := range seq {
		if !cmpFunc(k, v) {
			return false
		}
	}
	return true
}

func when[T any](cond bool, vTrue T, vFalse T) T {
	if cond {
		return vTrue
	}
	return vFalse
}
