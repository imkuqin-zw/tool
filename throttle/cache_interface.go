package throttle

type CatchThrottle interface {
	Attempts(key string) (int, error)
	Has(key string) (bool, error)
	ResetAttempts(key string) error
	Hit(key string, decayMinutes int) (int, error)
	availableIn(key string) (int, error)
}

type CatchSlidingWindow interface {
	Attempts(key string) (int, error)
}
