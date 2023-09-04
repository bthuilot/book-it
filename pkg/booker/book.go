package booker

import (
	"fmt"
	"github.com/bthuilot/book-it/pkg/resy"
	"github.com/bthuilot/book-it/pkg/util"
	"github.com/sirupsen/logrus"
	"time"
)

func (r *booker) Book(venueID int, partySize int, date time.Time, spread time.Duration, includeTypes ...string) (resy.ReservationSlot, error) {
	r.logger.Infof("booking reservation for venue %d on %s", venueID, date)
	venue, err := r.client.GetVenue(venueID)
	if err != nil {
		logrus.Errorf("error while retrieving venue with ID %d: %s", venueID, err)
		return resy.ReservationSlot{}, err
	}

	// TODO: can be done within the GetVenue method, since it can also return slots
	reservations, err := r.client.FindReservation(venue, date, partySize)
	if err != nil {
		r.logger.Errorf("error while retrieving reservations for venue %s at %s: %s", venue.Name, date, err)
		return resy.ReservationSlot{}, err
	}

	venueSlots := extractVenueSlots(reservations, venueID)

	slots := r.findApplicableSlots(venueSlots, date, spread, includeTypes...)

	for _, slot := range slots {
		if booked, bookingErr := r.bookReservation(partySize, date, slot.Config.Token); bookingErr != nil {
			r.logger.Warnf("unable to book reservation: %s", bookingErr)
		} else if booked {
			r.logger.Info("booked reservation!")
			return slot, nil
		}
	}
	return resy.ReservationSlot{}, fmt.Errorf("unable to book any reservation")
}

func extractVenueSlots(found resy.Search, venueID int) []resy.ReservationSlot {
	for _, venue := range found.Results.Venues {
		if venue.Venue.IDs.Resy == venueID {
			return venue.Slots
		}
	}
	return nil
}

func (r *booker) findApplicableSlots(venueSlots []resy.ReservationSlot, date time.Time, spread time.Duration, types ...string) (slots []resy.ReservationSlot) {
	for _, slot := range venueSlots {
		if timeDiff := AbsDuration(slot.Date.Start.Sub(date)); len(types) != 0 && !util.ListContains(types, slot.Config.Type) {
			r.logger.Debugf("skipping reservation as type %s is not in include types", slot.Config.Type)
			continue
		} else if timeDiff > spread {
			r.logger.Debugf("skipping reservation at %s, as diff of %d is greater than %d exists", slot.Date, timeDiff.Seconds(), spread.Seconds())
			continue
		}
		slots = append(slots, slot)
	}
	return
}

func (r *booker) bookReservation(partySize int, day time.Time, configToken string) (bool, error) {
	details, err := r.client.GetReservationDetails(partySize, day, configToken)
	if err != nil {
		return false, err
	}

	var paymentID int
	for _, payment := range details.User.PaymentMethods {
		if payment.IsDefault {
			paymentID = payment.ID
			break
		}
	}

	res, err := r.client.BookReservation(details.BookToken.Value, details.Cancellation.Fee != nil, paymentID)
	if err != nil {
		return false, err
	}
	r.logger.Infof("reservation token: %s", res.ResyToken)
	return res.ResyToken != "", err
}

func AbsDuration(duration time.Duration) time.Duration {
	if duration < 0 {
		return duration * -1
	}
	return duration
}
