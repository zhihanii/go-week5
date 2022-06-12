package go_week5

import (
	"sync"
	"time"
)

// Limiter 限流器
type Limiter struct {
	size  time.Duration // 窗口大小
	limit int64         // 最大限制数
	m     sync.Mutex

	curr Window
	prev Window
}

func NewLimiter(size time.Duration, limit int64) *Limiter {
	curr := NewWindow()
	prev := NewWindow()
	l := &Limiter{
		size:  size,
		limit: limit,
		curr:  curr,
		prev:  prev,
	}
	return l
}

func (l *Limiter) Size() time.Duration {
	return l.size
}

func (l *Limiter) Limit() int64 {
	l.m.Lock()
	defer l.m.Unlock()
	return l.limit
}

func (l *Limiter) SetLimit(limit int64) {
	l.m.Lock()
	defer l.m.Unlock()
	l.limit = limit
}

func (l *Limiter) Allow() bool {
	return l.AllowN(time.Now(), 1)
}

// AllowN 是否允许n个事件
func (l *Limiter) AllowN(now time.Time, n int64) bool {
	l.m.Lock()
	defer l.m.Unlock()

	l.advance(now)

	d := now.Sub(l.curr.Start())
	weight := float64(l.size-d) / float64(l.size)
	count := int64(weight*float64(l.prev.Count())) + l.curr.Count()

	if count+n > l.limit {
		return false
	}

	l.curr.AddCount(n)
	return true
}

func (l *Limiter) advance(now time.Time) {
	newCurrStart := now.Truncate(l.size)
	diffSize := newCurrStart.Sub(l.curr.Start()) / l.size
	if diffSize >= 1 {
		newPrevCount := int64(0)
		if diffSize == 1 {
			newPrevCount = l.curr.Count()
		}
		l.prev.Reset(newCurrStart.Add(-l.size), newPrevCount)
		l.curr.Reset(newCurrStart, 0)
	}
}
