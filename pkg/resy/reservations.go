package resy

import "time"

//
//type ReservationBooker interface {
//	BookAtMidnight(request ReservationRequest) (bool, error)
//}
//
//type reservationBooker struct {
//	client APIClient
//}
//
//func setHour(date time.Time, hour int) time.Time {
//	year, month, day := date.Date()
//	return time.Date(year, month, day, hour, 0, 0, 0, date.Location())
//}
//
//func (r reservationBooker) BookAtMidnight(request ReservationRequest) (bool, error) {
//	tomorrow := setHour(time.Now().Add(time.Hour*24), 0)
//	bookDate := tomorrow.Add(time.Hour * time.Duration(24*(request.DaysInAdvance)))
//	bookTarget := setHour(bookDate, request.Hour)
//	logrus.Infof("waiting till %s to book reservation for %s", tomorrow.String(), bookTarget.String())
//	//<-time.NewTimer(time.Until(tomorrow)).C // Wait
//	return r.findAndBookReservation(bookTarget, request, 50)
//}
//
//func getVenueSlots(found FindResponse, venueID int) []ReservationSlot {
//	for _, venue := range found.Results.Venues {
//		if venue.Venue.IDs.Resy == venueID {
//			return venue.Slots
//		}
//	}
//	return nil
//}
//
//func AbsDuration(duration time.Duration) time.Duration {
//	if duration < 0 {
//		return duration * -1
//	}
//	return duration
//}
//
//func (r reservationBooker) findAndBookReservation(day time.Time, request ReservationRequest, retryLimit int) (bool, error) {
//	logrus.Info("retrieving venue information")
//	venue, err := r.client.GetVenue(request.VenueId)
//	if err != nil {
//		return false, err
//	}
//	var (
//		innerErr   error
//		found      FindResponse
//		venueSlots []ReservationSlot
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
//		if AbsDuration(slot.Date.Start.Sub(day)) > time.Hour {
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
//
//func (r reservationBooker) bookReservation(partySize int, day time.Time, configToken string) (bool, error) {
//	details, err := r.client.GetReservationDetails(partySize, day, configToken)
//	if err != nil {
//		return false, err
//	}
//	res, err := r.client.BookReservations(details.BookToken.Value)
//	if err != nil {
//		return false, err
//	}
//	logrus.Infof("reservation token: %s", res.ResyToken)
//	return res.ResyToken != "", err
//}
//
//func NewReservationBooker(client APIClient) ReservationBooker {
//	return reservationBooker{
//		client: client,
//	}
//}
//
type ReservationTime struct {
	time.Time
}

//
//type ReservationRequest struct {
//	VenueId       int    `mapstructure:"venue_id"`
//	DaysInAdvance int    `mapstructure:"days_in_advance"`
//	Type          string `mapstructure:"type"`
//	PartySize     int    `mapstructure:"party_size"`
//	Hour          int    `mapstructure:"hour"`
//}
