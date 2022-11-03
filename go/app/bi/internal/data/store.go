package data

import "time"

// Click
type Click struct {
	Alias     string    `json:"alias"      bson:"alias"`
	IP        string    `json:"ip"         bson:"ip"`
	UA        string    `json:"ua"         bson:"ua"`
	Referer   string    `json:"referer"    bson:"referer"`
	Country   string    `json:"country"    bson:"country"`
	City      string    `json:"city"       bson:"city"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// IBIStore interface
type IBIStore interface {
	Open() error
	Close() error
	Create(click *Click) error
	Clicks(alias string) (uint64, error)
}
