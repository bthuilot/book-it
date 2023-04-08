package reserve

import (
	"fmt"
	"github.com/bthuilot/book-it/pkg/resy"
	"github.com/bthuilot/book-it/pkg/util"
	"github.com/sirupsen/logrus"
	"time"
)

type Booker interface {
	// TODO(change return to internal type)
	Book(venueID int, partySize int, date time.Time, spread time.Duration, includeTypes ...string) (resy.ReservationSlot, error)
}

type booker struct {
	client resy.APIClient
}

func setHour(date time.Time, hour int) time.Time {
	year, month, day := date.Date()
	return time.Date(year, month, day, hour, 0, 0, 0, date.Location())
}

func (r booker) Book(venueID int, partySize int, date time.Time, spread time.Duration, includeTypes ...string) (resy.ReservationSlot, error) {
	logrus.Infof("booking reservation for venue %d on %s", venueID, date)
	venue, err := r.client.GetVenue(venueID)
	if err != nil {
		logrus.Errorf("error while retrieving venue with ID %d: %s", venueID, err)
		return resy.ReservationSlot{}, err
	}

	reservations, err := r.client.FindReservation(venue, date, partySize)
	if err != nil {
		logrus.Errorf("error while retrieving reservations for venue %s at %s: %s", venue.Name, date, err)
		return resy.ReservationSlot{}, err
	}
	// TODO(change to use response venue)
	venueSlots := getVenueSlots(reservations, venueID)
	var slots []resy.ReservationSlot
	for _, slot := range venueSlots {
		if timeDiff := AbsDuration(slot.Date.Start.Sub(date)); len(includeTypes) != 0 && !util.ListContains(includeTypes, slot.Config.Type) {
			logrus.Debugf("skipping reservation as type %s is not in include types", slot.Config.Type)
			continue
		} else if timeDiff > spread {
			logrus.Debugf("skipping reservation at %s, as diff of %d is greater than %d exists", slot.Date, timeDiff.Seconds(), spread.Seconds())
			continue
		}
		slots = append(slots, slot)
	}
	for _, slot := range slots {
		if booked, bookingErr := r.bookReservation(partySize, date, slot.Config.Token); bookingErr != nil {
			logrus.Warnf("unable to book reservation: %s", bookingErr)
		} else if booked {
			logrus.Info("booked reservation!")
			return slot, nil
		}
	}
	return resy.ReservationSlot{}, fmt.Errorf("unable to book any reservation")
}

func getVenueSlots(found resy.FindResponse, venueID int) []resy.ReservationSlot {
	for _, venue := range found.Results.Venues {
		if venue.Venue.IDs.Resy == venueID {
			return venue.Slots
		}
	}
	return nil
}

func AbsDuration(duration time.Duration) time.Duration {
	if duration < 0 {
		return duration * -1
	}
	return duration
}

func (r booker) bookReservation(partySize int, day time.Time, configToken string) (bool, error) {
	details, err := r.client.GetReservationDetails(partySize, day, configToken)
	if err != nil {
		return false, err
	}
	res, err := r.client.BookReservation(details.BookToken.Value)
	if err != nil {
		return false, err
	}
	logrus.Infof("reservation token: %s", res.ResyToken)
	return res.ResyToken != "", err
}

func NewBooker(client resy.APIClient) Booker {
	return booker{
		client: client,
	}
}
