package resy

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"
)

// GetReservationDetails will return the details of a reservation for a given party size, day,
// and config token (from a previous FindReservation call).
func (c *client) GetReservationDetails(partySize int, day time.Time, configToken string) (v ReservationDetails, err error) {
	data := url.Values{}
	data.Set("party_size", strconv.Itoa(partySize))
	data.Set("day", day.Format("2006-01-02"))
	data.Set("config_id", configToken)
	reqURL, err := generateURL("3/details", data)
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
	return apiRequest[ReservationDetails](req)
}

// BookReservation will book a reservation for a given book token. If usePayment is true, the
// payment method with the given ID will be used. If usePayment is false, the default payment
// method will be used.
func (c *client) BookReservation(bookToken string, usePayment bool, paymentMethodID int) (v Reservation, err error) {
	data := url.Values{}
	data.Set("book_token", bookToken)
	if usePayment {
		payment, _ := json.Marshal(map[string]int{"id": paymentMethodID})
		data.Set("struct_payment_method", string(payment))
	}
	data.Set("source_id", "resy.com-venue-details")
	reqURL, err := generateURL("3/book", nil)
	if err != nil {
		return
	}
	additionalHeaders, err := c.generateAuthHeaders()
	if err != nil {
		return
	}
	req, err := generateRequest("POST", reqURL, data, additionalHeaders)
	if err != nil {
		return
	}
	return apiRequest[Reservation](req)
}

/*********
 * Types *
 *********/

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
	Amount float64 `json:"amount"`
}

type CancellationDetails struct {
	Fee *CancellationFeeDetails `json:"fee"`
}

type ReservationDetails struct {
	BookToken    ReservationBookToken `json:"book_token"`
	Payment      PaymentDetails       `json:"payment"`
	User         UserPayment          `json:"user"`
	Cancellation CancellationDetails  `json:"cancellation"`
}

type Reservation struct {
	ResyToken string `json:"resy_token"`
}
