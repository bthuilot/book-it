package resy

import "time"

/*
 * Auth endpoints responses
 */

type authPhone struct {
	MobileNumber string `json:"mobile_number"`
	Method       int    `json:"method"`
	Sent         bool   `json:"sent"`
}

type challengeResponse struct {
	MobileClaim struct {
		MobileNumber string `json:"mobile_number"`
		ClaimToken   string `json:"claim_token"`
		DateExpires  string `json:"date_expires"`
	} `json:"mobile_claim"`
	Challenge struct {
		ChallengeId string `json:"challenge_id"`
		FirstName   string `json:"first_name"`
		Message     string `json:"message"`
		Properties  []struct {
			Name    string `json:"name"`
			Type    string `json:"type"`
			Message string `json:"message"`
		} `json:"properties"`
	} `json:"challenge"`
}

type authResponse struct {
	Token string `json:"token"`
}

/*
 * Venue endpoints responses
 */

type Location struct {
	Lat  float64 `json:"latitude"`
	Long float64 `json:"longitude"`
}

type VenueIDs struct {
	Resy       int    `json:"resy"`
	Foursquare string `json:"foursquare"`
	Google     string `json:"google"`
}

type Venue struct {
	IDs      VenueIDs `json:"id"`
	Location Location `json:"location"`
}

type slotConfig struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

type reservationDate struct {
	time.Time
}

func (t *reservationDate) UnmarshalJSON(b []byte) (err error) {
	date, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(b), time.Local)
	if err != nil {
		return err
	}
	t.Time = date
	return
}

type reservationDuration struct {
	start reservationDate
	end   reservationDate
}

type reservationSlot struct {
	Config slotConfig          `json:"config"`
	Date   reservationDuration `json:"date"`
}

type foundVenue struct {
	Venue Venue             `json:"venue"`
	Slots []reservationSlot `json:"slots"`
}

type findResult struct {
	Venues []foundVenue `json:"venues"`
}

type findResponse struct {
	Results    findResult `json:"results"`
	GuestToken string     `json:"guest_token"`
}

/*
 * Reservation endpoints responses
 */

type reservationBookToken struct {
	Expires ReservationTime `json:"expires"`
	Value   string          `json:"value"`
}

type reservationDetailsResponse struct {
	BookToken reservationBookToken `json:"book_token"`
}

type reservationResponse struct {
	ResyToken string `json:"resy_token"`
}
