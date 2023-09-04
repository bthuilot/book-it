package resy

import (
	"net/url"
)

// passwordLength is the max length allowed for a password.
// when resy authenticates on their frontend only the first 25
// characters are sent, security I guess...
const passwordLength = 25

func (c *client) passwordAuth() (authStorage, error) {
	data := url.Values{}
	data.Set("email", c.authStorage.email)
	if len(c.authStorage.password) > passwordLength {
		c.authStorage.password = c.authStorage.password[:passwordLength]
	}
	data.Set("password", c.authStorage.password)

	reqURL, err := generateURL("3/auth/password", nil)
	if err != nil {
		return authStorage{}, err
	}

	req, err := generateRequest("POST", reqURL, data, nil)
	if err != nil {
		return authStorage{}, err
	}

	resp, err := apiRequest[Auth](req)
	if err != nil {
		return authStorage{}, err
	}
	return newAuthToken(resp.Token)
}

//
//func (c *client) smsAuth() (authStorage, error) {
//	data := url.Values{}
//	data.Set("mobile_number", c.phone)
//	data.Set("method", "sms")
//	if _, err := apiRequest[ValidationCode]("POST", "3/auth/mobile", "", url.Values{}, data); err != nil {
//		return authStorage{}, err
//	}
//	var code string
//	fmt.Println("enter code:")
//	if _, err := fmt.Scanf("%s", &code); err != nil {
//		return authStorage{}, err
//	}
//
//	data = url.Values{}
//	data.Set("mobile_number", c.phone)
//	data.Set("code", code)
//	challenge, err := apiRequest[Challenge]("POST", "3/auth/mobile", "", url.Values{}, data)
//	if err != nil {
//		return authStorage{}, err
//	}
//
//	data = url.Values{}
//	data.Set("em_address", c.email)
//	data.Set("challenge_id", challenge.Challenge.ChallengeId)
//	auth, err := apiRequest[authResponse]("POST", "3/auth/challenge", "", url.Values{}, data)
//	if err != nil {
//		return authStorage{}, err
//	}
//	fmt.Println(auth.Token)
//	token, err := jwt.Parse(auth.Token, func(token *jwt.Token) (interface{}, error) { return nil, nil })
//	if token == nil {
//		fmt.Println(token)
//		return authStorage{}, err
//	}
//
//	c.authStorage = authStorage{
//		raw: auth.Token,
//		jwt: token,
//	}
//	return authStorage{}, nil
//}

/*********
 * Types *
 *********/

type ValidationCode struct {
	MobileNumber string `json:"mobile_number"`
	Method       int    `json:"method"`
	Sent         bool   `json:"sent"`
}

type Challenge struct {
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

type Auth struct {
	Id                     int             `json:"id"`
	FirstName              string          `json:"first_name"`
	LastName               string          `json:"last_name"`
	MobileNumber           string          `json:"mobile_number"`
	EmAddress              string          `json:"em_address"`
	EmIsVerified           int             `json:"em_is_verified"`
	MobileNumberIsVerified int             `json:"mobile_number_is_verified"`
	IsActive               int             `json:"is_active"`
	ReferralCode           string          `json:"referral_code"`
	IsMarketable           int             `json:"is_marketable"`
	IsConcierge            int             `json:"is_concierge"`
	DateUpdated            int             `json:"date_updated"`
	DateCreated            int             `json:"date_created"`
	HasSetPassword         int             `json:"has_set_password"`
	NumBookings            int             `json:"num_bookings"`
	PaymentMethods         []PaymentMethod `json:"payment_methods"`
	ResySelect             int             `json:"resy_select"`
	ProfileImageUrl        string          `json:"profile_image_url"`
	PaymentMethodId        int             `json:"payment_method_id"`
	PaymentProviderId      int             `json:"payment_provider_id"`
	PaymentProviderName    string          `json:"payment_provider_name"`
	PaymentDisplay         string          `json:"payment_display"`
	Token                  string          `json:"token"`
	LegacyToken            string          `json:"legacy_token"`
}

type PaymentMethod struct {
	Id           int    `json:"id"`
	IsDefault    bool   `json:"is_default"`
	ProviderId   int    `json:"provider_id"`
	ProviderName string `json:"provider_name"`
	Display      string `json:"display"`
	Type         string `json:"type"`
	ExpYear      int    `json:"exp_year"`
	ExpMonth     int    `json:"exp_month"`
	IssuingBank  string `json:"issuing_bank"`
}
