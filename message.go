package townsita

import (
	"time"
)

const (
	MessageDraft = iota
	MessagePublished
	MessageDisabled
)

type MessageType struct {
	ID    string
	Title string
}

type Message struct {
	ID     string
	UserID string
	TypeID string

	Readers    int // Designated readers
	Completed  int // Number of times seen
	TargetHash string
	Latitude   float64
	Longitude  float64
	Radius     int

	Headline string
	Body     string
	Photo    string

	Status    int
	CreatedAd time.Time
	UpdatedAd time.Time
}
