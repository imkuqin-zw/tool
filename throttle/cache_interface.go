package throttle

type CatchThrottle interface {
	Attempts(key string) (int, error)
	Has(key string) (bool, error)
	ResetAttempts(key string) error
	Hit(key string, decayMinutes int) (int, error)
	availableIn(key string) (int, error)
}

type CatchSlidingWindow interface {
	GetLastTs(key string) (int64, error)
	Hit(key string, curWind, now, duration int64) error
}
