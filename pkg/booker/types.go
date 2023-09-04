package booker

import (
	"time"
)

// Request represents a request to book
// a reservation using Resy API
type Request struct {
	// VenueID is the ID for the venue to book
	VenueID int
	// PartySize is the amount of people to make the reservation for
	PartySize int
	// Type is the type of reservation to make (i.e. Booth, Bar, Outdoor, Indoor, etc.)
	// This value must exactly match the one listed on the Resy web UI
	Type string

	// TODO(refactor section below)

	// Date is the date and time to book
	Date time.Time
	// TimeSpread is how much spread to have when booking the time for the reservation
	// i.e. if the Date field's time is set to 10:00pm, and the TimeSpread is 30 minutes,
	// the booker will try to book anywhere from 9:30 to 10:30pm on the given date
	TimeSpread time.Duration
}

type Reservation struct {
}
