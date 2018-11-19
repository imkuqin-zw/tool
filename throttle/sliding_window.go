package throttle

import (
	"net/http"
	"time"
)

type SlidingWindow struct {
	catch       CatchSlidingWindow
	maxAttempts int //最大限制数量
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
			_, err = c.cleanWinds(key, c.getWindForClean(curWind, c.count-1, c.count))
			if err != nil {
				return false, err
			}
			return false, nil
		}
		lastWind := (lastTs / c.interval) % c.count
		winds := c.getWindForClean(curWind, lastWind, c.count)
		cur := 0
		if len(winds) > 0 {
			if cur, err = c.cleanWinds(key, winds); err != nil {
				return false, err
			}
		}
		if cur+1 > c.maxAttempts {
			return true, nil
		}
	}
	return false, nil
}

func (c *SlidingWindow) getWindForClean(curWind, lastWind, count int64) []interface{} {
	winds := make([]interface{}, 0)
	if curWind > lastWind {
		temp := curWind
		for temp <= lastWind {
			winds = append(winds, temp)
		}
	} else if curWind < lastWind {
		var temp int64
		for temp <= curWind {
			winds = append(winds, temp)
		}
		temp = lastWind + 1
		for temp < count {
			winds = append(winds, temp)
		}
	}
	return winds
}

func (c *SlidingWindow) cleanWinds(key string, winds []interface{}) (int, error) {
	decr, err := c.catch.DelWind(key, winds...)
	if err != nil {
		return 0, err
	}
	return c.catch.DecrCount(key, decr)
}
