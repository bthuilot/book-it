package resy

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Client interface {
	// GetVenue will return a Venue by its id
	GetVenue(id int) (v Venue, err error)
	// FindReservation will find all reservations for a given day
	FindReservation(venue Venue, day time.Time, partySize int) (found Search, err error)
	// GetReservationDetails will return the reservation details for a given reservation config token
	GetReservationDetails(partySize int, day time.Time, configToken string) (v ReservationDetails, err error)
	// BookReservation will book a reservation from its bookToken
	BookReservation(bookToken string, sendPaymentMethod bool, paymentMethodID int) (v Reservation, err error)
}

type client struct {
	authStorage authStorage
}

const (
	// apiURL is the base URL for all API requests
	apiURL = "https://api.resy.com"

	// apiKey is a token to be included with all API requests
	// As far as I can tell this is a hard coded value which can be found
	// be running `apiConfig.Config.apiKey` in the Resy web app console
	apiKey = "VbWk7s3L4KiK5fzlO7JD3Q5EYolJI7n5"
	// resultLimit
	resultLimit = 20
)

// NewClient will return a new Client with the given options
// (see the With* functions for available options)
func NewClient(opts ...Opts) (Client, error) {
	c := &client{}
	for _, opt := range opts {
		opt.apply(c)
	}
	var err error
	c.authStorage, err = c.passwordAuth()
	return c, err
}

/********************
 * Option Functions *
 ********************/

// Opts is an interface which can be used to set options on the client
// (see the With* functions for available options)
type Opts interface {
	apply(client *client)
}

// WithCredentialsOpts will return an Opts which will set the email and password
// for the client
func WithCredentialsOpts(email, password string) Opts {
	return withCredentialsOpts{
		Email:    email,
		Password: password,
	}
}

type withCredentialsOpts struct {
	Email    string
	Password string
}

func (co withCredentialsOpts) apply(client *client) {
	client.authStorage = authStorage{
		email:    co.Email,
		password: co.Password,
	}
}

/***********************************
 * Internal authentication storage *
 ***********************************/

type authStorage struct {
	email    string
	password string
	raw      string
	jwt      *jwt.Token
}

func (a authStorage) isValid() bool {
	if a.jwt == nil {
		return false
	}

	if claims, ok := a.jwt.Claims.(jwt.MapClaims); ok {
		return claims.Valid() == nil
	}

	return false
}

func newAuthToken(raw string) (authStorage, error) {
	token, err := jwt.Parse(raw, func(token *jwt.Token) (interface{}, error) { return nil, nil })
	if token == nil {
		return authStorage{}, err
	}

	return authStorage{
		raw: raw,
		jwt: token,
	}, nil
}

func (c *client) generateAuthHeaders() (res map[string]string, err error) {
	if !c.authStorage.isValid() {
		c.authStorage, err = c.passwordAuth()
	}
	res = map[string]string{
		"X-Resy-Auth-Token": c.authStorage.raw,
		//"X-Resy-Universal-Auth": authStorage,
	}
	return
}
