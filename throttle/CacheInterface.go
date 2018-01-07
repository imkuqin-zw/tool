package throttle

type CatchThrottle interface{
	Attempts(key string) (int, error)
	Has(key string) (bool, error)
	ResetAttempts(key string) error
	Hit(key string, decayMinutes int) int
	availableIn(key string) (int, error)
}

