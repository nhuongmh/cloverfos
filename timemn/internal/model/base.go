package model

import "time"

type Base struct {
	ID          uint64
	CreatedTime time.Time
	UpdatedTime time.Time
}
