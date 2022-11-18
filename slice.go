package slice

import (
	"errors"
	"fmt"
	"strings"
)

type Slice[T any] struct {
	array    *[]T // For edu purpose only
	length   int
	capacity int
}

func New[T any](elems ...T) Slice[T] {
	s := Make[T](len(elems))

	for i := 0; i < len(elems); i++ {
		s.Set(i, elems[i])
	}

	return s
}

// Make with only length (and same capacity): Make(length >= 0), with capacity too: Make(length >= 0, capacity >= length)
func Make[T any](size ...int) Slice[T] {
	length, capacity, err := extractMakeIndexes(size...)
	if err != nil {
		panic("slice.Make: " + err.Error())
	}

	array := make([]T, length, capacity)
	return Slice[T]{
		&array,
		length,
		capacity,
	}
}

func extractMakeIndexes(size ...int) (length, capacity int, err error) {
	// Check args count
	if len(size) == 0 {
		err = errors.New("missing size arguments")
		return
	}

	if len(size) > 2 {
		err = errors.New("too many arguments")
		return
	}

	// Init length & capacity
	length = size[0]
	capacity = length
	if len(size) == 2 {
		capacity = size[1]
	}

	// Check length & capacity
	if length < 0 {
		err = errors.New("negative slice length")
		return
	}

	if capacity < 0 {
		err = errors.New("negative slice capacity")
		return
	}

	if length > capacity {
		err = errors.New("slice length greater than capacity")
		return
	}
	return
}

func (s Slice[T]) IsNil() bool {
	return s.array == nil
}

func (s Slice[T]) Len() int {
	return s.length
}

func (s Slice[T]) Cap() int {
	return s.capacity
}

func (s Slice[T]) Get(idx int) T {
	if idx < 0 || idx >= s.Len() {
		panic("slice.Get: index out of range")
	}

	return (*s.array)[idx]
}

func (s Slice[T]) Set(idx int, val T) {
	if idx < 0 || idx >= len(*s.array) {
		panic("slice.Set: index out of range")
	}
	(*s.array)[idx] = val
}

// Sliced Strictly 2 cases: s[low:high] -> s.Sliced(low, high), s[low:high:maxCap] -> s.Sliced(low, high, maxCap)
func (s Slice[T]) Sliced(indexes ...int) Slice[T] {
	low, high, newCap, err := s.extractSlicedIndexes(indexes...)
	if err != nil {
		panic("slice.Sliced: " + err.Error())
	}

	array := (*s.array)[low:high:newCap]
	return Slice[T]{
		array:    &array,
		length:   high - low,
		capacity: newCap - low,
	}
}

func (s Slice[T]) extractSlicedIndexes(indexes ...int) (low, high, newCap int, err error) {
	if len(indexes) < 2 || len(indexes) > 3 {
		err = errors.New("invalid count of indexes")
		return
	}

	low = indexes[0]
	high = indexes[1]

	if low < 0 || high < 0 || low > high || high > s.Cap() {
		err = errors.New("index out of bound")
		return
	}

	newCap = s.Cap()
	if len(indexes) == 3 {
		newCap = indexes[2]
	}

	if newCap < high || newCap > s.Cap() {
		err = errors.New("index out of bound")
		return
	}

	return
}

func Append[T any](s Slice[T], elems ...T) Slice[T] {
	resLen := s.Len() + len(elems)

	var res Slice[T]
	if resLen <= s.Cap() {
		res = s.Sliced(0, resLen)
	} else {
		res = growSlice(s, resLen)
	}

	for i := s.Len(); i < resLen; i++ {
		res.Set(i, elems[i-s.Len()])
	}

	return res
}

func growSlice[T any](s Slice[T], newLen int) Slice[T] {
	newCap := nextSliceCapacity(newLen, s.Cap())

	newS := Make[T](newLen, newCap)
	for i := 0; i < s.Len(); i++ {
		newS.Set(i, s.Get(i))
	}

	return newS
}

func nextSliceCapacity(newLen, oldCap int) int {
	doubleCap := oldCap + oldCap
	if newLen > doubleCap {
		return newLen
	}

	const threshold = 1024

	// len(slice) x2 if < threshold
	if oldCap < threshold {
		return doubleCap
	}

	// len(slice) x1.25 if < threshold (or more if newLen > len(slice) x1.25)
	newCap := oldCap
	for {
		newCap += newCap / 4

		// uint for overflowing case
		if uint(newCap) >= uint(newLen) {
			break
		}
	}

	if newCap <= 0 {
		newCap = newLen
	}

	return newCap
}

func Copy[T any](dst, src Slice[T]) int {
	minLen := dst.Len()
	if src.Len() < minLen {
		minLen = src.Len()
	}

	for i := 0; i < minLen; i++ {
		dst.Set(i, src.Get(i))
	}

	return minLen
}

const (
	emptyOrNilSliceStr = "[]"
	elemDividerStr     = " "
)

func (s Slice[T]) String() string {
	if s.Len() == 0 {
		return emptyOrNilSliceStr
	}

	var sb strings.Builder
	sb.WriteString("[")
	for i := 0; i < s.Len()-1; i++ {
		sb.WriteString(fmt.Sprintf("%v", (*s.array)[i]))
		sb.WriteString(elemDividerStr)
	}
	sb.WriteString(fmt.Sprintf("%v", (*s.array)[s.Len()-1]))
	sb.WriteString("]")

	return sb.String()
}
