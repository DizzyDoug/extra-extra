package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Connector is the connector to Github. It handles authorization and rate limits.
type Connector struct {
	authToken         string
	remainingRequests int
}

// NewConnector is the constructor for Connector
func NewConnector(authToken string) Connector {
	return Connector{
		authToken: authToken,
	}
}

// Get handles authentication and ratze limits of github
func (c *Connector) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req, false)
}

func (c *Connector) getWithoutRateCheck(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req, true)
}

func (c *Connector) do(req *http.Request, skipRateLimitCheck bool) (*http.Response, error) {
	if !skipRateLimitCheck {
		if err := c.checkRateLimit(); err != nil {
			return nil, err
		}
	}

	req.Header.Add("Authorization", fmt.Sprintf("token %s", c.authToken))
	return http.DefaultClient.Do(req)
}

type rateLimitAPIResp struct {
	Resources struct {
		Core struct {
			Limit     int `json:"limit"`
			Remaining int `json:"remaining"`
			Reset     int `json:"reset"`
		} `json:"core"`
		Search struct {
			Limit     int `json:"limit"`
			Remaining int `json:"remaining"`
			Reset     int `json:"reset"`
		} `json:"search"`
		Graphql struct {
			Limit     int `json:"limit"`
			Remaining int `json:"remaining"`
			Reset     int `json:"reset"`
		} `json:"graphql"`
		IntegrationManifest struct {
			Limit     int `json:"limit"`
			Remaining int `json:"remaining"`
			Reset     int `json:"reset"`
		} `json:"integration_manifest"`
	} `json:"resources"`
	Rate struct {
		Limit     int `json:"limit"`
		Remaining int `json:"remaining"`
		Reset     int `json:"reset"`
	} `json:"rate"`
}

func (c *Connector) getRateLimit() (rlr rateLimitAPIResp, err error) {
	res, err := c.getWithoutRateCheck("https://api.github.com/rate_limit")
	if err != nil {
		return rlr, err
	}
	defer res.Body.Close()

	bz, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return rlr, err
	}

	err = json.Unmarshal(bz, &rlr)
	return rlr, err
}

func (c *Connector) checkRateLimit() error {
	if c.remainingRequests <= 0 {
		rl, err := c.getRateLimit()
		if err != nil {
			return err
		}
		if rl.Resources.Core.Remaining == 0 {
			return errors.New("Github API rate limit exceeded")
		}

		c.remainingRequests = rl.Resources.Core.Remaining - 1
		return nil
	}

	c.remainingRequests--
	return nil
}
