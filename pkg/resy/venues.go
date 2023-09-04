package resy

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

func (c *client) GetVenue(id int) (v Venue, err error) {
	data := url.Values{}
	data.Set("id", strconv.Itoa(id))
	reqURL, err := generateURL("3/venue", data)
	if err != nil {
		return
	}
	req, err := generateRequest("GET", reqURL, nil, nil)
	if err != nil {
		return
	}
	v, err = apiRequest[Venue](req)
	return
}

func (c *client) FindReservation(venue Venue, day time.Time, partySize int) (found Search, err error) {
	data := url.Values{}
	data.Set("lat", fmt.Sprintf("%f", venue.Location.Lat))
	data.Set("long", fmt.Sprintf("%f", venue.Location.Long))
	data.Set("party_size", strconv.Itoa(partySize))
	data.Set("day", day.Format("2006-01-02"))
	data.Set("venue_id", strconv.Itoa(venue.IDs.Resy))
	data.Set("limit", strconv.Itoa(resultLimit))
	reqURL, err := generateURL("4/find", data)
	if err != nil {
		return
	}
	additionalHeaders, err := c.generateAuthHeaders()
	if err != nil {
		return
	}
	req, err := generateRequest("GET", reqURL, nil, additionalHeaders)
	if err != nil {
		return
	}
	return apiRequest[Search](req)
}

/*********
 * Types *
 *********/

type Location struct {
	Lat  float64 `json:"latitude"`
	Long float64 `json:"longitude"`
}

type VenueID struct {
	Resy       int    `json:"resy"`
	Foursquare string `json:"foursquare"`
	Google     string `json:"google"`
}

type Venue struct {
	Name     string   `json:"name"`
	IDs      VenueID  `json:"id"`
	Location Location `json:"location"`
}

type SlotConfig struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

type ReservationDate struct {
	time.Time
}

func (t *ReservationDate) UnmarshalJSON(b []byte) (err error) {
	date, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(b), time.Local)
	if err != nil {
		return err
	}
	t.Time = date
	return
}

type ReservationDuration struct {
	Start ReservationDate
	End   ReservationDate
}

type ReservationSlot struct {
	Config SlotConfig          `json:"config"`
	Date   ReservationDuration `json:"date"`
}

type VenueSearchResult struct {
	Venue Venue             `json:"venue"`
	Slots []ReservationSlot `json:"slots"`
}

type VenueSearchResults struct {
	Venues []VenueSearchResult `json:"venues"`
}

type Search struct {
	Results    VenueSearchResults `json:"results"`
	GuestToken string             `json:"guest_token"`
}
