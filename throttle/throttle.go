package throttle

import (
	"net/http"
	"strings"
	"github.com/imkuqin-zw/tool/encoder"
	"strconv"
	"time"
)


type Throttle struct {
	catch CatchThrottle
	maxAttempts int
	decayMinutes int
	req *http.Request
	resp http.ResponseWriter
}

func NewThrottle(req *http.Request, resp http.ResponseWriter) *Throttle {
	return &Throttle{
		req: req,
		resp: resp,
		maxAttempts: 60,
		decayMinutes: 1,
	}
}

func (t *Throttle) Handle() bool {
	sign := t.getRequestSign()
	if t.tooManyAttempts(sign) {
		t.buildException(sign)
		return false
	}

	t.catch.Hit(sign, t.decayMinutes)
	t.AddHead(t.retriesLeft(sign, 0), 0)
	return true
}

func (t *Throttle) buildException(sign string) {
	retryAfter := t.catch.availableIn(sign)
	remainingAttempts := t.retriesLeft(sign, retryAfter)
	t.AddHead(remainingAttempts, retryAfter)
}

func (t *Throttle) retriesLeft(sign string, retryAfter int) int {
	if retryAfter == 0 {
		attempts := t.catch.Attempts(sign)
		return t.maxAttempts - attempts
	} else {
		return 0
	}
}

func (t *Throttle) AddHead(remainingAttempts, retryAfter int) {
	t.resp.Header().Set("X-RateLimit-Limit", strconv.Itoa(t.maxAttempts))
	t.resp.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remainingAttempts))
	if retryAfter != 0 {
		t.resp.WriteHeader(http.StatusTooManyRequests)
		t.resp.Header().Set("Retry-After", strconv.Itoa(retryAfter))
		availableAt := int(time.Now().Unix()) + retryAfter
		t.resp.Header().Set("X-RateLimit-Reset", strconv.Itoa(availableAt))
	}
	return
}

func (t *Throttle) getRequestSign() (sign string) {
	ip := strings.Split(t.req.RemoteAddr, ":")[0]
	domain := strings.Replace(t.req.Host, "http://", "", -1)
	domain = strings.Replace(domain, "https://", "", -1)
	sign, _ = encoder.Sha1(domain + "|" + ip)
	return
}

func (t *Throttle) tooManyAttempts(sign string) bool {
	if t.catch.Attempts(sign) > t.maxAttempts {
		if t.catch.Has(sign + ":timer") {
			return true
		}
		t.catch.resetAttempts(sign)
	}
	return false
}