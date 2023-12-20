package services

import "time"

const (
	MaxDurationFastQueryDB = time.Second * 2
	MaxDurationSlowQueryDB = time.Second * 5
)
