package redis

import "time"

// SetExParameter is the struct that defines parameters that are required for calling SetEx
type SetExParameter struct {
	//Data is the data that will be saved
	Data interface{}
	//ExpireDuration is the TTL of the key
	ExpireDuration time.Duration
	//IsTesting is the flag to define if the SetEx is called during the load test
	IsTesting bool
}
