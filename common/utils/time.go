package utils

import "time"

func init() {
	DefaultTimeService = new(SimpleTimeService)
}

var DefaultTimeService *SimpleTimeService

type TimeService interface {
	Now() time.Time
}

type SimpleTimeService struct{}

func (t *SimpleTimeService) Now() time.Time {
	return time.Now().UTC()
}

func (t *SimpleTimeService) Unix() int64 {
	return t.Now().Unix()
}

func (t *SimpleTimeService) UnixUint64() uint64 {
	return uint64(t.Unix())
}

func (t *SimpleTimeService) Nano() int64 {
	return t.Now().UnixNano()
}
