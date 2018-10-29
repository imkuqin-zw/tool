package throttle

import (
	"net/http"
	"time"
)

type SlidingWindow struct {
	catch       CatchSlidingWindow
	maxAttempts int //最大限制数量
	interval    int //窗口时间间隔(s)
	count       int //窗口数量
	req         *http.Request
	resp        http.ResponseWriter
}

func (c *SlidingWindow) Handle() (bool, error) {
	now := time.Now().Second()
	curWind := (now / c.interval) % c.count
	for i := 1; i < c.count; i++ {

	}
}
