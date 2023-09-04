package booker

import (
	"github.com/bthuilot/book-it/pkg/resy"
	"github.com/sirupsen/logrus"
	"time"
)

// Booker is an interface which can be used to book reservations
type Booker interface {
	// Book will book a reservation for the given venue, party size, and date
	// spread is the amount of seconds to "pad" the reservation selection time by
	// (i.e. if spread is 60, the reservation will be booked at a time between
	// date - 30 seconds  and date + 30 seconds). includeTypes is a list of
	// reservation types to include in the search (i.e. "standard", "bar", etc.)
	Book(venueID int, partySize int, date time.Time, spread time.Duration, includeTypes ...string) (resy.ReservationSlot, error)
}

// booker is the internal implementation of the Booker interface
type booker struct {
	client resy.Client
	logger *logrus.Logger
}

// NewBooker will create a new booker with the given client and options
// (see the Opt interface for available options)
func NewBooker(client resy.Client, options ...Opt) Booker {
	b := &booker{
		client: client,
		logger: logrus.New(),
	}

	b.logger.SetLevel(logrus.DebugLevel)

	for _, opt := range options {
		opt.apply(b)
	}
	return b
}

/***********
 * Options *
 ***********/

// Opt is an interface which can be used to set options on the booker
// (see the With* functions for available options)
type Opt interface {
	apply(booker *booker)
}

type withLogLevel struct {
	level logrus.Level
}

func (w withLogLevel) apply(booker *booker) {
	booker.logger.SetLevel(w.level)
}

// WithLogLevel will set the log level for the booker
func WithLogLevel(level logrus.Level) Opt {
	return withLogLevel{level: level}
}
