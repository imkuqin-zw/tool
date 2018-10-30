package throttle

import (
	"net/http"
	"time"
)

type SlidingWindow struct {
	catch       CatchSlidingWindow
	maxAttempts int64 //最大限制数量
	duration    int64
	interval    int64 //窗口时间间隔(s)
	count       int64 //窗口数量
	req         *http.Request
	resp        http.ResponseWriter
}

func (c *SlidingWindow) Handle(key string) (bool, error) {
	now := time.Now().Unix()
	curWind := now / int64(c.interval) % c.count
	tooMany, err := c.tooManyAttempts(key, curWind, now)
	if err != nil {
		return false, err
	}
	if tooMany {
		//c.buildException(sign)
		return false, nil
	}

	if err := c.catch.Hit(key, curWind, now, c.duration); err != nil {
		return false, err
	}
	return true, nil
}

func (c *SlidingWindow) tooManyAttempts(key string, curWind, now int64) (bool, error) {
	lastTs, err := c.catch.GetLastTs(key)
	if err != nil {
		return false, err
	}
	if lastTs != 0 {
		elapsed := now - lastTs
		if elapsed >= c.duration {
			//if err = c.catch.Reset(key, curWind, now, c.duration); err != nil {
			//	return false, err
			//}
			return false, nil
		}

		lastWind := (lastTs / c.interval) % c.count
		if curWind > lastWind {

		}
	}
}
