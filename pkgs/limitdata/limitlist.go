package limitdata

import (
	"sync"
)

const maxListSize = 60
const step = 1

type LimitLister[T any] interface {
	Add(T, ...func(T))
	Get() []T
	Select(...func(T))
}

type LimitList[T any] struct {
	data     []T
	maxLenth int
	step     int
	sync.Mutex
}

type opts[T any] func(*LimitList[T])

func SetMaxLenth[T any](n int) opts[T] {
	return func(l *LimitList[T]) {
		l.maxLenth = n
	}
}

func SetStep[T any](n int) opts[T] {
	return func(l *LimitList[T]) {
		l.step = n
	}
}

func NewLimitList[T any](otps ...opts[T]) *LimitList[T] {
	l := &LimitList[T]{
		maxLenth: maxListSize,
		step:     step,
	}
	for _, f := range otps {
		f(l)
	}
	return l
}

func (l *LimitList[T]) Add(a T, funcs ...func(T)) {
	l.Lock()
	defer l.Unlock()
	l.data = append(l.data, a)
	if len(l.data) >= l.maxLenth+1 {
		for _, f := range funcs {
			for _, i := range l.data[:step] {
				f(i)
			}
		}

		l.data = l.data[step:]
	}
}

func (l *LimitList[T]) Get() []T {
	return append([]T(nil), l.data...)
}

func (l *LimitList[T]) Select(funcs ...func(T)) {
	for _, f := range funcs {
		for _, i := range l.data {
			f(i)
		}
	}
}
