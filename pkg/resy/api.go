package resy

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type APIClient interface {
	// GetVenue will return a Venue by its id
	GetVenue(id int) (v Venue, err error)
	// FindReservation will find all reservations for a given day
	FindReservation(venue Venue, day time.Time, partySize int) (found FindResponse, err error)
	// GetReservationDetails will return the reservation details for a given reservation config token
	GetReservationDetails(partySize int, day time.Time, configToken string) (v ReservationDetailsResponse, err error)
	// BookReservation will book a reservation from its bookToken
	BookReservation(bookToken string) (v ReservationResponse, err error)
}

type client struct {
	email     string
	password  string
	authToken authToken
}

const (
	apiURL = "https://api.resy.com"

	// apiKey is a token to be included with all API requests
	// As far as I can tell this is a hard coded value found here
	// https://resy.com/modules/app.2624781f09c5841e389c.js - Line 14
	apiKey = "VbWk7s3L4KiK5fzlO7JD3Q5EYolJI7n5"
	// resultLimit
	resultLimit = 20
)

func NewClient() (APIClient, error) {
	c := &client{
		email:    viper.GetString("email"),
		password: viper.GetString("password"),
	}
	err := c.refreshToken()
	return c, err
}

func apiRequest[T any](method, endpoint, authToken string, queryParams, formParams url.Values) (resp T, err error) {
	var (
		httpReq *http.Request
		httpRes *http.Response
		body    []byte
	)
	fullURL := fmt.Sprintf("%s/%s", apiURL, endpoint)
	if len(queryParams) > 0 {
		fullURL = fmt.Sprintf("%s?%s", fullURL, queryParams.Encode())
	}
	if httpReq, err = http.NewRequest(method, fullURL, strings.NewReader(formParams.Encode())); err != nil {
		return
	}
	httpReq.Header.Add("Authorization", fmt.Sprintf("ResyAPI api_key=\"%s\"", apiKey))
	httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Add("Host", "api.resy.com")
	httpReq.Header.Add("Origin", "https://widgets.resy.com")

	if authToken != "" {
		httpReq.Header.Add("X-Resy-Auth-Token", authToken)
		httpReq.Header.Add("X-Resy-Universal-Auth", authToken)
	}
	logrus.Debugf("making '%s' API request to '%s'", method, fullURL)
	headers, _ := json.Marshal(httpReq.Header)
	logrus.Debugf("headers: %s", headers)
	logrus.Debugf("body: %s", formParams.Encode())
	if httpRes, err = http.DefaultClient.Do(httpReq); err != nil {
		logrus.Debugf("error sending request: %s", err)
		return
	} else if httpRes.StatusCode != http.StatusOK && httpRes.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(httpRes.Body)
		logrus.Debugf("API endpoint returned %d not 200: %s", httpRes.StatusCode, string(b))
		err = fmt.Errorf("non 200 status returned: %s", string(b))
		return
	}

	if body, err = io.ReadAll(httpRes.Body); err != nil {
		logrus.Fatalf("unable to read body: %s", err)
		return
	}
	logrus.Debugf("recieved response from API: %s", string(body))
	err = json.Unmarshal(body, &resp)
	logrus.Debugf("received response: %+v", resp)
	return
}

func (c *client) GetVenue(id int) (v Venue, err error) {
	data := url.Values{}
	data.Set("id", strconv.Itoa(id))
	v, err = apiRequest[Venue]("GET", "3/venue", "", data, url.Values{})
	return
}

func (c *client) FindReservation(venue Venue, day time.Time, partySize int) (found FindResponse, err error) {
	if !c.authToken.isValid() {
		logrus.Debug("auth token is invalid, refreshing")
		if err = c.refreshToken(); err != nil {
			err = fmt.Errorf("unbale to refresh auth token: %s", err)
			return
		}
	}

	data := url.Values{}
	data.Set("lat", fmt.Sprintf("%f", venue.Location.Lat))
	data.Set("long", fmt.Sprintf("%f", venue.Location.Long))
	data.Set("party_size", strconv.Itoa(partySize))
	data.Set("day", day.Format("2006-01-02"))
	data.Set("venue_id", strconv.Itoa(venue.IDs.Resy))
	data.Set("limit", strconv.Itoa(resultLimit))

	found, err = apiRequest[FindResponse]("GET", "4/find", c.authToken.raw, data, url.Values{})
	return
}

func (c *client) GetReservationDetails(partySize int, day time.Time, configToken string) (v ReservationDetailsResponse, err error) {
	if !c.authToken.isValid() { // TODO(check expiry)
		logrus.Debug("auth token is invalid, refreshing")
		if err = c.refreshToken(); err != nil {
			err = fmt.Errorf("unbale to refresh auth token: %s", err)
			return
		}
	}

	data := url.Values{}
	data.Set("party_size", strconv.Itoa(partySize))
	data.Set("day", day.Format("2006-01-02"))
	data.Set("config_id", configToken)

	v, err = apiRequest[ReservationDetailsResponse]("GET", "3/details", c.authToken.raw, data, url.Values{})
	return
}

func (c *client) BookReservation(bookToken string) (v ReservationResponse, err error) {
	if !c.authToken.isValid() { // TODO(check expiry)
		logrus.Debug("auth token is invalid, refreshing")
		if err = c.refreshToken(); err != nil {
			err = fmt.Errorf("unbale to refresh auth token: %s", err)
			return
		}
	}

	data := url.Values{}
	data.Set("book_token", bookToken)
	// TODO(reverse engineer how payment method is retrieved)
	data.Set("struct_payment_method", "{\"id\": 9313397}")
	data.Set("source_id", "resy.com-venue-details")
	v, err = apiRequest[ReservationResponse]("POST", "3/book", c.authToken.raw, url.Values{}, data)
	return
}
