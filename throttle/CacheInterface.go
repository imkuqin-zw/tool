package throttle

type CatchThrottle interface{
	Attempts(key string) int
	Has(key string) bool
	resetAttempts(key string)
	Hit(key string, decayMinutes int)
	availableIn(key string) int
}
