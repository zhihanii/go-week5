package go_week5

import "time"

type Window interface {
	Start() time.Time
	Stop()
	Count() int64
	AddCount(n int64)
	Reset(s time.Time, c int64)
}

// SlideWindow 滑动窗口
type SlideWindow struct {
	start int64 // 开始时间
	count int64 // 计数器
}

func NewWindow() Window {
	return &SlideWindow{}
}

func (w *SlideWindow) Start() time.Time {
	return time.Unix(0, w.start)
}

func (w *SlideWindow) Count() int64 {
	return w.count
}

func (w *SlideWindow) AddCount(n int64) {
	w.count += n
}

func (w *SlideWindow) Reset(s time.Time, c int64) {
	w.start = s.UnixNano()
	w.count = c
}

func (w *SlideWindow) Stop() {

}
