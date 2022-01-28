package model

import "time"

type Office struct {
	ID            string       `json:"id" bson:"id"`
	LightsOnTime  time.Time `json:"lights_on_time" bson:"lights_on_time"`
	LightsOffTime time.Time `json:"lights_off_time" bson:"lights_off_time"`
}
