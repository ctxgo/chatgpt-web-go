package aiclient

import (
	"math/rand"
	"sync"
	"time"
)

type scheduleType string

const RandomSchedule scheduleType = "random"
const RoundRobinSchedule scheduleType = "round_robin"

type Scheduler interface {
	NextIndex(max int) int
}

type RandomScheduler struct {
	rnd *rand.Rand
}

func NewRandomScheduler() *RandomScheduler {
	src := rand.NewSource(time.Now().UnixNano())
	return &RandomScheduler{
		rnd: rand.New(src),
	}
}

func (r *RandomScheduler) NextIndex(max int) int {
	return r.rnd.Intn(max)
}

type RoundRobinScheduler struct {
	index int
}

func NewRoundRobinScheduler() *RoundRobinScheduler {
	return &RoundRobinScheduler{}
}

func (r *RoundRobinScheduler) NextIndex(max int) int {
	defer func() { r.index = (r.index + 1) % max }()
	return r.index
}

type SchedulerFactory struct {
	schedulers map[string]Scheduler
	mu         sync.Mutex
}

func NewSchedulerFactory() *SchedulerFactory {
	return &SchedulerFactory{
		schedulers: make(map[string]Scheduler),
	}
}

// GetOrCreateScheduler 返回对应键的调度器，如果不存在则创建一个新的
func (f *SchedulerFactory) GetOrCreateScheduler(key string, schedulerType scheduleType) Scheduler {
	f.mu.Lock()
	defer f.mu.Unlock()

	if scheduler, exists := f.schedulers[key]; exists {
		return scheduler
	}

	var newScheduler Scheduler
	switch schedulerType {
	case "round_robin":
		newScheduler = NewRoundRobinScheduler()
	default:
		newScheduler = NewRandomScheduler()
	}

	f.schedulers[key] = newScheduler
	return newScheduler
}
