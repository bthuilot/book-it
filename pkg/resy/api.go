package resy

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func generateURL(endpoint string, queryParams url.Values) (string, error) {
	parsedURL, err := url.Parse(apiURL)
	if err != nil {
		return "", err
	}
	parsedURL = parsedURL.JoinPath(endpoint)
	parsedURL.RawQuery = queryParams.Encode()
	return parsedURL.String(), nil
}

func generateRequest(method, url string, formBody url.Values, additionalHeaders map[string]string) (req *http.Request, err error) {
	if req, err = http.NewRequest(method, url, strings.NewReader(formBody.Encode())); err != nil {
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("ResyAPI api_key=\"%s\"", apiKey))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Host", "api.resy.com")
	req.Header.Add("Origin", "https://widgets.resy.com")
	for k, v := range additionalHeaders {
		req.Header.Add(k, v)
	}

	return
}

func apiRequest[T any](r *http.Request) (apiRes T, err error) {
	log := logrus.WithFields(logrus.Fields{
		"method":   r.Method,
		"endpoint": r.URL.Path,
	})

	log.Infof("sending request to API: %s", r.URL.String())
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Errorf("error sending request: %s", err)
		return
	}

	defer res.Body.Close()
	switch res.StatusCode {
	case http.StatusOK, http.StatusCreated:
		if err = json.NewDecoder(res.Body).Decode(&apiRes); err != nil {
			log.Errorf("error decoding response: %s", err)
		}
	default:
		b, _ := io.ReadAll(res.Body)
		log.Errorf("API endpoint returned %d not 200: %s", res.StatusCode, string(b))
		err = fmt.Errorf("%w: non 200 status returned - '%s'", ErrAPI, string(b))
	}
	return
}
