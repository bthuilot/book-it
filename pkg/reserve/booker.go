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

//
//func (r reservationBooker) findAndBookReservation(day time.Time, request ReservationRequest, retryLimit int) (bool, error) {
//	logrus.Info("retrieving venue information")
//	venue, err := r.client.GetVenue(request.VenueId)
//	if err != nil {
//		return false, err
//	}
//	var (
//		innerErr   error
//		found      findResponse
//		venueSlots []reservationSlot
//	)
//	for retryLimit > 0 {
//		retryLimit--
//		logrus.Infof("attempt %d", retryLimit)
//		logrus.Info("finding reservations for venue")
//		if found, innerErr = r.client.FindReservation(venue, day, request.PartySize); innerErr != nil {
//			logrus.Errorf("unable to retry reservation: %s")
//			innerErr = nil
//			continue
//		}
//		if venueSlots = getVenueSlots(found, venue.IDs.Resy); len(venueSlots) != 0 {
//			logrus.Infof("found slots! %+v", venueSlots)
//			break
//		}
//		logrus.Info("no venue slots found")
//		time.Sleep(time.Second / 5)
//	}
//
//	// TODO(only doing this since i have no way of testing what the type is)
//	var typeAndTimeTokens []string
//	var timeTokens []string
//	var otherTokens []string
//
//	for _, slot := range venueSlots {
//		if AbsDuration(slot.Date.start.Sub(day)) > time.Hour {
//			otherTokens = append(otherTokens, slot.Config.Token)
//		}
//		if slot.Config.Type == request.Type {
//			typeAndTimeTokens = append(typeAndTimeTokens, slot.Config.Token)
//		} else {
//			timeTokens = append(timeTokens, slot.Config.Token)
//		}
//	}
//
//	for _, token := range append(typeAndTimeTokens, append(timeTokens, otherTokens...)...) {
//		if booked, err := r.bookReservation(request.PartySize, day, token); err != nil {
//			logrus.Warnf("unable to book reservation: %s", err)
//		} else if booked {
//			logrus.Info("booked reservation!")
//			return true, nil
//		}
//	}
//	return false, fmt.Errorf("no reservation booked")
//}

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
