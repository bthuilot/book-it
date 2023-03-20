package resy

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/url"
)

type authToken struct {
	raw string
	jwt *jwt.Token
}

func (a authToken) isValid() bool {
	if a.jwt == nil {
		return false
	}

	if claims, ok := a.jwt.Claims.(jwt.MapClaims); ok {
		return claims.Valid() == nil
	}

	return false
}

func newAuthToken(raw string) (authToken, error) {
	token, err := jwt.Parse(raw, func(token *jwt.Token) (interface{}, error) { return nil, nil })
	if token == nil {
		return authToken{}, err
	}

	return authToken{
		raw: raw,
		jwt: token,
	}, nil
}

func (c *client) refreshToken() (err error) {
	c.authToken, err = c.passwordAuth()
	return
}

func (c *client) passwordAuth() (authToken, error) {
	data := url.Values{}
	data.Set("email", c.email)
	data.Set("password", c.password[:25])
	resp, err := makeResyAPIRequest[authResponse]("POST", "3/auth/password", c.apiKey, "", url.Values{}, data)
	if err != nil {
		return authToken{}, err
	}
	return newAuthToken(resp.Token)
}

func (c *client) smsAuth() (authToken, error) {
	data := url.Values{}
	data.Set("mobile_number", c.phone)
	data.Set("method", "sms")
	if _, err := makeResyAPIRequest[authPhone]("POST", "3/auth/mobile", c.apiKey, "", url.Values{}, data); err != nil {
		return authToken{}, err
	}
	var code string
	fmt.Println("enter code:")
	if _, err := fmt.Scanf("%s", &code); err != nil {
		return authToken{}, err
	}

	data = url.Values{}
	data.Set("mobile_number", c.phone)
	data.Set("code", code)
	challenge, err := makeResyAPIRequest[challengeResponse]("POST", "3/auth/mobile", c.apiKey, "", url.Values{}, data)
	if err != nil {
		return authToken{}, err
	}

	data = url.Values{}
	data.Set("em_address", c.email)
	data.Set("challenge_id", challenge.Challenge.ChallengeId)
	auth, err := makeResyAPIRequest[authResponse]("POST", "3/auth/challenge", c.apiKey, "", url.Values{}, data)
	if err != nil {
		return authToken{}, err
	}
	fmt.Println(auth.Token)
	token, err := jwt.Parse(auth.Token, func(token *jwt.Token) (interface{}, error) { return nil, nil })
	if token == nil {
		fmt.Println(token)
		return authToken{}, err
	}

	c.authToken = authToken{
		raw: auth.Token,
		jwt: token,
	}
	return authToken{}, nil
}
