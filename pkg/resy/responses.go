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
	// TODO(name)
	Name     string   `json:"name"`
	IDs      VenueIDs `json:"id"`
	Location Location `json:"location"`
}

type SlotConfig struct {
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

type ReservationDuration struct {
	Start reservationDate
	End   reservationDate
}

type ReservationSlot struct {
	Config SlotConfig          `json:"config"`
	Date   ReservationDuration `json:"date"`
}

type FoundVenue struct {
	Venue Venue             `json:"venue"`
	Slots []ReservationSlot `json:"slots"`
}

type FindResult struct {
	Venues []FoundVenue `json:"venues"`
}

type FindResponse struct {
	Results    FindResult `json:"results"`
	GuestToken string     `json:"guest_token"`
}

/*
 * Reservation endpoints responses
 */

type ReservationTime struct {
	time.Time
}

type ReservationBookToken struct {
	Expires ReservationTime `json:"expires"`
	Value   string          `json:"value"`
}

type PaymentDetails struct {
}

type UserPayment struct {
	PaymentMethods []UserPaymentMethod
}

type UserPaymentMethod struct {
	ID        int    `json:"id"`
	IsDefault bool   `json:"is_default"`
	Type      string `json:"type"`
}

type CancellationFeeDetails struct {
	Amount int `json:"amount"`
}

type CancellationDetails struct {
	Fee *CancellationFeeDetails `json:"fee"`
}

type ReservationDetailsResponse struct {
	BookToken    ReservationBookToken `json:"book_token"`
	Payment      PaymentDetails       `json:"payment"`
	User         UserPayment          `json:"user"`
	Cancellation CancellationDetails  `json:"cancellation"`
}

type ReservationResponse struct {
	ResyToken string `json:"resy_token"`
}
