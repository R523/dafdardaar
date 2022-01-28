package model

import "time"

type Activity struct {
	User     string
	Office   string
	Datetime time.Time
	Type     string
}
