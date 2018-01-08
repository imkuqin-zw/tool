package throttle

import (
	"net/http"
	"strings"
	"github.com/imkuqin-zw/tool/encoder"
	"strconv"
	"time"
	"fmt"
)

type Throttle interface {
	Handle() (bool, error)
}

type throttle struct {
	catch CatchThrottle
	maxAttempts int
	decayMinutes int
	req *http.Request
	resp http.ResponseWriter
}

func NewThrottle(req *http.Request, resp http.ResponseWriter, cache CatchThrottle, config ...int) Throttle {
	instance := &throttle{
		catch: cache,
		req: req,
		resp: resp,
		maxAttempts: 60,
		decayMinutes: 1,
	}
	if len(config) > 0 {
		instance.maxAttempts = config[0]
		if len(config) > 1 {
			instance.decayMinutes = config[1]
		}
	}
	return instance
}

func (t *throttle) Handle() (bool, error) {
	sign := t.getRequestSign()
	fmt.Println(sign)
	tooMany, err := t.tooManyAttempts(sign)
	if err != nil {
		return false, err
	}
	if tooMany {
		t.buildException(sign)
		return false, nil
	}

	if _, err = t.catch.Hit(sign, t.decayMinutes); err != nil {
		return false, err
	}
	attempts, err := t.retriesLeft(sign, 0)
	if err != nil {
		return false, err
	}
	t.addHead(attempts, 0)
	return true, nil
}

func (t *throttle) buildException(sign string) (err error) {
	retryAfter, err := t.catch.availableIn(sign)
	if err != nil {
		return
	}
	remainingAttempts, err := t.retriesLeft(sign, retryAfter)
	if err != nil {
		return
	}
	t.addHead(remainingAttempts, retryAfter)
	return
}

func (t *throttle) retriesLeft(sign string, retryAfter int) (int, error) {
	if retryAfter == 0 {
		attempts, err := t.catch.Attempts(sign)
		return t.maxAttempts - attempts, err
	} else {
		return 0, nil
	}
}

func (t *throttle) addHead(remainingAttempts, retryAfter int) {
	t.resp.Header().Set("X-RateLimit-Limit", strconv.Itoa(t.maxAttempts))
	t.resp.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remainingAttempts))
	if retryAfter != 0 {
		t.resp.Header().Set("Retry-After", strconv.Itoa(retryAfter))
		availableAt := int(time.Now().Unix()) + retryAfter
		t.resp.Header().Set("X-RateLimit-Reset", strconv.Itoa(availableAt))
		t.resp.WriteHeader(http.StatusTooManyRequests)
	}
	return
}

func (t *throttle) getRequestSign() (sign string) {
	ip := strings.Split(t.req.RemoteAddr, ":")[0]
	domain := strings.Replace(t.req.Host, "http://", "", -1)
	domain = strings.Replace(domain, "https://", "", -1)
	sign, _ = encoder.Sha1(domain + "|" + ip)
	return
}

func (t *throttle) tooManyAttempts(sign string) (bool, error) {
	count, err := t.catch.Attempts(sign)
	if err != nil {
		return false, err
	}
	if count >= t.maxAttempts {
		exist, err := t.catch.Has(sign + ":timer")
		if  err != nil {
			return false, err
		}
		if exist {
			return true, nil
		}
		t.catch.ResetAttempts(sign)
	}
	return false, nil
}